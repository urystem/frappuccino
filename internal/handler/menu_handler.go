package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"frappuccino/internal/service"
	"frappuccino/models"
)

type menuHaldToService struct {
	menuServInt service.MenuServiceInter
}

func ReturnMenuHaldStruct(menuSerInt service.MenuServiceInter) *menuHaldToService {
	return &menuHaldToService{menuServInt: menuSerInt}
}

func (handMenu *menuHaldToService) GetMenus(w http.ResponseWriter, r *http.Request) {
	menus, err := handMenu.menuServInt.CollectMenus()
	if err != nil {
		slog.Error("Error getting all menus", "error", err)
		writeHttp(w, http.StatusInternalServerError, "get all", err.Error())
		return
	}

	err = bodyJsonStruct(w, menus, http.StatusOK)
	if err != nil {
		slog.Error("Get menus: cannot give body all menus")
		return
	}

	slog.Info("Get all menu list")
}

func (handMenu *menuHaldToService) GetMenuByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Get Menu: invalid id")
		writeHttp(w, http.StatusBadRequest, "ID", "Invalid id")
		return
	}

	menu, err := handMenu.menuServInt.TakeMenu(id)
	if err != nil {
		slog.Error("Get Menu: cannot get menu struct", "error", err)
		if err == models.ErrNotFound {
			writeHttp(w, http.StatusNotFound, "menu", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "get menu by id", err.Error())
		}
		return
	}

	if err = bodyJsonStruct(w, menu, http.StatusOK); err != nil {
		slog.Error("Get menu: cannot write struct to the body")
	}
}

func (handMenu *menuHaldToService) DelMenu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Del Menu:", "invalid id", err)
		writeHttp(w, http.StatusBadRequest, "ID", "Invalid id")
		return
	}

	menuDepends, err := handMenu.menuServInt.DelServiceMenuById(id)
	if err != nil {
		if err == models.ErrNotFound {
			slog.Error("Delete menu :", "by id", err)
			writeHttp(w, http.StatusNotFound, "menu", err.Error())
		} else {
			slog.Error("Delete menu by id", "unknown error", err)
			writeHttp(w, http.StatusInternalServerError, "delete menu", err.Error())
		}
		return
	}
	if menuDepends != nil {
		slog.Error("DELETE menu: found depend orders")
		err = bodyJsonStruct(w, menuDepends, http.StatusBadRequest)
		if err != nil {
			slog.Error("DELETE Invent: Error in decoder")
		}
		return
	}

	slog.Info("Deleted: ", " menu by id :", id)
	writeHttp(w, http.StatusNoContent, "", "")
}

func (h *menuHaldToService) PostMenu(w http.ResponseWriter, r *http.Request) {
	var menuStruct models.MenuItem
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("post the menu: content_Type must be application/json")
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&menuStruct); err != nil {
		slog.Error("incorrect input to post menu", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}

	if ings, err := h.menuServInt.CreateMenu(&menuStruct); err != nil {
		slog.Error("Post menu", "error", err)
		if err == models.InvalidIngs {
			err = bodyJsonStruct(w, ings, http.StatusBadRequest)
			if err != nil {
				slog.Error("DELETE Invent: Error in decoder")
			}
			return
		}
		writeHttp(w, http.StatusInternalServerError, "error post menu", err.Error())
		return
	}
	slog.Info("menu created: ", "success", menuStruct.ID)
	writeHttp(w, http.StatusCreated, "success", "menu created:")
}

// func (h *menuHaldToService) PutMenuById(w http.ResponseWriter, r *http.Request) {
// 	var menuStruct models.MenuItem
// 	if id := r.PathValue("id"); checkName(id) {
// 		slog.Error("Put Menu by id", "Invalid id ", id)
// 		writeHttp(w, http.StatusBadRequest, "id", "invalid")
// 	} else if r.Header.Get("Content-Type") != "application/json" {
// 		slog.Error("Put the menu: content_Type must be application/json")
// 		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
// 	} else if err := json.NewDecoder(r.Body).Decode(&menuStruct); err != nil {
// 		slog.Error("incorrect input to put menu", "error", err)
// 		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
// 	} else if err = checkMenuStruct(menuStruct, false); err != nil {
// 		slog.Error("Invalid Post menu struct: ", "error", err)
// 		writeHttp(w, http.StatusBadRequest, "menu struct", err.Error())
// 	} else if ings, err := h.menuServInt.PutServiceMenuById(&menuStruct, id); err != nil {
// 		slog.Error("Put menu by id", "error", err)
// 		if err == models.ErrNotFound {
// 			writeHttp(w, http.StatusNotFound, "menu", err.Error())
// 		} else if err == models.ErrNotFoundIngs {
// 			writeHttp(w, http.StatusNotFound, "inventory for menu", strings.Join(ings, ", ")+err.Error())
// 		} else {
// 			writeHttp(w, http.StatusInternalServerError, "error post menu", err.Error())
// 		}
// 	} else {
// 		slog.Info("Menu: ", "Updated Menu by id: ", id)
// 		writeHttp(w, http.StatusOK, "Updated Menu by id: ", id)
// 	}
// }

// func checkMenuStruct(menu models.MenuItem, isPost bool) error {
// 	// if isPost && checkName(menu.ID) {
// 	// 	return errors.New("invalid ID")
// 	// } else if checkName(menu.Name) {
// 	// 	return errors.New("invalid name of menu item")
// 	// } else if checkName(menu.Description) {
// 	// 	return errors.New("invalid description")
// 	// } else if menu.Price < 0 {
// 	// 	return errors.New("negative price")
// 	// } else if len(menu.Ingredients) == 0 {
// 	// 	return errors.New("none ingredients")
// 	// }
// 	// for _, v := range menu.Ingredients {
// 	// 	if checkName(v.IngredientID) {
// 	// 		return errors.New("invalid name of ingredient_id: " + v.IngredientID)
// 	// 	} else if v.Quantity <= 0 {
// 	// 		return errors.New("invalid quantity to: " + v.IngredientID)
// 	// 	}
// 	// }
// 	return nil
// }
