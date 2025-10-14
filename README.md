# stock-news-analysis

Миграция

Сначала ставим утилиту:

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


Путь: $GOPATH/bin/migrate (не забудь добавить его в $PATH).

migrate create -ext sql -dir news_parser/internal/infrastructure/db/migrations -seq create_table

migrate -path news_parser/internal/infrastructure/db/migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up 1
