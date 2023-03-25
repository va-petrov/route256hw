package postgres

import (
	"context"
	"route256/checkout/internal/service"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Cart struct {
	UserID int64  `db:"userid"`
	SKU    uint32 `db:"sku"`
	Count  uint16 `db:"count"`
}

type CartRepo struct {
	Pool *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewCartRepo(pool *pgxpool.Pool) *CartRepo {
	return &CartRepo{
		Pool: pool,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

const (
	tableCarts          = "carts"
	fieldCartItemUserID = "userid"
	fieldCartItemmSKU   = "sku"
	fieldCartItemCount  = "count"
)

var cartItemFields = []string{
	fieldCartItemUserID,
	fieldCartItemmSKU,
	fieldCartItemCount,
}

func (c CartRepo) GetCartItem(ctx context.Context, user int64, sku uint32) (*service.Item, error) {
	query := c.psql.Select(cartItemFields...).From(tableCarts).Where(sq.Eq{fieldCartItemUserID: user, fieldCartItemmSKU: sku})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var items []Cart
	if err := pgxscan.Select(ctx, c.Pool, &items, rawQuery, args...); err != nil {
		return nil, err
	}

	if len(items) > 0 {
		return &service.Item{
			SKU:   items[0].SKU,
			Count: items[0].Count,
		}, nil
	} else {
		return nil, nil
	}
}

func (c CartRepo) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	query := c.psql.Insert(tableCarts).Columns(cartItemFields...).Values(user, sku, count)
	query = query.Suffix("ON CONFLICT ("+fieldCartItemUserID+","+fieldCartItemmSKU+") DO UPDATE SET "+fieldCartItemCount+" = "+tableCarts+"."+fieldCartItemCount+" + ?", count)
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	if _, err := c.Pool.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}

func (c CartRepo) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	item, err := c.GetCartItem(ctx, user, sku)
	if err != nil {
		return err
	}

	if item == nil {
		return nil
	} else if item.Count <= count {
		query := c.psql.Delete(tableCarts).Where(sq.Eq{fieldCartItemUserID: user, fieldCartItemmSKU: sku})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := c.Pool.Exec(ctx, rawQuery, args...); err != nil {
			return err
		}
		return nil
	} else {
		query := c.psql.Update(tableCarts).Set(fieldCartItemCount, item.Count-count).Where(sq.Eq{fieldCartItemUserID: user, fieldCartItemmSKU: sku})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := c.Pool.Exec(ctx, rawQuery, args...); err != nil {
			return err
		}
		return nil
	}
}

func (c CartRepo) GetCart(ctx context.Context, user int64) ([]service.Item, error) {
	query := c.psql.Select(cartItemFields...).From(tableCarts).Where(sq.Eq{fieldCartItemUserID: user})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var items []Cart
	if err := pgxscan.Select(ctx, c.Pool, &items, rawQuery, args...); err != nil {
		return nil, err
	}

	if len(items) > 0 {
		result := make([]service.Item, len(items))
		for i, item := range items {
			result[i] = service.Item{
				SKU:   item.SKU,
				Count: item.Count,
			}
		}
		return result, nil
	} else {
		return nil, nil
	}
}

func (c CartRepo) CleanCart(ctx context.Context, user int64) error {
	query := c.psql.Delete(tableCarts).Where(sq.Eq{fieldCartItemUserID: user})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	if _, err := c.Pool.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}
