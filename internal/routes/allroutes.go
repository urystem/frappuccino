package router

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

func Allrouter(db *sqlx.DB) *http.ServeMux {
	muxRoot := http.NewServeMux()

	inventoryRouter := inventoryRouter(db)

	addPrefixToRouter("/inventory", muxRoot, inventoryRouter)

	menuMux := menuRouter(db)
	addPrefixToRouter("/menu", muxRoot, menuMux)

	orderMux := orderRouter(db)
	addPrefixToRouter("/orders", muxRoot, orderMux)

	reports := aggregationReportRouter(db)
	addPrefixToRouter("/reports", muxRoot, reports)

	return muxRoot
}

func addPrefixToRouter(prefix string, mux, child *http.ServeMux) {
	// "/" ті пайдалануға болмайт panic болады
	// "/" ПАЙДАЛАНУ ОПАСНО
	// mux.Handle("/reports/", http.StripPrefix("/reports", reportsMux))
	mux.Handle(prefix+"/", http.StripPrefix(prefix, child))
}
