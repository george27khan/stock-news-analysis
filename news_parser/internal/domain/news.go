package domain

import "time"

type Article struct {
	ArticleType         string
	Body                string
	CategoryIds         []int
	CountryIds          []int
	Id                  int
	Link                string
	NewsType            string
	NewsType1           string
	PublishedAt         time.Time
	SourceId            int
	SourceName          string
	SourceImage         string
	Title               string
	UpdatedAt           time.Time
	Provider            string
	Date                string
	WriterId            int
	FrontWriterName     string
	FrontMemberName     string
	WriterName          string
	WriterLink          string
	WriterImage         string
	UserFirstName       string
	UserLastName        string
	CompanyName         string
	CompanyNamePrepared string
}

type Pairs struct {
	PairChangeNumeric   string
	PairId              int
	PairLink            string
	PairType            string
	StockSymbol         string
	PairNameExport      string
	PairNameExportTrans string
	PairChangePercent   string
}
