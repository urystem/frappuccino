package handlers

import (
	"cafeteria/internal/handlers/middleware"
	"context"
	"log/slog"
	"net/http"
)

type TransactionService interface {
	TotalSales(ctx context.Context) (float32, error)
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
}

func (h *TransactionHandler) TotalSales(w http.ResponseWriter, r *http.Request) {
	total, err := h.Service.TotalSales(r.Context())
	if err != nil {

	}

	WriteJSON(w, 200, map[string]float32{"total_sales": total})
}
