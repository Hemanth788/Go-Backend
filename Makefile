postgres:
	docker run --name go-be-pg -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:18.3-alpine3.23

createdb:
	docker exec -it go-be-pg createdb --username=admin --owner=admin bank

dropdb: 
	docker exec -it go-be-pg dropdb --username=admin --owner=admin bank

migrateup:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

startPg:
	docker container start 4a5b248bd83f

stopPg:
	docker container stop 4a5b248bd83f

test: 
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen --build_flags=--mod=mod -package mockdb -destination db/mock/store.go go.com/go-backend/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc stopPg startPg server mock