package dto

import "time"

type ArticleDTO struct {
	ID         int64     `db:"id"`
	Source     string    `db:"source"`
	Category   string    `db:"category"`
	URL        string    `db:"url"`
	Data_json  string    `db:"data_json"`
	Created_dt time.Time `db:"created_dt"`
	Is_send    bool      `db:"is_send"`
}
