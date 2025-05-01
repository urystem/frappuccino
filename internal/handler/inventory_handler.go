package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"frappuccino/internal/service"
	"frappuccino/models"
)

type inventoryHandler struct {
	invSrv service.InventoryService
}

type inventoryHandlerInt interface {
	PostInventory(w http.ResponseWriter, r *http.Request)
	GetInventories(w http.ResponseWriter, r *http.Request)
	GetInventoryByID(w http.ResponseWriter, r *http.Request)
	PutInventory(w http.ResponseWriter, r *http.Request)
	DeleteInventory(w http.ResponseWriter, r *http.Request)
}

func NewInventoryHandler(service service.InventoryService) inventoryHandlerInt {
	return &inventoryHandler{invSrv: service}
}

func (handl *inventoryHandler) PostInventory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("Handler: post inventory -> content type is not json")
		writeHttp(w, http.StatusUnsupportedMediaType, "content type", "invalid")
		return
	}

	newInvent := new(models.Inventory)
	err := json.NewDecoder(r.Body).Decode(newInvent)
	if err != nil {
		slog.Error("Handler: post inventory -> Error decoding input JSON data", "error", err)
		writeHttp(w, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	err = handl.invSrv.CreateInventory(newInvent)
	if err != nil {
		slog.Error("Post inventory: ", newInvent.Name, err)
		code := http.StatusInternalServerError
		if errors.Is(err, models.ErrConflict) {
			code = http.StatusConflict
		} else if errors.Is(err, models.ErrBadInput) {
			code = http.StatusUnprocessableEntity
		}
		writeHttp(w, code, "Inventory", err.Error())
		return
	}

	writeHttp(w, http.StatusCreated, "inventory", "created")
	slog.Info("Post inventory: ", "success", newInvent.ID)
}

func (handl *inventoryHandler) GetInventories(w http.ResponseWriter, r *http.Request) {
	invents, err := handl.invSrv.CollectInventories()
	if err != nil {
		slog.Error("Can't get all inventory")
		writeHttp(w, http.StatusInternalServerError, "get all invents", err.Error())
		return
	}

	bodyJsonStruct(w, invents, http.StatusOK)
	slog.Info("Get", "inventories:", "succes")
}

func (handl *inventoryHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Get invent by id: ", "failed", err)
		writeHttp(w, http.StatusBadRequest, "invent", err.Error())
		return
	}

	invent, err := handl.invSrv.TakeInventory(id)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, models.ErrNotFound) {
			code = http.StatusNotFound
		}
		slog.Error("Get invent by id: ", "failed - ", err)
		writeHttp(w, code, "invent", err.Error())
		return
	}

	bodyJsonStruct(w, invent, http.StatusOK)
	slog.Info("get ", "inventory", "success")
}

func (handl *inventoryHandler) PutInventory(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")

	id, err := strconv.ParseUint(idPath, 10, 0)
	if err != nil {
		slog.Error("Put Invent: invalid parse id")
		writeHttp(w, http.StatusBadRequest, "id url", "invalid id")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("Put inventory: content type not json")
		writeHttp(w, http.StatusUnsupportedMediaType, "content type", "invalid")
		return
	}

	inv := new(models.Inventory)

	if err = json.NewDecoder(r.Body).Decode(inv); err != nil {
		slog.Error("Put Invent: Error in decoder")
		writeHttp(w, http.StatusBadRequest, "inventory", err.Error())
		return
	}

	inv.ID = id

	err = handl.invSrv.UpgradeInventory(inv)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, models.ErrBadInput) {
			code = http.StatusUnprocessableEntity
		} else if errors.Is(err, models.ErrNotFound) {
			code = http.StatusNotFound
		}
		slog.Error("Put inventory", "error", err)
		writeHttp(w, code, "inventory", err.Error())
		return
	}

	slog.Info("put inventory success", "id", inv.ID)
	writeHttp(w, http.StatusOK, "updated", idPath)
}

func (handl *inventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Del invent ", "failed parse to uint64", err)
		writeHttp(w, http.StatusBadRequest, "invent", err.Error())
		return
	}

	menus, err := handl.invSrv.RemoveInventory(id)
	if err != nil {
		slog.Error("Del invent", "error", err, "id = ", id)
		code := http.StatusInternalServerError
		if err == models.ErrNotFound {
			code = http.StatusNotFound
		}
		slog.Error("Del invent", "failed", err, "id = ", id)
		writeHttp(w, code, "invent", err.Error())
		return
	}

	if menus != nil {
		slog.Error("DELETE Invent:", "found depend menus id = ", id)
		bodyJsonStruct(w, menus, http.StatusFailedDependency)
		return
	}
	slog.Info("Deleted invent", "id", id)
	writeHttp(w, http.StatusNoContent, "", "")
}
