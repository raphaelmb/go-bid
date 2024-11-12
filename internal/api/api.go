package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/raphaelmb/go-bid/internal/services"
)

type Api struct {
	Router      *chi.Mux
	UserService services.UserService
}
