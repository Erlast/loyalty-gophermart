package components

import (
	"context"
	"fmt"

	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
)

var timeSleep = 1 * time.Minute
var percentFull int64 = 100

func OrderProcessing(ctx context.Context, db *pgxpool.Pool, log *zap.SugaredLogger) {
	for {
		query := `Select id FROM orders WHERE status=$1`
		rows, err := db.Query(ctx, query, helpers.StatusRegistered)
		if err != nil {
			log.Errorf("ошибка при попытке выбрать новые заказы: %v", err)
		}

		rules, err := FetchRewardRules(ctx, db)
		if err != nil {
			log.Error("не могу выбрать правила начислений")
		}

		for rows.Next() {
			var orderID int64
			err = rows.Scan(&orderID)
			if err != nil {
				log.Errorf("ошибка при получение заказа %v", err)
			}

			err = UpdateOrderStatus(ctx, db, orderID, helpers.StatusProcessing)
			if err != nil {
				log.Error("невозможно обновоить статус заказа")
			}

			products, err := FetchProducts(ctx, db, orderID)

			if err != nil {
				log.Error("не могу получить товары из заказа", err)
				err = UpdateOrderStatus(ctx, db, orderID, helpers.StatusInvalid)
				if err != nil {
					log.Error("невозможно обновоить статус заказа", err)
				}
			}

			points := make([]int64, len(products))

			for i, product := range products {
				for _, rule := range rules {
					if strings.Contains(product.Description, rule.Match) {
						switch rule.RewardType {
						case "%":
							points[i] += (product.Price * rule.Reward) / percentFull
						case "pt":
							points[i] += rule.Reward
						default:
							points[i] += 0
						}
					}
				}
			}

			err = SaveOrderPoints(ctx, db, orderID, points)
			if err != nil {
				log.Error("не могу сохранить информацию о заказе. ", err)
				err = UpdateOrderStatus(ctx, db, orderID, helpers.StatusInvalid)
				if err != nil {
					log.Error("невозможно обновоить статус заказа", err)
				}
			}
		}

		time.Sleep(timeSleep)
	}
}

func FetchRewardRules(ctx context.Context, db *pgxpool.Pool) ([]models.Goods, error) {
	var rules []models.Goods
	rows, err := db.Query(ctx, "SELECT match,reward, reward_type FROM accrual_rules")

	if err != nil {
		return nil, fmt.Errorf("can't get rules. %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Goods
		if err := rows.Scan(&r.Match, &r.Reward, &r.RewardType); err != nil {
			return nil, fmt.Errorf("can't parse rule. %w", err)
		}
		rules = append(rules, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't parse error row. %w", err)
	}
	return rules, nil
}

func FetchProducts(ctx context.Context, db *pgxpool.Pool, orderID int64) ([]models.Items, error) {
	rows, err := db.Query(ctx, "SELECT description, price FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, fmt.Errorf("can't get products. %w", err)
	}
	defer rows.Close()

	var products []models.Items
	for rows.Next() {
		var p models.Items
		if err := rows.Scan(&p.Description, &p.Price); err != nil {
			return nil, fmt.Errorf("can't parse product. %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't parse error row. %w", err)
	}
	return products, nil
}

func SaveOrderPoints(ctx context.Context, db *pgxpool.Pool, orderID int64, points []int64) error {
	var totalPoints int64
	for _, p := range points {
		totalPoints += p
	}

	_, err := db.Exec(
		ctx,
		"UPDATE orders SET status=$1,accrual=$2 where id=$3",
		helpers.StatusProcessed,
		totalPoints,
		orderID,
	)
	if err != nil {
		return fmt.Errorf("can't update orders. %w", err)
	}

	return nil
}

func UpdateOrderStatus(ctx context.Context, db *pgxpool.Pool, orderID int64, status string) error {
	_, err := db.Exec(ctx, "Update orders set status=$1 where id=$2", status, orderID)

	if err != nil {
		return fmt.Errorf("can't update order. %w", err)
	}
	return nil
}
