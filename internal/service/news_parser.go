package service

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"os"
	_ "stock-news-analysis/internal/domain"
	"stock-news-analysis/internal/service/browser"
	"time"
)

type NewsParser struct {
	Browser *browser.BrowserService
}

func NewNewsParser() *NewsParser {
	b := browser.NewBrowserChrome()
	b.Run()
	return &NewsParser{b}
}

func (p *NewsParser) Parse() error {
	url, err := p.Browser.GetWebSocketDebuggerURL()
	if err != nil {
		return err
	}
	parser := rod.New().
		ControlURL(url). //http://127.0.0.1:9222/json/version
		Timeout(60 * time.Second).
		MustConnect()
	defer parser.MustClose()
	//page := parser.MustPage("https://ru.investing.com/news/economy")
	page := parser.MustPage("https://ru.investing.com/news")
	err = proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}.Call(page)

	if err != nil {
		log.Fatalf("Failed to set user agent: %v", err)
	}

	page.MustWaitLoad()
	//fmt.Println()
	f, _ := os.Create("t.txt")
	f.WriteString(page.MustHTML())
	f.Close()
	//articles := page.MustElements("a[href*='/news/']")
	//articles := page.MustElements("div > a[href*='/news/']")
	//articles := page.MustElementsX(`//div/a[contains(@href, "/news/economy") and normalize-space(text()) != ""]`)
	//for i, a := range articles {
	//	fmt.Println(i)
	//	title := a.MustText()
	//	href, _ := a.Attribute("href")
	//	if href != nil {
	//		fmt.Printf("📰 %s\n🔗 %s\n\n", title, *href)
	//	}
	//}

	// Получаем содержимое скрипта с id="__NEXT_DATA__"
	//script := page.MustElement(`#__NEXT_DATA__`)
	//jsonText := script.MustText()
	//f, _ = os.Create("tt.txt")
	//f.WriteString(jsonText)
	//f.Close()
	//fmt.Println("Парсинг завершён")
	//// Парсим нужное поле, например: props.pageProps.article.title
	//result := gjson.Get(jsonText, "props.pageProps.state.newsStore._newsList")
	//
	//// Выводим результат
	//article := domain.Article{}
	////posts := make([]*domain.Article,0)
	//if result.Exists() {
	//	//json.Unmarshal([]byte(result.String(),&posts)
	//	//log.Println("Заголовок статьи:", result.)
	//	result.ForEach(func(_, value gjson.Result) bool {
	//		json.Unmarshal([]byte(value.String()), &article)
	//		fmt.Println(article)
	//	})
	//} else {
	//	log.Println("Поле не найдено")
	//}
	//fmt.Println(jsonText)
	return nil
}
