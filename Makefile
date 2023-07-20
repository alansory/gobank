createdb:
	docker exec -it pgsql-dev createdb --username=postgres --owner=postgres gobank
dropdb:
	docker exec -it pgsql-dev dropdb --username=postgres gobank
migrate_create:
	migrate create -ext sql -dir database/migration --seq $(file)
migrateup:
	migrate -database "postgres://postgres:secret@localhost:5432/gobank?sslmode=disable" -path database/migration up
migratedown:
	migrate -database "postgres://postgres:secret@localhost:5432/gobank?sslmode=disable" -path database/migration down
migrateup_last:
	migrate -database "postgres://postgres:secret@localhost:5432/gobank?sslmode=disable" -path database/migration up 1
migratedown_last:
	migrate -database "postgres://postgres:secret@localhost:5432/gobank?sslmode=disable" -path database/migration down 1
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go

.PHONY: createdb dropdb migrate_create migrateup migratedown migrateup_last migratedown_last sqlc test server