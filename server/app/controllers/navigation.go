package controllers

import (
	"dwiki/app/handlers"
	"dwiki/utils"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

type NavigationController struct {
	handler *handlers.NavigationHandler
}

func NewNavigationController(db *bun.DB) *NavigationController {
	return &NavigationController{handler: handlers.NewNavigationHandler(db)}
}

func (c NavigationController) GetAll(w http.ResponseWriter, req bunrouter.Request) error {
	return utils.Response(w, func() (*[]handlers.NavigationPage, error) {
		return c.handler.GetAll()
	})
}
