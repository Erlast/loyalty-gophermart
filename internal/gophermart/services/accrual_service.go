package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"

	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"

	"net/http"
	"time"
)

type AccrualService struct {
	logger               *zap.SugaredLogger
	Client               *http.Client
	AccrualSystemAddress string
}

func NewAccrualService(
	address string,
	logger *zap.SugaredLogger,
) *AccrualService {
	return &AccrualService{
		logger:               logger,
		AccrualSystemAddress: address,
		Client:               &http.Client{Timeout: 10 * time.Second}, //nolint:mnd // Timeout 10 секунд
	}
}

func (s *AccrualService) GetAccrualInfo(orderNumber string) (*models.AccrualResponse, error) {
	baseURL, err := url.Parse(s.AccrualSystemAddress)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать базовый адрес %s: %w", s.AccrualSystemAddress, err)
	}

	// Добавляем путь к базовому URL
	baseURL.Path = path.Join(baseURL.Path, "api/orders", orderNumber)

	resp, err := s.Client.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о начислениях с %s: %w", baseURL.String(), err)
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
}
