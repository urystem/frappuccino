package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type ordHandToService struct {
	orderService service.OrdServiceInter
}

func ReturnOrdHaldStruct(ordSerInt service.OrdServiceInter) *ordHandToService {
	return &ordHandToService{orderService: ordSerInt}
}

func (h *ordHandToService) PostOrder(w http.ResponseWriter, r *http.Request) {
	var orderStruct models.Order
	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("post the menu: content_Type must be application/json")
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&orderStruct); err != nil {
		slog.Error("incorrect input to post order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
		return
	}
	if err := checkOrdStruct(&orderStruct); err != nil {
		slog.Error("invalid order struct in body")
		writeHttp(w, http.StatusBadRequest, "invalid struct", err.Error())
		return
	}
	if notFoundsMenusOrInvents, err := h.orderService.PostServiceOrder(&orderStruct); err != nil {
		slog.Error("Failed to post order", "error", err)
		if err == models.ErrOrdNotFoundItem {
			writeHttp(w, http.StatusNotFound, "item", strings.Join(notFoundsMenusOrInvents, ", ")+err.Error())
		} else if err == models.ErrOrdNotEnough {
			writeHttp(w, http.StatusNotFound, "invents", strings.Join(notFoundsMenusOrInvents, ", "+err.Error()))
		} else {
			writeHttp(w, http.StatusInternalServerError, "Error post order", err.Error())
		}
	} else {
		slog.Info("order created by : " + orderStruct.CustomerName)
		writeHttp(w, http.StatusCreated, "succes", "order created by : "+orderStruct.CustomerName)
	}
}

func (h *ordHandToService) GetOrders(w http.ResponseWriter, r *http.Request) {
	if orders, err := h.orderService.GetServiceOrders(); err != nil {
		slog.Error("Get orders", "error", err)
		writeHttp(w, http.StatusInternalServerError, "get orders: ", err.Error())
	} else if err = bodyJsonStruct(w, orders); err != nil {
		slog.Error("Can't give all orders to body", "error", err)
	}
}

func (h *ordHandToService) GetOrdById(w http.ResponseWriter, r *http.Request) {
	if id := r.PathValue("id"); checkOdrId(id) {
		slog.Warn("Invalid id for get order")
		writeHttp(w, http.StatusBadRequest, "Invalid id", "Check the order id")
	} else if order, err := h.orderService.GetServiceOrdById(id); err != nil {
		slog.Error("Can't get order struct: ", "error", err)
		if err == models.ErrNotFound {
			writeHttp(w, http.StatusNotFound, "order", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "get order failed", err.Error())
		}
	} else if err = bodyJsonStruct(w, order); err != nil {
		slog.Error("Can't give order struct to body ", "error", err)
	}
}

func (h *ordHandToService) PutOrdById(w http.ResponseWriter, r *http.Request) {
	var orderStruct models.Order
	if id := r.PathValue("id"); checkOdrId(id) {
		slog.Warn("Invalid id for put a order")
		writeHttp(w, http.StatusBadRequest, "put order", "invalid id")
	} else if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("post the menu: content_Type must be application/json")
		writeHttp(w, http.StatusBadRequest, "content/type", "not json")
	} else if err := json.NewDecoder(r.Body).Decode(&orderStruct); err != nil {
		slog.Error("incorrect input to put a order", "error", err)
		writeHttp(w, http.StatusBadRequest, "input json", err.Error())
	} else if err = checkOrdStruct(&orderStruct); err != nil {
		slog.Error("invalid order struct in body")
		writeHttp(w, http.StatusBadRequest, "invalid struct", err.Error())
	} else if notFoundItemsOrInvents, err := h.orderService.PutServiceOrdById(&orderStruct, id); err != nil {
		slog.Error("incorrect input to put a order", "error", err)
		if err == models.ErrNotFound {
			writeHttp(w, http.StatusNotFound, "order", err.Error())
		} else if err == models.ErrOrdNotFoundItem {
			writeHttp(w, http.StatusNotFound, "order item(s)", strings.Join(notFoundItemsOrInvents, ", ")+err.Error())
		} else if err == models.ErrOrdNotEnough {
			writeHttp(w, http.StatusNotFound, "invents", strings.Join(notFoundItemsOrInvents, ", ")+err.Error())
		} else if err == models.ErrOrdStatusClosed {
			writeHttp(w, http.StatusBadRequest, "order", err.Error())
		} else {
			writeHttp(w, http.StatusInternalServerError, "update order", err.Error())
		}
	} else {
		slog.Info("order puted by : " + orderStruct.CustomerName)
		writeHttp(w, http.StatusCreated, "order puted by ", orderStruct.CustomerName)
	}
}

func (h *ordHandToService) DelOrdById(w http.ResponseWriter, r *http.Request) {
	if id := r.PathValue("id"); checkOdrId(id) {
		slog.Warn("Invalid id for get order")
		writeHttp(w, http.StatusBadRequest, "Invalid id", "Check the order id")
	} else if err := h.orderService.DelServiceOrdById(id); err != nil {
		slog.Error("Delete order error id: " + id)
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

func (h *ordHandToService) PostOrdCloseById(w http.ResponseWriter, r *http.Request) {
	if id := r.PathValue("id"); checkOdrId(id) {
		slog.Warn("Invalid id for get order")
		writeHttp(w, http.StatusBadRequest, "Invalid id", "Check the order id")
	} else if err := h.orderService.PostServiseOrdCloseById(id); err != nil {
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

func (h *ordHandToService) TotalSales(w http.ResponseWriter, r *http.Request) {
	if total, err := h.orderService.GetServiseTotalSales(); err != nil {
		slog.Error("Get total sales", "error", err)
		writeHttp(w, http.StatusInternalServerError, "failed to get total sales:", err.Error())
	} else {
		slog.Info("Succes", "Get total sales:", total)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"total_sales": total})
	}
}

func (h *ordHandToService) PopularItem(w http.ResponseWriter, r *http.Request) {
	if sortedItems, err := h.orderService.GetServicePopularItem(); err != nil {
		slog.Error("Error", "get popular items list:", err)
		writeHttp(w, http.StatusInternalServerError, "get popular items", err.Error())
	} else if err = bodyJsonStruct(w, sortedItems); err != nil {
		slog.Error("Error write sorted items to body")
	} else {
		slog.Info("get popular items success")
	}
}

func checkOrdStruct(ord *models.Order) error {
	if checkName(ord.CustomerName) {
		return errors.New("invalid name")
	}

	if len(ord.Items) == 0 {
		return errors.New("empty items")
	}

	if ord.ID != "" || ord.Status != "" || ord.CreatedAt != "" {
		return errors.New("you cannot give to other fields")
	}

	for _, v := range ord.Items {
		if checkName(v.ProductID) {
			return errors.New("invalid item name: " + v.ProductID)
		} else if v.Quantity <= 0 {
			return errors.New("Invalid quantity of item: " + v.ProductID)
		}
	}

	return nil
}


INSERT INTO inventory (id) VALUES(11);