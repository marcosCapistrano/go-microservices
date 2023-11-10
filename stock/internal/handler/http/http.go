package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"stocksapp/stock/internal/controller/stock"
)

// Handler defines a movie handler.
type Handler struct {
	ctrl *stock.Controller
}

// New creates a new movie HTTP handler.
func New(ctrl *stock.Controller) *Handler {
	return &Handler{ctrl}
}

// GetMovieDetails handles GET /movie requests.
func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, stock.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
