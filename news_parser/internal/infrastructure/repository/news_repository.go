package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"stock-news-analysis/news_parser/internal/domain"
	"stock-news-analysis/news_parser/internal/dto"
	"stock-news-analysis/news_parser/internal/usecase"
)

var (
	_ usecase.NewsParseRepository = (*newsRepository)(nil)
	_ usecase.NewsSendRepository  = (*newsRepository)(nil)
)

// NewPostgresPool отдельный конструктор для пула, т.к. он должен быть общим для всех репозиториев
func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("DB ping error:  %w", err)
	}
	return pool, nil
}

type newsRepository struct {
	pool *pgxpool.Pool
}

func NewNewsRepository(pool *pgxpool.Pool) *newsRepository {
	return &newsRepository{pool: pool}
}

func (r *newsRepository) Close() {
	r.pool.Close()
}

func (r *newsRepository) articleDTO(article domain.Article) (articleDTO dto.ArticleDTO, err error) {
	dataJson, err := json.Marshal(article)
	if err != nil {
		return dto.ArticleDTO{}, fmt.Errorf("NewsRepository.articleDTO error: %w", err)
	}
	articleDTO.Source = article.Provider
	articleDTO.Category = article.ArticleType
	articleDTO.ArticleId = article.Id
	articleDTO.Title = article.Title
	articleDTO.URL = article.Link
	articleDTO.DataJson = string(dataJson)
	articleDTO.PublishedAt = article.PublishedAt
	return articleDTO, nil
}

// Insert вставка статьи в news
func (r *newsRepository) Insert(ctx context.Context, article domain.Article) error {
	query := `INSERT INTO news(source, category, article_id, title, url, data_json, published_at) 
			  VALUES (@source, @category, @article_id, @title, @url, @data_json, @published_at)`
	articleDTO, err := r.articleDTO(article)
	if err != nil {
		return fmt.Errorf("NewsRepository.Insert error: %w", err)
	}
	args := pgx.NamedArgs{
		"source":       articleDTO.Source,
		"category":     articleDTO.Category,
		"article_id":   articleDTO.ArticleId,
		"title":        articleDTO.Title,
		"url":          articleDTO.URL,
		"data_json":    articleDTO.DataJson,
		"published_at": articleDTO.PublishedAt,
	}
	if _, err := r.pool.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("NewsRepository.Insert error: %w", err)
	}
	return nil
}

// InsertBatch вставка слайса статей в news
func (r *newsRepository) InsertBatch(ctx context.Context, articles []domain.Article) error {
	var errs []error
	query := `INSERT INTO news(source, category, article_id, title, url, data_json, published_at) 
			  VALUES (@source, @category, @article_id, @title, @url, @data_json, @published_at)`
	batch := &pgx.Batch{}
	for _, article := range articles {
		articleDTO, err := r.articleDTO(article)
		if err != nil {
			return fmt.Errorf("NewsRepository.InsertBatch error: %w", err)
		}
		args := pgx.NamedArgs{
			"source":       articleDTO.Source,
			"category":     articleDTO.Category,
			"article_id":   articleDTO.ArticleId,
			"title":        articleDTO.Title,
			"url":          articleDTO.URL,
			"data_json":    articleDTO.DataJson,
			"published_at": articleDTO.PublishedAt,
		}
		batch.Queue(query, args)
	}

	br := r.pool.SendBatch(ctx, batch)

	defer br.Close()
	for range articles {
		if _, err := br.Exec(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("NewsRepository.InsertBatch error: %w", errors.Join(errs...))
	}
	return nil
}

// GetLastArticleId получение ID последней статьи за день
func (r *newsRepository) GetLastArticleId(ctx context.Context, path string) (int, error) {
	var articleIdLast int
	query := `select article_id
			    from news 
			   where published_at > CURRENT_DATE
			     and url like $1
			   order by published_at desc limit 1`
	row := r.pool.QueryRow(ctx, query, path+"%")
	if err := row.Scan(&articleIdLast); err != nil {
		return 0, fmt.Errorf("NewsRepository.GetLastArticleUrl error: %w", err)
	}
	slog.Debug("GetLastArticleDt", "articleIdLast", articleIdLast)
	return articleIdLast, nil
}

// GetArticleToSend получение статей на отправку за день
func (r *newsRepository) GetArticleToSend(ctx context.Context) ([]string, error) {
	query := `select data_json
			    from news 
			   where published_at > CURRENT_DATE
			     and is_send = false`
	row, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository.GetArticleToSend r.pool.Query error: %w", err)
	}
	articlesJson := make([]string, 0)
	data := ""
	for row.Next() {
		if err = row.Scan(&data); err != nil {
			return nil, fmt.Errorf("repository.GetArticleToSend row.Scan(&data) error: %w", err)
		}
		articlesJson = append(articlesJson, data)
	}
	return articlesJson, nil
}
