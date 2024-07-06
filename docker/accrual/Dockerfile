FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /accrual ./cmd/accrual-system/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /accrual /accrual

EXPOSE 8090

ENTRYPOINT ["/accrual"]
