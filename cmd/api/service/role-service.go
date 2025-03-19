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

type IRoleService interface {
	// CRUD operations
	RoleNew(*exchange.RoleNewReq) (*dbs.Role, error)
	RoleSelect(*gin.Context, *exchange.RoleQuery) ([]*dbs.Role, *utils.PaginatorResources, error)
	RoleSelectByID(string) (*dbs.Role, error)
	RoleUpdateByID(*exchange.RoleUpdateReq) (*dbs.Role, error)
	RoleDeleteByID(string) error
}

type RoleService struct {
	ctx     context.Context
	queries *dbs.Queries
}

func NewRoleService(ctx context.Context, queries *dbs.Queries) IRoleService {
	return &RoleService{ctx: ctx, queries: queries}
}

func (rcv *RoleService) RoleNew(req *exchange.RoleNewReq) (*dbs.Role, error) {
	res, err := rcv.queries.RoleNew(context.Background(), req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}
	return res, nil
}

func (rcv *RoleService) RoleSelect(c *gin.Context, qry *exchange.RoleQuery) ([]*dbs.Role, *utils.PaginatorResources, error) {
	count, err := rcv.queries.RoleCount(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordCount, err)
	}
	paginator := utils.NewPaginator(c.Request, qry.Size, count)

	params := &dbs.RoleSelectParams{
		SqlLimit:  int32(qry.Size),
		SqlOffset: int32(paginator.Offset()),
		SqlOrder:  c.DefaultQuery("order", "id"),
	}
	res, err := rcv.queries.RoleSelect(context.Background(), params)
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

func (rcv *RoleService) RoleSelectByID(req string) (*dbs.Role, error) {
	res, err := rcv.queries.RoleSelectByID(context.Background(), req)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *RoleService) RoleUpdateByID(req *exchange.RoleUpdateReq) (*dbs.Role, error) {
	params := &dbs.RoleUpdateByIDParams{
		ID: req.ID, Name: req.Name,
	}
	res, err := rcv.queries.RoleUpdateByID(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordUpdate, err)
		}
	}
	return res, nil
}

func (rcv *RoleService) RoleDeleteByID(id string) error {
	if _, err := rcv.queries.RoleDeleteByID(context.Background(), id); err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return fmt.Errorf("%w: %v", common.ErrDBRecordDelete, err)
		}
	}
	return nil
}
