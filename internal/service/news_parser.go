package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"net/url"
	"os"
	_ "stock-news-analysis/internal/domain"
	"stock-news-analysis/internal/service/browser"
	"time"
)

type ShortArticle struct {
	ArticleType         string    `json:"article_type"`
	Body                string    `json:"body"`
	CreatedBy           int       `json:"created_by"`
	DomainId            int       `json:"domain_id"`
	Id                  int       `json:"id"`
	Link                string    `json:"link"`
	PublishedAt         time.Time `json:"published_at"`
	SourceExternalLink  string    `json:"source_external_link"`
	SourceId            int       `json:"source_id"`
	SourceImage         string    `json:"source_image"`
	Title               string    `json:"title"`
	UpdatedAt           time.Time `json:"updated_at"`
	CommentsCount       int       `json:"commentsCount"`
	Provider            string    `json:"provider"`
	Date                string    `json:"date"`
	WriterId            int       `json:"writerId"`
	FrontWriterName     string    `json:"front_writer_name"`
	FrontMemberName     string    `json:"front_member_name"`
	FrontWriterLink     string    `json:"front_writer_link"`
	WriterName          string    `json:"writerName"`
	WriterLink          string    `json:"writerLink"`
	WriterImage         string    `json:"writerImage"`
	UserFirstName       string    `json:"userFirstName"`
	UserLastName        string    `json:"userLastName"`
	WriterArticleCount  int       `json:"writerArticleCount"`
	CompanyName         string    `json:"company_name"`
	CompanyNamePrepared string    `json:"company_name_prepared"`
	CompanyHref         string    `json:"company_href"`
	CompanyUrl          string    `json:"company_url"`
	UserType            string    `json:"user_type"`
	AuthorLink          string    `json:"authorLink"`
	EditorId            int       `json:"editorId"`
	EditorFirstName     string    `json:"editorFirstName"`
	EditorLastName      string    `json:"editorLastName"`
	EditorLink          string    `json:"editorLink"`
	ImageCopyright      string    `json:"image_copyright"`
	RelatedImageBig     string    `json:"related_image_big"`
	ImagesSerialized    string    `json:"images_serialized"`
	ImageHref           string    `json:"imageHref"`
}
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
		Path{"Форекс", "/news/forex-news"},
		Path{"Сырьевые товары", "/news/commodities-news"},
		Path{"Фондовый рынок", "/news/stock-market-news"},
		Path{"Отчеты о доходах", "/news/earnings"},
		Path{"Рейтинги аналитиков", "/news/analyst-ratings"},
		Path{"Расшифровки", "/news/transcripts"},
		Path{"Экономпоказатели", "/news/economic-indicators"},
		Path{"Экономика", "/news/economy"},
		Path{"Крипто", "/news/cryptocurrency-news"},
		//Path{"Срочно", "/news/headlines"},
		//Path{"Pro News", "/news/pro"},
	}
	url, _ := url.Parse("https://ru.investing.com")
	return &NewsParser{b, newsPath, url}
}

func (p *NewsParser) Parse() error {
	url, err := p.Browser.GetWebSocketDebuggerURL() // запуск браузера
	if err != nil {
		return err
	}
	parser := rod.New().
		ControlURL(url). //http://127.0.0.1:9222/json/version
		Timeout(60 * time.Second).
		MustConnect()
	defer parser.MustClose()

	// проходим всевкладки новостей
	for _, path := range p.NewsPath {
		pageURL, err := p.URL.Parse(path.Href)
		if err != nil {
			return err
		}
		fmt.Println(pageURL.String())
		page := parser.MustPage(pageURL.String())
		err = proto.NetworkSetUserAgentOverride{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		}.Call(page)
		if err != nil {
			log.Fatalf("Failed to set user agent: %v", err)
		}

		//page.MustWaitLoad() //ждёт полной загрузки страницы, блокирующий вызов, который завершится, когда страница полностью загрузится
		time.Sleep(3 * time.Second)
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
		news := make([]ShortArticle, 0)
		//news := &ShortArticle{}
		if newsJson, err := json.Marshal(data["_news"]); err != nil {
			return err
		} else {
			fmt.Println(json.Unmarshal(newsJson, &news))
		}
		for _, shortArticle := range news {
			shortArticle.Link
		}
		// Выводим, например, часть данных
		fmt.Println("news:", news)
		return nil
	}
	return nil
}

func (p *NewsParser) Parse1() error {
	url, err := p.Browser.GetWebSocketDebuggerURL() // запуск браузера
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

	page.MustWaitLoad() //ждёт полной загрузки страницы, блокирующий вызов, который завершится, когда страница полностью загрузится

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
