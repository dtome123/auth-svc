package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
)

// PathRule defines which interceptors should be executed for gRPC methods matching Pattern or Patterns.
type PathRule struct {
	Pattern      string   // Single pattern (e.g. "/auth.v1.AuthorizationService/*" or exact method name)
	Patterns     []string // Multiple patterns
	Interceptors []grpc.UnaryServerInterceptor
}

// Matches checks if the method matches any pattern defined in the rule.
func (r PathRule) Matches(method string) bool {
	if r.Pattern != "" && matchPattern(method, r.Pattern) {
		return true
	}
	for _, p := range r.Patterns {
		if matchPattern(method, p) {
			return true
		}
	}
	return false
}

// NewPathBasedInterceptor returns a grpc.UnaryServerInterceptor that routes incoming RPCs
// to specific interceptors based on matching path rules.
func NewPathBasedInterceptor(rules ...PathRule) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var matched []grpc.UnaryServerInterceptor
		for _, rule := range rules {
			if rule.Matches(info.FullMethod) {
				matched = append(matched, rule.Interceptors...)
			}
		}

		if len(matched) == 0 {
			return handler(ctx, req)
		}

		return chainInterceptors(matched)(ctx, req, info, handler)
	}
}

func matchPattern(method, pattern string) bool {
	if prefix, ok := strings.CutSuffix(pattern, "*"); ok {
		return strings.HasPrefix(method, prefix)
	}
	return method == pattern
}

func chainInterceptors(interceptors []grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		currHandler := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := currHandler
			currHandler = func(c context.Context, r interface{}) (interface{}, error) {
				return interceptor(c, r, info, next)
			}
		}
		return currHandler(ctx, req)
	}
}
