-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS carts (
  userID bigint NOT NULL,
  sku integer NOT NULL,
  count smallint NOT NULL,
  PRIMARY KEY (userID, sku)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS carts;
-- +goose StatementEnd
