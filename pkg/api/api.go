package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kotfalya/hulk/pkg/config"
	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/host"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var restLog *log.Entry

type Rest struct {
	echo *echo.Echo
	host *host.Host
	cfg  *config.HTTPConfig
}

func NewRestServer(h *host.Host, cfg *config.HTTPConfig) *Rest {
	restLog = log.WithFields(log.Fields{
		"pkg":     "rest",
		"host-id": h.ID().HexL(crypto.IDLogLen),
	})

	rest := &Rest{
		echo: echo.New(),
		host: h,
		cfg:  cfg,
	}

	rest.echo.Use(middleware.Recover())
	rest.echo.Use(rest.setupRequest)
	rest.echo.Use(rest.registerJWT())

	rest.echo.POST("/login", rest.anonymousLogin)
	rest.echo.GET("/status", rest.hostStatus)
	rest.echo.GET("/nodes", rest.nodesList)
	rest.echo.GET("/nodes/:id", rest.nodesItem)
	rest.echo.GET("/connect", rest.hostConnect)

	admin := rest.echo.Group("/admin")
	admin.Use(rest.adminAuth())
	admin.POST("/net/create", rest.createNet)
	admin.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"admin": true,
		})
	})

	rest.echo.GET("/user", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"userId": UserFromContext(c).ID.HexL(crypto.IDLogLen),
			"type":   UserFromContext(c).Type,
		})
	})

	return rest
}

func (r *Rest) setupRequest(f echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		reqId, _ := uuid.NewRandom()
		logger := restLog.WithFields(logrus.Fields{
			"method": req.Method,
			"path":   req.URL.Path,
			"req_id": reqId.String(),
		})
		ctx.Set(loggerKey, logger)

		startTime := time.Now()
		defer func() {
			rsp := ctx.Response()
			logger.WithFields(logrus.Fields{
				"status_code":   rsp.Status,
				"runtime_micro": time.Since(startTime).Microseconds(),
			}).Info("Finished request")
		}()

		logger.WithFields(logrus.Fields{
			"user_agent":     req.UserAgent(),
			"content_length": req.ContentLength,
		}).Info("Starting request")

		err := f(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}

func (r *Rest) Listen() error {
	return r.echo.Start(r.host.Addr())
}
