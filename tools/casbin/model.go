package casbin

import "github.com/casbin/casbin/v2/model"

// newRBACModel builds the Casbin RBAC model used by the Enforcer.
//
// Model spec:
//
//	[request_definition]  r = sub, obj, act
//	[policy_definition]   p = sub, obj, act
//	[role_definition]     g = _, _
//	[policy_effect]       e = some(where (p.eft == allow))
//	[matchers]            m = (r.obj == p.obj || p.obj == "*")
//	                          && (r.act == p.act || p.act == "*")
//	                          && (r.sub == p.sub || g(r.sub, p.sub))
//
// Short-circuit string matchers for resource and action run first.
// Direct role equality (r.sub == p.sub) is evaluated before role hierarchy g().
func newRBACModel() model.Model {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `(r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*") && (r.sub == p.sub || g(r.sub, p.sub))`)
	return m
}
