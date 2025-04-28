package handler

import (
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
	FullTextSearchReport(w http.ResponseWriter, r *http.Request)
	PeriodOrderedItems(w http.ResponseWriter, r *http.Request)
	GetLeftOvers(w http.ResponseWriter, r *http.Request)
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

	numberOf, err := h.aggreService.NumberOfOrderedItemsService(startDate, endDate)
	if err != nil {
		slog.Error("Get number of sales", "error", err)
		writeHttp(w, http.StatusInternalServerError, "failed to get number sales:", err.Error())
		return
	}

	err = bodyJsonStruct(w, numberOf, http.StatusOK)
	if err != nil {
		slog.Error("Get number of sales", "error", err)
	} else {
		slog.Info("succes")
	}
}

func (h *aggregationHandler) FullTextSearchReport(w http.ResponseWriter, r *http.Request) {
	find := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")
	res, err := h.aggreService.Search(find, filter, minPrice, maxPrice)
	if err != nil {
		slog.Error("Get search", "error", err)
		return
	}
	err = bodyJsonStruct(w, res, http.StatusOK)
	if err != nil {
		slog.Error("Get search", "error", err)
	} else {
		slog.Info("succes")
	}
}

func (h *aggregationHandler) PeriodOrderedItems(w http.ResponseWriter, r *http.Request) {
	dayOrMonth := r.URL.Query().Get("period")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	if dayOrMonth == "" {
		writeHttp(w, http.StatusBadRequest, "", "period parameter is required")
		return
	}
	orderStats, err := h.aggreService.OrderedItemsPeriod(dayOrMonth, month, year)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = bodyJsonStruct(w, orderStats, http.StatusOK)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("succes")
}

func (h *aggregationHandler) GetLeftOvers(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")
	overs, err := h.aggreService.GetLeftOversService(sortBy, page, pageSize)
	if err != nil {
		slog.Error("overs", err)
		writeHttp(w, http.StatusBadRequest, "error", err.Error())
		return
	}
	err = bodyJsonStruct(w, overs, http.StatusOK)
	if err != nil {
		slog.Error("overs", err)
		writeHttp(w, http.StatusBadRequest, "error", err.Error())
		return
	}
	slog.Info("overs", "OK")
}
