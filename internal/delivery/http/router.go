package http

import "stock-news-analysis/internal/service"

func Run() {
	parser := service.NewNewsParser()
	parser.Parse()

}
