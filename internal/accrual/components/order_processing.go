package components

import (
	"context"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"math"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
)

var timeSleep = 2 * time.Second
var percentFull float32 = 100

type Storage interface {
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	SaveOrderItems(ctx context.Context, items models.OrderItem) error
	SaveGoods(ctx context.Context, goods models.Goods) error
	GetRegisteredOrders(ctx context.Context) ([]int64, error)
	FetchRewardRules(ctx context.Context) ([]models.Goods, error)
	UpdateOrderStatus(ctx context.Context, orderNumber int64, status string) error
	FetchProducts(ctx context.Context, orderID int64) ([]models.Items, error)
	SaveOrderPoints(ctx context.Context, orderID int64, points []float32) error
}

func OrderProcessing(ctx context.Context, store Storage, logger *zap.SugaredLogger) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("OrderProcessing stopped")
			return
		default:
			orders, err := store.GetRegisteredOrders(ctx)
			if err != nil {
				logger.Errorf("ошибка при попытке выбрать новые заказы: %v", err)
				return
			}

			rules, err := store.FetchRewardRules(ctx)
			if err != nil {
				logger.Error("не могу выбрать правила начислений")
				return
			}

			for _, orderID := range orders {
				err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusProcessing)
				if err != nil {
					logger.Error("невозможно обновоить статус заказа")
					return
				}

				products, err := store.FetchProducts(ctx, orderID)

				if err != nil {
					logger.Error("не могу получить товары из заказа", err)
					err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusInvalid)
					if err != nil {
						logger.Error("невозможно обновоить статус заказа", err)
						return
					}
					return
				}

				points := make([]float32, len(products))

				for i, product := range products {
					for _, rule := range rules {
						if strings.Contains(product.Description, rule.Match) {
							switch rule.RewardType {
							case "%":
								points[i] += float32(math.Round(float64(product.Price*rule.Reward/percentFull*100)) / 100)
							case "pt":
								points[i] += rule.Reward
							default:
								points[i] += 0
							}
						}
					}
				}

				err = store.SaveOrderPoints(ctx, orderID, points)
				if err != nil {
					logger.Error("не могу сохранить информацию о заказе. ", err)
					err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusInvalid)
					if err != nil {
						logger.Error("невозможно обновоить статус заказа", err)
					}
					return
				}
			}

			time.Sleep(timeSleep)
		}
	}
}
