package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"stock-news-analysis/news_parser/internal/domain"
	"stock-news-analysis/news_parser/internal/usecase"
	"time"
)

var p_ usecase.Browser = (*Browser)(nil)

type Browser struct {
	Info *domain.Browser
	CMD  *exec.Cmd
}

func NewBrowser() *Browser {
	port := 9222
	return &Browser{
		Info: &domain.Browser{
			Host:    "127.0.0.1",
			Port:    port,
			ExePath: `C:\Program Files\Google\Chrome\Application\chrome.exe`,
			RunArgs: []string{fmt.Sprintf("--remote-debugging-port=%d", port), `--user-data-dir=C:\temp\chrome-debug`}},
		CMD: nil,
	}
}

func (b *Browser) Run(ctx context.Context) error {
	b.CMD = exec.Command(b.Info.ExePath, b.Info.RunArgs...)
	err := b.CMD.Start()
	if err != nil {
		log.Fatalf("Ошибка запуска Chrome: %v", err)
		return fmt.Errorf("Ошибка запуска Chrome: %v", err)
	}
	time.Sleep(3 * time.Second) //ждем 3 секунды для открытия браузера
	log.Printf("Chrome запущен с PID %d", b.CMD.Process.Pid)
	return nil
}

func (b *Browser) Close() error {
	if b.CMD != nil && b.CMD.Process != nil {
		if err := b.CMD.Process.Kill(); err != nil {
			return fmt.Errorf("ошибка закрытия Chrome: %w", err)
		}
		log.Printf("Chrome с PID %d завершен", b.CMD.Process.Pid)
		b.CMD = nil
	}
	return nil
}

func (b *Browser) GetURL(ctx context.Context) (string, error) {
	var info map[string]interface{}
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/json/version", b.Info.Host, b.Info.Port))
	if err != nil {
		log.Fatalf("Ошибка при запросе: %v", err)
		return "", fmt.Errorf("Ошибка при запросе: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("Ошибка декодирования: %v", err)
		return "", fmt.Errorf("Ошибка декодирования: %v", err)
	}
	if url, ok := info["webSocketDebuggerUrl"]; ok {
		return url.(string), nil
	} else {
		return "", fmt.Errorf("")
	}
}
