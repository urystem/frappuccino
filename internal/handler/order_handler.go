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

type ordHandToService struct {
	orderService service.OrdServiceInter
}

type ordHandInt interface {
	GetOrders(w http.ResponseWriter, r *http.Request)
	GetOrderByID(w http.ResponseWriter, r *http.Request)
	DelOrderByID(w http.ResponseWriter, r *http.Request)
	PostOrder(w http.ResponseWriter, r *http.Request)
	PutOrderByID(w http.ResponseWriter, r *http.Request)
	PostOrdCloseById(w http.ResponseWriter, r *http.Request)
	BatchProcess(w http.ResponseWriter, r *http.Request)
}

func ReturnOrdHaldStruct(ordSerInt service.OrdServiceInter) ordHandInt {
	return &ordHandToService{orderService: ordSerInt}
}

func (h *ordHandToService) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.CollectOrders()
	if err != nil {
		slog.Error("Get orders", "error", err)
		writeHttp(w, http.StatusInternalServerError, "get orders: ", err.Error())
		return
	}
	bodyJsonStruct(w, orders, http.StatusOK)
	slog.Info("Get orders success")
}

func (h *ordHandToService) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Error("Invalid id for get order")
		writeHttp(w, http.StatusBadRequest, "Invalid id", "Check the order id")
		return
	}

	order, err := h.orderService.TakeOrder(id)
	if err != nil {
		slog.Error("Can't get order struct: ", "error", err)
		if errors.Is(err, models.ErrNotFound) {
			writeHttp(w, http.StatusNotFound, "get order:", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "get order:", err.Error())
		}
		return
	}
	bodyJsonStruct(w, order, http.StatusOK)
	slog.Error("Get order: cannot write struct to the body")
}

func (h *ordHandToService) DelOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Warn("Invalid id for get order")
		writeHttp(w, http.StatusBadRequest, "Invalid id", "Check the order id")
		return
	}

	err = h.orderService.RemoveOrder(id)
	if err != nil {
		slog.Error("Delete order: ", "error id:", err)
		if err == models.ErrNotFound {
			writeHttp(w, http.StatusNotFound, "order", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "order", err.Error())
		}
	} else {
		slog.Info("Order ", "deleted:", id)
		writeHttp(w, http.StatusNoContent, "", "")
	}
}

func (h *ordHandToService) PostOrder(w http.ResponseWriter, r *http.Request) {
	var orderStruct models.Order
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("post the menu: content_Type must be application/json")
		writeHttp(w, http.StatusUnsupportedMediaType, "content/type", "not json")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&orderStruct)
	if err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}

	err = h.orderService.CreateOrder(&orderStruct)
	if err == nil {
		slog.Info("order created by : " + orderStruct.CustomerName)
		writeHttp(w, http.StatusCreated, "succes", "order created by : "+orderStruct.CustomerName)
		return
	}

	slog.Error("Failed to post order", "error", err)

	if errors.Is(err, models.ErrBadInput) {
		writeHttp(w, http.StatusBadRequest, "Error post order", err.Error())
		return
	}

	if errors.Is(err, models.ErrBadInputItems) {
		bodyJsonStruct(w, orderStruct.Items, http.StatusBadRequest)
		return
	}

	if errors.Is(err, models.ErrOrderNotEnoughItems) {
		bodyJsonStruct(w, orderStruct.Items, http.StatusFailedDependency)
		return
	}

	if errors.Is(err, models.ErrNotFoundItems) {
		bodyJsonStruct(w, orderStruct.Items, http.StatusNotFound)
		return
	}

	writeHttp(w, http.StatusInternalServerError, "Error post order", err.Error())
}

func (h *ordHandToService) PutOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Warn("Invalid id for get order")
		writeHttp(w, http.StatusBadRequest, "Invalid id", "Check the order id")
		return
	}
	var orderStruct models.Order
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("post the menu: content_Type must be application/json")
		writeHttp(w, http.StatusUnsupportedMediaType, "content/type", "not json")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&orderStruct)
	if err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}
	err = h.orderService.UpgradeOrder(id, &orderStruct)
	if err == nil {
		writeHttp(w, http.StatusInternalServerError, "Put order", "succes")
		slog.Info("order updated: ", "success", id)
		return
	}

	slog.Error("Failed to put order", "error", err)

	if errors.Is(err, models.ErrBadInput) {
		writeHttp(w, http.StatusBadRequest, "Error put order", err.Error())
		return
	}

	if errors.Is(err, models.ErrBadInputItems) {
		bodyJsonStruct(w, orderStruct.Items, http.StatusBadRequest)
		return
	}

	if errors.Is(err, models.ErrNotFoundItems) {
		bodyJsonStruct(w, orderStruct.Items, http.StatusNotFound)
		return
	}

	if errors.Is(err, models.ErrOrderNotEnoughItems) {
		bodyJsonStruct(w, orderStruct.Items, http.StatusFailedDependency)
		return
	}

	writeHttp(w, http.StatusInternalServerError, "Error put order", err.Error())
}

func (h *ordHandToService) PostOrdCloseById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		slog.Warn("Invalid id for get order")
		writeHttp(w, http.StatusBadRequest, "Invalid id", "Check the order id")
		return
	}

	err = h.orderService.ShutOrder(id)
	if err != nil {
		slog.Error("Close order", "error id:", id)
		if errors.Is(err, models.ErrNotFound) {
			writeHttp(w, http.StatusNotFound, "order", err.Error())
		} else if errors.Is(err, models.ErrOrderStatusClosed) {
			writeHttp(w, http.StatusBadRequest, "order already", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "close order", err.Error())
		}
	} else {
		slog.Info("order closed", "id: ", id)
		writeHttp(w, http.StatusOK, "order", "closed")
	}
}

func (h *ordHandToService) BatchProcess(w http.ResponseWriter, r *http.Request) {
	var batch models.PostSomeOrders
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("batch: content_Type must be application/json")
		writeHttp(w, http.StatusUnsupportedMediaType, "content/type", "not json")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&batch)
	if err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}
	response, err := h.orderService.CreateSomeOrders(&batch)
	if err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "unknown erreke", err.Error())
	}
	slog.Info("Batch succes")
	bodyJsonStruct(w, response, http.StatusOK)
}
