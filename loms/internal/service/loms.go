package service

type Item struct {
	SKU   uint32
	Count uint16
}

type Order struct {
	Status string
	User   int64
	Items  []Item
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Stocks struct {
	Stocks []Stock
}

type Service struct {
}

func New() *Service {
	return &Service{}
}

var DummyOrder = Order{
	Status: "new",
	User:   1,
	Items: []Item{
		{
			SKU:   1076963,
			Count: 10},
		{
			SKU:   1148162,
			Count: 5,
		},
		{
			SKU:   1625903,
			Count: 1,
		},
	},
}

var DummyStocks = Stocks{
	Stocks: []Stock{
		{
			WarehouseID: 1,
			Count:       5,
		},
		{
			WarehouseID: 2,
			Count:       7,
		},
	},
}
