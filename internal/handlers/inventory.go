package handlers

import (
	"cafeteria/internal/handlers/middleware"
	"cafeteria/internal/models"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type InventoryService interface {
	GetAll(ctx context.Context) ([]*models.InventoryItem, error)
	GetByID(ctx context.Context, id int) (*models.InventoryItem, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, item *models.InventoryItem) error
	Insert(ctx context.Context, item *models.InventoryItem) error
	GetLeftovers(ctx context.Context, sortBy string, page, pageSize int) (models.LeftoversResponse, error)
}

type InventoryHandler struct {
	Service InventoryService
	Logger  *slog.Logger
}

func NewInventoryHandler(service InventoryService, logger *slog.Logger) *InventoryHandler {
	return &InventoryHandler{service, logger}
}

func (h *InventoryHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /inventory", middleware.Middleware(h.Insert))
	mux.HandleFunc("GET /inventory", middleware.Middleware(h.GetAll))
	mux.HandleFunc("GET /inventory/{id}", middleware.Middleware(h.GetById))
	mux.HandleFunc("PUT /inventory", middleware.Middleware(h.Update))
	mux.HandleFunc("DELETE /inventory/{id}", middleware.Middleware(h.Delete))
	mux.HandleFunc("GET /inventory/getLeftOvers", middleware.Middleware(h.GetLeftovers))
}

func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	inventory, err := h.Service.GetAll(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get all inventory items: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, err, "something went wrong")
		return
	}

	WriteJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		h.Logger.Error("Failed to get an inventory item: ", "error", "invalid error id")
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid inventory id"), "error")
		return
	}

	inventory, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("Failed to get an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("no such inventory item"), "something went wrong")
		return
	}

	WriteJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")

	id, err := strconv.Atoi(rawId)
	if err != nil {
		h.Logger.Error("Failed to delete an inventory item: ", "error", "invalid error id")
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid inventory id"), "error")
		return
	}

	if err := h.Service.Delete(r.Context(), id); err != nil {
		h.Logger.Error("Failed to delete an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("no such inventory item"), "something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	inventory := new(models.InventoryItem)
	if err := ParseJSON(r, inventory); err != nil {
		h.Logger.Error("Failed parsing json on an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect json format"), "something went wrong")
		return
	}

	if err := inventory.IsValid(); err != nil {
		h.Logger.Error("Failed parsing json on an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, err, "something went wrong")
		return
	}

	if err := h.Service.Update(r.Context(), inventory); err != nil {
		h.Logger.Error("Failed to delete an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("no such inventory item"), "something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) Insert(w http.ResponseWriter, r *http.Request) {
	inventory := new(models.InventoryItem)
	if err := ParseJSON(r, inventory); err != nil {
		h.Logger.Error("Failed parsing json on an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect json format"), "something went wrong")
		return
	}

	if err := inventory.IsValid(); err != nil {
		h.Logger.Error("Failed parsing json on an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, err, "something went wrong")
		return
	}

	if err := h.Service.Insert(r.Context(), inventory); err != nil {
		h.Logger.Error("Failed to insert an inventory item: ", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("item with such name already exists"), "something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) GetLeftovers(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	sortBy := r.URL.Query().Get("sortBy")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	// Get leftovers
	response, err := h.Service.GetLeftovers(r.Context(), sortBy, page, pageSize)
	if err != nil {
		h.Logger.Error("Failed to get leftovers", "error", err.Error())
		WriteError(w, http.StatusInternalServerError, err, "failed to get leftovers")
		return
	}

	WriteJSON(w, http.StatusOK, response)
}
