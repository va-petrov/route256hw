package postgres

import (
	"context"
	"route256/loms/internal/repository/postgres/tranman"
	"route256/loms/internal/service"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type LOMSRepo struct {
	tranman.QueryEngineProvider
	psql sq.StatementBuilderType
}

func NewLOMSRepo(provider tranman.QueryEngineProvider) *LOMSRepo {
	return &LOMSRepo{
		QueryEngineProvider: provider,
		psql:                sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

const (
	tableStocks                  = "stocks"
	fieldStockWarehouseID        = "warehouseid"
	fieldStockSKU                = "sku"
	fieldStockCount              = "count"
	tableReservations            = "reservations"
	fieldReservationsWarehouseID = "warehouseid"
	fieldReservationsSKU         = "sku"
	fieldReservationsCount       = "count"
	fieldReservationsActiveUntil = "active_until"
	fieldReservationsOrderID     = "orderid"
)

var StocksFields = []string{
	fieldStockWarehouseID,
	fieldStockSKU,
	fieldStockCount,
}

var ReservationsFields = []string{
	fieldReservationsWarehouseID,
	fieldReservationsSKU,
	fieldReservationsCount,
	fieldReservationsOrderID,
}

type Stock struct {
	WarehouseID int64  `db:"warehouseid"`
	SKU         uint32 `db:"sku"`
	Count       uint64 `db:"count"`
}

const (
	OrderStatusCancelled       = -2
	OrderStatusFailed          = -1
	OrderStatusNew             = 0
	OrderStatusAwaitingPayment = 1
	OrderStatusPayed           = 2
)

const (
	tableOrders         = "orders"
	fieldOrderOrderID   = "orderid"
	fieldOrderUserID    = "userid"
	fieldOrderStatus    = "status"
	fieldOrderCreatedAt = "created_at"
)

var OrdersFields = []string{
	fieldOrderOrderID,
	fieldOrderUserID,
	fieldOrderStatus,
}

type Order struct {
	OrderID int64 `db:"orderid"`
	UserID  int64 `db:"userid"`
	status  int16 `db:"status"`
}

const (
	tableOrdersItems        = "orders_items"
	fieldOrdersItemsOrderID = "orderid"
	fieldOrdersItemsSKU     = "sku"
	fieldOrdersItemsCount   = "count"
)

var OrdersItemsFields = []string{
	fieldOrdersItemsOrderID,
	fieldOrdersItemsSKU,
	fieldOrdersItemsCount,
}

type OrderItem struct {
	OrderID int64  `db:"orderid"`
	SKU     uint32 `db:"sku"`
	Count   uint16 `db:"count"`
}

func (L LOMSRepo) GetStocks(ctx context.Context, sku uint32, checkReservations bool) ([]service.Stock, error) {
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	var items []Stock
	if checkReservations {
		rawQuery := "SELECT s." + fieldStockWarehouseID + ", s." + fieldStockSKU + ", sum(s." + fieldStockCount + ") - COALESCE(sum(r." + fieldStockCount + "), 0) as count FROM " + tableStocks + " as s LEFT JOIN " + tableReservations + " as r ON r." + fieldReservationsWarehouseID + " = s." + fieldStockWarehouseID + " AND r." + fieldReservationsSKU + " = s." + fieldStockSKU + " AND r." + fieldReservationsActiveUntil + " > now() WHERE s." + fieldStockSKU + " = $1 GROUP BY s." + fieldStockWarehouseID + ", s." + fieldStockSKU
		if err := pgxscan.Select(ctx, db, &items, rawQuery, sku); err != nil {
			return nil, err
		}

	} else {
		query := L.psql.Select(StocksFields...).From(tableStocks).Where(sq.Eq{fieldStockSKU: sku})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return nil, err
		}
		if err := pgxscan.Select(ctx, db, &items, rawQuery, args...); err != nil {
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
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	query := L.psql.Select(fieldStockCount).From(tableStocks).Where(sq.Eq{fieldStockSKU: sku, fieldStockWarehouseID: warehouseID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	var stock int64
	if err := db.QueryRow(ctx, rawQuery, args...).Scan(&stock); err != nil {
		return err
	}
	stock -= int64(count)
	if stock < 0 {
		return service.ErrInsufficientStocks
	}
	if stock == 0 {
		query := L.psql.Delete(tableStocks).Where(sq.Eq{fieldStockSKU: sku, fieldStockWarehouseID: warehouseID})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := db.Exec(ctx, rawQuery, args...); err != nil {
			return err
		}
	} else {
		query := L.psql.Update(tableStocks).Set(fieldStockCount, stock).Where(sq.Eq{fieldStockSKU: sku, fieldStockWarehouseID: warehouseID})
		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}
		if _, err := db.Exec(ctx, rawQuery, args...); err != nil {
			return err
		}
	}
	return nil
}

func (L LOMSRepo) MakeReserve(ctx context.Context, orderID int64, sku uint32, warehouseID int64, count uint64) error {
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	query := L.psql.Insert(tableReservations).Columns(ReservationsFields...).Values(sku, warehouseID, orderID, count)
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	if _, err := db.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}

func (L LOMSRepo) GetReserves(ctx context.Context, orderID int64) ([]service.Stock, error) {
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	query := L.psql.Select(ReservationsFields...).From(tableReservations).Where(sq.Eq{fieldReservationsOrderID: orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var stocks []Stock
	if err := pgxscan.Select(ctx, db, &stocks, rawQuery, args...); err != nil {
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
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	query := L.psql.Delete(tableReservations).Where(sq.Eq{fieldReservationsOrderID: orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	if _, err := db.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}

func (L LOMSRepo) CreateOrder(ctx context.Context, order service.Order) (int64, error) {
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	query := L.psql.Insert(tableOrders).Columns(fieldOrderUserID).Values(order.User).Suffix("RETURNING " + fieldOrderOrderID)
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return -1, err
	}
	var orderID int64
	if err := db.QueryRow(ctx, rawQuery, args...).Scan(&orderID); err != nil {
		return -1, err
	}
	query = L.psql.Insert(tableOrdersItems).Columns(OrdersItemsFields...)
	for _, item := range order.Items {
		query = query.Values(orderID, item.SKU, item.Count)
	}
	rawQuery, args, err = query.ToSql()
	if err != nil {
		return -1, err
	}
	if _, err := db.Exec(ctx, rawQuery, args...); err != nil {
		return -1, err
	}
	return orderID, nil
}

func (L LOMSRepo) GetOrder(ctx context.Context, orderID int64) (*service.Order, error) {
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	query := L.psql.Select(OrdersFields...).From(tableOrders).Where(sq.Eq{fieldOrderOrderID: orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var order Order
	if err := db.QueryRow(ctx, rawQuery, args...).Scan(&order.OrderID, &order.UserID, &order.status); err != nil {
		return nil, err
	}
	query = L.psql.Select(fieldOrdersItemsSKU, fieldOrdersItemsCount).From(tableOrdersItems).Where(sq.Eq{fieldOrdersItemsOrderID: orderID})
	rawQuery, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	var items []OrderItem
	err = pgxscan.Select(ctx, db, &items, rawQuery, args...)
	if err != nil {
		return nil, err
	}
	result := service.Order{
		User:  order.UserID,
		Items: make([]service.Item, len(items)),
	}
	switch order.status {
	case OrderStatusCancelled:
		result.Status = service.OrderStatusCancelled
	case OrderStatusFailed:
		result.Status = service.OrderStatusFailed
	case OrderStatusNew:
		result.Status = service.OrderStatusNew
	case OrderStatusAwaitingPayment:
		result.Status = service.OrderStatusAwaitingPayment
	case OrderStatusPayed:
		result.Status = service.OrderStatusPayed
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
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	var orderStatus int16
	switch status {
	case service.OrderStatusCancelled:
		orderStatus = OrderStatusCancelled
	case service.OrderStatusFailed:
		orderStatus = OrderStatusFailed
	case service.OrderStatusNew:
		orderStatus = OrderStatusNew
	case service.OrderStatusAwaitingPayment:
		orderStatus = OrderStatusAwaitingPayment
	case service.OrderStatusPayed:
		orderStatus = OrderStatusPayed
	}
	query := L.psql.Update(tableOrders).Set(fieldOrderStatus, orderStatus).Where(sq.Eq{fieldOrderOrderID: orderID})
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	if _, err := db.Exec(ctx, rawQuery, args...); err != nil {
		return err
	}
	return nil
}

func (L LOMSRepo) CancelUnpayedOrders(ctx context.Context) error {
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	rawQuery := "UPDATE " + tableOrders + " SET " + fieldOrderStatus + " = -1 WHERE " + fieldOrderStatus + " = 0 and " + fieldOrderCreatedAt + " + interval '10 minutes' <= now()"
	if _, err := db.Exec(ctx, rawQuery); err != nil {
		return err
	}
	return nil
}

func (L LOMSRepo) DeleteStaleReservations(ctx context.Context) error {
	db := L.QueryEngineProvider.GetQueryEngine(ctx)
	rawQuery := "DELETE FROM " + tableReservations + " WHERE " + fieldReservationsActiveUntil + " <= now()"
	if _, err := db.Exec(ctx, rawQuery); err != nil {
		return err
	}
	return nil
}
