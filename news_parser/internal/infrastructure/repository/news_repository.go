package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stock-news-analysis/news_parser/internal/domain"
	"stock-news-analysis/news_parser/internal/dto"
	"stock-news-analysis/news_parser/internal/usecase"
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

func (r *NewsRepository) Insert(ctx context.Context, article domain.Article) error {
	query := "INSERT INTO news(source, category, url, data_json) VALUES (@source, @category, @url, @data_json)"
	articleDTO, err := r.articleDTO(article)
	if err != nil {
		return err
	}
	args := pgx.NamedArgs{
		"source":    articleDTO.Source,
		"category":  articleDTO.Category,
		"url":       articleDTO.URL,
		"data_json": articleDTO.Data_json,
	}
	if _, err := r.pool.Exec(ctx, query, args); err != nil {
		return err
	}
	return nil
}

func (r *NewsRepository) articleDTO(article domain.Article) (articleDTO *dto.ArticleDTO, err error) {
	data_json, err := json.Marshal(article)
	if err != nil {
		return nil, err
	}
	articleDTO.Source = article.Provider
	articleDTO.Category = article.ArticleType
	articleDTO.URL = article.Link
	articleDTO.Data_json = string(data_json)
	return articleDTO, nil
}
