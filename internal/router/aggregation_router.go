package router

import (
	"net/http"

	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"

	"github.com/jmoiron/sqlx"
)

func aggregationRouter(db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()
	dalAggreInter := dal.ReturnDulAggregationDB(db)
	serAggreInter := service.ReturnAggregationService(dalAggreInter)
	handAggre := handler.ReturnAggregationHandInter(serAggreInter)
	mux.HandleFunc("GET /total-sales", handAggre.TotalSales)
	mux.HandleFunc("GET /popular-items", handAggre.PopularItems)
	return mux
}
