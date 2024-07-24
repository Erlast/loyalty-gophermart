package components

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/opensearch"
)

var timeSleep = 5 * time.Second
var percentFull float32 = 100

func OrderProcessing(ctx context.Context, store storage.Storage, logger *opensearch.Logger) {
	for {
		orders, err := store.GetRegisteredOrders(ctx)
		if err != nil {
			logger.SendLog("error", fmt.Sprintf("ошибка при попытке выбрать новые заказы: %v", err))
		}

		rules, err := store.FetchRewardRules(ctx)
		if err != nil {
			logger.SendLog("error", "не могу выбрать правила начислений")
		}

		for _, orderID := range orders {
			err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusProcessing)
			if err != nil {
				logger.SendLog("error", "невозможно обновоить статус заказа")
			}

			products, err := store.FetchProducts(ctx, orderID)

			if err != nil {
				logger.SendLog("error", fmt.Sprintf("не могу получить товары из заказа: %v", err))
				err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusInvalid)
				if err != nil {
					logger.SendLog("error", fmt.Sprintf("невозможно обновоить статус заказа: %v", err))
				}
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
				logger.SendLog("error", fmt.Sprintf("не могу сохранить информацию о заказе: %v", err))
				err = store.UpdateOrderStatus(ctx, orderID, helpers.StatusInvalid)
				if err != nil {
					logger.SendLog("error", fmt.Sprintf("невозможно обновоить статус заказа: %v", err))
				}
			}
		}

		time.Sleep(timeSleep)
	}
}
