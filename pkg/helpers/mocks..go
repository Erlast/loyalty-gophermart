package helpers

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
)

type MockFormContext struct {
	mock.Mock
}

func NewMockFormContext() *MockFormContext {
	return &MockFormContext{}
}

func (m *MockFormContext) GetUserID(r *http.Request, logger *zap.SugaredLogger) (int64, error) {
	args := m.Called(r, logger)
	userID, ok := args.Get(0).(int64)
	if !ok {
		return 0, fmt.Errorf("expected int64, got %T", args.Get(0))
	}
	err := args.Error(1)
	return userID, err
}
