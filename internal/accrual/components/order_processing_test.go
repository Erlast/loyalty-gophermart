package components

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/handlers"
	"github.com/golang/mock/gomock"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"go.uber.org/zap/zaptest"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/stretchr/testify/assert"
)

func TestOrderProcessing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	store.EXPECT().GetRegisteredOrders(gomock.Any()).Return([]int64{1}, nil).AnyTimes()
	store.EXPECT().FetchRewardRules(gomock.Any()).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil).AnyTimes()
	store.EXPECT().UpdateOrderStatus(gomock.Any(), int64(1), helpers.StatusProcessing).Return(nil).AnyTimes()
	store.EXPECT().FetchProducts(gomock.Any(), int64(1)).Return([]models.Items{
		{Description: "test product", Price: 100.00},
	}, nil).AnyTimes()
	store.EXPECT().SaveOrderPoints(gomock.Any(), int64(1), []float32{10}).Return(nil).AnyTimes()

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)

	cancel()

	assert.True(t, true)
}

func TestOrderProcessing_ErrorGetOrders(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core).Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)

	store.EXPECT().GetRegisteredOrders(gomock.Any()).Return(nil, errors.New("database error"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1, "expected one log entry")
	assert.Equal(t, "ошибка при попытке выбрать новые заказы: database error", logs[0].Message)
}

func TestFetchRewardRules_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core).Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)

	store.EXPECT().GetRegisteredOrders(gomock.Any()).Return([]int64{1}, nil)
	store.EXPECT().FetchRewardRules(gomock.Any()).Return(nil, errors.New("database error"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1, "expected one log entry")
	assert.Equal(t, "не могу выбрать правила начислений", logs[0].Message)
}

func TestUpdateOrderStatus_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)

	store.EXPECT().GetRegisteredOrders(gomock.Any()).Return([]int64{1}, nil)
	store.EXPECT().FetchRewardRules(gomock.Any()).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.EXPECT().UpdateOrderStatus(
		ctx, int64(1),
		helpers.StatusProcessing,
	).Return(errors.New("failed to update order status"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
}

func TestFetchProducts_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core).Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)

	store.EXPECT().GetRegisteredOrders(gomock.Any()).Return([]int64{1}, nil)
	store.EXPECT().FetchRewardRules(gomock.Any()).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.EXPECT().UpdateOrderStatus(ctx, int64(1), helpers.StatusProcessing).Return(nil)
	store.EXPECT().FetchProducts(ctx, int64(1)).Return(nil, errors.New("failed to fetch products"))
	store.EXPECT().UpdateOrderStatus(ctx, int64(1), helpers.StatusInvalid).Return(nil)

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1, "expected one log entry")
	assert.Equal(t, "не могу получить товары из заказаfailed to fetch products", logs[0].Message)
}

func TestFetchProducts_ErrorUpdateStatus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)

	store.EXPECT().GetRegisteredOrders(gomock.Any()).Return([]int64{1}, nil)
	store.EXPECT().FetchRewardRules(gomock.Any()).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.EXPECT().UpdateOrderStatus(ctx, int64(1), helpers.StatusProcessing).Return(nil)
	store.EXPECT().FetchProducts(gomock.Any(), int64(1)).Return(
		nil,
		errors.New("failed to fetch products"),
	)
	store.EXPECT().UpdateOrderStatus(ctx, int64(1), helpers.StatusInvalid).Return(
		errors.New("failed to update order status"),
	)

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
}

func TestSaveOrderPoints(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)

	store.EXPECT().GetRegisteredOrders(ctx).Return([]int64{1}, nil).AnyTimes()
	store.EXPECT().FetchRewardRules(gomock.Any()).Return([]models.Goods{
		{Match: "encrypt", Reward: 10, RewardType: "%"},
		{Match: "simple", Reward: 5, RewardType: "pt"},
		{Match: "simple", Reward: 5, RewardType: "no"},
	}, nil).AnyTimes()
	store.EXPECT().UpdateOrderStatus(ctx, int64(1), helpers.StatusProcessing).Return(nil).AnyTimes()
	store.EXPECT().FetchProducts(gomock.Any(), int64(1)).Return([]models.Items{
		{Description: "encrypt data", Price: 100.00},
		{Description: "simple product", Price: 50.00},
	}, nil).AnyTimes()
	store.EXPECT().SaveOrderPoints(gomock.Any(), int64(1), []float32{10.00, 5.00}).Return(nil).AnyTimes()

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
}

func TestSaveOrderPoints_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := handlers.NewMockStorage(ctrl)

	store.EXPECT().GetRegisteredOrders(ctx).Return([]int64{1}, nil)
	store.EXPECT().FetchRewardRules(gomock.Any()).Return([]models.Goods{
		{Match: "encrypt", Reward: 10, RewardType: "%"},
	}, nil)
	store.EXPECT().UpdateOrderStatus(ctx, int64(1), helpers.StatusProcessing).Return(nil)
	store.EXPECT().FetchProducts(gomock.Any(), int64(1)).Return([]models.Items{
		{Description: "encrypt data", Price: 100.00},
	}, nil)
	store.EXPECT().SaveOrderPoints(gomock.Any(), int64(1), []float32{10.00}).Return(errors.New("failed to fetch products"))
	store.EXPECT().UpdateOrderStatus(ctx, int64(1), helpers.StatusInvalid).Return(
		errors.New("failed to update order status"),
	)

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
}
