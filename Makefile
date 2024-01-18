MIGRATIONS_DIR := /Users/ali/GoProject/testProject/migrations

vet:
	go vet ./testProject/...

build:
	go build -o ./bin/testProject ./cmd/main.go

clean:
	rm -rf ./bin

test:
	go test ./service/person_service_test.go

run: 
	docker-compose up

stop:
	docker-compose down


migrate-up:
	migrate -database "postgres://postgres:21509@localhost:5432/test?sslmode=disable" -path $(MIGRATIONS_DIR) up

mmigrate-down:
	migrate -database "postgres://postgres:052005@localhost:5436/test?sslmode=disable" -path $(MIGRATIONS_DIR) down
	
