package zaplog

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

// InitLogger создает и возвращает новый экземпляр zap SugaredLogger для разработки.
func InitLogger() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Error creating zap logger: %v", err)
	}
	Logger = zapLogger.Sugar()
}
