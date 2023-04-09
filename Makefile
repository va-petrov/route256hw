LOCAL_BIN:=$(CURDIR)/bin

build-all:
	cd checkout && GOOS=linux make build
	cd loms && GOOS=linux make build
	cd notifications && GOOS=linux make build

run-all: build-all
	docker compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd product-service && make precommit

install-go-deps-all:
	cd checkout && make install-go-deps
	cd loms && make install-go-deps
	cd product-service && make install-go-deps

get-go-deps-all:
	cd checkout && make get-go-deps
	cd loms && make get-go-deps
	cd product-service && make get-go-deps

vendor-proto-all:
	cd checkout && make vendor-proto
	cd loms && make vendor-proto
	cd product-service && make vendor-proto

generate-all:
	cd checkout && make generate
	cd loms && make generate
	cd product-service && make generate

migrations-status-all:
	cd checkout && make migrations-status
	cd loms && make migrations-status

.PHONY: logs
logs:
	mkdir -p logs/data
	touch logs/data/checkout.txt
	touch logs/data/loms.txt
	touch logs/data/notifications.txt
	touch logs/data/offsets.yaml
	sudo chmod -R 777 logs/data
	cd logs && sudo docker compose up
