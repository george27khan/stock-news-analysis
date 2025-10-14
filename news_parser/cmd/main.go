package main

import (
	"context"
	"stock-news-analysis/news_parser/internal/delivery/local"
	"stock-news-analysis/news_parser/internal/infrastructure/browser"
	"stock-news-analysis/news_parser/internal/infrastructure/parser/rod"
	"stock-news-analysis/news_parser/internal/infrastructure/repository"
	"stock-news-analysis/news_parser/internal/usecase"
)

// go tool pprof -http=:8083  cpu.out
// go tool pprof -http=:8083  mem.out
func main() {
	ctx := context.Background()
	pool, err := repository.NewPostgresPool(ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil
	}
	rep := repository.NewNewsRepository(pool)
	b := browser.NewBrowser()
	parser := rod.NewNewsParser()
	newsUsecase := usecase.NewNewsParserUseCase(b, parser)
	newsUsecase.Parse()
	// Профилирование CPU
	//fcpu, err := os.Create("cpu.out")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(fcpu)
	//defer pprof.StopCPUProfile()

	// --- твой код ---
	local.Run()
	// ----------------

	// Профилирование памяти
	//fmem, err := os.Create("mem.out")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.WriteHeapProfile(fmem)
	//fmem.Close()

}
