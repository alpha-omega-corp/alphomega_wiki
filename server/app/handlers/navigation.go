package handlers

import (
	"context"
	"dwiki/app/models"

	"github.com/uptrace/bun"
)

type NavigationPage struct {
	Id        int64             `json:"id"`
	Slug      string            `json:"slug"`
	Title     string            `json:"title"`
	Tags      []string          `json:"tags"`
	Published bool              `json:"published"`
	Pages     []*NavigationPage `json:"pages,omitempty"`
}

type NavigationHandler struct {
	db *bun.DB
}

func NewNavigationHandler(db *bun.DB) *NavigationHandler {
	return &NavigationHandler{db: db}
}

func (h *NavigationHandler) GetAll() (*[]NavigationPage, error) {
	ctx := context.Background()

	var pages []models.Page
	if err := h.db.NewSelect().
		Model(&pages).
		Where("p.parent_id IS NULL").
		Relation("Pages", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("content", "excerpt")
		}).
		ExcludeColumn("content", "excerpt").
		Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]NavigationPage, len(pages))
	for i, p := range pages {
		result[i] = toNavPage(p)
	}
	return &result, nil
}

func toNavPage(p models.Page) NavigationPage {
	np := NavigationPage{
		Id:        p.Id,
		Slug:      p.Slug,
		Title:     p.Title,
		Tags:      p.Tags,
		Published: p.Published,
	}
	for _, child := range p.Pages {
		sub := toNavPage(*child)
		np.Pages = append(np.Pages, &sub)
	}
	return np
}
