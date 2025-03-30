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

type MenuService interface {
	GetAll(ctx context.Context) ([]*models.MenuItem, error)
	GetByID(ctx context.Context, id int) (*models.MenuItem, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, item *models.MenuItem) error
	Insert(ctx context.Context, item *models.MenuItem) error
}

type MenuHandler struct {
	Service MenuService
	Logger  *slog.Logger
}

func NewMenuHandler(service MenuService, logger *slog.Logger) *MenuHandler {
	return &MenuHandler{service, logger}
}

func (h *MenuHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /menu", middleware.Middleware(h.Insert))
	mux.HandleFunc("GET /menu", middleware.Middleware(h.GetAll))
	mux.HandleFunc("GET /menu/{id}", middleware.Middleware(h.GetElementById))
	mux.HandleFunc("PUT /menu", middleware.Middleware(h.Update))
	mux.HandleFunc("DELETE /menu/{id}", middleware.Middleware(h.Delete))
}

func (h *MenuHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	menu, err := h.Service.GetAll(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get all menu items", "error", err.Error())
		WriteError(w, http.StatusInternalServerError, err, "something went wrong")
		return
	}
	WriteJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) GetElementById(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		h.Logger.Error("Invalid menu item ID", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid menu item ID"), "Invalid ID")
		return
	}
	menuItem, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("Failed to get menu item", "error", err.Error())
		WriteError(w, http.StatusNotFound, err, "Menu item not found")
		return
	}
	WriteJSON(w, http.StatusOK, menuItem)
}

func (h *MenuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		h.Logger.Error("Invalid menu item ID", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid menu item ID"), "Invalid ID")
		return
	}
	if err := h.Service.Delete(r.Context(), id); err != nil {
		h.Logger.Error("Failed to delete menu item", "error", err.Error())
		WriteError(w, http.StatusNotFound, err, "Menu item not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MenuHandler) Update(w http.ResponseWriter, r *http.Request) {
	menuItem := new(models.MenuItem)
	if err := json.NewDecoder(r.Body).Decode(menuItem); err != nil {
		h.Logger.Error("Failed to parse menu item JSON", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON format"), "Invalid JSON")
		return
	}
	if err := h.Service.Update(r.Context(), menuItem); err != nil {
		h.Logger.Error("Failed to update menu item", "error", err.Error())
		WriteError(w, http.StatusInternalServerError, err, "Failed to update menu item")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MenuHandler) Insert(w http.ResponseWriter, r *http.Request) {
	menuItem := new(models.MenuItem)
	if err := json.NewDecoder(r.Body).Decode(menuItem); err != nil {
		h.Logger.Error("Failed to parse menu item JSON", "error", err.Error())
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON format"), "Invalid JSON")
		return
	}
	if err := h.Service.Insert(r.Context(), menuItem); err != nil {
		h.Logger.Error("Failed to insert menu item", "error", err.Error())
		WriteError(w, http.StatusInternalServerError, err, "Failed to insert menu item")
		return
	}
	w.WriteHeader(http.StatusCreated)
}
