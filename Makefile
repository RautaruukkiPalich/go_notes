

run:
	go mod tidy 
	go run ./cmd/app/main.go

lint:
	golangci-lint run ./...

swagger:
	swag init -g cmd/app/main.go
	swag fmt

makemigrations:
	migrate create -ext sql -dir migrations $(name)

migratetables:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/go_notes?sslmode=disable" $(mode)