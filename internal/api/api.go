package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/raphaelmb/go-bid/internal/services"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserService
	ProductService services.ProductService
	BidService     services.BidService
	Sessions       *scs.SessionManager
	WsUpgrader     websocket.Upgrader
}
