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
	var serviceInventInter service.InventoryService = service.NewInventoryService(dalInventInter)
	handInv := handler.NewInventoryHandler(serviceInventInter)
	// mux.Handle("", mux1)
	mux.HandleFunc("POST /inventory", handInv.PostInventory)
	mux.HandleFunc("GET /inventory", handInv.GetInventories)
	mux.HandleFunc("GET /inventory/{id}", handInv.GetInventoryByID)
	mux.HandleFunc("PUT /inventory/{id}", handInv.PutInventory)
	mux.HandleFunc("DELETE /inventory/{id}", handInv.DeleteInventory)
	// mux.HandleFunc("PUT /inventory", handInv.PutAllIng)

	// // setup pathfile to dulmenu struct and build to handfunc
	var dalMenuInter dal.MenuDalInter = dalCore
	var menuSerInter service.MenuServiceInter = service.ReturnMenuSerStruct(dalMenuInter)
	handMenu := handler.ReturnMenuHaldStruct(menuSerInter)
	mux.HandleFunc("POST /menu", handMenu.PostMenu)
	mux.HandleFunc("GET /menu", handMenu.GetMenus)
	mux.HandleFunc("GET /menu/{id}", handMenu.GetMenuByID)
	// mux.HandleFunc("PUT /menu/{id}", menuHand.PutMenuById)
	mux.HandleFunc("DELETE /menu/{id}", handMenu.DelMenu)

	// // setup pathfiles to dulorder struct and build to handlfunc
	// var dalOrdInter dal.OrderDalInter = dal.ReturnOrdDalStruct(*dir + PathFiles[0])
	// var ordSer service.OrdServiceInter = service.ReturnOrdSerStruct(dalOrdInter, dalMenuInter, dalInventInter)
	// ordHand := handler.ReturnOrdHaldStruct(ordSer)
	// mux.HandleFunc("POST /orders", ordHand.PostOrder)
	// mux.HandleFunc("GET /orders", ordHand.GetOrders)
	// mux.HandleFunc("GET /orders/{id}", ordHand.GetOrdById)
	// mux.HandleFunc("PUT /orders/{id}", ordHand.PutOrdById)
	// mux.HandleFunc("DELETE /orders/{id}", ordHand.DelOrdById)
	// mux.HandleFunc("POST /orders/{id}/close", ordHand.PostOrdCloseById)
	// mux.HandleFunc("GET /reports/total-sales", ordHand.TotalSales)
	// mux.HandleFunc("GET /reports/popular-items", ordHand.PopularItem)
	return mux
}
