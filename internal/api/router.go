package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
func NewRouter(a *AccountHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/health", HandleHealth)
	router.Route("/accounts" , func (r chi.Router)  {
		r.Post("/", a.CreateAccount)
	})
	return router
}