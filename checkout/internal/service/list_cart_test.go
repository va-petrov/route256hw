package service

import (
	"context"
	lomsClient "route256/checkout/internal/clients/lomsclient"
	lomsClientMocks "route256/checkout/internal/clients/lomsclient/mocks"
	productsClient "route256/checkout/internal/clients/productsclient"
	productsClientMocks "route256/checkout/internal/clients/productsclient/mocks"
	cartRepo "route256/checkout/internal/repository/postgres"
	cartRepoMocks "route256/checkout/internal/repository/postgres/mocks"
	"route256/checkout/internal/service/model"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	type cartRepoMockFunc func(mc *minimock.Controller) cartRepo.CartRepo
	type lomsClientMockFunc func(mc *minimock.Controller) lomsClient.Client
	type productsClientMockFunc func(mc *minimock.Controller) productsClient.Client

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		items = []model.Item{
			{SKU: 5097510, Count: 10},
			{SKU: 1076963, Count: 20},
		}

		itemsInfo = []model.Product{
			{Name: "Двенадцатый день Рождества", Price: 2300},
			{Name: "Теория нравственных чувств | Смит Адам", Price: 3379},
		}

		cart = model.Cart{
			Items: []model.CartItem{
				{
					SKU: 5097510, Count: 10,
					Name: "Двенадцатый день Рождества", Price: 2300,
				},
				{
					SKU: 1076963, Count: 20,
					Name: "Теория нравственных чувств | Смит Адам", Price: 3379,
				},
			},
			TotalPrice: 90580,
		}

		userID = gofakeit.Int64()

		cartRepoError       = errors.New("carts db")
		productServiceError = errors.New("some product service error")
	)

	tests := []struct {
		name               string
		args               args
		want               model.Cart
		err                error
		cartRepoMock       cartRepoMockFunc
		lomsClientMock     lomsClientMockFunc
		productsClientMock productsClientMockFunc
	}{
		{
			name: "positive test",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: cart,
			err:  nil,
			cartRepoMock: func(mc *minimock.Controller) cartRepo.CartRepo {
				mock := cartRepoMocks.NewCartRepoMock(mc)
				mock.GetCartMock.Expect(ctx, userID).Return(items, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) lomsClient.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				return mock
			},
			productsClientMock: func(mc *minimock.Controller) productsClient.Client {
				mock := productsClientMocks.NewClientMock(mc)
				mock.GetProductsInfoMock.Set(func(ctx context.Context, items []model.CartItem) error {
					for i := range items {
						items[i].Name = itemsInfo[i].Name
						items[i].Price = itemsInfo[i].Price
					}
					return nil
				})
				return mock
			},
		},
		{
			name: "cart repo error",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: model.Cart{},
			err:  cartRepoError,
			cartRepoMock: func(mc *minimock.Controller) cartRepo.CartRepo {
				mock := cartRepoMocks.NewCartRepoMock(mc)
				mock.GetCartMock.Expect(ctx, userID).Return(nil, errors.New("some cartRepo error"))
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) lomsClient.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				return mock
			},
			productsClientMock: func(mc *minimock.Controller) productsClient.Client {
				mock := productsClientMocks.NewClientMock(mc)
				mock.GetProductsInfoMock.Set(func(ctx context.Context, items []model.CartItem) error {
					for i := range items {
						items[i].Name = itemsInfo[i].Name
						items[i].Price = itemsInfo[i].Price
					}
					return nil
				})
				return mock
			},
		},
		{
			name: "product service error",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: model.Cart{},
			err:  productServiceError,
			cartRepoMock: func(mc *minimock.Controller) cartRepo.CartRepo {
				mock := cartRepoMocks.NewCartRepoMock(mc)
				mock.GetCartMock.Expect(ctx, userID).Return(items, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) lomsClient.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				return mock
			},
			productsClientMock: func(mc *minimock.Controller) productsClient.Client {
				mock := productsClientMocks.NewClientMock(mc)
				mock.GetProductsInfoMock.Return(productServiceError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(tt.lomsClientMock(mc), tt.productsClientMock(mc), tt.cartRepoMock(mc))

			res, err := service.ListCart(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}
