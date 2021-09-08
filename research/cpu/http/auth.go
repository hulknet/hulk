package http

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kotfalya/hulk/pkg/utils"
	"github.com/kotfalya/hulk/research/cpu/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type JWTClaims struct {
	ID   string    `json:"id"`
	Type TokenType `json:"type"`
	jwt.StandardClaims
}

func (jc *JWTClaims) ToService() (*Service, error) {
	id, err := types.FromHex(jc.ID)
	if err != nil {
		return nil, err
	}
	return &Service{id, jc.Type}, nil
}

type Service struct {
	ID   types.ID
	Type TokenType
}

type TokenType int

const (
	AdminTokenType TokenType = iota
	AuthTypeType
	NetTokenType
	CpuTokenType
	StoreTokenType
	UnknownTokenType
)

func TokenTypeFromString(token string) TokenType {
	switch token {
	case "admin":
		return AdminTokenType
	case "net":
		return NetTokenType
	case "store":
		return StoreTokenType
	case "cpu":
		return CpuTokenType
	case "auth":
		return AuthTypeType
	default:
		return UnknownTokenType
	}
}

func TokenToString(token TokenType) string {
	switch token {
	case AdminTokenType:
		return "admin"
	case NetTokenType:
		return "net"
	case StoreTokenType:
		return "store"
	case CpuTokenType:
		return "cpu"
	case AuthTypeType:
		return "auth"
	default:
		return "unknown"
	}
}

func GenerateToken(key interface{}, tokenType TokenType) (string, error) {
	claims := &JWTClaims{
		utils.GenerateSHA().Hex(),
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return t, nil
}

func RegisterJWT(key interface{}) echo.MiddlewareFunc {
	jwtCfg := middleware.JWTConfig{
		Claims:        &JWTClaims{},
		ContextKey:    tokenKey,
		SigningKey:    key,
		SigningMethod: "ES256",
		SuccessHandler: func(ctx echo.Context) {
			logger := LoggerFromContext(ctx)
			data := ctx.Get(tokenKey).(*jwt.Token)
			claims := data.Claims.(*JWTClaims)
			service, err := claims.ToService()
			if err != nil {
				logger.WithError(err).Error("Failed to get service from JWT token")
				return
			}

			ctx.Set(serviceKey, service)
		},
	}
	return middleware.JWTWithConfig(jwtCfg)
}
