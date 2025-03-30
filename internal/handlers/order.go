package handlers

import (
	"cafeteria/internal/handlers/middleware"
	"cafeteria/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type OrderService interface {
	GetAll(ctx context.Context) ([]*models.Order, error)
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, order *models.Order) error
	Insert(ctx context.Context, order *models.Order) error
	ProcessBatchOrders(ctx context.Context, orders []models.BatchOrder) ([]models.OrderResult, models.BatchSummary, error)
}

type OrderHandler struct {
	Service OrderService
	Logger  *slog.Logger
}

func NewOrderHandler(service OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{service, logger}
}

func (h *OrderHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", middleware.Middleware(h.Insert))
	mux.HandleFunc("GET /orders", middleware.Middleware(h.GetAll))
	mux.HandleFunc("GET /orders/{id}", middleware.Middleware(h.GetElementById))
	mux.HandleFunc("PUT /orders", middleware.Middleware(h.Update))
	mux.HandleFunc("DELETE /orders/{id}", middleware.Middleware(h.Delete))
	mux.HandleFunc("POST /orders/batch-process", middleware.Middleware(h.ProcessBatchOrders))
}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.Service.GetAll(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get all orders: ", "error", err)
		WriteError(w, http.StatusBadRequest, err, "something went wrong")
		return
	}

	WriteJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetElementById(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		h.Logger.Error("Failed to get an order: ", "error", "invalid order id")
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order id"), "error")
		return
	}

	order, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("Failed to get an order: ", "error", err)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("no such order"), "something went wrong")
		return
	}

	WriteJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		h.Logger.Error("Failed to delete an order: ", "error", "invalid order id")
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order id"), "error")
		return
	}

	if err := h.Service.Delete(r.Context(), id); err != nil {
		h.Logger.Error("Failed to delete an order: ", "error", err)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("no such order"), "something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	order := new(models.Order)
	if err := ParseJSON(r, order); err != nil {
		h.Logger.Error("Failed parsing JSON for an order: ", "error", err)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect JSON format"), "something went wrong")
		return
	}
	if err := h.Service.Update(r.Context(), order); err != nil {
		h.Logger.Error("Failed to update an order: ", "error", err)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("no such order"), "something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) Insert(w http.ResponseWriter, r *http.Request) {
	order := new(models.Order)
	if err := ParseJSON(r, order); err != nil {
		h.Logger.Error("Failed parsing JSON for an order: ", "error", err)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect JSON format"), "something went wrong")
		return
	}

	if err := h.Service.Insert(r.Context(), order); err != nil {
		h.Logger.Error("Failed to insert an order: ", "error", err)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("error inserting order"), "something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) ProcessBatchOrders(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Orders []models.BatchOrder `json:"orders"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, err, "invalid request body")
		return
	}

	results, summary, err := h.Service.ProcessBatchOrders(r.Context(), request.Orders)
	if err != nil {
		h.Logger.Error("Failed to process batch orders", "error", err)
		WriteError(w, http.StatusInternalServerError, err, "failed to process orders")
		return
	}

	response := map[string]interface{}{
		"processed_orders": results,
		"summary":          summary,
	}

	WriteJSON(w, http.StatusOK, response)
}
