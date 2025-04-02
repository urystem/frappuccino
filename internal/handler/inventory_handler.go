package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"frappuccino/internal/service"
	"frappuccino/models"
)

type inventoryHandler struct {
	invSrv service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *inventoryHandler {
	return &inventoryHandler{invSrv: service}
}

func (handl *inventoryHandler) PostInventory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("Put Menu: content type not json")
		writeHttp(w, http.StatusBadRequest, "content type", "invalid")
		return
	}

	newInvent := new(models.Inventory)
	err := json.NewDecoder(r.Body).Decode(newInvent)
	if err != nil {
		slog.Error("Error decoding input JSON data", "error", err)
		writeHttp(w, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	err = handl.invSrv.CreateInventory(newInvent)
	if err != nil {
		slog.Error("Post inventory: "+newInvent.Name, "error", err)
		writeHttp(w, http.StatusInternalServerError, "Inventory", err.Error())
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

	err = bodyJsonStruct(w, invents, http.StatusOK)
	if err != nil {
		slog.Error("Get Invents: Cannot write struct to body")
		return
	}
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
		slog.Error("Get invent by id: ", "failed - ", err)
		writeHttp(w, http.StatusInternalServerError, "invent", err.Error())
		return
	}

	if err = bodyJsonStruct(w, invent, http.StatusOK); err != nil {
		slog.Error("Get Invent: Cannot write struct to body", "id: ", "")
		return
	}
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
		writeHttp(w, http.StatusBadRequest, "content type", "invalid")
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
		slog.Error("Put inventory", "error", err)
		writeHttp(w, http.StatusInternalServerError, "inventory", err.Error())
		return
	}

	slog.Info("put inventory success", "id", inv.ID)
	writeHttp(w, http.StatusOK, "updated", idPath)
}

func (handl *inventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Get invent by id: ", "failed", err)
		writeHttp(w, http.StatusBadRequest, "invent", err.Error())
		return
	}

	menus, err := handl.invSrv.RemoveInventory(id)
	if err != nil {
		if err == models.ErrNotFound {
			slog.Error("Del invent", "not found", err)
			writeHttp(w, http.StatusNotFound, "invent", err.Error())
		} else {
			slog.Error("Del invent", "failed", err)
			writeHttp(w, http.StatusInternalServerError, "invent", err.Error())
		}
		return
	}

	if menus != nil {
		slog.Error("DELETE Invent: found depend menus")
		err = bodyJsonStruct(w, menus, http.StatusBadRequest)
		if err != nil {
			slog.Error("DELETE Invent: Error in decoder")
		}
		return
	}
	slog.Info("Deleted invent", "id", id)
	writeHttp(w, http.StatusNoContent, "", "")
}

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
