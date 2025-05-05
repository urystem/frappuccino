package router

import (
	"net/http"

	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"

	"github.com/jmoiron/sqlx"
)

func menuRouter(db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()

	var dalMenuInter dal.MenuDalInter = dal.ReturnDalMenuCore(db)
	menuSerInter := service.ReturnMenuSerStruct(dalMenuInter)
	handMenu := handler.ReturnMenuHaldStruct(menuSerInter)

	mux.HandleFunc("GET /", handMenu.GetMenus)
	mux.HandleFunc("GET /{id}", handMenu.GetMenuByID)
	mux.HandleFunc("DELETE /{id}", handMenu.DelMenu)
	mux.HandleFunc("POST /", handMenu.PostMenu)
	mux.HandleFunc("PUT /{id}", handMenu.PutMenuByID)
	mux.HandleFunc("GET /history", handMenu.GetHistory)
	return mux
}
