// Package casbin provides an in-memory, thread-safe RBAC policy engine
// built on top of Casbin. It is designed as a standalone IAM (Identity & Access Management)
// library with clean interface specifications, configuration structs, and Redis Pub/Sub synchronization.
//
// Pattern:
//
//	// 1. Configure
//	cfg := casbin.Config{
//	    SyncChannel: "auth:casbin:sync",
//	    SyncTimeout: 10 * time.Second,
//	}
//
//	// 2. Initialize (returns IAM / Enforcer interface)
//	iam, err := casbin.New(cfg, casbin.WithRedis(redisClient))
//	if err != nil { ... }
//
//	// 3. Bulk load policies
//	err = iam.LoadPolicies(groupings, policies)
//
//	// 4. Check permissions (in-memory, sub-microsecond)
//	allowed, err := iam.Enforce("user_123", "documents", "read")
//
//	// 5. Dynamic mutations
//	iam.SetUserRoles("user_123", []string{"role_admin"})
//	iam.SetRolePolicies("role_admin", []casbin.Policy{
//	    {Resource: "documents", Action: "*"},
//	})
//
//	// 6. Cross-instance sync
//	iam.PublishSync(ctx)
package casbin
