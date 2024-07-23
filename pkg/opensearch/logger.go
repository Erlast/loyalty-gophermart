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

func NewOpenSearchLogger() {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Адрес вашего OpenSearch кластера
		},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Создание zap логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Пример логгирования сообщения
	logMessage := map[string]string{
		"message": "Hello, OpenSearch!",
		"level":   "info",
	}

	// Сериализация сообщения в JSON
	jsonMessage, err := json.Marshal(logMessage)
	if err != nil {
		logger.Fatal("Failed to serialize log message", zap.Error(err))
	}

	// Отправка сообщения в OpenSearch
	req := opensearchapi.IndexRequest{
		Index:      "logs", // Имя индекса
		DocumentID: "",     // Если пусто, то будет автоматически создано
		Body:       strings.NewReader(string(jsonMessage)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		logger.Fatal("Failed to send log message to OpenSearch", zap.Error(err))
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Fatal("Error indexing document", zap.String("status", res.Status()))
	} else {
		logger.Info("Document indexed successfully", zap.String("result", res.String()))
	}
}
