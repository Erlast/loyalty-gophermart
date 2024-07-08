package components

import (
	"context"
	"github.com/Erlast/loyalty-gophermart.git/internal/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
	"time"
)

var timeSleep = 1 * time.Minute

func OrderProcessing(ctx context.Context, db *pgxpool.Pool, log *zap.SugaredLogger) {
	for {
		query := `Select id FROM orders WHERE status=$1`
		rows, err := db.Query(ctx, query, storage.StatusRegistered)
		if err != nil {
			log.Errorf("ошибка при удалении мягко удалённых записей: %v", err)
		}

		rules, err := FetchRewardRules(ctx, db)
		if err != nil {
			log.Error("не могу выбрать правила начислений")
		}

		for rows.Next() {
			var id int64
			err = rows.Scan(&id)
			if err != nil {
				log.Errorf("ошибка при получение заказа %v", err)
			}

			_, err := db.Exec(ctx, "Update orders set status=$1 where id=$2", storage.StatusProcessing, id)

			if err != nil {
				log.Errorf("ошибка при обработке заказа %v", err)
			}

			products, err := FetchProducts(ctx, db, id)

			if err != nil {
				log.Error("не могу получить товары из заказа")
			}
			points := make([]int64, len(products))

			for i, product := range products {
				for _, rule := range rules {
					if strings.Contains(product.Description, rule.Match) {
						switch rule.RewardType {
						case "%":
							points[i] += (product.Price * rule.Reward) / 100
						case "pt":
							points[i] += rule.Reward
						default:
							points[i] += 0
						}
					}
				}
			}
			err = SaveOrderPoints(ctx, db, id, points)
			if err != nil {
				log.Error("не могу сохранить информацию о заказе")
			}

		}

		time.Sleep(timeSleep)
	}
}

func FetchRewardRules(ctx context.Context, db *pgxpool.Pool) ([]models.Goods, error) {
	rows, err := db.Query(ctx, "SELECT match,reward, reward_type FROM accrual_rules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.Goods
	for rows.Next() {
		var r models.Goods
		if err := rows.Scan(&r.Match, &r.Reward, &r.RewardType); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}

func FetchProducts(ctx context.Context, db *pgxpool.Pool, orderID int64) ([]models.Items, error) {
	rows, err := db.Query(ctx, "SELECT desciption, price FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Items
	for rows.Next() {
		var p models.Items
		if err := rows.Scan(&p.Description, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func SaveOrderPoints(ctx context.Context, db *pgxpool.Pool, orderID int64, points []int64) error {
	var totalPoints int64
	for _, p := range points {
		totalPoints += p
	}

	_, err := db.Exec(ctx, "UPDATE orders SET status=$1,accrual=$2 where order_id=$3", storage.StatusProcessed, totalPoints, orderID)
	if err != nil {
		return err
	}

	return nil
}
