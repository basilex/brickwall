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

type IUserController interface {
	common.IController

	// CRUD operations
	UserSelect(*gin.Context)
	UserSelectByID(*gin.Context)
	UserSelectByUsername(*gin.Context)
	UserDeleteByID(*gin.Context)

	// Business operations
	UserUpdateCredentialsByID(*gin.Context)
	UserUpdateIsBlockedByID(*gin.Context)
	UserUpdateIsCheckedByID(*gin.Context)
}

type UserController struct {
	ctx         context.Context
	group       *gin.RouterGroup
	userService service.IUserService
}

func NewUserController(ctx context.Context, grp *gin.RouterGroup) IUserController {
	serviceManager := ctx.Value(common.KeyServiceManager).(service.IServiceManager)

	return &UserController{
		ctx: ctx, group: grp, userService: serviceManager.UserService(),
	}
}

func (rcv *UserController) Register() {
	// CRUD operations
	rcv.group.GET("/user", rcv.UserSelect)
	rcv.group.GET("/user/:id", rcv.UserSelectByID)
	rcv.group.PUT("/user/section/credentials", rcv.UserUpdateCredentialsByID)
	rcv.group.PUT("/user/section/is_blocked", rcv.UserUpdateIsBlockedByID)
	rcv.group.PUT("/user/section/is_checked", rcv.UserUpdateIsCheckedByID)
	rcv.group.DELETE("/user/:id", rcv.UserDeleteByID)
}

func (rcv *UserController) UserSelect(c *gin.Context) {
	qry := &exchange.UserQuery{}

	if err := c.ShouldBindQuery(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, pag, err := rcv.userService.UserSelect(c, qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"data": res, "paginator": pag}),
	)
}

func (rcv *UserController) UserSelectByID(c *gin.Context) {
	uri := &exchange.UserUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.userService.UserSelectByID(uri.ID)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *UserController) UserSelectByUsername(c *gin.Context) {
	uri := &exchange.UserUriUsername{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.userService.UserSelectByUsername(uri.Username)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *UserController) UserDeleteByID(c *gin.Context) {
	uri := &exchange.UserUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	if err := rcv.userService.UserDeleteByID(uri.ID); err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"message": "no data"}),
	)
}

// Buseness operations
func (rcv *UserController) UserUpdateCredentialsByID(c *gin.Context) {
	req := &exchange.UserUpdateCredentialsReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.userService.UserUpdateCredentialsByID(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *UserController) UserUpdateIsBlockedByID(c *gin.Context) {
	req := &exchange.UserUpdateIsBlockedByIDReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.userService.UserUpdateIsBlockedByID(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *UserController) UserUpdateIsCheckedByID(c *gin.Context) {
	req := &exchange.UserUpdateIsCheckedByIDReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.userService.UserUpdateIsCheckedByID(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}
