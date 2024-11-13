package api

import (
	"errors"
	"net/http"

	"github.com/raphaelmb/go-bid/internal/jsonutils"
	"github.com/raphaelmb/go-bid/internal/services"
	"github.com/raphaelmb/go-bid/internal/usecase/user"
)

const AUTHENTICATED_USER_ID = "AuthenticatedUserID"

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.CreateUserReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			jsonutils.EncodeJSON(w, r, http.StatusConflict, map[string]any{
				"error": "email or username already exists",
			})
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.LoginUserReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid email or password",
			})
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	api.Sessions.Put(r.Context(), AUTHENTICATED_USER_ID, id)

	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "logged in successfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "something went wrong",
		})
		return
	}

	api.Sessions.Remove(r.Context(), AUTHENTICATED_USER_ID)

	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "logged out successfully",
	})
}
