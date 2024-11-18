package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphaelmb/go-bid/internal/jsonutils"
	"github.com/raphaelmb/go-bid/internal/usecase/product"
)

func (api *Api) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[product.CreateProductReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), AUTHENTICATED_USER_ID).(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	id, err := api.ProductService.CreateProduct(
		r.Context(),
		userID,
		data.ProductName,
		data.Description,
		data.Baseprice,
		data.AuctionEnd,
	)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"message": "product created successfully",
		"id":      id,
	})
}
