package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"brickwall/cmd/api/exchange"
	"brickwall/internal/common"
	"brickwall/internal/storage/dbs"
	"brickwall/internal/utils"
)

type ICurrencyService interface {
	// CRUD operations
	CurrencyNew(*exchange.CurrencyNewReq) (*dbs.Currency, error)
	CurrencySelect(*gin.Context, *exchange.CurrencyQuery) ([]*dbs.Currency, *utils.PaginatorResources, error)
	CurrencySelectByID(string) (*dbs.Currency, error)
	CurrencyUpdateByID(*exchange.CurrencyUpdateReq) (*dbs.Currency, error)
	CurrencyDeleteByID(string) error

	// Business logic operations
	CurrencyCountrySelect(*gin.Context, *exchange.CurrencyQuery) ([]*dbs.CurrencyCountrySelectRow, *utils.PaginatorResources, error)
}

type CurrencyService struct {
	ctx     context.Context
	queries *dbs.Queries
}

func NewCurrencyService(ctx context.Context, queries *dbs.Queries) ICurrencyService {
	return &CurrencyService{ctx: ctx, queries: queries}
}

func (rcv *CurrencyService) CurrencyNew(req *exchange.CurrencyNewReq) (*dbs.Currency, error) {
	res, err := rcv.queries.CurrencyNew(context.Background(), &dbs.CurrencyNewParams{
		Name: req.Name, Code: req.Code, NumCode: req.NumCode, Symbol: req.Symbol,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}
	return res, nil
}

func (rcv *CurrencyService) CurrencySelect(c *gin.Context, qry *exchange.CurrencyQuery) ([]*dbs.Currency, *utils.PaginatorResources, error) {
	count, err := rcv.queries.CurrencyCount(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordCount, err)
	}
	paginator := utils.NewPaginator(c.Request, qry.Size, count)

	params := &dbs.CurrencySelectParams{
		SqlLimit:  int32(qry.Size),
		SqlOffset: int32(paginator.Offset()),
		SqlOrder:  c.DefaultQuery("order", "id"),
	}
	res, err := rcv.queries.CurrencySelect(context.Background(), params)
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

func (rcv *CurrencyService) CurrencySelectByID(id string) (*dbs.Currency, error) {
	res, err := rcv.queries.CurrencySelectByID(context.Background(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *CurrencyService) CurrencyUpdateByID(req *exchange.CurrencyUpdateReq) (*dbs.Currency, error) {
	params := &dbs.CurrencyUpdateByIDParams{
		ID: req.ID, Name: req.Name, Code: req.Code, NumCode: req.NumCode, Symbol: req.Symbol,
	}
	res, err := rcv.queries.CurrencyUpdateByID(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordUpdate, err)
		}
	}
	return res, nil
}

func (rcv *CurrencyService) CurrencyDeleteByID(id string) error {
	if _, err := rcv.queries.CurrencyDeleteByID(context.Background(), id); err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return fmt.Errorf("%w: %v", common.ErrDBRecordDelete, err)
		}
	}
	return nil
}

func (rcv *CurrencyService) CurrencyCountrySelect(c *gin.Context, qry *exchange.CurrencyQuery) ([]*dbs.CurrencyCountrySelectRow, *utils.PaginatorResources, error) {
	count, err := rcv.queries.CurrencyCount(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordCount, err)
	}
	paginator := utils.NewPaginator(c.Request, qry.Size, count)

	params := &dbs.CurrencyCountrySelectParams{
		SqlLimit:  int32(qry.Size),
		SqlOffset: int32(paginator.Offset()),
		SqlOrder:  c.DefaultQuery("order", "id"),
	}
	res, err := rcv.queries.CurrencyCountrySelect(context.Background(), params)
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
