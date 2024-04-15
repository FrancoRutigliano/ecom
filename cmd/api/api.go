package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/FrancoRutigliano/ecom/service/cart"
	"github.com/FrancoRutigliano/ecom/service/order"
	"github.com/FrancoRutigliano/ecom/service/product"
	"github.com/FrancoRutigliano/ecom/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	//product
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	//order
	orderStore := order.NewStore(s.db)

	//cart
	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
