package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stock-news-analysis/news_parser/internal/domain"
	"stock-news-analysis/news_parser/internal/dto"
	"stock-news-analysis/news_parser/internal/usecase"
	"time"
)

var _ usecase.NewsRepository = (*NewsRepository)(nil)

// отдельный конструктор для пула, т.к. он должен быть общим для всех репозиториев
func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("DB ping error", err)
	}
	return pool, nil
}

type NewsRepository struct {
	pool *pgxpool.Pool
}

func NewNewsRepository(pool *pgxpool.Pool) *NewsRepository {

	return &NewsRepository{pool: pool}
}

func (r *NewsRepository) Close() {
	r.pool.Close()
}

func (r *NewsRepository) articleDTO(article domain.Article) (articleDTO dto.ArticleDTO, err error) {
	data_json, err := json.Marshal(article)
	if err != nil {
		return dto.ArticleDTO{}, fmt.Errorf("NewsRepository.articleDTO error: %w", err)
	}
	articleDTO.Source = article.Provider
	articleDTO.Category = article.ArticleType
	articleDTO.URL = article.Link
	articleDTO.Data_json = string(data_json)
	articleDTO.Published_at = article.PublishedAt
	return articleDTO, nil
}

// Insert вставка статьи в news
func (r *NewsRepository) Insert(ctx context.Context, article domain.Article) error {
	query := "INSERT INTO news(source, category, url, data_json, published_at) VALUES (@source, @category, @url, @data_json, @published_at)"
	articleDTO, err := r.articleDTO(article)
	if err != nil {
		return fmt.Errorf("NewsRepository.Insert error: %w", err)
	}
	args := pgx.NamedArgs{
		"source":       articleDTO.Source,
		"category":     articleDTO.Category,
		"url":          articleDTO.URL,
		"data_json":    articleDTO.Data_json,
		"published_at": articleDTO.Published_at,
	}
	if _, err := r.pool.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("NewsRepository.Insert error: %w", err)
	}
	return nil
}

// InsertBatch вставка слайса статей в news
func (r *NewsRepository) InsertBatch(ctx context.Context, articles []domain.Article) error {
	var errs []error
	query := "INSERT INTO news(source, category, url, data_json,published_at) VALUES (@source, @category, @url, @data_json, @published_at)"
	batch := &pgx.Batch{}
	for _, article := range articles {
		articleDTO, err := r.articleDTO(article)
		if err != nil {
			return fmt.Errorf("NewsRepository.InsertBatch error: %w", err)
		}
		args := pgx.NamedArgs{
			"source":       articleDTO.Source,
			"category":     articleDTO.Category,
			"url":          articleDTO.URL,
			"data_json":    articleDTO.Data_json,
			"published_at": articleDTO.Published_at,
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

// GetLastArticleDt получение последней статьи
func (r *NewsRepository) GetLastArticleDt(ctx context.Context, path string) (time.Time, error) {
	var lastDt time.Time
	query := `select published_at 
			    from news 
			   where published_at > CURRENT_DATE
			     and url like $1
			   order by published_at desc limit 1`
	row := r.pool.QueryRow(ctx, query, path+"%")
	if err := row.Scan(&lastDt); err != nil {
		return time.Time{}, fmt.Errorf("NewsRepository.GetLastArticleUrl error: %w", err)
	}
	return lastDt, nil
}
