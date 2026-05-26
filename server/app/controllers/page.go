package controllers

import (
	"dwiki/app/handlers"
	"dwiki/app/models"
	"dwiki/utils"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

type ArticleController struct {
	handler *handlers.ArticleHandler
}

type GetOneBySlug struct {
	Slug string `json:"slug"`
}

func NewArticleController(db *bun.DB) *ArticleController {
	return &ArticleController{handler: handlers.NewArticleHandler(db)}
}

func (c ArticleController) Create(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*models.Page, error) {
		data := utils.GetBody[models.Page](w, req)
		return c.handler.Create(data)
	})
}

func (c ArticleController) GetOne(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*models.Page, error) {
		data := utils.GetParams[GetOneBySlug](w, req)
		return c.handler.GetOne(data.Slug)
	})
}

func (c ArticleController) GetAll(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*[]models.Page, error) {
		return c.handler.GetAll()
	})
}
