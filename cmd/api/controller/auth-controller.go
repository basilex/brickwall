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

type IAuthController interface {
	common.IController

	AuthSignup(*gin.Context)
	AuthSignin(*gin.Context)
	AuthSignout(*gin.Context)
}

type AuthController struct {
	ctx         context.Context
	group       *gin.RouterGroup
	authService service.IAuthService
	userService service.IUserService
}

func NewAuthController(ctx context.Context, grp *gin.RouterGroup) IAuthController {
	serviceManager := ctx.Value(common.KeyServiceManager).(service.IServiceManager)

	return &AuthController{
		ctx:         ctx,
		group:       grp,
		authService: serviceManager.AuthService(),
		userService: serviceManager.UserService(),
	}
}

func (rcv *AuthController) Register() {
	rcv.group.POST("/auth/signup", rcv.AuthSignup)
	rcv.group.POST("/auth/signin", rcv.AuthSignin)
	rcv.group.POST("/auth/signout", rcv.AuthSignout)
}

func (rcv *AuthController) AuthSignup(c *gin.Context) {
	qry := &exchange.AuthSignupReq{}

	if err := c.ShouldBindJSON(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.authService.Signup(qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *AuthController) AuthSignin(c *gin.Context) {
	qry := &exchange.AuthSigninReq{}

	if err := c.ShouldBindJSON(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.authService.Signin(qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *AuthController) AuthSignout(c *gin.Context) {
	ok, err := rcv.authService.Signout()
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(ok))
}
