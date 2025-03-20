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

type ICurrencyController interface {
	common.IController

	// CRUD operations
	CurrencyNew(*gin.Context)
	CurrencySelect(*gin.Context)
	CurrencySelectByID(*gin.Context)
	CurrencyUpdateByID(*gin.Context)
	CurrencyDeleteByID(*gin.Context)

	// Business logic operations
	CurrencyCountrySelect(*gin.Context)
}
type CurrencyController struct {
	ctx             context.Context
	group           *gin.RouterGroup
	currencyService service.ICurrencyService
}

func NewCurrencyController(ctx context.Context, grp *gin.RouterGroup) ICurrencyController {
	serviceManager := ctx.Value(common.KeyServiceManager).(service.IServiceManager)

	return &CurrencyController{
		ctx: ctx, group: grp, currencyService: serviceManager.CurrencyService(),
	}
}

func (rcv *CurrencyController) Register() {
	// CRUD operations
	rcv.group.POST("/currency", rcv.CurrencyNew)
	rcv.group.GET("/currency", rcv.CurrencySelect)
	rcv.group.GET("/currency/:id", rcv.CurrencySelectByID)
	rcv.group.PUT("/currency", rcv.CurrencyUpdateByID)
	rcv.group.DELETE("/currency/:id", rcv.CurrencyDeleteByID)

	// Business logic operations
	rcv.group.GET("/currency-country", rcv.CurrencyCountrySelect)
}

func (rcv *CurrencyController) CurrencyNew(c *gin.Context) {
	req := &exchange.CurrencyNewReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.currencyService.CurrencyNew(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusCreated, common.NewResponse(res))
}

func (rcv *CurrencyController) CurrencySelect(c *gin.Context) {
	qry := &exchange.CurrencyQuery{}

	if err := c.ShouldBindQuery(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, pag, err := rcv.currencyService.CurrencySelect(c, qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"data": res, "paginator": pag}),
	)
}

func (rcv *CurrencyController) CurrencySelectByID(c *gin.Context) {
	uri := &exchange.CurrencyUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.currencyService.CurrencySelectByID(uri.ID)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *CurrencyController) CurrencyUpdateByID(c *gin.Context) {
	req := &exchange.CurrencyUpdateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.currencyService.CurrencyUpdateByID(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *CurrencyController) CurrencyDeleteByID(c *gin.Context) {
	uri := &exchange.CurrencyUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	if err := rcv.currencyService.CurrencyDeleteByID(uri.ID); err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"message": "no data"}),
	)
}

// Business operations
func (rcv *CurrencyController) CurrencyCountrySelect(c *gin.Context) {
	qry := &exchange.CurrencyQuery{}

	if err := c.ShouldBindQuery(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, pag, err := rcv.currencyService.CurrencyCountrySelect(c, qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"data": res, "paginator": pag}),
	)
}
