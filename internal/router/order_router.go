package router

import (
	"net/http"

	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"

	"github.com/jmoiron/sqlx"
)

func orderRouter(db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()
	var dalOrdInter dal.OrderDalInter = dal.ReturnDulOrderDB(db)
	var serOrderInter service.OrdServiceInter = service.ReturnOrdSerStruct(dalOrdInter)
	handOrd := handler.ReturnOrdHaldStruct(serOrderInter)

	mux.HandleFunc("GET /", handOrd.GetOrders)
	mux.HandleFunc("GET /{id}", handOrd.GetOrderByID)
	mux.HandleFunc("DELETE /{id}", handOrd.DelOrderByID)
	mux.HandleFunc("POST /", handOrd.PostOrder)
	mux.HandleFunc("PUT /{id}", handOrd.PutOrderByID)
	mux.HandleFunc("POST /{id}/close", handOrd.PostOrdCloseById)

	return mux
}
