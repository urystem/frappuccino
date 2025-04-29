package handler

import (
	"encoding/json"
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
	err = bodyJsonStruct(w, orders, http.StatusOK)
	if err != nil {
		slog.Error("Can't give all orders to body", "error", err)
		return
	}
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
		writeHttp(w, http.StatusInternalServerError, "get order:", err.Error())
		return
	}

	if err = bodyJsonStruct(w, order, http.StatusOK); err != nil {
		slog.Error("Get order: cannot write struct to the body")
	}
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
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&orderStruct)
	if err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}
	// if err := checkOrdStruct(&orderStruct); err != nil {
	// 	slog.Error("invalid order struct in body")
	// 	writeHttp(w, http.StatusBadRequest, "invalid struct", err.Error())
	// 	return
	// }
	err = h.orderService.CreateOrder(&orderStruct)
	if err != nil {
		slog.Error("Failed to post order", "error", err)
		if err == models.ErrInputOrder || err == models.ErrOrderItems {
			err = bodyJsonStruct(w, orderStruct.Items, http.StatusBadRequest)
			if err != nil {
				slog.Error("Post order: Error in decoder")
			}
		} else {
			writeHttp(w, http.StatusInternalServerError, "Error post order", err.Error())
		}
	} else {
		slog.Info("order created by : " + orderStruct.CustomerName)
		writeHttp(w, http.StatusCreated, "succes", "order created by : "+orderStruct.CustomerName)
	}
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
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&orderStruct)
	if err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}
	err = h.orderService.UpgradeOrder(id, &orderStruct)
	if err != nil {
		slog.Error("Failed to put order", "error", err)
		if err == models.ErrInputOrder || err == models.ErrOrderItems {
			err = bodyJsonStruct(w, orderStruct.Items, http.StatusBadRequest)
			if err != nil {
				slog.Error("put order:", " Error in decoder", err)
			}
		} else {
			writeHttp(w, http.StatusInternalServerError, "Error put order", err.Error())
		}
	} else {
		writeHttp(w, http.StatusInternalServerError, "Put order", "succes")
		slog.Info("order updated: ", "success", id)
	}
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
		if err == models.ErrNotFound {
			writeHttp(w, http.StatusNotFound, "order", err.Error())
		} else if err == models.ErrOrdStatusClosed {
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
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&batch)
	if err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}
	h.orderService.CreateSomeOrders(&batch)
}
