package services

import (
	"encoding/json"
	"fmt"
	"gofermart/internal/gofermart/models"
	"net/http"
	"time"
)

type AccrualService struct {
	AccrualSystemAddress string
	Client               *http.Client
}

func NewAccrualService(address string) *AccrualService {
	return &AccrualService{
		AccrualSystemAddress: address,
		Client:               &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *AccrualService) GetAccrualInfo(orderNumber string) (*models.AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", s.AccrualSystemAddress, orderNumber)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get accrual info: status %d", resp.StatusCode)
	}

	var accrualResp models.AccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
		return nil, err
	}

	return &accrualResp, nil
}
