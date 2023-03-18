-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks (warehouseid, sku, count) VALUES (2, 1076963, 79) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 1076963, 99) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (2, 1076963, 92) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 1076963, 50) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (2, 1076963, 66) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 1076963, 26) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 1076963, 59) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 1148162, 25) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (8, 1148162, 24) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (5, 1148162, 57) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 1148162, 26) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 1148162, 18) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (6, 1625903, 96) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 1625903, 12) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (9, 1625903, 81) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (9, 1625903, 68) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 1625903, 88) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 2956315, 18) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 2956315, 47) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (6, 2956315, 27) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (5, 2956315, 4) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 2956315, 11) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 2956315, 16) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 2956315, 67) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 2956315, 96) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (6, 2958025, 45) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 2958025, 38) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 2958025, 36) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (2, 2958025, 74) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (9, 3596599, 55) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 3596599, 37) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (6, 3596599, 74) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (2, 3596599, 90) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (5, 3618852, 26) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 3618852, 22) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 3618852, 35) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (8, 3618852, 62) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 3618852, 88) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 3618852, 16) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 3618852, 99) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (6, 4288068, 47) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 4288068, 85) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (6, 4288068, 70) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (2, 4288068, 84) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 4288068, 69) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 4288068, 4) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 4288068, 64) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 4288068, 90) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (9, 4465995, 14) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (9, 4465995, 74) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 4465995, 67) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 4465995, 83) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (9, 4487693, 69) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 4669069, 10) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 4678816, 63) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 4687693, 7) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (3, 4687693, 20) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (2, 4687693, 4) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (9, 4687693, 70) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (4, 4687693, 39) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 4687693, 2) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (1, 4687693, 20) ON CONFLICT (warehouseid, sku) DO NOTHING;
INSERT INTO stocks (warehouseid, sku, count) VALUES (7, 5097510, 24) ON CONFLICT (warehouseid, sku) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stocks;
-- +goose StatementEnd
