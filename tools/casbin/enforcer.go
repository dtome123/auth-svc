package casbin

import (
	"context"
	"fmt"

	casbinLib "github.com/casbin/casbin/v2"
)

// impl implements the IAM and Enforcer interfaces.
type impl struct {
	cfg   Config
	inner *casbinLib.SyncedEnforcer
	sync  *syncer
}

// New creates and initializes a new IAM enforcer instance conforming to the IAM / Enforcer spec.
//
// Example:
//
//	iam, err := casbin.New(casbin.DefaultConfig(), casbin.WithRedis(redisClient))
func New(cfg Config, opts ...Option) (IAM, error) {
	cfg = cfg.withDefaults()

	inner, err := casbinLib.NewSyncedEnforcer(newRBACModel())
	if err != nil {
		return nil, fmt.Errorf("casbin: failed to initialize enforcer: %w", err)
	}

	instance := &impl{
		cfg:   cfg,
		inner: inner,
	}

	for _, opt := range opts {
		opt(instance)
	}

	return instance, nil
}

// NewIAM is an alias for New.
func NewIAM(cfg Config, opts ...Option) (IAM, error) {
	return New(cfg, opts...)
}

// ---------------------------------------------------------------------------
// Query
// ---------------------------------------------------------------------------

// Enforce checks whether (subject, resource, action) is permitted.
// Matches are performed strictly in RAM with sub-microsecond latency.
func (i *impl) Enforce(subject, resource, action string) (bool, error) {
	return i.inner.Enforce(subject, resource, action)
}

// EnforceRoles checks whether any of the given roles is permitted to perform (resource, action).
// Short-circuits immediately on the first allowed role.
func (i *impl) EnforceRoles(roles []string, resource, action string) (bool, error) {
	for _, role := range roles {
		allowed, err := i.inner.Enforce(role, resource, action)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// Bulk Load
// ---------------------------------------------------------------------------

// LoadPolicies atomically replaces all in-memory grouping and permission rules.
func (i *impl) LoadPolicies(groupings, policies [][]string) error {
	i.inner.ClearPolicy()

	if len(groupings) > 0 {
		if _, err := i.inner.AddGroupingPolicies(groupings); err != nil {
			return fmt.Errorf("casbin: failed to load grouping policies: %w", err)
		}
	}
	if len(policies) > 0 {
		if _, err := i.inner.AddPolicies(policies); err != nil {
			return fmt.Errorf("casbin: failed to load permission policies: %w", err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Dynamic Mutations
// ---------------------------------------------------------------------------

// SetUserRoles replaces all role assignments for a user.
func (i *impl) SetUserRoles(userID string, roleIDs []string) error {
	if _, err := i.inner.RemoveFilteredGroupingPolicy(0, userID); err != nil {
		return fmt.Errorf("casbin: remove user roles: %w", err)
	}
	if len(roleIDs) == 0 {
		return nil
	}
	groupings := make([][]string, len(roleIDs))
	for idx, rID := range roleIDs {
		groupings[idx] = []string{userID, rID}
	}
	if _, err := i.inner.AddGroupingPolicies(groupings); err != nil {
		return fmt.Errorf("casbin: add user roles: %w", err)
	}
	return nil
}

// SetRolePolicies replaces all permission policies for a role.
func (i *impl) SetRolePolicies(roleID string, policies []Policy) error {
	if _, err := i.inner.RemoveFilteredPolicy(0, roleID); err != nil {
		return fmt.Errorf("casbin: remove role policies: %w", err)
	}
	if len(policies) == 0 {
		return nil
	}
	rules := make([][]string, len(policies))
	for idx, p := range policies {
		rules[idx] = []string{roleID, p.Resource, p.Action}
	}
	if _, err := i.inner.AddPolicies(rules); err != nil {
		return fmt.Errorf("casbin: add role policies: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Sync
// ---------------------------------------------------------------------------

// PublishSync broadcasts a policy-reload event across all service instances via Redis Pub/Sub.
func (i *impl) PublishSync(ctx context.Context) {
	if i.sync != nil {
		i.sync.publish(ctx)
	}
}

// NotifySync is an alias for PublishSync.
func (i *impl) NotifySync(ctx context.Context) {
	i.PublishSync(ctx)
}

// OnSyncEvent registers a callback that is invoked when a remote reload event is received.
func (i *impl) OnSyncEvent(fn func(ctx context.Context) error) {
	if i.sync == nil {
		return
	}
	i.sync.onEvent = fn
	go i.sync.subscribe(context.Background(), i.cfg)
}
