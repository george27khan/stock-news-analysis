package usecase

import "context"

type NewsSendRepository interface {
	GetArticleToSend(ctx context.Context) ([]string, error)
}

type newsSender struct {
	Repository NewsSendRepository
}

func NewNewsSender() *newsSender {
	return &newsSender{}
}

func (s *newsSender) SendArticles(ctx context.Context) error {
	articles, err := s.Repository.GetArticleToSend(ctx)
	if err != nil {
		return err
	}

}
