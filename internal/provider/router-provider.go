package provider

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"

	"brickwall/internal/common"
)

type GinLoggerAdapter struct{}

func (rcv *GinLoggerAdapter) Write(p []byte) (n int, err error) {
	slog.Info(string(p))
	return len(p), nil
}

type IRouterProvider interface {
	Init() IRouterProvider
	Engine() *gin.Engine
}

type RouterProvider struct {
	ctx    context.Context
	engine *gin.Engine
}

func NewRouterProvider(ctx context.Context) IRouterProvider {
	return &RouterProvider{ctx: ctx}
}

func (rcv *RouterProvider) Init() IRouterProvider {
	gin.DefaultWriter = &GinLoggerAdapter{}
	gin.DefaultErrorWriter = &GinLoggerAdapter{}

	rcv.engine = gin.Default()
	rcv.engine.Use(rcv.cors(rcv.ctx))

	return rcv
}

func (rcv *RouterProvider) Engine() *gin.Engine {
	return rcv.engine
}

func (rcv *RouterProvider) cors(ctx context.Context) gin.HandlerFunc {
	env := rcv.ctx.Value(common.KeyEnvProvider).(IEnvProvider)

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", env.GetString("CORS_ALLOW_ORIGIN", DefCorsAllowOrigin))
		c.Writer.Header().Set("Access-Control-Allow-Methods", env.GetString("CORS_ALLOW_METHODS", DefCorsAllowMethods))
		c.Writer.Header().Set("Access-Control-Allow-Headers", env.GetString("CORS_ALLOW_HEADERS", DefCorsAllowHeaders))
		c.Writer.Header().Set("Access-Control-Expose-Headers", env.GetString("CORS_EXPOSE_HEADERS", DefCorsExposeHeaders))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", env.GetString("CORS_ALLOW_CREDENTIALS", DefCorsAllowCredentials))
		c.Writer.Header().Set("Access-Control-Max-Age", env.GetString("CORS_MAX_AGE", DefCorsMaxAge))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
