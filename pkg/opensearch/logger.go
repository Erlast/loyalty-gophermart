package opensearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
	"time"
)

type LogMessage struct {
	Message  string `json:"message"`
	Level    string `json:"level"`
	DateTime string `json:"dateTime"`
}

type Logger struct {
	client *opensearch.Client
	Logger *zap.Logger
	index  string
}

func NewOpenSearchLogger() (*Logger, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "admin",
		Password: "yourStrongPassword123!",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	logger, _ := zap.NewProduction()
	return &Logger{
		client: client,
		Logger: logger,
		index:  "logs",
	}, nil
}

func (l *Logger) SendLog(level string, message string) {
	logMessage := LogMessage{
		Message:  message,
		Level:    level,
		DateTime: time.Now().Format(time.RFC3339),
	}

	jsonMessage, err := json.Marshal(logMessage)
	if err != nil {
		l.Logger.Fatal("Failed to serialize log message", zap.Error(err))
	}

	req := opensearchapi.IndexRequest{
		Index:      l.index,
		DocumentID: "",
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
