package rod

import (
	"encoding/json"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"net/url"
	"os"
	dm "stock-news-analysis/news_parser/internal/domain"
	"stock-news-analysis/news_parser/internal/usecase"
	"strconv"
	"time"
)

var _ usecase.NewsParser = (*NewsParser)(nil)

type NewsParser struct {
	Browser     *rod.Browser
	articleChan chan *dm.Article
}

func NewNewsParser() *NewsParser {
	return &NewsParser{rod.New(), make(chan *dm.Article, 10)}
}

func (p *NewsParser) Run(url string, timeout time.Duration) {
	p.Browser = p.Browser.ControlURL(url).Timeout(timeout).MustConnect()
}

func (p *NewsParser) Close() {
	p.Browser.MustClose()
}

//func (p *NewsParser) Parse() error {
//	articleChan := make(chan *dm.Article)
//	//start := time.Now() // запоминаем время старта
//	//
//	//url, err := p.Browser.GetWebSocketDebuggerURL() // запуск браузера
//	//if err != nil {
//	//	return err
//	//}
//	////запуск браузера
//	//browser := rod.New().
//	//	ControlURL(url). //http://127.0.0.1:9222/json/version
//	//	Timeout(60 * time.Second).
//	//	MustConnect()
//	//defer browser.MustClose()
//	wg := &sync.WaitGroup{}
//	// проходим всевкладки новостей
//	for _, path := range p.NewsPath {
//		pageURL, err := p.URL.Parse(path.Href)
//		if err != nil {
//			return err
//		}
//		wg.Add(1)
//		go p.ParseArticleInfo(browser, pageURL.String(), wg, articleChan)
//	}
//	go func() {
//		wg.Wait()
//		close(articleChan)
//	}()
//
//	wgArticle := &sync.WaitGroup{}
//	for article := range articleChan {
//		articleURL, err := p.URL.Parse(article.Link)
//		if err != nil {
//			return err
//		}
//		wgArticle.Add(1)
//		p.ParseArticle(browser, articleURL.String(), wgArticle)
//		return nil
//	}
//	wg.Wait()
//
//	elapsed := time.Since(start) // считаем разницу
//	fmt.Printf("Время выполнения: %s\n", elapsed)
//	return nil
//
//}

func (p *NewsParser) Parse(url *url.URL) ([]dm.Article, error) {
	resArticle := make([]dm.Article, 0)
	start := time.Now() // запоминаем время старта
	shortArticleInfo, err := p.parseNewsInfo(url)
	if err != nil {
		return nil, err
	}
	for _, shortArticle := range shortArticleInfo {
		articleURL, err := url.Parse(shortArticle.Link)
		if err != nil {
			return nil, err
		}
		article, err := p.parseArticleInfo(articleURL)
		if err != nil {
			return nil, err
		}
		resArticle = append(resArticle, *article)
		fmt.Println(resArticle)
		//return nil, nil
	}

	elapsed := time.Since(start) // считаем разницу
	fmt.Printf("Время выполнения: %s\n", elapsed)
	return resArticle, nil

}

func (p *NewsParser) parseNewsInfo(url *url.URL) ([]*dm.Article, error) {
	data, err := p.getPageNewsInfo(url) //получаем информацию о постах
	if err != nil {
		return nil, err
	}
	// преобразуем массив интерфейсов в json и обратно в структуру массива ShortArticle
	shortArticleInfo := make([]*dm.Article, 0)
	if newsJson, err := json.Marshal(data["_news"]); err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(newsJson, &shortArticleInfo); err != nil {
			return nil, err
		}
	}
	return shortArticleInfo, nil
}

func (p *NewsParser) parseArticleInfo(url *url.URL) (*dm.Article, error) {
	data, err := p.getPageNewsInfo(url) //получаем информацию о постах
	if err != nil {
		return nil, err
	}
	fmt.Println(url.String())
	// преобразуем массив интерфейсов в json и обратно в структуру массива Article
	article := &dm.Article{}
	if articleJson, err := json.Marshal(data["_article"]); err != nil {
		return nil, err
	} else {
		//fmt.Println(string(articleJson))
		if err = json.Unmarshal(articleJson, article); err != nil {
			return nil, err
		}
	}
	//fmt.Println(article)
	// вытаскиваем информации по парам для статьи
	pairs := make([]dm.Pairs, 0)
	data = data["_relatedPairs"].(map[string]interface{})
	if pairsJson, err := json.Marshal(data[strconv.Itoa(article.Id)]); err != nil {
		return nil, err
	} else {
		//fmt.Println(string(pairsJson))
		if err = json.Unmarshal(pairsJson, &pairs); err != nil {
			return nil, err
		}
	}
	article.Pairs = pairs
	return article, nil
}

func (p *NewsParser) getPageNewsInfo(url *url.URL) (map[string]interface{}, error) {
	//fmt.Println(url.String())
	page := p.Browser.MustPage(url.String())
	defer page.MustClose()
	err := proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}.Call(page)
	if err != nil {
		log.Fatalf("Failed to set user agent: %v", err)
		return nil, err
	}
	page.MustWaitLoad() //ждёт полной загрузки страницы, блокирующий вызов, который завершится, когда страница полностью загрузится
	//fmt.Println(url)
	f, _ := os.Create("t.txt")
	f.WriteString(page.MustHTML())
	f.Close()
	// Находим JSON описывающий данные
	dataJSON := page.MustElement(`script#__NEXT_DATA__`).MustText()
	var data map[string]interface{}
	// разбираем в структуру
	err = json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
		return nil, err
	}
	// доходим до нужного атрибута где лежит массив ленты статей
	data = data["props"].(map[string]interface{})
	data = data["pageProps"].(map[string]interface{})
	data = data["state"].(map[string]interface{})
	data = data["newsStore"].(map[string]interface{})

	return data, nil
}
