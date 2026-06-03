package handlers

import (
	"context"
	"fmt"
	"wiki/app/models"
	"wiki/utils"

	"github.com/uptrace/bun"
)

type PageHandler struct {
	repository utils.Repository[models.Page]
}

func NewPageHandler(db *bun.DB) *PageHandler {
	return &PageHandler{repository: utils.NewRepository[models.Page](db)}
}

func (h *PageHandler) GetOne(id int64) (*models.Page, error) {
	return h.repository.GetOne(context.Background(), id)
}

func (h *PageHandler) GetOneBySlug(slug string) (*models.Page, error) {
	return h.repository.GetOneBy(context.Background(), "slug", slug)
}

func (h *PageHandler) GetAll() (*[]models.Page, error) {
	// Bun relations only nest to a fixed depth, so instead we fetch every page
	// in a single query and assemble the parent/child tree in memory. This
	// nests pages as deep as the data goes.
	all, err := h.repository.GetAll(context.Background(), func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.ExcludeColumn("content", "excerpt").OrderExpr("id ASC")
	})
	if err != nil {
		return nil, err
	}

	// Index every page by id and clear any relation slice so we control nesting.
	byID := make(map[int64]*models.Page, len(*all))
	for i := range *all {
		p := &(*all)[i]
		p.Pages = nil
		byID[p.Id] = p
	}

	// Attach each page to its parent; pages without a (known) parent are roots.
	roots := make([]*models.Page, 0)
	for i := range *all {
		p := &(*all)[i]
		if p.ParentID != nil {
			if parent, ok := byID[*p.ParentID]; ok {
				parent.Pages = append(parent.Pages, p)
				continue
			}
		}
		roots = append(roots, p)
	}

	result := make([]models.Page, len(roots))
	for i, r := range roots {
		result[i] = *r
	}

	return &result, nil
}

func (h *PageHandler) Create(page *models.Page) (*models.Page, error) {
	fmt.Println(page)
	return h.repository.Create(context.Background(), page)
}

func (h *PageHandler) Update(page *models.Page) (*models.Page, error) {
	return h.repository.Update(context.Background(), page)
}

func (h *PageHandler) Delete(id int64) error {
	return h.repository.Delete(context.Background(), id)
}
