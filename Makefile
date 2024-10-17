postgres:
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

migrateup:
	migrate -path db/migrations -database "postgres://postgres:6Wi9N3ucDgJK8tXT4gdq@go-bank.c7cia4imqgo7.ap-south-1.rds.amazonaws.com:5432/bank" -verbose up 

migrateup1:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgres://postgres:6Wi9N3ucDgJK8tXT4gdq@go-bank.c7cia4imqgo7.ap-south-1.rds.amazonaws.com:5432/bank" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen  -package mockdb -destination db/mock/store.go github.com/ayushrakesh/go-bank/db/sqlc Store

.PHONY: postgres migrateup migratedown sqlc test server mock migrateup1 migratedown1