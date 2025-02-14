package router

import (
	"net/http"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

var PathFiles [3]string = [3]string{"/orders.json", "/menu_items.json", "/inventory.json"}

func Allrouter(dir *string) *http.ServeMux {
	mux := http.NewServeMux()

	// setup pathfile to dulinvent and build to handfunc
	var dalInventInter dal.InventoryDataAccess = dal.NewInventoryRepository(*dir + PathFiles[2])
	var serviceInventInter service.InventoryService = service.NewInventoryService(dalInventInter)
	handInv := handler.NewInventoryHandler(serviceInventInter)
	mux.HandleFunc("POST /inventory", handInv.CreateInventory)
	mux.HandleFunc("GET /inventory", handInv.GetAllInventory)
	mux.HandleFunc("GET /inventory/{id}", handInv.GetSpecificInventory)
	mux.HandleFunc("PUT /inventory/{id}", handInv.UpdateInventory)
	mux.HandleFunc("DELETE /inventory/{id}", handInv.DeleteInventory)
	mux.HandleFunc("PUT /inventory", handInv.PutAllIng)

	// setup pathfile to dulmenu struct and build to handfunc
	var dalMenuInter dal.MenuDalInter = dal.ReturnMenuDalStruct(*dir + PathFiles[1])
	var menuSer service.MenuServiceInter = service.ReturnMenuSerStruct(dalMenuInter, dalInventInter)
	menuHand := handler.ReturnMenuHaldStruct(menuSer)
	mux.HandleFunc("POST /menu", menuHand.PostMenu)
	mux.HandleFunc("GET /menu", menuHand.GetAllMenus)
	mux.HandleFunc("GET /menu/{id}", menuHand.GetMenuById)
	mux.HandleFunc("PUT /menu/{id}", menuHand.PutMenuById)
	mux.HandleFunc("DELETE /menu/{id}", menuHand.DelMenuById)

	// setup pathfiles to dulorder struct and build to handlfunc
	var dalOrdInter dal.OrderDalInter = dal.ReturnOrdDalStruct(*dir + PathFiles[0])
	var ordSer service.OrdServiceInter = service.ReturnOrdSerStruct(dalOrdInter, dalMenuInter, dalInventInter)
	ordHand := handler.ReturnOrdHaldStruct(ordSer)
	mux.HandleFunc("POST /orders", ordHand.PostOrder)
	mux.HandleFunc("GET /orders", ordHand.GetOrders)
	mux.HandleFunc("GET /orders/{id}", ordHand.GetOrdById)
	mux.HandleFunc("PUT /orders/{id}", ordHand.PutOrdById)
	mux.HandleFunc("DELETE /orders/{id}", ordHand.DelOrdById)
	mux.HandleFunc("POST /orders/{id}/close", ordHand.PostOrdCloseById)
	mux.HandleFunc("GET /reports/total-sales", ordHand.TotalSales)
	mux.HandleFunc("GET /reports/popular-items", ordHand.PopularItem)
	return mux
}
