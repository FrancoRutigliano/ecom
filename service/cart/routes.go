package cart

import (
	"fmt"
	"net/http"

	"github.com/FrancoRutigliano/ecom/service/auth"
	typePayload "github.com/FrancoRutigliano/ecom/types"
	"github.com/FrancoRutigliano/ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store      typePayload.ProductStore
	orderstore typePayload.OrderStore
	userStore  typePayload.UserStore
}

func NewHandler(store typePayload.ProductStore,
	orderStore typePayload.OrderStore,
	userStore typePayload.UserStore) *Handler {
	return &Handler{store: store,
		orderstore: orderStore,
		userStore:  userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods("POST")
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIdFromContext(r.Context())

	var cart typePayload.CartCheckoutPayload
	// parse
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, err)
		return
	}

	//validate the payload
	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// obtaining the productsId
	productIDs, err := getCartItemsIDs(cart.Item)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, err)
		return
	}

	//get products
	products, err := h.store.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(products, cart.Item, userId)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"total_price": totalPrice,
		"order":       orderID,
	})

}
