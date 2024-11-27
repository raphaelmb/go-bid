package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/raphaelmb/go-bid/internal/jsonutils"
	"github.com/raphaelmb/go-bid/internal/services"
)

func (api *Api) HandleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductId := chi.URLParam(r, "product_id")

	productId, err := uuid.Parse(rawProductId)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{
			"message": "invalid product id - must be valid uuid",
		})
		return
	}

	_, err = api.ProductService.GetProductById(r.Context(), productId)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJSON(w, r, http.StatusNotFound, map[string]any{
				"message": "no product with given id",
			})
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), AUTHENTICATED_USER_ID).(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"message": "could not upgrade connection to websocket protocol",
		})
		return
	}
}
