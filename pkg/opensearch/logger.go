package opensearch

import (
	"context"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
	"log"
	"strings"
)

type LogMessage struct {
	Message string `json:"message"`
	Level   string `json:"level"`
}

// Logger представляет логгер с клиентом OpenSearch
type Logger struct {
	client *opensearch.Client
	Logger *zap.Logger
	index  string
}

func NewOpenSearchLogger() (*Logger, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Создание zap логгера
	logger, _ := zap.NewProduction()
	return &Logger{
		client: client,
		Logger: logger,
		index:  "logs",
	}, nil
}

func (l *Logger) SendLog(level string, message string) {
	logMessage := LogMessage{
		Message: message,
		Level:   level,
	}

	// Сериализация сообщения в JSON
	jsonMessage, err := json.Marshal(logMessage)
	if err != nil {
		l.Logger.Fatal("Failed to serialize log message", zap.Error(err))
	}

	// Отправка сообщения в OpenSearch
	req := opensearchapi.IndexRequest{
		Index:      l.index, // Имя индекса
		DocumentID: "",      // Если пусто, то будет автоматически создано
		Body:       strings.NewReader(string(jsonMessage)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), l.client)
	if err != nil {
		l.Logger.Fatal("Failed to send log message to OpenSearch", zap.Error(err))
	}
	defer res.Body.Close()

	if res.IsError() {
		l.Logger.Fatal("Error indexing document", zap.String("status", res.Status()))
	} else {
		l.Logger.Info("Document indexed successfully", zap.String("result", res.String()))
	}
}
