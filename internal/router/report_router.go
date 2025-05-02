package router

import (
	"net/http"

	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"

	"github.com/jmoiron/sqlx"
)

func aggregationReportRouter(db *sqlx.DB) *http.ServeMux {
	dalAggreInter := dal.ReturnDulAggregationDB(db)
	serAggreInter := service.ReturnAggregationService(dalAggreInter)
	handAggre := handler.ReturnAggregationHandInter(serAggreInter)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /total-sales", handAggre.TotalSales)
	mux.HandleFunc("GET /popular-items", handAggre.PopularItems)
	mux.HandleFunc("GET /search", handAggre.FullTextSearchReport)
	mux.HandleFunc("GET /orderedItemsByPeriod", handAggre.PeriodOrderedItems)
	mux.HandleFunc("GET /getLeftOvers", handAggre.GetLeftOvers)
	mux.HandleFunc("GET /numberOfOrderedItems", handAggre.NumberOfOrderedItems)
	return mux
}
