package user

import (
	"fmt"
	"net/http"

	"github.com/FrancoRutigliano/ecom/config"
	"github.com/FrancoRutigliano/ecom/service/auth"
	typePayload "github.com/FrancoRutigliano/ecom/types"
	"github.com/FrancoRutigliano/ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store typePayload.UserStore
}

func NewHandler(store typePayload.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")

	//admin route
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload typePayload.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	// need the secret coming from the env var's
	secret := []byte(config.Envs.JWTSecret)
	// Generate the JWT
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload typePayload.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, err)
	}

	// validate the Payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//checkear si existe el usuario
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteErr(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exist", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, err)
		return
	}

	// si no existe lo creamos
	err = h.store.CreateUser(typePayload.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteErr(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)

}
