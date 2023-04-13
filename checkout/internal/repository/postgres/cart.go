package postgres

//go:generate sh -c "mkdir -p mocks && rm -rf mocks/cart_repo_minimock.go"
//go:generate minimock -i CartRepo -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"route256/checkout/internal/service/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
)

type CartRepo interface {
	GetCartItem(ctx context.Context, user int64, sku uint32) (*model.Item, error)
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error
	GetCart(ctx context.Context, user int64) ([]model.Item, error)
	CleanCart(ctx context.Context, user int64) error
}

type Cart struct {
	UserID int64  `db:"userid"`
	SKU    uint32 `db:"sku"`
	Count  uint16 `db:"count"`
}

type cartRepo struct {
	Pool *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewCartRepo(pool *pgxpool.Pool) CartRepo {
	return &cartRepo{
		Pool: pool,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

const (
	tableCarts          = "carts"
	fieldCartItemUserID = "userid"
	fieldCartItemSKU    = "sku"
	fieldCartItemCount  = "count"
)

var cartItemFields = []string{
	fieldCartItemUserID,
	fieldCartItemSKU,
	fieldCartItemCount,
}

func (c cartRepo) GetCartItem(ctx context.Context, user int64, sku uint32) (*model.Item, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CartRepo.GetCartItem")
	defer span.Finish()

	span.SetTag("userID", user)
	span.SetTag("SKU", sku)

	query := c.psql.Select(cartItemFields...).From(tableCarts).Where(sq.Eq{fieldCartItemUserID: user, fieldCartItemSKU: sku})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var items []Cart
	if err := pgxscan.Select(ctx, c.Pool, &items, rawQuery, args...); err != nil {
		return nil, err
	}

	if len(items) > 0 {
		return &model.Item{
			SKU:   items[0].SKU,
			Count: items[0].Count,
		}, nil
	} else {
		return nil, nil
	}
}

const (
	addToCartQuerySuffix = "ON CONFLICT (" + fieldCartItemUserID + "," + fieldCartItemSKU + ") DO UPDATE SET " + fieldCartItemCount + " = " + tableCarts + "." + fieldCartItemCount + " + ?"
)

func (c cartRepo) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CartRepo.AddToCart")
	defer span.Finish()

	span.SetTag("userID", user)
	span.SetTag("SKU", sku)
	span.SetTag("count", count)

	query := c.psql.Insert(tableCarts).Columns(cartItemFields...).Values(user, sku, count)
	query = query.Suffix(addToCartQuerySuffix, count)
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	if _, err := c.Pool.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}

func (c cartRepo) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CartRepo.DeleteFromCart")
	defer span.Finish()

	span.SetTag("userID", user)
	span.SetTag("SKU", sku)
	span.SetTag("count", count)

	item, err := c.GetCartItem(ctx, user, sku)
	if err != nil {
		return err
	}

	if item == nil {
		return nil
	} else if item.Count <= count {
		query := c.psql.Delete(tableCarts).Where(sq.Eq{fieldCartItemUserID: user, fieldCartItemSKU: sku})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := c.Pool.Exec(ctx, rawQuery, args...); err != nil {
			return err
		}
		return nil
	} else {
		query := c.psql.Update(tableCarts).Set(fieldCartItemCount, item.Count-count).Where(sq.Eq{fieldCartItemUserID: user, fieldCartItemSKU: sku})
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

func (c cartRepo) GetCart(ctx context.Context, user int64) ([]model.Item, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CartRepo.GetCart")
	defer span.Finish()

	span.SetTag("userID", user)

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
		result := make([]model.Item, len(items))
		for i, item := range items {
			result[i] = model.Item{
				SKU:   item.SKU,
				Count: item.Count,
			}
		}
		return result, nil
	} else {
		return nil, nil
	}
}

func (c cartRepo) CleanCart(ctx context.Context, user int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CartRepo.CleanCart")
	defer span.Finish()

	span.SetTag("userID", user)

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
