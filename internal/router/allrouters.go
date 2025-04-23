package router

import (
	"net/http"

	"frappuccino/internal/router/aggregation" // for mux

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

	aggregations := aggregation.NewAggregationRouter(db)

	reports := aggregations.AggregationReportRouter()
	addPrefixToRouter("/reports", muxRoot, reports)

	asRoot := aggregations.AggregationsAsRootMux()
	addPrefixToRouter("", muxRoot, asRoot)

	return muxRoot
}

func addPrefixToRouter(prefix string, mux, child *http.ServeMux) {
	// "/" ті пайдалануға болмайт panic болады
	// "/" ПАЙДАЛАНУ ОПАСНО
	// mux.Handle("/reports/", http.StripPrefix("/reports", reportsMux))
	mux.Handle(prefix+"/", http.StripPrefix(prefix, child))
}
