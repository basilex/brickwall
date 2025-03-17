package provider

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"

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
	rcv.engine.Use(cors(rcv.ctx))
	return rcv
}

func (rcv *RouterProvider) Engine() *gin.Engine {
	return rcv.engine
}

func cors(ctx context.Context) gin.HandlerFunc {
	cli := ctx.Value(common.KeyCommand).(*cli.Command)

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cli.String("cors-allow-origin"))
		c.Writer.Header().Set("Access-Control-Allow-Methods", cli.String("cors-allow-methods"))
		c.Writer.Header().Set("Access-Control-Allow-Headers", cli.String("cors-allow-headers"))
		c.Writer.Header().Set("Access-Control-Expose-Headers", cli.String("cors-expose-headers"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", cli.String("cors-allow-credentials"))
		c.Writer.Header().Set("Access-Control-Max-Age", cli.String("cors-max-age"))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
