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

type menuHandToService struct {
	menuServInt service.MenuServiceInter
}

type menuHandInt interface {
	GetMenus(w http.ResponseWriter, r *http.Request)
	GetMenuByID(w http.ResponseWriter, r *http.Request)
	DelMenu(w http.ResponseWriter, r *http.Request)
	PostMenu(w http.ResponseWriter, r *http.Request)
	PutMenuByID(w http.ResponseWriter, r *http.Request)
}

func ReturnMenuHaldStruct(menuSerInt service.MenuServiceInter) menuHandInt {
	return &menuHandToService{menuServInt: menuSerInt}
}

func (handMenu *menuHandToService) GetMenus(w http.ResponseWriter, r *http.Request) {
	menus, err := handMenu.menuServInt.CollectMenus()
	if err != nil {
		slog.Error("Error getting all menus", "error", err)
		writeHttp(w, http.StatusInternalServerError, "get all", err.Error())
		return
	}

	bodyJsonStruct(w, menus, http.StatusOK)
	slog.Info("Get all menu list")
}

func (handMenu *menuHandToService) GetMenuByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Get Menu: invalid id")
		writeHttp(w, http.StatusBadRequest, "ID", "Invalid id")
		return
	}

	menu, err := handMenu.menuServInt.TakeMenu(id)
	if err != nil {
		slog.Error("Get Menu: cannot get menu struct", "error", err)
		if errors.Is(err, models.ErrNotFound) {
			writeHttp(w, http.StatusNotFound, "menu", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "get menu by id", err.Error())
		}
		return
	}

	bodyJsonStruct(w, menu, http.StatusOK)
}

func (handMenu *menuHandToService) DelMenu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Del Menu:", "invalid id", err)
		writeHttp(w, http.StatusBadRequest, "ID", "Invalid id")
		return
	}

	menuDepends, err := handMenu.menuServInt.DelServiceMenuById(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
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
		bodyJsonStruct(w, menuDepends, http.StatusFailedDependency)
		return
	}

	slog.Info("Deleted: ", " menu by id :", id)
	writeHttp(w, http.StatusNoContent, "", "")
}

func (handMenu *menuHandToService) PostMenu(w http.ResponseWriter, r *http.Request) {
	var menuStruct models.MenuItem
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("post the menu: content_Type must be application/json")
		writeHttp(w, http.StatusUnsupportedMediaType, "content/type", "not json")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&menuStruct)
	if err != nil {
		slog.Error("incorrect input to post menu", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}

	err = handMenu.menuServInt.CreateMenu(&menuStruct)
	if err == nil {
		slog.Info("menu created: ", "success", menuStruct.ID)
		writeHttp(w, http.StatusCreated, "success", "menu created:")
		return
	}

	slog.Error("Post menu", "error", err)

	//ErrBadInputItems дегеннің ішінде ErrBadInput бар
	//
	if errors.Is(err, models.ErrBadInputItems) {
		bodyJsonStruct(w, menuStruct.Ingredients, http.StatusUnprocessableEntity)
		return
	}
	if errors.Is(err, models.ErrBadInput) {
		writeHttp(w, http.StatusUnprocessableEntity, "failed", "menu input:")
		return
	}

	if errors.Is(err, models.ErrNotFoundItems) {
		bodyJsonStruct(w, menuStruct.Ingredients, http.StatusNotFound)
		return
	}
	writeHttp(w, http.StatusInternalServerError, "failed", "post menu:")
}

func (handMenu *menuHandToService) PutMenuByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Put Menu by id", "Invalid id ", id)
		writeHttp(w, http.StatusBadRequest, "id", "invalid")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("Put the menu: content_Type must be application/json")
		writeHttp(w, http.StatusUnsupportedMediaType, "content/type", "not json")
		return
	}

	var menuStruct models.MenuItem

	if err := json.NewDecoder(r.Body).Decode(&menuStruct); err != nil {
		slog.Error("incorrect input to put menu", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}

	menuStruct.ID = id
	err = handMenu.menuServInt.UpgradeMenu(&menuStruct)
	if err == nil {
		slog.Info("Menu: ", "Updated Menu by id: ", id)
		writeHttp(w, http.StatusOK, "Updated Menu by id: ", "")
		return
	}

	slog.Error("Put menu by id", "error", err)

	if errors.Is(err, models.ErrBadInputItems) {
		bodyJsonStruct(w, menuStruct.Ingredients, http.StatusBadRequest)
		return
	}

	if errors.Is(err, models.ErrBadInput) {
		writeHttp(w, http.StatusBadRequest, "error put menu", err.Error())
		return
	}

	if errors.Is(err, models.ErrNotFound) {
		writeHttp(w, http.StatusNotFound, "error put menu", err.Error())
		return
	}

	if errors.Is(err, models.ErrNotFoundItems) {
		bodyJsonStruct(w, menuStruct.Ingredients, http.StatusNotFound)
		return
	}
	writeHttp(w, http.StatusInternalServerError, "error put menu", err.Error())
}
