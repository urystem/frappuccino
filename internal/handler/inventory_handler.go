package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type inventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *inventoryHandler {
	return &inventoryHandler{inventoryService: service}
}

func (h *inventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var newInvent models.InventoryItem
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("Put Menu: content type not json")
		writeHttp(w, http.StatusBadRequest, "content type", "invalid")
	} else if err := json.NewDecoder(r.Body).Decode(&newInvent); err != nil {
		slog.Error("Error decoding input JSON data", "error", err)
		writeHttp(w, http.StatusBadRequest, "Invalid input", err.Error())
	} else if err = h.inventoryService.InsertInventory(&newInvent); err != nil {
		slog.Error("Post inventory: "+newInvent.Name, "error", err)
		if err == models.ErrConflict {
			writeHttp(w, http.StatusConflict, "Inventory ", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "Inventory", err.Error())
		}
	} else {
		slog.Info("Post inventory: ", "success", newInvent.Name)
		writeHttp(w, http.StatusCreated, "inventory", "created")
	}
}

// func (h *inventoryHandler) GetAllInventory(w http.ResponseWriter, r *http.Request) {
// 	if invents, err := h.inventoryService.GetAllInventory(); err != nil {
// 		slog.Error("Can't get all inventory")
// 		writeHttp(w, http.StatusInternalServerError, "get all invents", err.Error())
// 	} else if err := bodyJsonStruct(w, invents); err != nil {
// 		slog.Error("Get Invents: Cannot write struct to body")
// 	}
// }

// func (h *inventoryHandler) GetSpecificInventory(w http.ResponseWriter, r *http.Request) {
// 	if id := r.PathValue("id"); checkName(id) {
// 		slog.Error("Get InventBYid: Invalid id")
// 		writeHttp(w, http.StatusBadRequest, "", "")
// 	} else if invent, err := h.inventoryService.GetSpecificInventory(id); err != nil {
// 		slog.Error("Get invent by id: ", "failed - ", err)
// 		if err == models.ErrNotFound {
// 			writeHttp(w, http.StatusNotFound, "inventory", err.Error())
// 		} else {
// 			writeHttp(w, http.StatusInternalServerError, "invent", err.Error())
// 		}
// 	} else if err = bodyJsonStruct(w, invent); err != nil {
// 		slog.Error("Get Invent: Cannot write struct to body", "id: ", id)
// 	}
// }

// func (h *inventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
// 	var newInvent models.InventoryItem
// 	if id := r.PathValue("id"); checkName(id) {
// 		slog.Error("Put Invent: nvalid id in url")
// 		writeHttp(w, http.StatusBadRequest, "id url", "invalid id")
// 	} else if r.Header.Get("Content-Type") != "application/json" {
// 		slog.Error("Put Menu: content type not json")
// 		writeHttp(w, http.StatusBadRequest, "content type", "invalid")
// 	} else if err := json.NewDecoder(r.Body).Decode(&newInvent); err != nil {
// 		slog.Error("Put Invent: Error in decoder")
// 		writeHttp(w, http.StatusBadRequest, "inventory", err.Error())
// 	} else if err = checkInventStruct(&newInvent, false); err != nil {
// 		slog.Error("Put invent: ", "invalid struct", err)
// 		writeHttp(w, http.StatusBadRequest, "invalid struct", err.Error())
// 	} else if err = h.inventoryService.UpdateInventory(id, &newInvent); err != nil {
// 		slog.Error("Put inventory", "error", err)
// 		if err == models.ErrNotFound {
// 			writeHttp(w, http.StatusNotFound, "inventory", err.Error())
// 		} else {
// 			writeHttp(w, http.StatusInternalServerError, "inventory", err.Error())
// 		}
// 	} else {
// 		slog.Info("put inventory success", "id", newInvent.IngredientID)
// 		writeHttp(w, http.StatusOK, "updated", newInvent.IngredientID)
// 	}
// }

// func (h *inventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
// 	if id := r.PathValue("id"); checkName(id) {
// 		slog.Error("Del invent", "invalid", "id")
// 		writeHttp(w, http.StatusBadRequest, "Invalid", "id")
// 	} else if err := h.inventoryService.DeleteInventory(id); err != nil {
// 		slog.Error("Del invent", "failed", err)
// 		if err == models.ErrNotFound {
// 			writeHttp(w, http.StatusNotFound, "invent", err.Error())
// 		} else {
// 			writeHttp(w, http.StatusInternalServerError, "invent", err.Error())
// 		}
// 	} else {
// 		slog.Info("Deleted invent", "id", id)
// 		writeHttp(w, http.StatusNoContent, "", "")
// 	}
// }

// func (h *inventoryHandler) PutAllIng(w http.ResponseWriter, r *http.Request) {
// 	var invents []models.InventoryItem
// 	if err := json.NewDecoder(r.Body).Decode(&invents); err != nil {
// 		slog.Error("Wrong json or error with decode the body", "error", err)
// 		writeHttp(w, http.StatusBadRequest, "decode th json", err.Error())
// 		return
// 	}
// 	for _, v := range invents {
// 		if err := checkInventStruct(&v, false); err != nil {
// 			slog.Error(v.IngredientID + ": wrong struct")
// 			writeHttp(w, http.StatusBadRequest, "put some invents:"+v.IngredientID, err.Error())
// 			return
// 		}
// 	}

// 	if ings, err := h.inventoryService.PutAllInvets(invents); err != nil {
// 		if err == models.ErrNotFoundIngs {
// 			writeHttp(w, http.StatusNotFound, "inventory ", strings.Join(ings, ", ")+err.Error())
// 		} else {
// 			writeHttp(w, http.StatusInternalServerError, "put some invets", err.Error())
// 		}
// 	} else {
// 		slog.Info("Invents updated successfully")
// 		writeHttp(w, http.StatusOK, "invents", "updated")
// 	}
// }
