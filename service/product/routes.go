package product

import (
	"net/http"

	typePayload "github.com/FrancoRutigliano/ecom/types"
	"github.com/FrancoRutigliano/ecom/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store typePayload.ProductStore
}

func NewHandler(store typePayload.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleCreateProduct).Methods("GET")
	router.HandleFunc("/products", h.handleCreateProduct).Methods("POST")
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, ps)
}
