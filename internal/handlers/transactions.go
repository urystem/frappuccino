package handlers

import (
	"cafeteria/internal/handlers/middleware"
	"cafeteria/internal/models"
	"context"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TransactionService interface {
	TotalSales(ctx context.Context) (float32, error)
	PopularItems(ctx context.Context) (models.JSONB, error)
	NumberOfOrderedItems(ctx context.Context, start, end time.Time) (models.JSONB, error)
	SearchOrders(ctx context.Context, q, filter string, low, high float32) (models.JSONB, error)
	OrderedItemsByPeriod(ctx context.Context, period, month string, year int) (models.JSONB, error)
}

type TransactionHandler struct {
	Service TransactionService
	Logger  *slog.Logger
}

func NewTransactionHandler(service TransactionService, logger *slog.Logger) *TransactionHandler {
	return &TransactionHandler{service, logger}
}

func (h *TransactionHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /total-sales", middleware.Middleware(h.TotalSales))
	mux.HandleFunc("GET /popular-items", middleware.Middleware(h.PopularItems))
	mux.HandleFunc("GET /numberOfOrderedItems", middleware.Middleware(h.NumberOfOrderedItems))
	mux.HandleFunc("GET /search", middleware.Middleware(h.Search))
	mux.HandleFunc("GET /orderedItemsByPeriod", middleware.Middleware(h.OrderedItemsByPeriod))
}

func (h *TransactionHandler) TotalSales(w http.ResponseWriter, r *http.Request) {
	total, err := h.Service.TotalSales(r.Context())
	if err != nil {
		h.Logger.Error("Error while fetching popular items", "error", err)
		WriteError(w, http.StatusBadRequest, err, "something went wrong")
		return
	}

	WriteJSON(w, 200, map[string]float32{"total_sales": total})
}

func (h *TransactionHandler) PopularItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.PopularItems(r.Context())
	if err != nil {
		h.Logger.Error("Error while fetching popular items", "error", err)
		WriteError(w, http.StatusBadRequest, err, "something went wrong")
		return
	}

	WriteJSON(w, http.StatusOK, items)
}

func (h *TransactionHandler) NumberOfOrderedItems(w http.ResponseWriter, r *http.Request) {
	start, end := r.FormValue("startDate"), r.FormValue("endDate")

	layout := "2006-01-02"
	s, _ := time.Parse(layout, start)
	e, _ := time.Parse(layout, end)

	items, err := h.Service.NumberOfOrderedItems(r.Context(), s, e)
	if err != nil {
		h.Logger.Error("Error while fetching popular items", "error", err)
		WriteError(w, http.StatusBadRequest, err, "something went wrong")
		return
	}

	WriteJSON(w, http.StatusOK, items)
}

func (h *TransactionHandler) Search(w http.ResponseWriter, r *http.Request) {
	q, minPrice, maxPrice := r.FormValue("q"), r.FormValue("minPrice"), r.FormValue("maxPrice")
	if len(strings.Trim(q, " ")) == 0 {
		WriteError(w, http.StatusBadRequest, errors.New("search query cannot be empty"), "error")
		return
	}

	minFloat, _ := strconv.ParseFloat(minPrice, 32)
	maxFloat, _ := strconv.ParseFloat(maxPrice, 32)

	minFloat = max(0, minFloat)
	maxFloat = max(0, maxFloat)
	if maxFloat == 0 {
		maxFloat = math.MaxFloat32
	}

	filter := strings.Trim(r.FormValue("filter"), " ")
	if filter != "menu" && filter != "orders" {
		filter = "all"
	}

	body, err := h.Service.SearchOrders(r.Context(), q, filter, float32(minFloat), float32(maxFloat))
	if err != nil {
		h.Logger.Error("Failed whlile making search", "error", err)
		WriteError(w, http.StatusInternalServerError, err, "something went wrong")
		return
	}

	WriteJSON(w, http.StatusOK, body)
}

func (h *TransactionHandler) OrderedItemsByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period != "day" && period != "month" {
		WriteError(w, http.StatusBadRequest, errors.New("invalid period parameter"), "period must be either 'day' or 'month'")
		return
	}

	var month string
	var year int
	var err error

	if period == "day" {
		month = r.URL.Query().Get("month")
		if month == "" {
			WriteError(w, http.StatusBadRequest, errors.New("month parameter is required for day period"), "month parameter is required")
			return
		}
	}

	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		year = time.Now().Year()
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err, "invalid year parameter")
			return
		}
	}

	items, err := h.Service.OrderedItemsByPeriod(r.Context(), period, month, year)
	if err != nil {
		h.Logger.Error("Error while fetching ordered items by period", "error", err)
		WriteError(w, http.StatusInternalServerError, err, "something went wrong")
		return
	}

	response := map[string]interface{}{
		"period":       period,
		"orderedItems": items,
	}

	if period == "day" {
		response["month"] = month
	} else {
		response["year"] = year
	}

	WriteJSON(w, http.StatusOK, response)
}
