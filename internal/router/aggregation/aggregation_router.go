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
	AggregationReportRouter() *http.ServeMux // report
	AggregationsAsRootMux() *http.ServeMux   // /
	// inventory
}

// constructor
func NewAggregationRouter(db *sqlx.DB) AggregationRouters {
	dalAggreInter := dal.ReturnDulAggregationDB(db)
	serAggreInter := service.ReturnAggregationService(dalAggreInter)
	handAggre := handler.ReturnAggregationHandInter(serAggreInter)
	return &aggregationRoute{aggreHandler: handAggre}
}

// report
func (agg *aggregationRoute) AggregationReportRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /total-sales", agg.aggreHandler.TotalSales)
	mux.HandleFunc("GET /popular-items", agg.aggreHandler.PopularItems)
	mux.HandleFunc("GET /search", agg.aggreHandler.FullTextSearchReport)
	mux.HandleFunc("GET /orderedItemsByPeriod", agg.aggreHandler.PeriodOrderedItems)
	mux.HandleFunc("GET /getLeftOvers", agg.aggreHandler.GetLeftOvers)
	return mux
}

// like only at
// As Root
func (agg *aggregationRoute) AggregationsAsRootMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /numberOfOrderedItems", agg.aggreHandler.NumberOfOrderedItems)
	return mux
}
