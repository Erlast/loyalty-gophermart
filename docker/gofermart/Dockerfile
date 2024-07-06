FROM golang:1.22 as builder

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../.. .

RUN go build -o /gofermart ./cmd/gofermart/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /gofermart /gofermart

EXPOSE 8080

ENTRYPOINT ["/gofermart"]
