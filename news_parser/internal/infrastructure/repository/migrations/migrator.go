package migrations

//go:generate migrate -path ./ -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable down -all

//go:generate migrate -path ./ -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

//migrate create -ext sql -dir news_parser/internal/infrastructure/repository/migrations -seq create_table
//migrate -path news_parser/internal/infrastructure/repository/migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down -all

//migrate -path news_parser/internal/infrastructure/repository/migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
