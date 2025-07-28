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
	api.router.Route("/api", func(r1 chi.Router) {
		r1.Use(middleware.Logger)
		r1.Route("/auth", func(r2 chi.Router) {
			r2.Post("/register", api.registerUser)
			r2.Post("/login", api.loginUser)
		})
		r1.Route("/users", func(r2 chi.Router) {
			r2.Use(api.accessTokenMiddleware(api.userService.JwtSecret))
			r2.Get("/me", api.getUser)
		})

	})

}

func (api *API) Start() error {
	return http.ListenAndServe(":9999", api.router)
}
