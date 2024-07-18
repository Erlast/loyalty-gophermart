package components

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
)

func TestOrderProcessing(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

	store.On("GetRegisteredOrders", mock.Anything).Return([]int64{1}, nil)
	store.On("FetchRewardRules", mock.Anything).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)
	store.On("UpdateOrderStatus", mock.Anything, int64(1), helpers.StatusProcessing).Return(nil)
	store.On("FetchProducts", mock.Anything, int64(1)).Return([]models.Items{
		{Description: "test product", Price: 100.00},
	}, nil)
	store.On("SaveOrderPoints", mock.Anything, int64(1), []float32{10}).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go OrderProcessing(ctx, store, newLogger)

	time.Sleep(2 * time.Second)

	store.AssertExpectations(t)
}

func TestFetchRewardRules(t *testing.T) {
	store := &storage.MockStorage{}
	ctx := context.Background()

	store.On("FetchRewardRules", ctx).Return([]models.Goods{
		{Match: "test", Reward: 10, RewardType: "%"},
	}, nil)

	rules, err := store.FetchRewardRules(ctx)

	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.Equal(t, "test", rules[0].Match)
	assert.Equal(t, int64(10), rules[0].Reward)
	assert.Equal(t, "%", rules[0].RewardType)

	store.AssertExpectations(t)
}

func TestFetchProducts(t *testing.T) {
	store := &storage.MockStorage{}
	ctx := context.Background()

	store.On("FetchProducts", ctx, int64(1)).Return([]models.Items{
		{Description: "test product", Price: 100.00},
	}, nil)

	products, err := store.FetchProducts(ctx, int64(1))

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "test product", products[0].Description)
	assert.Equal(t, float32(100), products[0].Price)

	store.AssertExpectations(t)
}

func TestSaveOrderPoints(t *testing.T) {
	store := &storage.MockStorage{}
	ctx := context.Background()

	store.On("SaveOrderPoints", ctx, int64(1), []float32{10}).Return(nil)

	err := store.SaveOrderPoints(ctx, int64(1), []float32{10})

	assert.NoError(t, err)

	store.AssertExpectations(t)
}

func TestUpdateOrderStatus(t *testing.T) {
	store := &storage.MockStorage{}
	ctx := context.Background()

	store.On("UpdateOrderStatus", ctx, int64(1), helpers.StatusProcessed).Return(nil)

	err := store.UpdateOrderStatus(ctx, int64(1), helpers.StatusProcessed)

	assert.NoError(t, err)

	store.AssertExpectations(t)
}
