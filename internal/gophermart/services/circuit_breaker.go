package services

import (
	"time"

	"github.com/sony/gobreaker"
)

func NewCircuitBreaker(name string, timeout time.Duration) *gobreaker.CircuitBreaker {
	// Настройка параметров Circuit Breaker
	var st gobreaker.Settings
	st.Name = "HTTP API"
	st.Timeout = 10 * time.Second // время охлаждения
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests > 3 && failureRatio >= 0.6
	}

	return gobreaker.NewCircuitBreaker(st)
}
