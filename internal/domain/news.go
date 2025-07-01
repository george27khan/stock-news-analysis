package domain

type Article struct {
	ArticleID       int    `json:"article_ID"`
	Title           string `json:"title"`
	ShortTitle      string `json:"shortTitle"`
	Href            string `json:"href"`
	ImageHref       string `json:"imageHref"`
	Provider        string `json:"provider"`
	ProviderHref    string `json:"providerHref"`
	Date            string `json:"date"`
	CommentsCounter int    `json:"commentsCounter"`
	Snippet         string `json:"snippet"`
	OpenInNewTab    bool   `json:"openInNewTab"`
	MediumImageHref string `json:"mediumImageHref"`
	SmallImageHref  string `json:"smallImageHref"`
	NewsType        string `json:"news_type"`
	Category        string `json:"category"`
	CategoryLink    string `json:"categoryLink"`
	ArticleType     string `json:"article_type"`
	Pairs           []struct {
		StockSymbol       string `json:"stock_symbol"`
		PairLink          string `json:"pair_link"`
		PairName          string `json:"pair_name"`
		PairNameSynonym   string `json:"pair_name_synonym"`
		PairChangePercent string `json:"pair_change_percent"`
	} `json:"pairs"`
}
