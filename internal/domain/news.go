package domain

import "time"

type Article struct {
	ArticleType         string    `json:"article_type"`
	Body                string    `json:"body"`
	CategoryIds         []int     `json:"category_ids"`
	CountryIds          []int     `json:"country_ids"`
	Id                  int       `json:"id"`
	Link                string    `json:"link"`
	NewsType            string    `json:"news_type"`
	NewsType1           string    `json:"newsType"`
	PublishedAt         time.Time `json:"published_at"`
	SourceId            int       `json:"source_id"`
	SourceName          string    `json:"source_name"`
	SourceImage         string    `json:"source_image"`
	Title               string    `json:"title"`
	UpdatedAt           time.Time `json:"updated_at"`
	Provider            string    `json:"provider"`
	Date                string    `json:"date"`
	WriterId            int       `json:"writerId"`
	FrontWriterName     string    `json:"front_writer_name"`
	FrontMemberName     string    `json:"front_member_name"`
	WriterName          string    `json:"writerName"`
	WriterLink          string    `json:"writerLink"`
	WriterImage         string    `json:"writerImage"`
	UserFirstName       string    `json:"userFirstName"`
	UserLastName        string    `json:"userLastName"`
	CompanyName         string    `json:"company_name"`
	CompanyNamePrepared string    `json:"company_name_prepared"`
}

type Pairs struct {
	PairChangeNumeric   string `json:"pair_change_numeric"`
	PairId              int    `json:"pair_id"`
	PairLink            string `json:"pair_link"`
	PairType            string `json:"pair_type"`
	StockSymbol         string `json:"stock_symbol"`
	PairNameExport      string `json:"pair_name_export"`
	PairNameExportTrans string `json:"pair_name_export_trans"`
	PairChangePercent   string `json:"pair_change_percent"`
}
