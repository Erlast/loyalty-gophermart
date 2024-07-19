package zaplog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger создает и возвращает новый экземпляр zap SugaredLogger для логирования в файл.
func InitLogger() *zap.SugaredLogger {
	// Указываем путь к файлу лога
	logFile := "logfile.log"

	// Открываем файл для логирования
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zap.L().Fatal("Failed to open log file", zap.Error(err))
	}

	// Настраиваем Encoder для логирования
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Форматирование времени
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Создаем Core для логирования в файл с уровнем логирования DebugLevel
	core := zapcore.NewCore(encoder, zapcore.AddSync(file), zapcore.DebugLevel)

	// Создаем Logger
	zapLogger := zap.New(core)
	SugaredLogger := zapLogger.Sugar()

	return SugaredLogger
}
