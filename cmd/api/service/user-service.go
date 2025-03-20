package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"brickwall/cmd/api/exchange"
	"brickwall/internal/common"
	"brickwall/internal/utils"

	"brickwall/internal/storage/dbs"
)

type IUserService interface {
	// CRUD operations
	UserNew(*exchange.UserNewReq) (*dbs.User, error)
	UserSelect(*gin.Context, *exchange.UserQuery) ([]*dbs.User, *utils.PaginatorResources, error)
	UserSelectByID(string) (*dbs.User, error)
	UserSelectByUsername(string) (*dbs.User, error)
	UserUpdateCredentialsByID(*exchange.UserUpdateCredentialsReq) (*dbs.User, error)
	UserUpdateIsBlockedByID(*exchange.UserUpdateIsBlockedByIDReq) (*dbs.User, error)
	UserUpdateIsCheckedByID(*exchange.UserUpdateIsCheckedByIDReq) (*dbs.User, error)
	UserUpdateVisitedAtByID(*exchange.UserUpdateVisitedAtByIDReq) (*dbs.User, error)
	UserDeleteByID(string) error
}

type UserService struct {
	ctx     context.Context
	queries *dbs.Queries
}

func NewUserService(ctx context.Context, queries *dbs.Queries) IUserService {
	return &UserService{ctx: ctx, queries: queries}
}

func (rcv *UserService) UserNew(req *exchange.UserNewReq) (*dbs.User, error) {
	passwordCrypted, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	res, err := rcv.queries.UserNew(context.Background(), &dbs.UserNewParams{
		Username: req.Username,
		Password: string(passwordCrypted),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}
	return res, nil
}

func (rcv *UserService) UserSelect(c *gin.Context, qry *exchange.UserQuery) ([]*dbs.User, *utils.PaginatorResources, error) {
	count, err := rcv.queries.UserCount(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", common.ErrDBRecordCount, err)
	}
	paginator := utils.NewPaginator(c.Request, qry.Size, count)

	params := &dbs.UserSelectParams{
		SqlLimit:  int32(qry.Size),
		SqlOffset: int32(paginator.Offset()),
		SqlOrder:  c.DefaultQuery("order", "id"),
	}
	res, err := rcv.queries.UserSelect(context.Background(), params)
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

func (rcv *UserService) UserSelectByID(req string) (*dbs.User, error) {
	res, err := rcv.queries.UserSelectByID(context.Background(), req)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *UserService) UserSelectByUsername(req string) (*dbs.User, error) {
	res, err := rcv.queries.UserSelectByUsername(context.Background(), req)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *UserService) UserUpdateCredentialsByID(req *exchange.UserUpdateCredentialsReq) (*dbs.User, error) {
	passwordCrypted, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	params := &dbs.UserUpdateCredentialsByIDParams{
		ID: req.ID, Username: req.Username, Password: string(passwordCrypted),
	}
	res, err := rcv.queries.UserUpdateCredentialsByID(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordUpdate, err)
		}
	}
	return res, nil
}

func (rcv *UserService) UserUpdateIsBlockedByID(req *exchange.UserUpdateIsBlockedByIDReq) (*dbs.User, error) {
	params := &dbs.UserUpdateIsBlockedByIDParams{
		ID:        req.ID,
		IsBlocked: req.IsBlocked,
	}
	res, err := rcv.queries.UserUpdateIsBlockedByID(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *UserService) UserUpdateIsCheckedByID(req *exchange.UserUpdateIsCheckedByIDReq) (*dbs.User, error) {
	params := &dbs.UserUpdateIsCheckedByIDParams{
		ID:        req.ID,
		IsChecked: req.IsChecked,
	}
	res, err := rcv.queries.UserUpdateIsCheckedByID(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *UserService) UserUpdateVisitedAtByID(req *exchange.UserUpdateVisitedAtByIDReq) (*dbs.User, error) {
	params := &dbs.UserUpdateVisitedAtByIDParams{
		ID:        req.ID,
		VisitedAt: req.VisitedAt,
	}
	res, err := rcv.queries.UserUpdateVisitedAtByID(context.Background(), params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	return res, nil
}

func (rcv *UserService) UserDeleteByID(id string) error {
	if _, err := rcv.queries.UserDeleteByID(context.Background(), id); err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return fmt.Errorf("%w: %v", common.ErrDBRecordDelete, err)
		}
	}
	return nil
}
