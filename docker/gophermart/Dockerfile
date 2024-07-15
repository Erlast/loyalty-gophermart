FROM golang:1.22 as builder

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../.. .

RUN go build -o /gophermart ./cmd/gophermart/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /gophermart /gophermart

EXPOSE 8080

ENTRYPOINT ["/gophermart"]
