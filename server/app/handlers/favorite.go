package handlers

import (
	"context"
	"wiki/app/models"
	"wiki/utils"

	"github.com/uptrace/bun"
)

type FavoriteHandler struct {
	repository utils.Repository[models.Favorite]
	db         *bun.DB
}

func NewFavoriteHandler(db *bun.DB) *FavoriteHandler {
	return &FavoriteHandler{repository: utils.NewRepository[models.Favorite](db), db: db}
}

// ListByUser returns the favorites owned by a user, newest first.
func (h *FavoriteHandler) ListByUser(userSub string) (*[]models.Favorite, error) {
	return h.repository.GetAll(context.Background(), func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("user_sub = ?", userSub).OrderExpr("created_at DESC")
	})
}

// Add stars a page for the user. It is idempotent: starring an already-favorited
// page is a no-op (relies on the (user_sub, page_id) unique constraint).
func (h *FavoriteHandler) Add(userSub string, pageID int64) (*models.Favorite, error) {
	fav := &models.Favorite{UserSub: userSub, PageID: pageID}
	_, err := h.db.NewInsert().
		Model(fav).
		On("CONFLICT (user_sub, page_id) DO NOTHING").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return fav, nil
}

// Remove unstars a page for the user.
func (h *FavoriteHandler) Remove(userSub string, pageID int64) error {
	_, err := h.db.NewDelete().
		Model((*models.Favorite)(nil)).
		Where("user_sub = ? AND page_id = ?", userSub, pageID).
		Exec(context.Background())
	return err
}