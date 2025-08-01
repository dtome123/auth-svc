package types

type RouteScope int32

const (
	RouteScopePublic RouteScope = iota
	RouteScopeAuthenticated
	RouteScopeAuthorized
)
