package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"

	"github.com/sony/gobreaker"

	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"

	"net/http"
)

type AccrualService struct { //nolint:govet //4кб не стоят этого
	AccrualSystemAddress string
	logger               *zap.SugaredLogger
	Circuit              *gobreaker.CircuitBreaker
}

func NewAccrualService(
	logger *zap.SugaredLogger,
	circuit *gobreaker.CircuitBreaker,
	address string,
) *AccrualService {
	return &AccrualService{
		logger:               logger,
		AccrualSystemAddress: address,
		Circuit:              circuit,
	}
}

func (s *AccrualService) GetAccrualInfo(orderNumber string) (*models.AccrualResponse, error) {
	baseURL, err := url.Parse(s.AccrualSystemAddress)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать базовый адрес %s: %w", s.AccrualSystemAddress, err)
	}

	// Добавляем путь к базовому URL
	baseURL.Path = path.Join(baseURL.Path, "api/orders", orderNumber)

	accrualResp, err := s.Circuit.Execute(func() (interface{}, error) {
		resp, err := http.Get(baseURL.String())
		if err != nil {
			return nil, fmt.Errorf("error request to accrual service: %w", err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				s.logger.Errorf("не удалось закрыть body ответа от %s: %w", baseURL.String(), err)
			}
		}()

		if resp.StatusCode == http.StatusNoContent {
			return nil, errors.New("нет содержимого")
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("не удалось получить информацию о начислениях: статус %d", resp.StatusCode)
		}

		var accrualResp models.AccrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
			return nil, fmt.Errorf("не удалось разобрать информацию о начислениях с %s: %w", baseURL.String(), err)
		}
		return &accrualResp, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error accrual info: %w", err)
	}

	ar, ok := accrualResp.(*models.AccrualResponse)
	if !ok {
		return nil, errors.New("response is not accrual info")
	}
	return ar, nil
}
