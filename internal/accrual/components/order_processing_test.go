package components

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
)

func TestOrderProcessing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := &storage.MockStorage{}
	logger := zaptest.NewLogger(t).Sugar()

	store.On("GetRegisteredOrders", mock.Anything).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.On("UpdateOrderStatus", mock.Anything, int64(1), helpers.StatusProcessing).Return(nil)
	store.On("FetchProducts", mock.Anything, int64(1)).Return([]models.Items{
		{Description: "test product", Price: 100.00},
	}, nil)
	store.On("SaveOrderPoints", mock.Anything, int64(1), []float32{10}).Return(nil)

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)

	cancel()

	store.AssertExpectations(t)
	assert.True(t, true)
}

func TestOrderProcessing_ErrorGetOrders(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core).Sugar()

	store := &storage.MockStorage{}

	store.On("GetRegisteredOrders", mock.Anything).Return(nil, errors.New("database error"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1, "expected one log entry")
	assert.Equal(t, "ошибка при попытке выбрать новые заказы: err database error", logs[0].Message)

	store.AssertExpectations(t)
}

func TestFetchRewardRules_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core).Sugar()

	store := &storage.MockStorage{}

	store.On("GetRegisteredOrders", mock.Anything).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return(nil, errors.New("database error"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1, "expected one log entry")
	assert.Equal(t, "не могу выбрать правила начислений", logs[0].Message)

	store.AssertExpectations(t)
}

func TestUpdateOrderStatus_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	store := &storage.MockStorage{}

	store.On("GetRegisteredOrders", mock.Anything).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.On(
		"UpdateOrderStatus",
		ctx, int64(1),
		helpers.StatusProcessing,
	).Return(errors.New("failed to update order status"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()

	store.AssertExpectations(t)
}

func TestFetchProducts_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core).Sugar()

	store := &storage.MockStorage{}

	store.On("GetRegisteredOrders", mock.Anything).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusProcessing).Return(nil)
	store.On("FetchProducts", ctx, int64(1)).Return(nil, errors.New("failed to fetch products"))
	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusInvalid).Return(nil)

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1, "expected one log entry")
	assert.Equal(t, "не могу получить товары из заказаerr failed to fetch products", logs[0].Message)

	store.AssertExpectations(t)
}

func TestFetchProducts_ErrorUpdateStatus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	store := &storage.MockStorage{}

	store.On("GetRegisteredOrders", mock.Anything).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusProcessing).Return(nil)
	store.On("FetchProducts", mock.Anything, int64(1)).Return(nil, errors.New("failed to fetch products"))
	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusInvalid).Return(errors.New("failed to update order status"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()

	store.AssertExpectations(t)
}

func TestSaveOrderPoints(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	store := &storage.MockStorage{}

	store.On("GetRegisteredOrders", ctx).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return([]models.Goods{
		{Match: "encrypt", Reward: 10, RewardType: "%"},
		{Match: "simple", Reward: 5, RewardType: "pt"},
		{Match: "simple", Reward: 5, RewardType: "no"},
	}, nil)
	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusProcessing).Return(nil)
	store.On("FetchProducts", mock.Anything, int64(1)).Return([]models.Items{
		{Description: "encrypt data", Price: 100.00},
		{Description: "simple product", Price: 50.00},
	}, nil)
	store.On("SaveOrderPoints", mock.Anything, int64(1), []float32{10.00, 5.00}).Return(nil)

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()

	store.AssertExpectations(t)
}

func TestSaveOrderPoints_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample().Sugar()

	store := &storage.MockStorage{}

	store.On("GetRegisteredOrders", ctx).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return([]models.Goods{
		{Match: "encrypt", Reward: 10, RewardType: "%"},
	}, nil)
	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusProcessing).Return(nil)
	store.On("FetchProducts", mock.Anything, int64(1)).Return([]models.Items{
		{Description: "encrypt data", Price: 100.00},
	}, nil)
	store.On("SaveOrderPoints", mock.Anything, int64(1), []float32{10.00}).Return(errors.New("failed to fetch products"))
	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusInvalid).Return(errors.New("failed to update order status"))

	go func() {
		OrderProcessing(ctx, store, logger)
	}()

	time.Sleep(2 * time.Second)
	cancel()

	store.AssertExpectations(t)
}
