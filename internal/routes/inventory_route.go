package router

import (
	"net/http"

	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"

	"github.com/jmoiron/sqlx"
)

func inventoryRouter(db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()

	var dalInventInter dal.InventoryDataAccess = dal.ReturnDalInvCore(db)
	serviceInventInter := service.ReturnInventorySerInt(dalInventInter)
	handInvInt := handler.NewInventoryHandler(serviceInventInter)
	
	mux.HandleFunc("POST /", handInvInt.PostInventory)
	mux.HandleFunc("GET /", handInvInt.GetInventories)
	mux.HandleFunc("GET /{id}", handInvInt.GetInventoryByID)
	mux.HandleFunc("PUT /{id}", handInvInt.PutInventory)
	mux.HandleFunc("DELETE /{id}", handInvInt.DeleteInventory)
	mux.HandleFunc("GET /history", handInvInt.GetInventoryHistory)
	mux.HandleFunc("GET /reorder", handInvInt.GetReorderInventories)
	return mux
}
