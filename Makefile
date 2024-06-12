postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

migrateup:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose up 

migratedown:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres migrateup migratedown