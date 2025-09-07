package main

import (
	"log"
	"os"
	"runtime/pprof"
	"stock-news-analysis/news_parser/internal/delivery/http"
)

// go tool pprof -http=:8083  cpu.out
// go tool pprof -http=:8083  mem.out
func main() {

	// Профилирование CPU
	fcpu, err := os.Create("cpu.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(fcpu)
	defer pprof.StopCPUProfile()

	// --- твой код ---
	http.Run()
	// ----------------

	// Профилирование памяти
	fmem, err := os.Create("mem.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(fmem)
	fmem.Close()

}
