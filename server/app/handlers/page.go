package handlers

import (
	"context"
	"dwiki/app/models"
	"dwiki/utils"

	"github.com/uptrace/bun"
)

type ArticleHandler struct {
	repository utils.Repository[models.Page]
}

func NewArticleHandler(db *bun.DB) *ArticleHandler {
	return &ArticleHandler{repository: utils.NewRepository[models.Page](db)}
}

func (h *ArticleHandler) Create(page *models.Page) (*models.Page, error) {
	return h.repository.Create(context.Background(), page)
}

func (h *ArticleHandler) GetOne(slug string) (*models.Page, error) {
	return h.repository.GetOneBy(context.Background(), "slug", slug)
}

func (h *ArticleHandler) GetAll() (*[]models.Page, error) {
	return h.repository.GetAll(context.Background(), nil)
}
