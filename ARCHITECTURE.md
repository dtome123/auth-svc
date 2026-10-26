# AI Agent Context & Architecture Index (`auth-svc`)

> **Note for AI Coding Agents**: This document serves as an explicit context map of the `auth-svc` repository to accelerate codebase comprehension, eliminate repetitive directory scanning, and enforce architecture consistency when planning or implementing changes.

---

## 1. Project Overview & Tech Stack

* **Service Name**: `auth-svc`
* **Language & Runtime**: Go (Golang)
* **Core Responsibilities**: Authentication (JWT management, session tracking) and Authorization (RBAC/ABAC permission evaluation).
* **Protocols**:
  * **gRPC**: Primary inter-service communication protocol for internal microservices.
  * **HTTP REST**: Public/Edge endpoints for token management.
* **Storage & Caching**:
  * **MongoDB**: Persistent database (`sessions`, `roles`, `permissions`, `assignments`).
  * **Flexible Cache Layer**: Supports 3 caching options configured via `caching.type`:
    * `none`: Disabled cache (No-Op).
    * `memory`: Thread-safe In-Memory RAM cache.
    * `redis`: External Redis cache.
* **Key Dependencies**:
  * `github.com/dtome123/auth-sdk`: Proprietary SDK for JWT verification/signing (RSA/HMAC) and replay protection.
  * `github.com/dtome123/go-mongo-generic`: Mongo DB wrapper.
  * `github.com/redis/go-redis/v9`: Redis client (initialized conditionally).

---

## 2. Architectural Design (Hexagonal / Ports and Adapters)

The project strictly adheres to **Ports and Adapters (Hexagonal Architecture)**. Code within `internal/services` contains pure business logic and does not depend on transport layers or specific DB drivers.

```text
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                 PRIMARY PORTS                                   │
│    gRPC Handlers (internal/port/grpc/)    │   REST Handlers (internal/port/rest/)   │
└────────────────────────────────────────┬────────────────────────────────────────┘
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                             CORE SERVICE LAYER                                  │
│             services.Service  -->  auth.AuthorizationService                    │
└────────────────────────────────────────┬────────────────────────────────────────┘
                                         │
                                         ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                               SECONDARY ADAPTERS                                │
│  MongoDB Repo (internal/adapters/mongodb/) │ Redis Cache (internal/adapters/redis)│
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Comprehensive File & Package Map

### 📍 Root & Entrypoint
* `main.go`: Application bootstrap. Loads config, initializes MongoDB and Redis connections, sets up M2M verifiers/signers and per-client user token verifiers/signers (`ClientEntry`), and invokes `port.NewServer()`.
* `config/config.go` & `config/config.yaml`: Configuration structures for MongoDB DSN, Redis credentials, OAuth server keys, default user token settings (`default_user_token_config`), and M2M/User client definitions (`clients`).

### 🔌 Tầng Port (Primary Adapters - Input)
* `internal/port/new.go`: Server supervisor (`Server.Run()`) launching gRPC and REST servers concurrently. Initializes server-level OAuth verifiers and signers.
* `internal/port/grpc/`: gRPC request handlers.
  * `token.go`, `refresh.go`, `sign.go`: Token lifecycle handlers.
  * `check.go`, `check_permissions.go`, `get_user_permission.go`, `get_user_role.go`, `has_permission.go`: Authorization check endpoints.
  * `create_role.go`, `update_role.go`, `list_roles.go`, `assign_role_to_user.go`: RBAC management endpoints.
* `internal/port/rest/`: HTTP REST request handlers (`token.go`, `new.go`).

### 🧠 Core Business Logic (Application Layer)
* `internal/services/service.go`: Orchestrator linking DB repositories to domain services.
* `internal/services/auth/`: Central authorization and authentication engine (`AuthorizationService`).
  * `token.go`: Handles M2M OAuth Client Credentials token generation using Server OAuth keys. Supports wildcard `"*"` audience authorization checks.
  * `sign.go` & `refresh.go`: Handles End-User authentication and token refresh using client-specific (or fallback default) `UserSigner` and custom `AccessTokenTTL` / `RefreshTokenTTL`.
  * `has_permission.go`, `check_permissions.go`, `get_user_permission.go`: Permission resolution engine with Redis cache integration.
  * `assign_role.go`, `create_role.go`, `update_role.go`, `upsert_permission.go`: RBAC mutations.

### 💾 Adapters (Secondary Adapters - Output & Persistence)
* `internal/adapters/cache/`: Unified caching module containing interface (`interface.go`) and implementations (`noop.go`, `memory.go`, `redis.go`).
* `internal/adapters/redis/`:
  * `authen/`: Redis persistence for user sessions (`new.go`, `session.go`).
* `internal/adapters/mongodb/`:
  * `author/`: MongoDB persistence for RBAC domain models (`create_role.go`, `list_roles.go`, `batch_upsert_permissions.go`, `upsert_assignment.go`).

### 📦 Domain Models & DTOs
* `internal/domain/`:
  * `session.go`: `Session` domain model (`user_id`, `device_id`, `access_token_hash`, `refresh_token_hash`, `ttl`).
  * `role.go`: `Role` domain model (`name`, `resource_actions`, `permission_ids`, `scope`, `description`).
  * `permission.go`: `Permission` domain model (`name`, `domain`, `resource`, `actions`).
  * `assignment.go`: User role assignment domain model.
* `internal/types/`: DTOs, OAuth client entries (`ClientEntry` storing `Verifier`, `Signer`, `UserSigner`, `UserVerifier`, and custom TTLs), route scopes, auth mechanism types (`AuthM2MTypeRSA`, `AuthM2MTypeHMAC`).
* `internal/utils/`: Helper utilities (`hash.go`, `token.go`, `time.go`, `page.go`).

---

## 4. Key Execution Workflows

### 🔑 Token Issuance & Refresh
1. Request arrives at `port/grpc/token.go`, `sign.go`, or `refresh.go`.
2. Evaluated by `services/auth/`.
   * **M2M Tokens (`token.go`)**: Validates client assertion and signs token via `svc.serverSigner`.
   * **User Tokens (`sign.go` & `refresh.go`)**: Resolves calling client's `UserSigner` and dynamic `accessTTL` / `refreshTTL`.
3. Session and token hashes are persisted to Redis with TTL via `adapters/redis/authen/session.go`.

### 🛡️ Permission Check & Caching Policy
1. Request arrives at `port/grpc/check_permissions.go`.
2. Handled by `services/auth/check_permissions.go`.
3. **Cache Look-aside Pattern**:
   * Query Redis via `adapters/redis/author`.
   * **On Cache Hit**: Return permissions immediately.
   * **On Cache Miss**: Query MongoDB via `adapters/mongodb/author`, calculate user permissions/roles, populate Redis with TTL, and return result.

---

## 5. Guidelines for AI Agents Modifying Code

1. **Maintain Layer Separation**: Do not import MongoDB or Redis packages directly in `internal/port/` or `internal/services/auth/`. All data access must pass through repository interfaces.
2. **Preserve Error Handling**: Panic recovery is handled in `internal/port/new.go`. Domain errors should be cleanly returned as typed errors or gRPC status codes.
3. **Configuration Inheritance**: When adding client configuration fields, ensure validation rules in `config/config.go` properly handle inheritance from `default_user_token_config`.
4. **Update Index On Structural Changes**: Keep this document updated whenever changing configuration layouts, domain models, or key workflows.
