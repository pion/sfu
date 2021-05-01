package jwt

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthConfig auth config
type AuthConfig struct {
	Enabled         bool     `mapstructure:"enabled"`
	Key             string   `mapstructure:"key"`
	KeyType         string   `mapstructure:"key_type"`
	AllowedServices []string `mapstructure:"allowed_services"`
}

// KeyFunc auth key types
func (a AuthConfig) KeyFunc(t *jwt.Token) (interface{}, error) {
	// nolint: gocritic
	switch a.KeyType {
	//TODO: add more support for keytypes here
	default:
		return []byte(a.Key), nil
	}
}

// claims custom claims type for jwt
type CustomClaims struct {
	UID string `json:"uid"`
	SID string `json:"sid"`
	*jwt.StandardClaims
}

func GetClaim(ctx context.Context, ac *AuthConfig) (*CustomClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "valid JWT token required")
	}

	token, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "valid JWT token required")
	}

	jwtToken, err := jwt.ParseWithClaims(token[0], &CustomClaims{}, ac.KeyFunc)

	if claims, ok := jwtToken.Claims.(*CustomClaims); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, status.Errorf(codes.Unauthenticated, "valid JWT token required: %v", err)
}
