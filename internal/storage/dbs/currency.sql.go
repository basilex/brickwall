// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: currency.sql

package dbs

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const currencyCount = `-- name: CurrencyCount :one
select count(*) from currency
`

// CurrencyCount
//
//	select count(*) from currency
func (q *Queries) CurrencyCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, currencyCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const currencyCountrySelect = `-- name: CurrencyCountrySelect :many
select cr.id, cr.name, cr.code, cr.num_code, cr.symbol, cr.created_at, cr.updated_at,
  coalesce(
    json_agg(
      jsonb_build_object(
        'id',         cn.id,
        'name',       cn.name,
        'iso2',       cn.iso2,
        'iso3',       cn.iso3,
        'num_code',   cn.num_code, 
        'created_at', cn.created_at,
        'updated_at', cn.updated_at
      )
    ) filter (where cn.id is not null), '[]') as countries
  from currency cr
  left join country_currency cc on cr.id = cc.currency_id
  left join country cn ON cc.country_id = cn.id
 group by cr.id, cr.name
 order by $1::text
 limit  $3 offset $2
`

type CurrencyCountrySelectParams struct {
	SqlOrder  string `json:"sql_order"`
	SqlOffset int32  `json:"sql_offset"`
	SqlLimit  int32  `json:"sql_limit"`
}

type CurrencyCountrySelectRow struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Code      string           `json:"code"`
	NumCode   int16            `json:"num_code"`
	Symbol    string           `json:"symbol"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Countries interface{}      `json:"countries"`
}

// CurrencyCountrySelect
//
//	select cr.id, cr.name, cr.code, cr.num_code, cr.symbol, cr.created_at, cr.updated_at,
//	  coalesce(
//	    json_agg(
//	      jsonb_build_object(
//	        'id',         cn.id,
//	        'name',       cn.name,
//	        'iso2',       cn.iso2,
//	        'iso3',       cn.iso3,
//	        'num_code',   cn.num_code,
//	        'created_at', cn.created_at,
//	        'updated_at', cn.updated_at
//	      )
//	    ) filter (where cn.id is not null), '[]') as countries
//	  from currency cr
//	  left join country_currency cc on cr.id = cc.currency_id
//	  left join country cn ON cc.country_id = cn.id
//	 group by cr.id, cr.name
//	 order by $1::text
//	 limit  $3 offset $2
func (q *Queries) CurrencyCountrySelect(ctx context.Context, arg *CurrencyCountrySelectParams) ([]*CurrencyCountrySelectRow, error) {
	rows, err := q.db.Query(ctx, currencyCountrySelect, arg.SqlOrder, arg.SqlOffset, arg.SqlLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*CurrencyCountrySelectRow
	for rows.Next() {
		var i CurrencyCountrySelectRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Code,
			&i.NumCode,
			&i.Symbol,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Countries,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const currencyDeleteByID = `-- name: CurrencyDeleteByID :one
delete from currency c where c.id = $1 returning id
`

// CurrencyDeleteByID
//
//	delete from currency c where c.id = $1 returning id
func (q *Queries) CurrencyDeleteByID(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRow(ctx, currencyDeleteByID, id)
	err := row.Scan(&id)
	return id, err
}

const currencyNew = `-- name: CurrencyNew :one
insert into currency (
	name, code, num_code, symbol
) values (
	$1, $2, $3, $4
) returning id, name, code, num_code, symbol, created_at, updated_at
`

type CurrencyNewParams struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	NumCode int16  `json:"num_code"`
	Symbol  string `json:"symbol"`
}

// CurrencyNew
//
//	insert into currency (
//		name, code, num_code, symbol
//	) values (
//		$1, $2, $3, $4
//	) returning id, name, code, num_code, symbol, created_at, updated_at
func (q *Queries) CurrencyNew(ctx context.Context, arg *CurrencyNewParams) (*Currency, error) {
	row := q.db.QueryRow(ctx, currencyNew,
		arg.Name,
		arg.Code,
		arg.NumCode,
		arg.Symbol,
	)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.NumCode,
		&i.Symbol,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const currencySelect = `-- name: CurrencySelect :many
select id, name, code, num_code, symbol, created_at, updated_at
  from currency c
 order by $1::text
 limit $3 offset $2
`

type CurrencySelectParams struct {
	SqlOrder  string `json:"sql_order"`
	SqlOffset int32  `json:"sql_offset"`
	SqlLimit  int32  `json:"sql_limit"`
}

// CurrencySelect
//
//	select id, name, code, num_code, symbol, created_at, updated_at
//	  from currency c
//	 order by $1::text
//	 limit $3 offset $2
func (q *Queries) CurrencySelect(ctx context.Context, arg *CurrencySelectParams) ([]*Currency, error) {
	rows, err := q.db.Query(ctx, currencySelect, arg.SqlOrder, arg.SqlOffset, arg.SqlLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Currency
	for rows.Next() {
		var i Currency
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Code,
			&i.NumCode,
			&i.Symbol,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const currencySelectByID = `-- name: CurrencySelectByID :one
select id, name, code, num_code, symbol, created_at, updated_at from currency c where c.id = $1
`

// CurrencySelectByID
//
//	select id, name, code, num_code, symbol, created_at, updated_at from currency c where c.id = $1
func (q *Queries) CurrencySelectByID(ctx context.Context, id string) (*Currency, error) {
	row := q.db.QueryRow(ctx, currencySelectByID, id)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.NumCode,
		&i.Symbol,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const currencyUpdateByID = `-- name: CurrencyUpdateByID :one
update currency
   set name = $1, code = $2, num_code = $3, symbol = $4
 where id = $5
       returning id, name, code, num_code, symbol, created_at, updated_at
`

type CurrencyUpdateByIDParams struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	NumCode int16  `json:"num_code"`
	Symbol  string `json:"symbol"`
	ID      string `json:"id"`
}

// CurrencyUpdateByID
//
//	update currency
//	   set name = $1, code = $2, num_code = $3, symbol = $4
//	 where id = $5
//	       returning id, name, code, num_code, symbol, created_at, updated_at
func (q *Queries) CurrencyUpdateByID(ctx context.Context, arg *CurrencyUpdateByIDParams) (*Currency, error) {
	row := q.db.QueryRow(ctx, currencyUpdateByID,
		arg.Name,
		arg.Code,
		arg.NumCode,
		arg.Symbol,
		arg.ID,
	)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.NumCode,
		&i.Symbol,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
