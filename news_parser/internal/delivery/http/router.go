package http

import (
	"stock-news-analysis/news_parser/internal/service"
)

func Run() {
	parser := service.NewNewsParser()
	parser.Parse()

}
