package browser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type Browser struct {
	host    string
	port    int
	exePath string
	runArgs []string
}

func NewBrowserChrome() *Browser {
	port := 9222
	return &Browser{"127.0.0.1",
		port,
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		[]string{fmt.Sprintf("--remote-debugging-port=%d", port), `--user-data-dir=C:\temp\chrome-debug`}}
}

func (c *Browser) Run() error {
	cmd := exec.Command(c.exePath, c.runArgs...)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Ошибка запуска Chrome: %v", err)
		return fmt.Errorf("Ошибка запуска Chrome: %v", err)
	}
	time.Sleep(3 * time.Second) //ждем 3 секунды для открытия браузера
	log.Printf("Chrome запущен с PID %s", cmd.String())
	log.Printf("Chrome запущен с PID %d", cmd.Process.Pid)
	return nil

}

func (c *Browser) GetWebSocketDebuggerURL() (string, error) {
	var info map[string]interface{}
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/json/version", c.host, c.port))
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
