package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"brickwall/cmd/api/exchange"
	"brickwall/internal/common"
	"brickwall/internal/utils"

	"brickwall/internal/storage/dbs"
)

type ICountryService interface {
	// CRUD operations
	CountryNew(*exchange.CountryNewReq) (*dbs.Country, error)
	CountrySelect(*gin.Context, *exchange.CountryQuery) ([]*dbs.Country, *utils.PaginatorResources, error)
	CountrySelectByID(string) (*dbs.Country, error)
	CountryUpdateByID(*exchange.CountryUpdateReq) (*dbs.Country, error)
	CountryDeleteByID(string) error

	// Business logic operations
	CountryCurrencySelect(*gin.Context, *exchange.CountryQuery) ([]*dbs.CountryCurrencySelectRow, *utils.PaginatorResources, error)
}

type CountryService struct {
	ctx     context.Context
	queries *dbs.Queries
}

func NewCountryService(ctx context.Context, queries *dbs.Queries) ICountryService {
	return &CountryService{ctx: ctx, queries: queries}
}

func (rcv *CountryService) CountryNew(req *exchange.CountryNewReq) (*dbs.Country, error) {
	res, err := rcv.queries.CountryNew(context.Background(), &dbs.CountryNewParams{
		Name: req.Name, Iso2: req.Iso2, Iso3: req.Iso3, NumCode: req.NumCode,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}
	return res, nil
}

func (rcv *CountryService) CountrySelect(c *gin.Context, qry *exchange.CountryQuery) ([]*dbs.Country, *utils.PaginatorResources, error) {
	count, err := rcv.queries.CountryCount(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordCount, err)
	}
	paginator := utils.NewPaginator(c.Request, qry.Size, count)

	params := &dbs.CountrySelectParams{
		SqlLimit:  int32(qry.Size),
		SqlOffset: int32(paginator.Offset()),
		SqlOrder:  c.DefaultQuery("order", "id"),
	}
	res, err := rcv.queries.CountrySelect(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	pag := utils.NewPaginatorBuilder(paginator).Build()
	return res, pag, nil
}

func (rcv *CountryService) CountrySelectByID(id string) (*dbs.Country, error) {
	res, err := rcv.queries.CountrySelectByID(context.Background(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *CountryService) CountryUpdateByID(req *exchange.CountryUpdateReq) (*dbs.Country, error) {
	params := &dbs.CountryUpdateByIDParams{
		ID: req.ID, Name: req.Name, Iso2: req.Iso2, Iso3: req.Iso3, NumCode: req.NumCode,
	}
	res, err := rcv.queries.CountryUpdateByID(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordUpdate, err)
		}
	}
	return res, nil
}

func (rcv *CountryService) CountryDeleteByID(id string) error {
	if _, err := rcv.queries.CountryDeleteByID(context.Background(), id); err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return fmt.Errorf("%w: %v", common.ErrDBRecordDelete, err)
		}
	}
	return nil
}

func (rcv *CountryService) CountryCurrencySelect(c *gin.Context, qry *exchange.CountryQuery) ([]*dbs.CountryCurrencySelectRow, *utils.PaginatorResources, error) {
	count, err := rcv.queries.CountryCount(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordCount, err)
	}
	paginator := utils.NewPaginator(c.Request, qry.Size, count)

	params := &dbs.CountryCurrencySelectParams{
		SqlLimit:  int32(qry.Size),
		SqlOffset: int32(paginator.Offset()),
		SqlOrder:  c.DefaultQuery("order", "id"),
	}
	res, err := rcv.queries.CountryCurrencySelect(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	pag := utils.NewPaginatorBuilder(paginator).Build()
	return res, pag, nil
}
