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
	dalCore := dal.ReturnRepoCore(db)
	var dalInventInter dal.InventoryDataAccess = dalCore
	serviceInventInter := service.ReturnInventorySerInt(dalInventInter) // basqasha
	handInvInt := handler.NewInventoryHandler(serviceInventInter)
	// mux.Handle("", mux1)
	mux.HandleFunc("POST /inventory", handInvInt.PostInventory)
	mux.HandleFunc("GET /inventory", handInvInt.GetInventories)
	mux.HandleFunc("GET /inventory/{id}", handInvInt.GetInventoryByID)
	mux.HandleFunc("PUT /inventory/{id}", handInvInt.PutInventory)
	mux.HandleFunc("DELETE /inventory/{id}", handInvInt.DeleteInventory)
	// mux.HandleFunc("PUT /inventory", handInv.PutAllIng)

	// // setup pathfile to dulmenu struct and build to handfunc
	var dalMenuInter dal.MenuDalInter = dalCore
	menuSerInter := service.ReturnMenuSerStruct(dalMenuInter)
	handMenu := handler.ReturnMenuHaldStruct(menuSerInter)

	mux.HandleFunc("GET /menu", handMenu.GetMenus)
	mux.HandleFunc("GET /menu/{id}", handMenu.GetMenuByID)
	mux.HandleFunc("DELETE /menu/{id}", handMenu.DelMenu)
	mux.HandleFunc("POST /menu", handMenu.PostMenu)
	mux.HandleFunc("PUT /menu/{id}", handMenu.PutMenuByID)

	// // setup pathfiles to dulorder struct and build to handlfunc
	var dalOrdInter dal.OrderDalInter = dalCore
	var serOrderInter service.OrdServiceInter = service.ReturnOrdSerStruct(dalOrdInter)
	handOrd := handler.ReturnOrdHaldStruct(serOrderInter)

	mux.HandleFunc("GET /orders", handOrd.GetOrders)
	mux.HandleFunc("GET /orders/{id}", handOrd.GetOrderByID)
	// ordHand := handler.ReturnOrdHaldStruct(ordSer)
	// mux.HandleFunc("POST /orders", ordHand.PostOrder)
	// mux.HandleFunc("PUT /orders/{id}", ordHand.PutOrdById)
	// mux.HandleFunc("DELETE /orders/{id}", ordHand.DelOrdById)
	// mux.HandleFunc("POST /orders/{id}/close", ordHand.PostOrdCloseById)
	// mux.HandleFunc("GET /reports/total-sales", ordHand.TotalSales)
	// mux.HandleFunc("GET /reports/popular-items", ordHand.PopularItem)
	return mux
}
