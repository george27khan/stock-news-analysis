package dto

import (
	"time"
)

type ArticleDTO struct {
	ID          int       `db:"id"`
	Source      string    `db:"source"`
	Category    string    `db:"category"`
	ArticleId   int       `db:"article_id"`
	Title       string    `db:"title"`
	URL         string    `db:"url"`
	DataJson    string    `db:"data_json"`
	CreatedDt   time.Time `db:"created_dt"`
	PublishedAt time.Time `db:"published_at"`
	IsSend      bool      `db:"is_send"`
}
