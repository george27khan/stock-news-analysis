package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Chrome struct {
	host    string
	port    int
	exePath string
	runArgs []string
}

func newChrome() *Chrome {
	port := 9222
	return &Chrome{"127.0.0.1",
		port,
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		[]string{fmt.Sprintf("--remote-debugging-port=%d", port), `--user-data-dir=C:\temp\chrome-debug`}}
}

func (c *Chrome) runChrome() {
	cmd := exec.Command(c.exePath, c.runArgs...)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Ошибка запуска Chrome: %v", err)
	}
	log.Printf("Chrome запущен с PID %s", cmd.String())
	log.Printf("Chrome запущен с PID %d", cmd.Process.Pid)

}

func (c *Chrome) GetWebSocketDebuggerURL() string {
	var info map[string]interface{}
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/json/version", c.host, c.port))
	if err != nil {
		log.Fatalf("Ошибка при запросе: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("Ошибка декодирования: %v", err)
	}
	if url, ok := info["webSocketDebuggerUrl"]; ok {
		return strings.TrimPrefix(url.(string), fmt.Sprintf("ws://%s:%d/devtools/browser/", c.host, c.port)) //ws://127.0.0.1:9222/devtools/browser/0449ffed-6677-4c9d-8a58-9ae34236ebbd
	} else {
		return ""
	}
}

func main() {
	browser := newChrome()
	browser.runChrome()
	fmt.Println(browser.GetWebSocketDebuggerURL())
}

func main1() {
	browser := rod.New().
		ControlURL("ws://127.0.0.1:9222/devtools/browser/ce94cb7d-5d69-4693-9343-ad6794c86647"). //http://127.0.0.1:9222/json/version

		Timeout(60 * time.Second).
		MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://www.investing.com/news")

	// Правильная установка User-Agent
	err := proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}.Call(page)

	if err != nil {
		log.Fatalf("Failed to set user agent: %v", err)
	}

	page.MustWaitLoad()
	fmt.Println(page.MustHTML())
	articles := page.MustElements("a[href*='/news/']")

	for i, a := range articles {
		fmt.Println(i)
		title := a.MustText()
		href, _ := a.Attribute("href")
		if href != nil {
			fmt.Printf("📰 %s\n🔗 https://www.investing.com%s\n\n", title, *href)
		}
	}

	fmt.Println("Парсинг завершён")
}
