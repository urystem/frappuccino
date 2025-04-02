package service

import (
	"frappuccino/internal/dal"
	"frappuccino/models"
)

type ordServiceToDal struct {
	ordDalInt    dal.OrderDalInter
	menuDalInt   dal.MenuDalInter
	inventDalInt dal.InventoryDataAccess
}

type OrdServiceInter interface {
	GetServiceOrders() ([]models.Order, error)
	PostServiceOrder(*models.Order) ([]string, error)
	GetServiceOrdById(string) (*models.Order, error)
	PutServiceOrdById(*models.Order, string) ([]string, error)
	DelServiceOrdById(string) error
	PostServiseOrdCloseById(string) error
	GetServiseTotalSales() (float64, error)
	GetServicePopularItem() ([]models.OrderItem, error)
}

func ReturnOrdSerStruct(ordInter dal.OrderDalInter, menuInt dal.MenuDalInter, invIntDal dal.InventoryDataAccess) *ordServiceToDal {
	return &ordServiceToDal{ordDalInt: ordInter, menuDalInt: menuInt, inventDalInt: invIntDal}
}

func (ser *ordServiceToDal) GetServiceOrders() ([]models.Order, error) {
	return ser.ordDalInt.ReadOrdersDal()
}

// func (ser *ordServiceToDal) PostServiceOrder(ord *models.Order) ([]string, error) {
// 	if notFounds, err := ser.checkItemWithMenu(ord.Items); err != nil {
// 		return nil, err
// 	} else if len(notFounds) != 0 {
// 		return notFounds, models.ErrOrdNotFoundItem
// 	} else if notEnough, err := ser.updateInventory(ord.Items, true); err != nil {
// 		return notEnough, err
// 	}
// 	orders, err := ser.ordDalInt.ReadOrdersDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(orders) == 0 {
// 		ord.ID = "order1"
// 	} else if last := orders[len(orders)-1].ID; len(last) > 5 {
// 		if n, err := strconv.Atoi(last[5:]); err != nil {
// 			return nil, err
// 		} else {
// 			ord.ID = "order" + strconv.Itoa(n+1)
// 		}
// 	} else {
// 		return nil, errors.New("metadata order id error")
// 	}
// 	ord.Status, ord.CreatedAt = "open", time.Now().Format(time.RFC3339)
// 	return nil, ser.ordDalInt.WriteOrderDal(append(orders, *ord))
// }

// func (ser *ordServiceToDal) GetServiceOrdById(id string) (*models.Order, error) {
// 	orders, err := ser.ordDalInt.ReadOrdersDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, v := range orders {
// 		if v.ID == id {
// 			return &v, nil
// 		}
// 	}
// 	return nil, models.ErrNotFound
// }

// func (ser *ordServiceToDal) PutServiceOrdById(ord *models.Order, id string) ([]string, error) {
// 	orders, err := ser.ordDalInt.ReadOrdersDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i, v := range orders {
// 		if v.ID == id {
// 			if v.Status != "open" {
// 				return nil, models.ErrOrdStatusClosed
// 			} else if notFounds, err := ser.checkItemWithMenu(ord.Items); err != nil {
// 				return nil, err
// 			} else if len(notFounds) != 0 {
// 				return notFounds, models.ErrOrdNotFoundItem
// 			} else if _, err := ser.updateInventory(v.Items, false); err != nil { // reset invents with old orders items
// 				return nil, err
// 			} else if notEnough, err := ser.updateInventory(ord.Items, true); err != nil { // then update invents with new order items
// 				if err == models.ErrOrdNotEnough { // if not enough
// 					if _, er := ser.updateInventory(v.Items, true); er != nil { // reset stay with old order items
// 						return nil, errors.New("metadata changed fatal error") // if metadata changed when checking new updated order
// 					}
// 					return notEnough, err
// 				}
// 				return nil, err
// 			}
// 			orders[i].CustomerName, orders[i].Items = ord.CustomerName, ord.Items
// 			return nil, ser.ordDalInt.WriteOrderDal(orders)
// 		}
// 	}
// 	return nil, models.ErrNotFound
// }

// func (ser *ordServiceToDal) DelServiceOrdById(id string) error {
// 	ords, err := ser.ordDalInt.ReadOrdersDal()
// 	if err != nil {
// 		return err
// 	}
// 	for i, v := range ords {
// 		if v.ID == id {
// 			if v.Status == "open" {
// 				if _, err = ser.updateInventory(v.Items, false); err != nil {
// 					return err
// 				}
// 			}
// 			return ser.ordDalInt.WriteOrderDal(append(ords[:i], ords[i+1:]...))
// 		}
// 	}
// 	return models.ErrNotFound
// }

// func (ser *ordServiceToDal) PostServiseOrdCloseById(id string) error {
// 	orders, err := ser.ordDalInt.ReadOrdersDal()
// 	if err != nil {
// 		return err
// 	}
// 	for i, v := range orders {
// 		if v.ID == id {
// 			if orders[i].Status != "open" {
// 				return models.ErrOrdStatusClosed
// 			}
// 			orders[i].Status = "closed"
// 			return ser.ordDalInt.WriteOrderDal(orders)
// 		}
// 	}
// 	return models.ErrNotFound
// }

// func (ser *ordServiceToDal) GetServiseTotalSales() (float64, error) {
// 	totalMenus, err := ser.returnMapSelled(false)
// 	if err != nil {
// 		return 0, err
// 	}
// 	menus, err := ser.menuDalInt.ReadMenuDal()
// 	if err != nil {
// 		return 0, err
// 	}
// 	var ansTotal float64
// 	for _, menu := range menus {
// 		if quanti, exists := totalMenus[menu.ID]; exists {
// 			ansTotal += float64(quanti) * menu.Price
// 		}
// 	}
// 	return ansTotal, nil
// }

// func (ser *ordServiceToDal) GetServicePopularItem() ([]models.OrderItem, error) {
// 	totalMenus, err := ser.returnMapSelled(true)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var ordItems []models.OrderItem
// 	for itemId, quantity := range totalMenus {
// 		ordItems = append(ordItems, models.OrderItem{ProductID: itemId, Quantity: quantity})
// 	}
// 	for i := 0; i < len(ordItems)-1; i++ {
// 		for j := i + 1; j < len(ordItems); j++ {
// 			if ordItems[i].Quantity < ordItems[j].Quantity {
// 				ordItems[i], ordItems[j] = ordItems[j], ordItems[i]
// 			}
// 		}
// 	}
// 	return ordItems, nil
// }

// func (ser *ordServiceToDal) returnMapSelled(isPopular bool) (map[string]int, error) {
// 	orders, err := ser.ordDalInt.ReadOrdersDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	totalMenus := make(map[string]int)
// 	for _, order := range orders {
// 		if isPopular || order.Status != "open" {
// 			for _, item := range order.Items {
// 				totalMenus[item.ProductID] += item.Quantity
// 			}
// 		}
// 	}
// 	return totalMenus, nil
// }

// func (ser *ordServiceToDal) checkItemWithMenu(items []models.OrderItem) ([]string, error) {
// 	menus, err := ser.menuDalInt.ReadMenuDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var notFounds []string
// 	for _, item := range items {
// 		var isHere bool
// 		for _, menu := range menus {
// 			if item.ProductID == menu.ID {
// 				isHere = true
// 				break
// 			}
// 		}
// 		if !isHere {
// 			notFounds = append(notFounds, item.ProductID)
// 		}
// 	}
// 	return notFounds, nil
// }

// func (ser *ordServiceToDal) updateInventory(menuesOrd []models.OrderItem, minus bool) ([]string, error) {
// 	menus, err := ser.menuDalInt.ReadMenuDal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	invetory, err := ser.inventDalInt.ReadInventory()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var notEnough []string
// 	for _, v := range menuesOrd {
// 		invents := ser.returnMenu(menus, v.ProductID) // search in menu and return her invents
// 		if invents == nil {                           // this menu is not found in menues
// 			return nil, errors.New("after checking not found menu") // ---FATAL ERROR
// 		}
// 		for _, invent := range invents { // in needed invents
// 			invent.Quantity *= float64(v.Quantity)
// 			if err := ser.changeCountInventsTempBase(invetory, invent.IngredientID, invent.Quantity, minus); err != nil {
// 				notEnough = append(notEnough, err.Error())
// 			}
// 		}
// 	}
// 	if len(notEnough) != 0 {
// 		return notEnough, models.ErrOrdNotEnough
// 	}
// 	return nil, ser.inventDalInt.WriteInventory(invetory)
// }

// func (ser *ordServiceToDal) returnMenu(menus []models.MenuItem, menuId string) []models.MenuItemIngredient {
// 	for _, v := range menus {
// 		if v.ID == menuId {
// 			return v.Ingredients
// 		}
// 	}
// 	return nil // it is impossible, because it is checked before and it is must be in here
// }

// func (ser *ordServiceToDal) changeCountInventsTempBase(invertory []models.InventoryItem, nameInvent string, quanti float64, minus bool) error {
// 	for i, v := range invertory {
// 		if v.IngredientID == nameInvent {
// 			if minus {
// 				invertory[i].Quantity -= quanti // алып тастаймыз минусқа кетсе де, үйткені кейбірі минусқа кетпеуі мүмкін, сосын тек минусқа кеткендерін терңп жүрмеу үшін - бәрі біріңғай алынған дұрыс. Кейін бәрін біріңғай қайтара салуға оңай болады
// 				if invertory[i].Quantity < 0 {
// 					return errors.New(v.IngredientID + ": Required " + strconv.FormatFloat(quanti, 'f', 2, 64) + " Available: " + strconv.FormatFloat(invertory[i].Quantity+quanti, 'f', 2, 64))
// 				}
// 			} else {
// 				invertory[i].Quantity += quanti
// 			}
// 		}
// 	}
// 	return nil
// }
