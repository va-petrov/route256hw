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

func TestPurchase(t *testing.T) {
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

		userID = gofakeit.Int64()

		items = []model.Item{
			{SKU: 5097510, Count: 10},
			{SKU: 1076963, Count: 20},
		}

		order = model.Order{
			User:  userID,
			Items: items,
		}

		orderID = gofakeit.Int64()

		cartRepoError  = errors.New("getting cart from db")
		cartCleanError = errors.New("cleaning cart")
		lomsError      = errors.New("creating order")
	)

	tests := []struct {
		name               string
		args               args
		want               int64
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
			want: orderID,
			err:  nil,
			cartRepoMock: func(mc *minimock.Controller) cartRepo.CartRepo {
				mock := cartRepoMocks.NewCartRepoMock(mc)
				mock.GetCartMock.Expect(ctx, userID).Return(items, nil)
				mock.CleanCartMock.Expect(ctx, userID).Return(nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) lomsClient.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, order).Return(orderID, nil)
				return mock
			},
			productsClientMock: func(mc *minimock.Controller) productsClient.Client {
				mock := productsClientMocks.NewClientMock(mc)
				return mock
			},
		},
		{
			name: "empty cart test",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: -1,
			err:  ErrEmptyCart,
			cartRepoMock: func(mc *minimock.Controller) cartRepo.CartRepo {
				mock := cartRepoMocks.NewCartRepoMock(mc)
				mock.GetCartMock.Expect(ctx, userID).Return(nil, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) lomsClient.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				return mock
			},
			productsClientMock: func(mc *minimock.Controller) productsClient.Client {
				mock := productsClientMocks.NewClientMock(mc)
				return mock
			},
		},
		{
			name: "cart repo error",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: -1,
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
				return mock
			},
		},
		{
			name: "loms service error",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: -1,
			err:  lomsError,
			cartRepoMock: func(mc *minimock.Controller) cartRepo.CartRepo {
				mock := cartRepoMocks.NewCartRepoMock(mc)
				mock.GetCartMock.Expect(ctx, userID).Return(items, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) lomsClient.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				mock.CreateOrderMock.Return(-1, lomsError)
				return mock
			},
			productsClientMock: func(mc *minimock.Controller) productsClient.Client {
				mock := productsClientMocks.NewClientMock(mc)
				return mock
			},
		},
		{
			name: "cart clean error",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: -1,
			err:  cartCleanError,
			cartRepoMock: func(mc *minimock.Controller) cartRepo.CartRepo {
				mock := cartRepoMocks.NewCartRepoMock(mc)
				mock.GetCartMock.Expect(ctx, userID).Return(items, nil)
				mock.CleanCartMock.Expect(ctx, userID).Return(cartCleanError)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) lomsClient.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, order).Return(orderID, nil)
				return mock
			},
			productsClientMock: func(mc *minimock.Controller) productsClient.Client {
				mock := productsClientMocks.NewClientMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(tt.lomsClientMock(mc), tt.productsClientMock(mc), tt.cartRepoMock(mc))

			res, err := service.Purchase(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}
