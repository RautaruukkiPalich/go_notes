

run:
	go mod tidy 
	go run ./cmd/app/main.go

makemigrations:
	migrate create -ext sql -dir migrations $(name)

migratetables:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/go_notes?sslmode=disable" $(mode)