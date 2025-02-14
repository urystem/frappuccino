package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type menuHaldToService struct {
	menuServInt service.MenuServiceInter
}

func ReturnMenuHaldStruct(menuSerInt service.MenuServiceInter) *menuHaldToService {
	return &menuHaldToService{menuServInt: menuSerInt}
}

func (h *menuHaldToService) PostMenu(w http.ResponseWriter, r *http.Request) {
	var menuStruct models.MenuItem
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("post the menu: content_Type must be application/json")
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
	} else if err := json.NewDecoder(r.Body).Decode(&menuStruct); err != nil {
		slog.Error("incorrect input to post menu", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
	} else if err = checkMenuStruct(menuStruct, true); err != nil {
		slog.Error("Invalid Post menu struct: ", "error", err)
		writeHttp(w, http.StatusBadRequest, "invalid struct of input json", err.Error())
	} else if ings, err := h.menuServInt.PostServiceMenu(&menuStruct); err != nil {
		slog.Error("Post menu", "error", err)
		if err == models.ErrConflict {
			writeHttp(w, http.StatusConflict, "menu", err.Error())
		} else if err == models.ErrNotFoundIngs {
			writeHttp(w, http.StatusNotFound, "inventory ", strings.Join(ings, ", ")+err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "error post menu", err.Error())
		}
	} else {
		slog.Info("menu created: " + menuStruct.ID)
		writeHttp(w, http.StatusCreated, "success", "menu created: "+menuStruct.ID)
	}
}

func (h *menuHaldToService) GetAllMenus(w http.ResponseWriter, r *http.Request) {
	if menus, err := h.menuServInt.GetServiceMenus(); err != nil {
		slog.Error("Error getting all menus", "error", err)
		writeHttp(w, http.StatusInternalServerError, "get all", err.Error())
	} else if err = bodyJsonStruct(w, menus); err != nil {
		slog.Error("Get menus: cannot give body all menus")
	} else {
		slog.Info("Get all menu list")
	}
}

func (h *menuHaldToService) GetMenuById(w http.ResponseWriter, r *http.Request) {
	if idname := r.PathValue("id"); checkName(idname) {
		slog.Error("Get Menu: invalid id")
		writeHttp(w, http.StatusBadRequest, "ID", "Invalid id")
	} else if menu, err := h.menuServInt.GetServiceMenuById(idname); err != nil {
		slog.Error("Get Menu: cannot get menu struct", "error", err)
		if err == models.ErrNotFound {
			writeHttp(w, http.StatusNotFound, "menu", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "get menu by id", err.Error())
		}
	} else if err = bodyJsonStruct(w, menu); err != nil {
		slog.Error("Get menu: cannot write struct to the body")
	}
}

func (h *menuHaldToService) PutMenuById(w http.ResponseWriter, r *http.Request) {
	var menuStruct models.MenuItem
	if id := r.PathValue("id"); checkName(id) {
		slog.Error("Put Menu by id", "Invalid id ", id)
		writeHttp(w, http.StatusBadRequest, "id", "invalid")
	} else if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("Put the menu: content_Type must be application/json")
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
	} else if err := json.NewDecoder(r.Body).Decode(&menuStruct); err != nil {
		slog.Error("incorrect input to put menu", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
	} else if err = checkMenuStruct(menuStruct, false); err != nil {
		slog.Error("Invalid Post menu struct: ", "error", err)
		writeHttp(w, http.StatusBadRequest, "menu struct", err.Error())
	} else if ings, err := h.menuServInt.PutServiceMenuById(&menuStruct, id); err != nil {
		slog.Error("Put menu by id", "error", err)
		if err == models.ErrNotFound {
			writeHttp(w, http.StatusNotFound, "menu", err.Error())
		} else if err == models.ErrNotFoundIngs {
			writeHttp(w, http.StatusNotFound, "inventory for menu", strings.Join(ings, ", ")+err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "error post menu", err.Error())
		}
	} else {
		slog.Info("Menu: ", "Updated Menu by id: ", id)
		writeHttp(w, http.StatusOK, "Updated Menu by id: ", id)
	}
}

func (h *menuHaldToService) DelMenuById(w http.ResponseWriter, r *http.Request) {
	if idname := r.PathValue("id"); checkName(idname) {
		slog.Error("Del Menu: invalid id")
		writeHttp(w, http.StatusBadRequest, "ID", "Invalid id")
	} else if err := h.menuServInt.DelServiceMenuById(idname); err != nil {
		if err == models.ErrNotFound {
			slog.Error("Delete menu by id:", idname, err)
			writeHttp(w, http.StatusNotFound, "menu", err.Error())
		} else {
			slog.Error("Delete menu by id", "unknown error", err)
			writeHttp(w, http.StatusInternalServerError, "delete menu", err.Error())
		}
	} else {
		slog.Info("Deleted: ", " menu by id :", idname)
		writeHttp(w, http.StatusNoContent, "", "")
	}
}

func checkMenuStruct(menu models.MenuItem, isPost bool) error {
	if isPost && checkName(menu.ID) {
		return errors.New("invalid ID")
	} else if checkName(menu.Name) {
		return errors.New("invalid name of menu item")
	} else if checkName(menu.Description) {
		return errors.New("invalid description")
	} else if menu.Price < 0 {
		return errors.New("negative price")
	} else if len(menu.Ingredients) == 0 {
		return errors.New("none ingredients")
	}
	for _, v := range menu.Ingredients {
		if checkName(v.IngredientID) {
			return errors.New("invalid name of ingredient_id: " + v.IngredientID)
		} else if v.Quantity <= 0 {
			return errors.New("invalid quantity to: " + v.IngredientID)
		}
	}
	return nil
}
