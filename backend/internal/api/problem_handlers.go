package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kar1mov-u/LeetClone/internal/models"
)

func (api *API) CreateProblem(w http.ResponseWriter, r *http.Request) {

	//extract input payload to struct
	var problemData models.CreateProblem

	err := json.NewDecoder(r.Body).Decode(&problemData)
	if err != nil {
		RespondWithErr(w, "error on decoding json", http.StatusBadRequest)
		return
	}

	id, err := api.problemService.CreateProblem(r.Context(), problemData)

	if err != nil {
		RespondWithErr(w, fmt.Sprintf("error on creating problem:%v", err), http.StatusInternalServerError)
		return
	}
	RespondWithJson(w, map[string]int{"problem_id": id}, http.StatusOK)
}
