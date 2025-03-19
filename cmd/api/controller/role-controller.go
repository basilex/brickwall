package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"brickwall/cmd/api/exchange"
	"brickwall/cmd/api/service"
	"brickwall/internal/common"
)

type IRoleController interface {
	common.IController

	// CRUD operations
	RoleNew(*gin.Context)
	RoleSelect(*gin.Context)
	RoleSelectByID(*gin.Context)
	RoleUpdateByID(*gin.Context)
	RoleDeleteByID(*gin.Context)
}

type RoleController struct {
	ctx         context.Context
	group       *gin.RouterGroup
	RoleService service.IRoleService
}

func NewRoleController(ctx context.Context, grp *gin.RouterGroup) IRoleController {
	serviceManager := ctx.Value(common.KeyServiceManager).(service.IServiceManager)

	return &RoleController{
		ctx: ctx, group: grp, RoleService: serviceManager.RoleService(),
	}
}

func (rcv *RoleController) Register() {
	// CRUD operations
	rcv.group.POST("/role", rcv.RoleNew)
	rcv.group.GET("/role", rcv.RoleSelect)
	rcv.group.GET("/role/:id", rcv.RoleSelectByID)
	rcv.group.PUT("/role", rcv.RoleUpdateByID)
	rcv.group.DELETE("/role/:id", rcv.RoleDeleteByID)
}

func (rcv *RoleController) RoleNew(c *gin.Context) {
	req := &exchange.RoleNewReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.RoleService.RoleNew(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusCreated, common.NewResponse(res))
}

func (rcv *RoleController) RoleSelect(c *gin.Context) {
	qry := &exchange.RoleQuery{}

	if err := c.ShouldBindQuery(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, pag, err := rcv.RoleService.RoleSelect(c, qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"data": res, "paginator": pag}),
	)
}

func (rcv *RoleController) RoleSelectByID(c *gin.Context) {
	uri := &exchange.RoleUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.RoleService.RoleSelectByID(uri.ID)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *RoleController) RoleUpdateByID(c *gin.Context) {
	req := &exchange.RoleUpdateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.RoleService.RoleUpdateByID(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *RoleController) RoleDeleteByID(c *gin.Context) {
	uri := &exchange.RoleUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	if err := rcv.RoleService.RoleDeleteByID(uri.ID); err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"message": "no data"}),
	)
}
