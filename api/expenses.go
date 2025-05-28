package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ayikoandrew/server/types"
)

func (s *Server) uploadExpenses(w http.ResponseWriter, r *http.Request) error {
	expense := new(types.Expense)
	if err := json.NewDecoder(r.Body).Decode(expense); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, expense)
}

func (s *Server) retriveExpenses(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "OK")
}
