
DB_URL=postgres://postgres:1111@localhost:5432/product_service?sslmode=disable
migrate create -ext sql -dir db/migrations -seq create_users_table
# Применить миграции
migrate-up:
	migrate -database "$env:POSTGRESQL_URL" -source "file://db/migrations" up


# Откатить миграции
migrate-down:
	migrate -database "$(DB_URL)" -source "file://db/migrations" down 1

gen:
	protoc -I proto proto/product.proto --go_out=./gen --go_opt=paths=source_relative --go-grpc_out=./gen --go-grpc_opt=paths=source_relative