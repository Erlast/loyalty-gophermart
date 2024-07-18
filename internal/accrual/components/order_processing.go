package components

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
)

var timeSleep = 5 * time.Second
var percentFull int64 = 100

func OrderProcessing(ctx context.Context, store storage.Storage, logger *zap.SugaredLogger) {
	for {
		orders, err := store.GetRegisteredOrders(ctx)
		if err != nil {
			logger.Errorf("ошибка при попытке выбрать новые заказы: %v", err)
		}

		rules, err := store.FetchRewardRules(ctx)
		if err != nil {
			logger.Error("не могу выбрать правила начислений")
		}

		for _, orderID := range orders {
			err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusProcessing)
			if err != nil {
				logger.Error("невозможно обновоить статус заказа")
			}

			products, err := store.FetchProducts(ctx, orderID)

			if err != nil {
				logger.Error("не могу получить товары из заказа", err)
				err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusInvalid)
				if err != nil {
					logger.Error("невозможно обновоить статус заказа", err)
				}
			}

			points := make([]float64, len(products))

			for i, product := range products {
				for _, rule := range rules {
					if strings.Contains(product.Description, rule.Match) {
						switch rule.RewardType {
						case "%":
							points[i] += float64((product.Price * float32(rule.Reward)) / float32(percentFull))
						case "pt":
							points[i] += float64(rule.Reward)
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
			}
		}

		time.Sleep(timeSleep)
	}
}
