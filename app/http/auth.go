package http

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hulknet/hulk/app/types"
)

type JWTClaims struct {
	ID   string    `json:"id"`
	Type TokenType `json:"type"`
	jwt.StandardClaims
}

func (jc *JWTClaims) ToService() (*Service, error) {
	id, err := types.ID256FromHex(jc.ID)
	if err != nil {
		return nil, err
	}
	return &Service{id, jc.Type}, nil
}

type Service struct {
	ID   types.ID256
	Type TokenType
}

type TokenType int

const (
	AdminTokenType TokenType = iota
	AuthTypeType
	NetTokenType
	CpuTokenType
	DiskTokenType
	MemoryTokenType
	UnknownTokenType
)

func TokenTypeFromString(token string) TokenType {
	switch token {
	case "admin":
		return AdminTokenType
	case "net":
		return NetTokenType
	case "disk":
		return DiskTokenType
	case "memory":
		return MemoryTokenType
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
	case DiskTokenType:
		return "disk"
	case MemoryTokenType:
		return "memory"
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
		types.ID256ToHex(types.GenerateSHA()),
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
