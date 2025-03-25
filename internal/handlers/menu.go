package handlers

import (
	"cafeteria/internal/handlers/middleware"
	"cafeteria/internal/helpers"
	"cafeteria/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type MenuService interface {
	GetAll(ctx context.Context) ([]models.Menu, error)
	GetElementById(ctx context.Context, id int) (models.Menu, error)
	Delete(ctx context.Context, id int) error
	Put(ctx context.Context, item models.Menu) error
	Post(ctx context.Context, item models.Menu) error
}

type MenuHandler struct {
	Service MenuService
	Logger  *slog.Logger
}

func NewMenuHandler(service MenuService, logger *slog.Logger) *MenuHandler {
	return &MenuHandler{service, logger}
}

func (h *MenuHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /menu", middleware.Middleware(h.Post))
	mux.HandleFunc("POST /menu/", middleware.Middleware(h.Post))

	mux.HandleFunc("GET /menu", middleware.Middleware(h.GetAll))
	mux.HandleFunc("GET /menu/", middleware.Middleware(h.GetAll))

	mux.HandleFunc("GET /menu/{id}", middleware.Middleware(h.GetElementById))
	mux.HandleFunc("GET /menu/{id}/", middleware.Middleware(h.GetElementById))

	mux.HandleFunc("PUT /menu/{id}", middleware.Middleware(h.Put))
	mux.HandleFunc("PUT /menu/{id}/", middleware.Middleware(h.Put))

	mux.HandleFunc("DELETE /menu/{id}", middleware.Middleware(h.Delete))
	mux.HandleFunc("DELETE /menu/{id}/", middleware.Middleware(h.Delete))
}

func (h *MenuHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.GetAll(r.Context())
	if err != nil {
		h.Logger.Error(err.Error())
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "Failed to fetch menu items")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.Logger.Error(err.Error())
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

	h.Logger.Info("menu items were fetched")
}

func (h *MenuHandler) GetElementById(w http.ResponseWriter, r *http.Request) {
	h.handleRequestWithID(w, r, func(ctx context.Context, id int) error {
		item, err := h.Service.GetElementById(ctx, id)
		if err != nil {
			h.Logger.Error(err.Error())
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		h.Logger.Info("Fetched an menu item", slog.Int("id", id))
		return json.NewEncoder(w).Encode(item)
	})
}

func (h *MenuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.handleRequestWithID(w, r, func(ctx context.Context, id int) error {
		if err := h.Service.Delete(ctx, id); err != nil {
			h.Logger.Error(err.Error())
			return err
		}

		w.WriteHeader(204)
		h.Logger.Info("menu item was deleted", slog.Int("id", id))
		return nil
	})
}

func (h *MenuHandler) Put(w http.ResponseWriter, r *http.Request) {
	h.handleRequestWithID(w, r, func(ctx context.Context, i int) error {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			h.Logger.Error(fmt.Sprintf("error reading request body: %v", err))
			return err
		}
		defer r.Body.Close()

		var item models.Menu
		if err := json.Unmarshal(data, &item); err != nil {
			h.Logger.Error(fmt.Sprintf("error unmarshalling menu item: %v", err))
			return err
		}

		if err := h.Service.Put(r.Context(), item); err != nil {
			h.Logger.Error(fmt.Sprintf("error updating menu item: %v", err))
			return err
		}

		w.WriteHeader(http.StatusOK)
		h.Logger.Info("menu item was updated")
		return nil
	})
}

func (h *MenuHandler) Post(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("error reading request body: %v", err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	var item models.Menu
	if err := json.Unmarshal(data, &item); err != nil {
		h.Logger.Error(fmt.Sprintf("error unmarshalling menu item: %v", err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := h.Service.Post(r.Context(), item); err != nil {
		h.Logger.Error(fmt.Sprintf("error creating menu item: %v", err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.Logger.Info("new menu item was added")
}

func (h *MenuHandler) handleRequestWithID(w http.ResponseWriter, r *http.Request, handler func(context.Context, int) error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid menu ID", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("invalid id: %v", err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := handler(r.Context(), id); err != nil {
		h.Logger.Error(err.Error())
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "No such menu item")
	}
}
