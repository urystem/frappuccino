package aggregation

import (
	"net/http"

	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"

	"github.com/jmoiron/sqlx"
)

type aggregationRoute struct {
	aggreHandler handler.AggregationHandInter
}

type AggregationRouters interface {
	AggregationReportRouter() *http.ServeMux
	AggregationsAsRootMux() *http.ServeMux
}

// constructor
func NewAggregationRouter(db *sqlx.DB) AggregationRouters {
	dalAggreInter := dal.ReturnDulAggregationDB(db)
	serAggreInter := service.ReturnAggregationService(dalAggreInter)
	handAggre := handler.ReturnAggregationHandInter(serAggreInter)
	return &aggregationRoute{aggreHandler: handAggre}
}

// report
func (aggreRoute *aggregationRoute) AggregationReportRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /total-sales", aggreRoute.aggreHandler.TotalSales)
	mux.HandleFunc("GET /popular-items", aggreRoute.aggreHandler.PopularItems)
	mux.HandleFunc("GET /search", aggreRoute.aggreHandler.FullTextSearchReport)
	return mux
}

// like only at
// As Root
func (aggreRoute *aggregationRoute) AggregationsAsRootMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /numberOfOrderedItems", aggreRoute.aggreHandler.NumberOfOrderedItems)
	return mux
}
