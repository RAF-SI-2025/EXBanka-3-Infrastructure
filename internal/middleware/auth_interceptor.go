package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/config"
	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"
	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/util"
)

// publicMethods are RPC methods that do not require a JWT
var publicMethods = map[string]bool{
	"/auth.v1.AuthService/Login":                true,
	"/auth.v1.AuthService/ActivateAccount":      true,
	"/auth.v1.AuthService/RequestPasswordReset": true,
	"/auth.v1.AuthService/ResetPassword":        true,
}

// requiredPermissions maps each protected RPC to its minimum required permission.
// Employees with the "admin" permission bypass all of these checks.
var requiredPermissions = map[string]string{
	"/employee.v1.EmployeeService/CreateEmployee":            models.PermEmployeeCreate,
	"/employee.v1.EmployeeService/ListEmployees":             models.PermEmployeeRead,
	"/employee.v1.EmployeeService/GetEmployee":               models.PermEmployeeRead,
	"/employee.v1.EmployeeService/UpdateEmployee":            models.PermEmployeeUpdate,
	"/employee.v1.EmployeeService/SetEmployeeActive":         models.PermEmployeeActivate,
	"/employee.v1.EmployeeService/UpdateEmployeePermissions": models.PermEmployeePermissions,
	"/employee.v1.EmployeeService/GetAllPermissions":         models.PermEmployeeRead,
	"/client.v1.ClientService/CreateClient":                  models.PermAdmin,
	"/client.v1.ClientService/GetClient":                     models.PermAdmin,
	"/client.v1.ClientService/ListClients":                   models.PermAdmin,
	"/client.v1.ClientService/UpdateClient":                  models.PermAdmin,
	"/client.v1.ClientService/UpdateClientPermissions":       models.PermAdmin,
}

// claimsContextKey is the unexported context key for JWT claims
type claimsContextKey struct{}

// ClaimsKey is used to store/retrieve *util.Claims in a context
var ClaimsKey = claimsContextKey{}

// GetClaimsFromContext retrieves JWT claims stored by the auth interceptor
func GetClaimsFromContext(ctx context.Context) (*util.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*util.Claims)
	return claims, ok
}

// AuthInterceptor validates JWT tokens and checks per-RPC permissions
func AuthInterceptor(cfg *config.Config) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip auth for public endpoints
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")
		claims, err := util.ParseToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		if claims.TokenType != "access" {
			return nil, status.Error(codes.Unauthenticated, "wrong token type: expected access token")
		}

		// Permission check — admins bypass all restrictions
		if requiredPerm, exists := requiredPermissions[info.FullMethod]; exists {
			isAdmin := util.HasPermission(claims, models.PermAdmin)
			hasPerm := util.HasPermission(claims, requiredPerm)
			if !isAdmin && !hasPerm {
				return nil, status.Errorf(codes.PermissionDenied,
					"permission %q required", requiredPerm)
			}
		}

		ctx = context.WithValue(ctx, ClaimsKey, claims)
		return handler(ctx, req)
	}
}
