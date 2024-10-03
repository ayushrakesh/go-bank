postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

migrateup:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose up 

migratedown:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen  -package mockdb -destination db/mock/store.go github.com/ayushrakesh/go-bank/db/sqlc Store

.PHONY: postgres migrateup migratedown sqlc test server