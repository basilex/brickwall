package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"brickwall/cmd/api/service"
	"brickwall/internal/common"
)

type IAuxController interface {
	common.IController

	Index(*gin.Context)
	Health(*gin.Context)
	Metadata(*gin.Context)
}

type AuxController struct {
	ctx        context.Context
	group      *gin.RouterGroup
	auxService service.IAuxService
}

func NewAuxController(ctx context.Context, grp *gin.RouterGroup) IAuxController {
	serviceManager := ctx.Value(common.KeyServiceManager).(service.IServiceManager)

	return &AuxController{
		ctx: ctx, group: grp, auxService: serviceManager.AuxService(),
	}
}

func (rcv *AuxController) Register() {
	rcv.group.GET("/aux", rcv.Index)
	rcv.group.GET("/aux/health", rcv.Health)
	rcv.group.GET("/aux/metadata", rcv.Metadata)
}

// @Summary     Greetings
// @Description Return the platformm greetings
// @Tags        aux
// @Accept      json
// @Produce     json
// @Success     200
// @Router      /aux [get]
func (rcv *AuxController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, common.NewResponse(rcv.auxService.Index()))
}

// @Summary     Health check
// @Description Checks the health
// @Tags        aux
// @Accept      json
// @Produce     json
// @Success     200
// @Router      /aux/health [get]
func (rcv *AuxController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, common.NewResponse(rcv.auxService.Health()))
}

// @Summary     Binary metadata
// @Description Obtain the binary metadata
// @Tags        aux
// @Accept      json
// @Produce     json
// @Success     200
// @Router      /aux/metadata [get]
func (rcv *AuxController) Metadata(c *gin.Context) {
	c.JSON(http.StatusOK, common.NewResponse(rcv.auxService.Metadata()))
}
