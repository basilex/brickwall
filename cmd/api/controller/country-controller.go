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

type ICountryController interface {
	common.IController

	// CRUD operations
	CountryNew(*gin.Context)
	CountrySelect(*gin.Context)
	CountrySelectByID(*gin.Context)
	CountryUpdateByID(*gin.Context)
	CountryDeleteByID(*gin.Context)

	// Business logic operations
	CountryCurrencySelect(*gin.Context)
}
type CountryController struct {
	ctx            context.Context
	group          *gin.RouterGroup
	countryService service.ICountryService
}

func NewCountryController(ctx context.Context, grp *gin.RouterGroup) ICountryController {
	serviceManager := ctx.Value(common.KeyServiceManager).(service.IServiceManager)

	return &CountryController{
		ctx: ctx, group: grp, countryService: serviceManager.CountryService(),
	}
}

func (rcv *CountryController) Register() {
	// CRUD operations
	rcv.group.POST("/country", rcv.CountryNew)
	rcv.group.GET("/country", rcv.CountrySelect)
	rcv.group.GET("/country/:id", rcv.CountrySelectByID)
	rcv.group.PUT("/country", rcv.CountryUpdateByID)
	rcv.group.DELETE("/country/:id", rcv.CountryDeleteByID)

	// Business logic operations
	rcv.group.GET("/country-currency", rcv.CountryCurrencySelect)
}

func (rcv *CountryController) CountryNew(c *gin.Context) {
	req := &exchange.CountryNewReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.countryService.CountryNew(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusCreated, common.NewResponse(res))
}

func (rcv *CountryController) CountrySelect(c *gin.Context) {
	qry := &exchange.CountryQuery{}

	if err := c.ShouldBindQuery(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, pag, err := rcv.countryService.CountrySelect(c, qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"data": res, "paginator": pag}),
	)
}

func (rcv *CountryController) CountrySelectByID(c *gin.Context) {
	uri := &exchange.CountryUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.countryService.CountrySelectByID(uri.ID)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *CountryController) CountryUpdateByID(c *gin.Context) {
	req := &exchange.CountryUpdateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, err := rcv.countryService.CountryUpdateByID(req)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(res))
}

func (rcv *CountryController) CountryDeleteByID(c *gin.Context) {
	uri := &exchange.CountryUriID{}

	if err := c.ShouldBindUri(uri); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	if err := rcv.countryService.CountryDeleteByID(uri.ID); err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"message": "no data"}),
	)
}

// Business operations
func (rcv *CountryController) CountryCurrencySelect(c *gin.Context) {
	qry := &exchange.CountryQuery{}

	if err := c.ShouldBindQuery(&qry); err != nil {
		c.JSON(common.ErrMapper(fmt.Errorf("%w: %v", common.ErrReqBindJson, err)))
		return
	}
	res, pag, err := rcv.countryService.CountryCurrencySelect(c, qry)
	if err != nil {
		c.JSON(common.ErrMapper(err))
		return
	}
	c.JSON(http.StatusOK, common.NewResponse(
		gin.H{"data": res, "paginator": pag}),
	)
}
