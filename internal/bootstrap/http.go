package bootstrap

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ramdhanrizki/bytecode-api/configs"
	"github.com/ramdhanrizki/bytecode-api/internal/identity"
	"github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

func NewHTTPServer(cfg configs.Config, logger *zap.Logger, identityModule *identity.Module) *http.Server {
	if cfg.App.Env != "production" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger))

	healthHandler := func(c *gin.Context) {
		response.Success(c, http.StatusOK, "service is healthy", gin.H{
			"status":    "ok",
			"name":      cfg.App.Name,
			"env":       cfg.App.Env,
			"timestamp": time.Now().UTC(),
		})
	}

	router.GET("/health", healthHandler)
	router.GET("/api/v1/health", healthHandler)

	if identityModule != nil {
		identityModule.RegisterRoutes(router)
	}

	return &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func requestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		logger.Info("http request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(startedAt)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
