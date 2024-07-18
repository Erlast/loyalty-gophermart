GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.57.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:

lint:
	# | jq > ./golangci-lint/report.json
	golangci-lint run --fix -c .golangci.yml > golangci-lint/report-unformatted.json

lint-clean:
	sudo rm -rf ./golangci-lint

test:
	go test ./...

# Параметры контейнера и образа
CONTAINER_NAME=gophermart_test
IMAGE=postgres:16.3
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=gophermart
VOLUME_NAME=gophermart_data

# Команда для запуска контейнера PostgreSQL
db_g:
	docker run -d \
        --name $(CONTAINER_NAME) \
        -e POSTGRES_USER=$(POSTGRES_USER) \
        -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
        -e POSTGRES_DB=$(POSTGRES_DB) \
        -p 5432:5432 \
        -v $(VOLUME_NAME):/var/lib/postgresql/data \
        $(IMAGE)

# Параметры контейнера и образа
CONTAINER_NAME=accrual_test
IMAGE=postgres:16.3
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=accrual
VOLUME_NAME=accrual_data

# Команда для запуска контейнера PostgreSQL
db_a:
	docker run -d \
        --name $(CONTAINER_NAME) \
        -e POSTGRES_USER=$(POSTGRES_USER) \
        -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
        -e POSTGRES_DB=$(POSTGRES_DB) \
        -p 54321:5432 \
        -v $(VOLUME_NAME):/var/lib/postgresql/data \
        $(IMAGE)


gophermarttest "-test.v" "-test.run=^TestGophermart$" "-gophermart-binary-path=cmd/gophermart/gophermart.exe" "-gophermart-host=localhost" "-gophermart-port=8080"  "-accrual-binary-path=cmd\accrual\accrual.exe" "-gophermart-database-uri=postgres://user:password@localhost:5432/gophermart?sslmode=disable" "-accrual-database-uri=postgres://user:password@localhost:54321/accrual?sslmode=disable" "-accrual-host=localhost" "-accrual-port=8081"