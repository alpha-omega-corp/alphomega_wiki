package controllers

import (
	"net/http"
	"wiki/app/handlers"
	"wiki/app/middleware"
	"wiki/app/models"
	"wiki/utils"

	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

type FavoriteController struct {
	handler *handlers.FavoriteHandler
}

func NewFavoriteController(db *bun.DB) *FavoriteController {
	return &FavoriteController{handler: handlers.NewFavoriteHandler(db)}
}

// favoriteBody is the JSON payload for POST /favorite.
type favoriteBody struct {
	PageID int64 `json:"page_id"`
}

// List returns the calling user's favorites. The user is taken from the
// verified access token (set on the context by the auth middleware), never from
// the request body.
func (c FavoriteController) List(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*[]models.Favorite, error) {
		sub := middleware.UserSubFromContext(req.Context())
		return c.handler.ListByUser(sub)
	})
}

// Add stars the page in the body for the calling user.
func (c FavoriteController) Add(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*models.Favorite, error) {
		sub := middleware.UserSubFromContext(req.Context())
		body := utils.GetBody[favoriteBody](w, req)
		return c.handler.Add(sub, body.PageID)
	})
}

// Remove unstars the page named by the :id path param for the calling user.
func (c FavoriteController) Remove(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*models.Favorite, error) {
		sub := middleware.UserSubFromContext(req.Context())
		params := utils.GetParams[GetOneById](w, req)
		return nil, c.handler.Remove(sub, params.Id)
	})
}