postgres:
	docker run --name test-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=rahasia -d postgres:12-alpine

postgres-stop:
	docker stop test-postgres

postgres-start:
	docker start test-postgres

createdb:
	docker exec -it test-postgres createdb --username=root --owner=root simple_bank 

dropdb:
	docker exec -it test-postgres dropdb simple_bank

sqlc:
	sqlc generate

migrateup:
	migrate -path db/migration -database "postgresql://root:rahasia@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:rahasia@localhost:5432/simple_bank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc