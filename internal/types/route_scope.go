package types

type RouteScope string

const (
	RouteScopeAuthenticated RouteScope = "AUTHENTICATED"
	RouteScopeAuthorized    RouteScope = "AUTHORIZED"
)
