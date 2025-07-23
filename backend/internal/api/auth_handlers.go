package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
	"github.com/kar1mov-u/LeetClone/internal/models"
	"github.com/kar1mov-u/LeetClone/internal/services"
)

func (api *API) registerUser(w http.ResponseWriter, r *http.Request) {
	var userData models.UserRegister
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		RespondWithErr(w, "Invalid data", http.StatusBadRequest)
		return
	}

	if len(userData.Username) == 0 || len(userData.Email) == 0 || len(userData.Password) == 0 {
		RespondWithErr(w, "Fields cannot be empty", http.StatusBadRequest)
		return
	}

	if _, err = mail.ParseAddress(userData.Email); err != nil {
		RespondWithErr(w, "Invalid Email", http.StatusBadRequest)
		return
	}

	userID, err := api.userService.RegisterUser(r.Context(), userData)
	if err != nil {
		if errors.Is(err, services.EmailTakenErr) {
			RespondWithErr(w, "Email is already in use", http.StatusConflict)
			return
		} else if errors.Is(err, services.UsernameTakenErr) {
			RespondWithErr(w, "UserName is already in use", http.StatusConflict)
			return
		} else {
			RespondWithErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	RespondWithJson(w, map[string]uuid.UUID{"userID": userID}, 200)
	//do some shit

}

func RespondWithJson(w http.ResponseWriter, data any, code int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error on responding JSON: %v", err)
		http.Error(w, "Error on sedning marshalled JSON", http.StatusInternalServerError)
		return
	}
}

func RespondWithErr(w http.ResponseWriter, errorMessage string, code int) {
	RespondWithJson(w, map[string]string{"error": errorMessage}, code)
}
