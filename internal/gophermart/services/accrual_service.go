package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"

	"net/http"
	"time"
)

type AccrualService struct {
	Client               *http.Client
	AccrualSystemAddress string
}

func NewAccrualService(address string) *AccrualService {
	return &AccrualService{
		AccrualSystemAddress: address,
		Client:               &http.Client{Timeout: 10 * time.Second}, //nolint:mnd // Timeout 10 секунд
	}
}

func (s *AccrualService) GetAccrualInfo(orderNumber string) (*models.AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", s.AccrualSystemAddress, orderNumber)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get accrual info from %s: %w", url, err)
	}
	defer resp.Body.Close() //nolint:errcheck // later change

	if resp.StatusCode == http.StatusNoContent {
		return nil, errors.New("no content")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get accrual info: status %d", resp.StatusCode)
	}

	var accrualResp models.AccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
		return nil, fmt.Errorf("could not parse accrual info from %s: %w", url, err)
	}

	return &accrualResp, nil
}
