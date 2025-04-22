package router

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

func Allrouter(db *sqlx.DB) *http.ServeMux {
	muxRoot := http.NewServeMux()

	inventoryRouter := InventoryRouter(db)
	addPrefixToRouter(muxRoot, inventoryRouter, "/inventory")

	menuMux := menuRouter(db)
	addPrefixToRouter(muxRoot, menuMux, "/menu")

	orderMux := orderRouter(db)
	addPrefixToRouter(muxRoot, orderMux, "/orders")

	reportsMux := aggregationRouter(db)
	addPrefixToRouter(muxRoot, reportsMux, "/reports")

	
	return muxRoot
}

func addPrefixToRouter(mux, child *http.ServeMux, prefix string) {
	// "/" ті пайдалануға болмайт panic болады
	// "/" ПАЙДАЛАНУ ОПАСНО
	// mux.Handle("/reports/", http.StripPrefix("/reports", reportsMux))
	mux.Handle(prefix+"/", http.StripPrefix(prefix, child))
}
