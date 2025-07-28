package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (api *API) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		RespondWithErr(w, "invalid userID", http.StatusBadRequest)
		return
	}
	userData, err := api.userService.GetUserByID(r.Context(), userID)
	RespondWithJson(w, userData, 200)
}

func getUserID(r *http.Request) (uuid.UUID, error) {
	userIdStr := r.Context().Value("userID").(string)
	userID, err := uuid.Parse(userIdStr)
	return userID, err
}
