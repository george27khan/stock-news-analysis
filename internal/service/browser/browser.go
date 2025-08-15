package browser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"stock-news-analysis/internal/domain"
	"time"
)

type BrowserService struct {
	Browser *domain.Browser
}

func NewBrowserChrome() *BrowserService {
	port := 9222
	return &BrowserService{&domain.Browser{"127.0.0.1",
		port,
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		[]string{fmt.Sprintf("--remote-debugging-port=%d", port), `--user-data-dir=C:\temp\chrome-debug`}}}
}

func (b *BrowserService) Run() error {
	cmd := exec.Command(b.Browser.ExePath, b.Browser.RunArgs...)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Ошибка запуска Chrome: %v", err)
		return fmt.Errorf("Ошибка запуска Chrome: %v", err)
	}
	time.Sleep(3 * time.Second) //ждем 3 секунды для открытия браузера
	log.Printf("Chrome запущен с PID %d", cmd.Process.Pid)
	return nil
}

func (b *BrowserService) GetWebSocketDebuggerURL() (string, error) {
	var info map[string]interface{}
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/json/version", b.Browser.Host, b.Browser.Port))
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
