package handlers

import (
	"cafeteria/internal/handlers/middleware"
	"cafeteria/internal/models"
	"context"
	"log/slog"
	"net/http"
	"time"
)

type TransactionService interface {
	TotalSales(ctx context.Context) (float32, error)
	PopularItems(ctx context.Context) (models.JSONB, error)
	NumberOfOrderedItems(ctx context.Context, start, end time.Time) (models.JSONB, error)
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
