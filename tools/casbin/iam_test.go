package casbin

import (
	"sync"
	"testing"
)

func TestIAM_BasicAndWildcards(t *testing.T) {
	iam, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("failed to create IAM enforcer: %v", err)
	}

	groupings := [][]string{
		{"alice", "admin"},
		{"bob", "editor"},
		{"charlie", "viewer"},
	}

	policies := [][]string{
		{"admin", "*", "*"},
		{"editor", "articles", "read"},
		{"editor", "articles", "write"},
		{"viewer", "articles", "read"},
	}

	if err := iam.LoadPolicies(groupings, policies); err != nil {
		t.Fatalf("failed to load policies: %v", err)
	}

	cases := []struct {
		sub, obj, act string
		expected      bool
	}{
		// Admin with wildcard (*, *)
		{"alice", "articles", "read", true},
		{"alice", "articles", "write", true},
		{"alice", "articles", "delete", true},
		{"alice", "billing", "manage", true},

		// Editor permissions
		{"bob", "articles", "read", true},
		{"bob", "articles", "write", true},
		{"bob", "articles", "delete", false},
		{"bob", "billing", "read", false},

		// Viewer permissions
		{"charlie", "articles", "read", true},
		{"charlie", "articles", "write", false},
		{"charlie", "articles", "delete", false},

		// Unknown user
		{"stranger", "articles", "read", false},
	}

	for _, c := range cases {
		allowed, err := iam.Enforce(c.sub, c.obj, c.act)
		if err != nil {
			t.Errorf("Enforce(%q, %q, %q) returned error: %v", c.sub, c.obj, c.act, err)
		}
		if allowed != c.expected {
			t.Errorf("Enforce(%q, %q, %q) = %v, expected %v", c.sub, c.obj, c.act, allowed, c.expected)
		}
	}
}

func TestIAM_DynamicMutations(t *testing.T) {
	iam, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("failed to create IAM enforcer: %v", err)
	}

	// 1. Initial state
	allowed, _ := iam.Enforce("user1", "reports", "export")
	if allowed {
		t.Fatal("expected false for empty policies")
	}

	// 2. Set role policies
	err = iam.SetRolePolicies("analyst", []Policy{
		{Resource: "reports", Action: "view"},
		{Resource: "reports", Action: "export"},
	})
	if err != nil {
		t.Fatalf("SetRolePolicies failed: %v", err)
	}

	// 3. Assign role to user
	err = iam.SetUserRoles("user1", []string{"analyst"})
	if err != nil {
		t.Fatalf("SetUserRoles failed: %v", err)
	}

	// 4. Verify permission granted
	allowed, _ = iam.Enforce("user1", "reports", "export")
	if !allowed {
		t.Error("expected user1 to have reports:export")
	}

	// 5. Update user roles (remove analyst, assign guest)
	err = iam.SetUserRoles("user1", []string{"guest"})
	if err != nil {
		t.Fatalf("SetUserRoles failed: %v", err)
	}

	allowed, _ = iam.Enforce("user1", "reports", "export")
	if allowed {
		t.Error("expected user1 to no longer have reports:export")
	}
}

func TestIAM_ConcurrentEnforce(t *testing.T) {
	iam, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("failed to create IAM enforcer: %v", err)
	}

	_ = iam.LoadPolicies(
		[][]string{{"userA", "roleA"}},
		[][]string{{"roleA", "resA", "actA"}},
	)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				allowed, err := iam.Enforce("userA", "resA", "actA")
				if err != nil || !allowed {
					t.Errorf("concurrent enforce failed: allowed=%v err=%v", allowed, err)
				}
			}
		}()
	}
	wg.Wait()
}

func TestIAM_EnforceRoles(t *testing.T) {
	iam, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("failed to create IAM enforcer: %v", err)
	}

	// Load only role -> permission policies (no users in Casbin)
	policies := [][]string{
		{"role_admin", "*", "*"},
		{"role_editor", "articles", "write"},
		{"role_viewer", "articles", "read"},
	}

	if err := iam.LoadPolicies(nil, policies); err != nil {
		t.Fatalf("failed to load policies: %v", err)
	}

	// User with [role_viewer] trying to write -> false
	allowed, err := iam.EnforceRoles([]string{"role_viewer"}, "articles", "write")
	if err != nil || allowed {
		t.Errorf("expected false, got allowed=%v, err=%v", allowed, err)
	}

	// User with [role_viewer, role_editor] trying to write -> true
	allowed, err = iam.EnforceRoles([]string{"role_viewer", "role_editor"}, "articles", "write")
	if err != nil || !allowed {
		t.Errorf("expected true, got allowed=%v, err=%v", allowed, err)
	}

	// User with empty roles -> false
	allowed, err = iam.EnforceRoles([]string{}, "articles", "read")
	if err != nil || allowed {
		t.Errorf("expected false for empty roles, got allowed=%v, err=%v", allowed, err)
	}
}
