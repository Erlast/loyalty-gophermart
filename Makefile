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
CONTAINER_NAME_GF=gophermart_test
IMAGE_GF=postgres:16.3
POSTGRES_USER_GF=user
POSTGRES_PASSWORD_GF=password
POSTGRES_DB_GF=gophermart

# Команда для запуска контейнера PostgreSQL
db_gophermart:
	docker run -d \
        --name $(CONTAINER_NAME_GF) \
        -e POSTGRES_USER=$(POSTGRES_USER_GF) \
        -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD_GF) \
        -e POSTGRES_DB=$(POSTGRES_DB_GF) \
        -p 5432:5432 \
        $(IMAGE_GF)

# Параметры контейнера и образа
CONTAINER_NAME_AC=accrual_test
IMAGE_AC=postgres:16.3
POSTGRES_USER_AC=user
POSTGRES_PASSWORD_AC=password
POSTGRES_DB_AC=accrual

# Команда для запуска контейнера PostgreSQL
db_accrual:
	docker run -d \
        --name $(CONTAINER_NAME_AC) \
        -e POSTGRES_USER=$(POSTGRES_USER_AC) \
        -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD_AC) \
        -e POSTGRES_DB=$(POSTGRES_DB_AC) \
        -p 54321:5432 \
        $(IMAGE_AC)

accrual_test_db_rm:
	docker rm -f accrual_test
gophermart_test_db_rm:
	docker rm -f gophermart_test

accrual_compile:
	go build -o cmd/accrual/accrual.exe ./cmd/accrual

gophermart_compile:
	go build -o cmd/gophermart/gophermart.exe ./cmd/gophermart

# удалить тестовые базы данных
rm_db_tests: accrual_test_db_rm gophermart_test_db_rm

# скомпилировать бинарники
compile_bins: accrual_compile gophermart_compile


#	gophermarttest "-test.v" "-test.run=^TestGophermart$" "-gophermart-binary-path=cmd/gophermart/gophermart.exe" "-gophermart-host=localhost" "-gophermart-port=8080"  "-accrual-binary-path=cmd\accrual\accrual.exe" "-gophermart-database-uri=postgres://user:password@localhost:5432/gophermart?sslmode=disable" "-accrual-database-uri=postgres://user:password@localhost:54321/accrual?sslmode=disable" "-accrual-host=localhost" "-accrual-port=8081"

#	gophermarttest "-test.v" "-test.run=^TestGophermart$" "-gophermart-binary-path=cmd/gophermart/gophermart.exe" "-gophermart-host=localhost" "-gophermart-port=8080"  "-accrual-binary-path=cmd\accrual\accrual.exe" "-gophermart-database-uri=postgres://user:password@localhost:5432/gophermart?sslmode=disable" "-accrual-database-uri=postgres://user:password@localhost:54321/accrual?sslmode=disable" "-accrual-host=localhost" "-accrual-port=8081"

run_gophermart:
	go run .\cmd\gophermart\main.go -a=localhost:8008 -d=postgres://user:password@localhost:5432/gophermart?sslmode=disable -r=localhost:8081

run_accrual:
	go run .\cmd\accrual\main.go -a=localhost:8081 -d=postgres://user:password@localhost:5432/accrual?sslmode=disable