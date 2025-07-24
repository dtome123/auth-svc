package types

type RouteScope string

const (
	RouteScopePublic        RouteScope = "public"
	RouteScopeAuthenticated RouteScope = "authenticated"
	RouteScopeAuthorized    RouteScope = "authorized"
)
