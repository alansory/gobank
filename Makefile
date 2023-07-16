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
build:
	@go build -o bin/gobank
run: build
	@./bin/gobank
test:
	go test -v -cover -short ./...