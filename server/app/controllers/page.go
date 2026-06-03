package controllers

import (
	"fmt"
	"net/http"
	"wiki/app/handlers"
	"wiki/app/models"
	"wiki/utils"

	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

type PageController struct {
	handler *handlers.PageHandler
}

type GetOneBySlug struct {
	Slug string `json:"slug"`
}

func NewPageHandler(db *bun.DB) *PageController {
	return &PageController{handler: handlers.NewPageHandler(db)}
}

func (c PageController) GetOne(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*models.Page, error) {
		data := utils.GetParams[GetOneBySlug](w, req)
		return c.handler.GetOneBySlug(data.Slug)
	})
}

func (c PageController) GetAll(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*[]models.Page, error) {
		return c.handler.GetAll()
	})
}

func (c PageController) Create(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*models.Page, error) {
		data := utils.GetBody[models.Page](w, req)
		fmt.Println(data)
		return c.handler.Create(data)
	})
}

func (c PageController) Update(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*models.Page, error) {
		data := utils.GetBody[models.Page](w, req)
		fmt.Println(data)
		return c.handler.Update(data)
	})
}
