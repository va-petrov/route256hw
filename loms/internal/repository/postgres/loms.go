package postgres

import (
	"context"
	"route256/loms/internal/service"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LOMSRepo struct {
	Pool *pgxpool.Pool
}

func NewLOMSRepo(pool *pgxpool.Pool) *LOMSRepo {
	return &LOMSRepo{
		Pool: pool,
	}
}

type Stock struct {
	WarehouseID int64  `db:"warehouseid"`
	SKU         uint32 `db:"sku"`
	Count       uint64 `db:"count"`
}

type Order struct {
	OrderID int64 `db:"orderid"`
	UserID  int64 `db:"userid"`
	status  int16 `db:"status"`
}

type OrderItem struct {
	OrderID int64  `db:"orderid"`
	SKU     uint32 `db:"sku"`
	Count   uint16 `db:"count"`
}

func (L LOMSRepo) GetStocks(ctx context.Context, sku uint32, checkReservations bool) ([]service.Stock, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var items []Stock
	if checkReservations {
		rawQuery := "SELECT s.warehouseid, s.sku, sum(s.count) - COALESCE(sum(r.count), 0) as count FROM stocks as s LEFT JOIN reservations as r ON r.warehouseid = s.warehouseid AND r.sku = s.sku AND r.active_until > now() WHERE s.sku = $1 GROUP BY s.warehouseid, s.sku"
		if err := pgxscan.Select(ctx, L.Pool, &items, rawQuery, sku); err != nil {
			return nil, err
		}

	} else {
		query := psql.Select("warehouseid", "sku", "count").From("stocks").Where(sq.Eq{"sku": sku})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return nil, err
		}
		if err := pgxscan.Select(ctx, L.Pool, &items, rawQuery, args...); err != nil {
			return nil, err
		}
	}

	if len(items) > 0 {
		result := make([]service.Stock, len(items))
		for i, item := range items {
			result[i] = service.Stock{
				WarehouseID: item.WarehouseID,
				Count:       item.Count,
			}
		}
		return result, nil
	} else {
		return nil, nil
	}
}

func (L LOMSRepo) ShipStock(ctx context.Context, sku uint32, warehouseID int64, count uint16) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("count").From("stocks").Where(sq.Eq{"sku": sku, "warehouseid": warehouseID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	var stock int64
	if err := L.Pool.QueryRow(ctx, rawQuery, args...).Scan(&stock); err != nil {
		return err
	}
	stock -= int64(count)
	if stock < 0 {
		return service.ErrInsufficientStocks
	}
	if stock == 0 {
		query := psql.Delete("stocks").Where(sq.Eq{"sku": sku, "warehouseid": warehouseID})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := L.Pool.Exec(ctx, rawQuery, args...); err != nil {
			return err
		}
	} else {
		query := psql.Update("stocks").Set("count", stock).Where(sq.Eq{"sku": sku, "warehouseid": warehouseID})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := L.Pool.Exec(ctx, rawQuery, args...); err != nil {
			return err
		}
	}
	return nil
}

func (L LOMSRepo) MakeReserve(ctx context.Context, orderID int64, sku uint32, warehouseID int64, count uint64) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("reservations").Columns("sku", "warehouseid", "orderid", "count").Values(sku, warehouseID, orderID, count)
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	if _, err := L.Pool.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}

func (L LOMSRepo) GetReserves(ctx context.Context, orderID int64) ([]service.Stock, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("sku", "warehouseid", "count").From("reservations").Where(sq.Eq{"orderid": orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var stocks []Stock
	if err := pgxscan.Select(ctx, L.Pool, &stocks, rawQuery, args...); err != nil {
		return nil, err
	}
	result := make([]service.Stock, len(stocks))
	for i, stock := range stocks {
		result[i] = service.Stock{
			SKU:         stock.SKU,
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		}
	}
	return result, nil
}

func (L LOMSRepo) CancelReservationsForOrder(ctx context.Context, orderID int64) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("reservations").Where(sq.Eq{"orderid": orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	if _, err := L.Pool.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}

func (L LOMSRepo) CreateOrder(ctx context.Context, order service.Order) (int64, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("orders").Columns("userid").Values(order.User).Suffix("RETURNING orderid")
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return -1, err
	}
	var orderID int64
	if err := L.Pool.QueryRow(ctx, rawQuery, args...).Scan(&orderID); err != nil {
		return -1, err
	}
	query = psql.Insert("orders_items").Columns("orderid", "sku", "count")
	for _, item := range order.Items {
		query = query.Values(orderID, item.SKU, item.Count)
	}
	rawQuery, args, err = query.ToSql()
	if err != nil {
		return -1, err
	}
	if _, err := L.Pool.Exec(ctx, rawQuery, args...); err != nil {
		return -1, err
	}
	return orderID, nil
}

func (L LOMSRepo) GetOrder(ctx context.Context, orderID int64) (*service.Order, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("orderid", "userid", "status").From("orders").Where(sq.Eq{"orderid": orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var order Order
	if err := L.Pool.QueryRow(ctx, rawQuery, args...).Scan(&order.OrderID, &order.UserID, &order.status); err != nil {
		return nil, err
	}
	query = psql.Select("sku", "count").From("orders_items").Where(sq.Eq{"orderid": orderID})
	rawQuery, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	var items []OrderItem
	err = pgxscan.Select(ctx, L.Pool, &items, rawQuery, args...)
	if err != nil {
		return nil, err
	}
	result := service.Order{
		User:  order.UserID,
		Items: make([]service.Item, len(items)),
	}
	/* 0 - created, 1 - awaiting payment 2 - payed, -1 - failed, -2 - cancelled */
	switch order.status {
	case -2:
		result.Status = "cancelled"
	case -1:
		result.Status = "failed"
	case 0:
		result.Status = "new"
	case 1:
		result.Status = "awaiting payment"
	case 2:
		result.Status = "payed"
	}
	for i, item := range items {
		result.Items[i] = service.Item{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}
	return &result, nil
}

func (L LOMSRepo) SetStatusOrder(ctx context.Context, orderID int64, status string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var orderStatus int16
	switch status {
	case "cancelled":
		orderStatus = -2
	case "failed":
		orderStatus = -1
	case "new":
		orderStatus = 0
	case "awaiting payment":
		orderStatus = 1
	case "payed":
		orderStatus = 2
	}
	query := psql.Update("orders").Set("status", orderStatus).Where(sq.Eq{"orderid": orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	if _, err := L.Pool.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}
