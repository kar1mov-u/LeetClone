package api

import (
	"net/http"

	"github.com/kar1mov-u/LeetClone/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	userService *services.UserService
	router      *chi.Mux
}

func New(userService *services.UserService) *API {
	router := chi.NewRouter()
	api := API{router: router, userService: userService}
	api.setRoutes()
	return &api
}

func (api *API) setRoutes() {
	//set routes of the api
	api.router.Use(middleware.Logger)

	api.router.Post("/api/user/register", api.registerUser)
}

func (api *API) Start() error {
	return http.ListenAndServe(":9999", api.router)
}
