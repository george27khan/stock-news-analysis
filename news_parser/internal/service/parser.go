package service

type NewsParser interface {
	Parse() error
}

type NewsSender interface {
	Send() error
}

type ArticeService struct {
	parser *NewsParser
	sender *NewsSender
}

func NewArticeService(parser *NewsParser, sender *NewsSender) *ArticeService {
	return &ArticeService{parser, sender}
}
