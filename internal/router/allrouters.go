package router

import (
	"net/http"

	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"

	"github.com/jmoiron/sqlx"
)

func Allrouter(db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()
	// setup pathfile to dulinvent and build to handfunc
	var dalInventInter dal.InventoryDataAccess = dal.ReturnDalInvCore(db)
	serviceInventInter := service.ReturnInventorySerInt(dalInventInter)
	handInvInt := handler.NewInventoryHandler(serviceInventInter)
	// mux.Handle("", mux1)
	mux.HandleFunc("POST /inventory", handInvInt.PostInventory)
	mux.HandleFunc("GET /inventory", handInvInt.GetInventories)
	mux.HandleFunc("GET /inventory/{id}", handInvInt.GetInventoryByID)
	mux.HandleFunc("PUT /inventory/{id}", handInvInt.PutInventory)
	mux.HandleFunc("DELETE /inventory/{id}", handInvInt.DeleteInventory)
	// mux.HandleFunc("PUT /inventory", handInv.PutAllIng)

	// // setup pathfile to dulmenu struct and build to handfunc
	var dalMenuInter dal.MenuDalInter = dal.ReturnDalMenuCore(db)
	menuSerInter := service.ReturnMenuSerStruct(dalMenuInter)
	handMenu := handler.ReturnMenuHaldStruct(menuSerInter)

	mux.HandleFunc("GET /menu", handMenu.GetMenus)
	mux.HandleFunc("GET /menu/{id}", handMenu.GetMenuByID)
	mux.HandleFunc("DELETE /menu/{id}", handMenu.DelMenu)
	mux.HandleFunc("POST /menu", handMenu.PostMenu)
	mux.HandleFunc("PUT /menu/{id}", handMenu.PutMenuByID)

	// // setup pathfiles to dulorder struct and build to handlfunc
	var dalOrdInter dal.OrderDalInter = dal.ReturnDulOrderCore(db)
	var serOrderInter service.OrdServiceInter = service.ReturnOrdSerStruct(dalOrdInter)
	handOrd := handler.ReturnOrdHaldStruct(serOrderInter)

	mux.HandleFunc("GET /orders", handOrd.GetOrders)
	mux.HandleFunc("GET /orders/{id}", handOrd.GetOrderByID)
	mux.HandleFunc("DELETE /orders/{id}", handOrd.DelOrderByID)
	mux.HandleFunc("POST /orders", handOrd.PostOrder)
	mux.HandleFunc("PUT /orders/{id}", handOrd.PutOrderByID)
	mux.HandleFunc("POST /orders/{id}/close", handOrd.PostOrdCloseById)
	// mux.HandleFunc("GET /reports/total-sales", ordHand.TotalSales)
	// mux.HandleFunc("GET /reports/popular-items", ordHand.PopularItem)
	return mux
}
