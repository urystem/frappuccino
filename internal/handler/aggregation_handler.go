package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"frappuccino/internal/service"
)

type aggregationHandler struct {
	aggreService service.AggregationServiceInter
}

type AggregationHandInter interface {
	TotalSales(w http.ResponseWriter, r *http.Request)
	PopularItems(w http.ResponseWriter, r *http.Request)
	NumberOfOrderedItems(w http.ResponseWriter, r *http.Request)
}

func ReturnAggregationHandInter(aggreSer service.AggregationServiceInter) AggregationHandInter {
	return &aggregationHandler{aggreService: aggreSer}
}

func (h *aggregationHandler) TotalSales(w http.ResponseWriter, r *http.Request) {
	if total, err := h.aggreService.SumOrder(); err != nil {
		slog.Error("Get total sales", "error", err)
		writeHttp(w, http.StatusInternalServerError, "failed to get total sales:", err.Error())
	} else {
		slog.Info("Succes", "Get total sales:", total)
		err = bodyJsonStruct(w, struct {
			Total_sales float64 // `json: "total_sales"`
		}{total}, http.StatusOK)
		if err != nil {
			slog.Error("")
		} else {
			slog.Info("succes")
		}
	}
}

func (h *aggregationHandler) PopularItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.aggreService.PopularItems()
	if err != nil {
		slog.Error("Get popular sales", "error", err)
		writeHttp(w, http.StatusInternalServerError, "failed to get popular sales:", err.Error())
		return
	}

	err = bodyJsonStruct(w, popularItems, http.StatusOK)
	if err != nil {
		slog.Error("Get popular sales", "error", err)
	} else {
		slog.Info("succes")
	}
}

func (h *aggregationHandler) NumberOfOrderedItems(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	fmt.Println(startDate, endDate)
	h.aggreService.NumberOfOrderedItemsService(startDate, endDate)
}
