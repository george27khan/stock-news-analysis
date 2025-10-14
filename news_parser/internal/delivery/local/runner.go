package local

import (
	"stock-news-analysis/news_parser/internal/infrastructure/browser"
	"stock-news-analysis/news_parser/internal/infrastructure/parser/rod"
	"stock-news-analysis/news_parser/internal/usecase"
)

func Run() {
	b := browser.NewBrowser()
	p := rod.NewNewsParser()
	parser := usecase.NewNewsParserUseCase(b, p)
	parser.Parse()

}
