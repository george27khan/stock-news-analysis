package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"net/url"
	"os"
	dm "stock-news-analysis/news_parser/internal/domain"
	"stock-news-analysis/news_parser/internal/service/browser"
	"strconv"
	"sync"
	"time"
)

type Path struct {
	Title string
	Href  string
}

type NewsParser struct {
	Browser  *browser.BrowserService
	NewsPath []Path
	URL      *url.URL
}

func NewNewsParser() *NewsParser {
	b := browser.NewBrowserChrome()
	b.Run()
	newsPath := []Path{
		//Path{"Форекс", "/news/forex-news"},
		//Path{"Сырьевые товары", "/news/commodities-news"},
		Path{"Фондовый рынок", "/news/stock-market-news"},
		//Path{"Отчеты о доходах", "/news/earnings"},
		//Path{"Рейтинги аналитиков", "/news/analyst-ratings"},
		//Path{"Расшифровки", "/news/transcripts"},
		//Path{"Экономпоказатели", "/news/economic-indicators"},
		//Path{"Экономика", "/news/economy"},
		//Path{"Крипто", "/news/cryptocurrency-news"},
		//Path{"Срочно", "/news/headlines"},
		//Path{"Pro News", "/news/pro"},
	}
	url, _ := url.Parse("https://ru.investing.com")
	return &NewsParser{b, newsPath, url}
}

func (p *NewsParser) Parse() error {
	articleChan := make(chan *dm.Article)
	start := time.Now() // запоминаем время старта

	url, err := p.Browser.GetWebSocketDebuggerURL() // запуск браузера
	if err != nil {
		return err
	}
	//запуск браузера
	browser := rod.New().
		ControlURL(url). //http://127.0.0.1:9222/json/version
		Timeout(60 * time.Second).
		MustConnect()
	defer browser.MustClose()
	wg := &sync.WaitGroup{}
	// проходим всевкладки новостей
	for _, path := range p.NewsPath {
		pageURL, err := p.URL.Parse(path.Href)
		if err != nil {
			return err
		}
		wg.Add(1)
		go p.ParseArticleInfo(browser, pageURL.String(), wg, articleChan)
	}
	go func() {
		wg.Wait()
		close(articleChan)
	}()

	wgArticle := &sync.WaitGroup{}
	for article := range articleChan {
		articleURL, err := p.URL.Parse(article.Link)
		if err != nil {
			return err
		}
		wgArticle.Add(1)
		p.ParseArticle(browser, articleURL.String(), wgArticle)
		return nil
	}
	wg.Wait()

	elapsed := time.Since(start) // считаем разницу
	fmt.Printf("Время выполнения: %s\n", elapsed)
	return nil

}

func (p *NewsParser) ParseArticleInfo(browser *rod.Browser, url string, wg *sync.WaitGroup, articleChan chan *dm.Article) error {
	defer wg.Done()
	page := browser.MustPage(url)
	defer page.MustClose()
	err := proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}.Call(page)
	if err != nil {
		log.Fatalf("Failed to set user agent: %v", err)
	}
	page.MustWaitLoad() //ждёт полной загрузки страницы, блокирующий вызов, который завершится, когда страница полностью загрузится
	//time.Sleep(3 * time.Second)
	f, _ := os.Create("t.txt")
	f.WriteString(page.MustHTML())
	f.Close()

	// Находим JSON описывающий ленту статей
	scriptText := page.MustElement(`script#__NEXT_DATA__`).MustText()
	// Преобразуем в map[string]interface{}
	var data map[string]interface{}
	// разбираем в json
	err = json.Unmarshal([]byte(scriptText), &data)
	if err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
	}
	// доходим до нужного атрибута где лежит массив ленты статей
	data = data["props"].(map[string]interface{})
	data = data["pageProps"].(map[string]interface{})
	data = data["state"].(map[string]interface{})
	data = data["newsStore"].(map[string]interface{})

	// преобразуем массив интерфейсов в json и обратно в структуру массива ShortArticle
	news := make([]dm.Article, 0)
	//news := &ShortArticle{}
	if newsJson, err := json.Marshal(data["_news"]); err != nil {
		return err
	} else {
		fmt.Println(string(newsJson))
		return nil
		if err = json.Unmarshal(newsJson, &news); err != nil {
			return err
		}
	}
	for _, shortArticle := range news {
		articleChan <- &shortArticle
	}
	return nil
}

func (p *NewsParser) ParseArticle(browser *rod.Browser, url string, wg *sync.WaitGroup) error {
	defer wg.Done()
	page := browser.MustPage(url)
	defer page.MustClose()
	err := proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}.Call(page)
	if err != nil {
		log.Fatalf("Failed to set user agent: %v", err)
	}
	page.MustWaitLoad() //ждёт полной загрузки страницы, блокирующий вызов, который завершится, когда страница полностью загрузится
	fmt.Println(url)
	f, _ := os.Create("t.txt")
	f.WriteString(page.MustHTML())
	f.Close()

	// Находим JSON описывающий статью
	scriptText := page.MustElement(`script#__NEXT_DATA__`).MustText()
	// Преобразуем в map[string]interface{}
	var data map[string]interface{}
	// разбираем в json
	err = json.Unmarshal([]byte(scriptText), &data)
	if err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
	}
	// доходим до нужного атрибута где лежит массив ленты статей
	data = data["props"].(map[string]interface{})
	data = data["pageProps"].(map[string]interface{})
	data = data["state"].(map[string]interface{})
	data = data["newsStore"].(map[string]interface{})

	// преобразуем массив интерфейсов в json и обратно в структуру массива ShortArticle
	article := &dm.Article{}
	if articleJson, err := json.Marshal(data["_article"]); err != nil {
		return err
	} else {
		fmt.Println(string(articleJson))
		if err = json.Unmarshal(articleJson, article); err != nil {
			return err
		}
	}
	// вытаскиваем информации по парам для статьи
	pairs := make([]dm.Pairs, 0)
	data = data["_relatedPairs"].(map[string]interface{})
	if pairsJson, err := json.Marshal(data[strconv.Itoa(article.Id)]); err != nil {
		return err
	} else {
		fmt.Println(string(pairsJson))
		if err = json.Unmarshal(pairsJson, &pairs); err != nil {
			return err
		}
	}
	fmt.Println(pairs)
	return nil
}
