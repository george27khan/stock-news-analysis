package usecase

import (
	"context"
	"fmt"
	"net/url"
	"stock-news-analysis/news_parser/internal/domain"
	"time"
)

type path struct {
	Title string
	Href  string
}

// var _ h.ArticleUseCase = (*ArticleUseCase)(nil)
type Browser interface {
	Run() error
	GetURL() (string, error)
}

type NewsParser interface {
	Run(url string, timeout time.Duration)
	Close()
	Parse(url *url.URL) ([]domain.Article, error)
}

type NewsRepository interface {
	Insert(ctx context.Context, article domain.Article) error
	Close()
}

type NewsParserUseCase struct {
	Browser    Browser
	Parser     NewsParser
	NewsPath   []path
	URL        *url.URL
	Repository NewsRepository
}

func NewNewsParserUseCase(b Browser, p NewsParser, r NewsRepository) *NewsParserUseCase {

	newsPath := []path{
		//PathСырьевые товары", "/news/commodities-news"},
		path{"Фондовый рынок", "/news/stock-market-news"},
		//		//Path{"Отчеты о доходах", "/news/earnings"},
		//		//Path{"Рейтинги аналитиков", "/news/analyst-ratings"},
		//		//Path{"Рас{"Форекс", "/news/forex-news"},
		//Path{"шифровки", "/news/transcripts"},
		//Path{"Экономпоказатели", "/news/economic-indicators"},
		//Path{"Экономика", "/news/economy"},
		//Path{"Крипто", "/news/cryptocurrency-news"},
		//Path{"Срочно", "/news/headlines"},
		//Path{"Pro News", "/news/pro"},
	}
	url, _ := url.Parse("https://ru.investing.com")
	return &NewsParserUseCase{b, p, newsPath, url, r}
}

func (p *NewsParserUseCase) Parse() error {
	start := time.Now()                     // запоминаем время старта
	if err := p.Browser.Run(); err != nil { // запуск браузера
		return err
	}
	url, err := p.Browser.GetURL() // получаем адрес браузера
	if err != nil {
		return err
	}
	//fmt.Println("url ", url)
	p.Parser.Run(url, 60*time.Second) //запуск парсера в браузере
	defer p.Parser.Close()

	// проходим по всем новостным страницам сайта
	for _, path := range p.NewsPath {
		pageURL, err := p.URL.Parse(path.Href) //формируем URL
		if err != nil {
			return err
		}
		articles, err := p.Parser.Parse(pageURL) // запускаем парсинг, получаем массив статей
		if err != nil {
			return err
		}
		articles = articles
		//fmt.Println(articles)
	}
	// нужно реализовать отправку статей в кафку
	elapsed := time.Since(start) // считаем разницу
	fmt.Printf("Время выполнения: %s\n", elapsed)
	return nil

}

//func (p *NewsParser) ParseArticleInfo(browser *rod.Browser, url string, wg *sync.WaitGroup, articleChan chan *dm.Article) error {
//	defer wg.Done()
//	page := browser.MustPage(url)
//	defer page.MustClose()
//	err := proto.NetworkSetUserAgentOverride{
//		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
//	}.Call(page)
//	if err != nil {
//		log.Fatalf("Failed to set user agent: %v", err)
//	}
//	page.MustWaitLoad() //ждёт полной загрузки страницы, блокирующий вызов, который завершится, когда страница полностью загрузится
//	//time.Sleep(3 * time.Second)
//	f, _ := os.Create("t.txt")
//	f.WriteString(page.MustHTML())
//	f.Close()
//
//	// Находим JSON описывающий ленту статей
//	scriptText := page.MustElement(`script#__NEXT_DATA__`).MustText()
//	// Преобразуем в map[string]interface{}
//	var data map[string]interface{}
//	// разбираем в json
//	err = json.Unmarshal([]byte(scriptText), &data)
//	if err != nil {
//		log.Fatal("Ошибка парсинга JSON:", err)
//	}
//	// доходим до нужного атрибута где лежит массив ленты статей
//	data = data["props"].(map[string]interface{})
//	data = data["pageProps"].(map[string]interface{})
//	data = data["state"].(map[string]interface{})
//	data = data["newsStore"].(map[string]interface{})
//
//	// преобразуем массив интерфейсов в json и обратно в структуру массива ShortArticle
//	news := make([]dm.Article, 0)
//	//news := &ShortArticle{}
//	if newsJson, err := json.Marshal(data["_news"]); err != nil {
//		return err
//	} else {
//		fmt.Println(string(newsJson))
//		return nil
//		if err = json.Unmarshal(newsJson, &news); err != nil {
//			return err
//		}
//	}
//	for _, shortArticle := range news {
//		articleChan <- &shortArticle
//	}
//	return nil
//}
//
//func (p *NewsParser) ParseArticle(browser *rod.Browser, url string, wg *sync.WaitGroup) error {
//	defer wg.Done()
//	page := browser.MustPage(url)
//	defer page.MustClose()
//	err := proto.NetworkSetUserAgentOverride{
//		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
//	}.Call(page)
//	if err != nil {
//		log.Fatalf("Failed to set user agent: %v", err)
//	}
//	page.MustWaitLoad() //ждёт полной загрузки страницы, блокирующий вызов, который завершится, когда страница полностью загрузится
//	fmt.Println(url)
//	f, _ := os.Create("t.txt")
//	f.WriteString(page.MustHTML())
//	f.Close()
//
//	// Находим JSON описывающий статью
//	scriptText := page.MustElement(`script#__NEXT_DATA__`).MustText()
//	// Преобразуем в map[string]interface{}
//	var data map[string]interface{}
//	// разбираем в json
//	err = json.Unmarshal([]byte(scriptText), &data)
//	if err != nil {
//		log.Fatal("Ошибка парсинга JSON:", err)
//	}
//	// доходим до нужного атрибута где лежит массив ленты статей
//	data = data["props"].(map[string]interface{})
//	data = data["pageProps"].(map[string]interface{})
//	data = data["state"].(map[string]interface{})
//	data = data["newsStore"].(map[string]interface{})
//
//	// преобразуем массив интерфейсов в json и обратно в структуру массива ShortArticle
//	article := &dm.Article{}
//	if articleJson, err := json.Marshal(data["_article"]); err != nil {
//		return err
//	} else {
//		fmt.Println(string(articleJson))
//		if err = json.Unmarshal(articleJson, article); err != nil {
//			return err
//		}
//	}
//	// вытаскиваем информации по парам для статьи
//	pairs := make([]dm.Pairs, 0)
//	data = data["_relatedPairs"].(map[string]interface{})
//	if pairsJson, err := json.Marshal(data[strconv.Itoa(article.Id)]); err != nil {
//		return err
//	} else {
//		fmt.Println(string(pairsJson))
//		if err = json.Unmarshal(pairsJson, &pairs); err != nil {
//			return err
//		}
//	}
//	fmt.Println(pairs)
//	return nil
//}
