package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserClaims struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
	jwt.StandardClaims
}

func (jc *UserClaims) ToUser() (*User, error) {
	id, err := crypto.FromHex(jc.ID)
	if err != nil {
		return nil, err
	}
	return &User{id, jc.Type}, nil
}

type User struct {
	ID   crypto.ID
	Type int
}

const (
	ClaimsAnonymousType = iota
	ClaimsPublicType
	ClaimsNodeType
)

func (r *Rest) anonymousLogin(ctx echo.Context) error {
	claims := &UserClaims{
		utils.GenerateSHA().Hex(),
		ClaimsAnonymousType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(r.cfg.Secret)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (r *Rest) registerJWT() echo.MiddlewareFunc {
	jwtCfg := middleware.JWTConfig{
		Claims:     &UserClaims{},
		ContextKey: tokenKey,
		SigningKey: r.cfg.Secret,
		SuccessHandler: func(ctx echo.Context) {
			logger := LoggerFromContext(ctx)
			data := ctx.Get(tokenKey).(*jwt.Token)
			claims := data.Claims.(*UserClaims)
			user, err := claims.ToUser()
			if err != nil {
				logger.WithError(err).Error("Failed to get user from JWT token")
				return
			}
			ctx.Set(userKey, user)
		},
		Skipper: middleware.Skipper(func(ctx echo.Context) bool {
			return ctx.Path() == "/login" || strings.HasPrefix(ctx.Path(), "/admin")
		}),
	}
	return middleware.JWTWithConfig(jwtCfg)
}

func (r *Rest) adminAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == r.cfg.Username && password == r.cfg.Password {
			return true, nil
		}
		return false, nil
	})
}
