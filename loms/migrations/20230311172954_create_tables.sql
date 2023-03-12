-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks
(
    sku           int4 NOT NULL,
    warehouseID   int8 NOT NULL,
    count         int8  NOT NULL,
    CONSTRAINT stocks_pk
        PRIMARY KEY (sku, warehouseID)
);

CREATE TABLE IF NOT EXISTS reservations
(
    sku           int4 NOT NULL,
    warehouseID   int8  NOT NULL,
    active_until  timestamp DEFAULT now() + interval '10 minutes',
    orderID       int8 NOT NULL,
    count         int8  NOT NULL,
    CONSTRAINT reservations_pk
        PRIMARY KEY (sku, warehouseID, active_until, orderID)
);

CREATE INDEX IF NOT EXISTS reservations_orderid_index
    ON reservations (orderid);

CREATE TABLE IF NOT EXISTS orders
(
    orderID          int8 NOT NULL,
    userID           int8 NOT NULL,
    status           int2 NOT NULL DEFAULT 0, /* 0 - created, 1 - awaiting payment 2 - payed, -1 - failed, -2 - cancelled */
    created_at  timestamp DEFAULT now(),
    CONSTRAINT orders_pk
        PRIMARY KEY (orderID)
);

CREATE TABLE IF NOT EXISTS orders_items
(
    orderID         int8 NOT NULL,
    sku             int4 NOT NULL,
    count           int2 NOT NULL,
    CONSTRAINT orders_items_pk
        PRIMARY KEY (orderID, sku)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS reservations;
DROP INDEX IF EXISTS reservations_orderid_index;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS orders_items;

-- +goose StatementEnd
