package village

import (
	"github.com/gin-gonic/gin"
	"github.com/timshannon/bolthold"
	"gopkg.in/mgo.v2/bson"
)

// Returns the service of the worker that should represent the information
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getOurServiceWorker(c *gin.Context, sTypeID bson.ObjectId, class string, method string, operation string) (Service, bool) {
	services := VS.getAllSpecificServicesFromLocalVillage(c, sTypeID, "h_database", "getOurServiceWorker", "services := VS.getAllSpecificServicesFromLocalVillage")
	user := VS.getUserFomCookie(c, "h_database", "getOurServiceWorker", "user := VS.getUserFomCookie")
	var serWithAccess []Service
	for _, ser := range services {
		if VS.hasAccess(c, user.ID, ser.ID, "h_database", "getOurServiceWorker", "user := VS.getUserFomCookie") {
			serWithAccess = append(serWithAccess, ser)
		}
	}
	if len(serWithAccess) == 0 {
		return Service{}, false
	}
	return serWithAccess[0], true
}

// This method Returns all Sales from the system
func (VS *Server) getAllSalesFromTheSystem(c *gin.Context, class string, method string, operation string) []Sale {
	var sales []Sale
	err := VS.DataBase.Find(&sales, nil)
	VS.checkOperation(c, class, method, operation, err)
	return sales
}

// This method Returns all Payments from the system
func (VS *Server) getAllPaymentsFromTheSystem(c *gin.Context, class string, method string, operation string) []Payment {
	var payments []Payment
	err := VS.DataBase.Find(&payments, nil)
	VS.checkOperation(c, class, method, operation, err)
	return payments
}

// This method Returns all Service Purchases from the system
func (VS *Server) getAllStockFromTheSystem(c *gin.Context, class string, method string, operation string) []Stock {
	var workshopsP []Stock
	err := VS.DataBase.Find(&workshopsP, nil)
	VS.checkOperation(c, class, method, operation, err)
	return workshopsP
}

// This method Returns all Services from the village given
func (VS *Server) getAllServicesFromVillage(c *gin.Context, idVillage bson.ObjectId, class string, method string, operation string) []Service {
	var workshops []Service
	err := VS.DataBase.Find(&workshops, bolthold.Where("VillageID").Eq(idVillage))
	VS.checkOperation(c, class, method, operation, err)
	return workshops
}

// This method finds all the Worker Orders related with a specific Worker and Service
func (VS *Server) getWorkerOrderByWorkerIDAndVideoID(c *gin.Context, workerID bson.ObjectId, videoID bson.ObjectId, class string, method string, operation string) []WorkerOrder {
	var orders []WorkerOrder
	err := VS.DataBase.Find(&orders, bolthold.Where("WorkerID").Eq(workerID).And("ProductID").Eq(videoID))
	VS.checkOperation(c, class, method, operation, err)
	// if len(orders)==0 {VS.sendError(c, class, method, operation, err)}
	return orders
}

// This method finds the ServiceOrder in database by his ID and returns it.
func (VS *Server) getServiceOrdersByVideoCourse(c *gin.Context, id bson.ObjectId, class string, method string, operation string) []ServiceOrder {
	var pOrder []ServiceOrder
	err := VS.DataBase.Find(&pOrder, bolthold.Where("VideoCourseOrItemID").Eq(id))
	VS.checkOperation(c, class, method, operation, err)
	return pOrder
}

// The first return is about if the service already exist, and the second is about if the service is active
func (VS *Server) existAccess(c *gin.Context, userID bson.ObjectId, serviceID bson.ObjectId) bool {
	var acc []Access
	err := VS.DataBase.Find(&acc, bolthold.Where("IDservice").Eq(serviceID).And("IDuser").Eq(userID).And("IsActive").Eq(true))
	VS.checkOperation(c, "Database", "existAccess", "find", err)
	if len(acc) == 0 {
		return false
	}
	return true
}

// Returns the specific Access from the Service and User provided
func (VS *Server) getSpecificAccessByServiceAndUser(c *gin.Context, userID bson.ObjectId, serviceID bson.ObjectId) Access {
	var acc []Access
	err := VS.DataBase.Find(&acc, bolthold.Where("IDservice").Eq(serviceID).And("IDuser").Eq(userID).And("IsActive").Eq(true))
	VS.checkOperation(c, "Database", "getSpecificAccessByServiceAndUser", "find", err)
	if len(acc) == 0 {
		return Access{}
	}
	return acc[0]
}

// This method returns all the Audits of the User into a Date range
func (VS *Server) getAuditsByUserAndDate(c *gin.Context, from int64, to int64, userID bson.ObjectId, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, bolthold.Where("IDuser").Eq(userID).And("Date").Gt(from).And("Date").Lt(to))
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns all the Audits of the Villages into a range of Dates
func (VS *Server) getAuditsByVillageAndDate(c *gin.Context, from int64, to int64, IDVillage bson.ObjectId, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, bolthold.Where("IDVillage").Eq(IDVillage).And("Date").Gt(from).And("Date").Lt(to))
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns all the Audits into a range of Dates
func (VS *Server) getAuditsByDates(c *gin.Context, from int64, to int64, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, bolthold.Where("Date").Gt(from).And("Date").Lt(to))
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns all the Audits of the Service into a range of Dates
func (VS *Server) getAuditsByServiceAndDate(c *gin.Context, from int64, to int64, IDservice bson.ObjectId, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, bolthold.Where("IDservice").Eq(IDservice).And("Date").Gt(from).And("Date").Lt(to))
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns all the Categories of our system that haveaccess to the serviceID
func (VS *Server) getAllCategoriesFromServiceID(c *gin.Context, serviceTypeID bson.ObjectId, class string, method string, operation string) []Category {
	allCategories := VS.getAllCategories(c, "functionality", "getAllItemsByService", "allCategories := VS.getAllCategories")
	// We clean and take only the categories with access to the Item
	var categoriesWithItem []Category
	for _, category := range allCategories {
		if ContainsID(category.ServicesAccess, serviceTypeID) {
			categoriesWithItem = append(categoriesWithItem, category)
		}
	}
	return allCategories
}

// This method returns all the Items from the Category given
func (VS *Server) getAllItemsWithStockByCategory(c *gin.Context, categoryID bson.ObjectId, serviceID bson.ObjectId, class string, method string, operation string) []Item {
	var items []Item
	err := VS.DataBase.Find(&items, bolthold.Where("IDCategory").Eq(categoryID))
	VS.checkOperation(c, "No Material found", method, operation, err)
	var itemsWithStock []Item
	for _, n := range items {
		_, b := VS.getStockByItemAndService(c, n.ID, serviceID, "database", "getAllItemsWithStockByCategory", "if VS.getStockByItemAndService")
		if b {
			itemsWithStock = append(itemsWithStock, n)
		}
	}
	return itemsWithStock
}

// This method returns all the Deliveries from the Service given
func (VS *Server) getAllDeliveriesByServiceType(c *gin.Context, serviceTypeID bson.ObjectId) []Delivery {
	var dPack []Delivery
	err := VS.DataBase.Find(&dPack, bolthold.Where("Type").Eq(serviceTypeID))
	VS.checkOperation(c, "", "", "ya todo", err)
	return dPack
}

// This method returns all the deliveries that are sent from the Service given and are marked as Sent
func (VS *Server) getAllDeliveriesFromServiceIDSent(c *gin.Context, serviceID bson.ObjectId) []Delivery {
	var dPack []Delivery
	err := VS.DataBase.Find(&dPack, bolthold.Where("IDserviceEmitter").Eq(serviceID).And("IsSent").Eq(true))
	VS.checkOperation(c, "", "", "ya todo", err)
	return dPack
}

// This method returns all the Deliveries from the Service given
func (VS *Server) getAllDeliveriesByServiceIDReceived(c *gin.Context, serviceID bson.ObjectId) []Delivery {
	var dPack []Delivery
	err := VS.DataBase.Find(&dPack, bolthold.Where("IDserviceReceiver").Eq(serviceID))
	VS.checkOperation(c, "", "", "ya todo", err)
	return dPack
}

// This method returns all the DeliveryPack from the Service given
func (VS *Server) getAllDeliveriesPackNOTCheckedByServiceID(c *gin.Context, serviceID bson.ObjectId) []Delivery {
	var dPack []Delivery
	err := VS.DataBase.Find(&dPack, bolthold.Where("IDservice").Eq(serviceID).And("IsComplete").Eq(false))
	VS.checkOperation(c, "", "", "ya todo", err)
	return dPack
}

// This method returns all the DeliveryPack from the Service given
func (VS *Server) getAllDeliveriesPackCheckedByServiceID(c *gin.Context, serviceID bson.ObjectId) []Delivery {
	var dPack []Delivery
	err := VS.DataBase.Find(&dPack, bolthold.Where("IDservice").Eq(serviceID).And("IsComplete").Eq(true))
	VS.checkOperation(c, "", "", "ya todo", err)
	return dPack
}

// This method returns all the DeliveryItem from the Delivery given
func (VS *Server) getAllDeliveryPack(c *gin.Context) []Delivery {
	var dPack []Delivery
	VS.DataBase.Find(&dPack, nil)
	return dPack
}

// This method finds the DeliveryPack in the database by his ID and returns it, also returns true if found and false if not
func (VS *Server) existDeliveryPackByID(c *gin.Context, id bson.ObjectId) bool {
	var dPack []Delivery
	VS.DataBase.Find(&dPack, bolthold.Where("ID").Eq(id))
	if len(dPack) == 0 {
		return false
	}
	return true
}

// This method finds and return the User in database by his ID
func (VS *Server) getTodoCheckedByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) ToDoChecked {
	var todo ToDoChecked
	err := VS.DataBase.Get(id, &todo)
	VS.checkOperation(c, class, method, operation, err)
	return todo
}

// This method returns all the Villages of our System
func (VS *Server) getAllToDoFromUser(c *gin.Context, class string, method string, operation string) []ToDo {
	u := VS.getUserFomCookie(c, class, method, operation)
	var todo []ToDo
	err := VS.DataBase.Find(&todo, bolthold.Where("IDuser").Eq(u.ID))
	VS.checkOperation(c, class, method, operation, err)
	return todo
}

/***********************  STOCK ***********************/
// This method returns all the Villages of our System
func (VS *Server) getStockByItemAndService(c *gin.Context, idItem bson.ObjectId, IDservice bson.ObjectId, class string, method string, operation string) (Stock, bool) {
	var stock []Stock
	VS.DataBase.Find(&stock, bolthold.Where("VideoCourseOrItemID").Eq(idItem).And("ServiceIDLocation").Eq(IDservice))
	if len(stock) != 0 {
		return stock[0], true
	}
	return Stock{}, false
}

// This method returns is exist the stock To Deliver from the delivery
// The second element return true if it was sucessfully and false in the opposite case
func (VS *Server) getStockByDeliveryAndItemOrVideoID(c *gin.Context, stockToDeliverID bson.ObjectId, idDelivery bson.ObjectId, class string, method string, operation string) (Stock, bool) {
	var stockResult []Stock
	VS.DataBase.Find(&stockResult, bolthold.Where("VideoCourseOrItemID").Eq(stockToDeliverID).And("ServiceIDLocation").Eq(idDelivery))
	if len(stockResult) != 0 {
		return stockResult[0], true
	}
	return Stock{}, false
}

// This method returns all the Villages of our System
func (VS *Server) getPDFByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) PDF {
	var rep PDF
	err := VS.DataBase.Get(id, &rep)
	VS.checkOperation(c, class, method, operation, err)
	return rep
}

// This method returns all the Tools of our workshop
func (VS *Server) getChecksFromItemID(c *gin.Context, itemID bson.ObjectId, class string, method string, operation string) []CheckStock {
	var check []CheckStock
	err := VS.DataBase.Find(&check, bolthold.Where("IDInstance").Eq(itemID))
	VS.checkOperation(c, class, method, operation, err)
	return check
}

// functionality

// ContainsService tells whether arrayString a contains string x
func ContainsService(a []ServiceType, x ServiceType) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// This method cleans the array of ServiceOrder and leave only the orders that are not complete
func cleanServiceOrdersWithoutWindow(array []ServiceOrder) []ServiceOrder {
	array = revertServiceOrder(array)
	var ordersClean []ServiceOrder
	for i := 0; i < len(array); i++ {
		if array[i].Quantity > array[i].AlreadyMade {
			ordersClean = append(ordersClean, array[i])
		}
	}
	return ordersClean
}

// This method return the Workshop Order to make the changes in the already Done
func lastWorkshopOrder(array []ServiceOrder) ServiceOrder {
	var ordersClean []ServiceOrder
	for i := 0; i < len(array); i++ {
		if array[i].Quantity > array[i].AlreadyMade {
			ordersClean = append(ordersClean, array[i])
		}
	}
	return ordersClean[0]
}

// // This method give us all the items that his category haveaccess to the service provided
// func (VS *Server) getAllItemsForRentWithStock(c *gin.Context, service Service) []ItemsPicker {
// 	allCategories := VS.getAllCategoriesOnlyForRent(c, "functionality", "getAllItemsByService", "allCategories := VS.getAllCategories")
// 	// We clean and take only the categories with access to the Item
// 	var categoriesWithItem []Category
// 	for _, category := range allCategories {
// 		if ContainsService(category.ServicesAccess, service.Type) {
// 			categoriesWithItem = append(categoriesWithItem, category)
// 		}
// 	}
// 	// Now we take all the items from the categories
// 	var itemsHTML []ItemsPicker
// 	var item ItemsPicker
// 	for _, category := range categoriesWithItem {
// 		item.CategoryName = category.Name
// 		item.Stocks = VS.itemsToHTML(c, VS.getAllItemsWithStockByCategory(c, category.ID, service.ID, "functionality", "getAllItemsByService", "item.Items = VS.getAllItemsByCategory"))
// 		itemsHTML = append(itemsHTML, item)
// 	}
// 	return itemsHTML
// }

// // This method cleans the array of WorkerOrder and leave only the orders that are not complete and with the window period time
// func (VS *Server) divideStocksByType(c *gin.Context, stocks []Stock) ([]Stock, []Stock, []Stock) {
// 	tpe := VS.getCategoryTypeByName(c, string(CategoryTypeTools), "database", "getAllCategoriesOnlyForSell", "tpe := VS.getCategoryTypeByName")
// 	var stocksMaterial []Stock
// 	var stocksProduct []Stock
// 	var stocksTool []Stock
// 	i := Item{}
// 	for _, s := range stocks {
// 		item, _ := VS.getItemByID(c, s.VideoCourseOrItemID, "functionality", "divideStocksByType", "item := VS.getItemByID", false)
// 		if item == i {
// 			stocksProduct = append(stocksProduct, s)
// 		} else {
// 			cat := VS.getCategoryByID(c, item.IDCategory, "functionality", "divideStocksByType", "item := VS.getItemByID")
// 			if cat.Type == tpe.ID {
// 				stocksTool = append(stocksTool, s)
// 			} else {
// 				stocksMaterial = append(stocksMaterial, s)
// 			}

// 		}
// 	}
// 	return stocksMaterial, stocksProduct, stocksTool
// }

// // This handler shows the Specific services to choose
// func (VS *Server) showChooseService(serviceType ServiceType) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		role := getRole(c)
// 		var services []Service
// 		if role == Admin {
// 			services = VS.getAllSpecificServiceFromSystem(c, serviceType, "h_stock", "showStockService", "ser := VS.getAllSpecificServicesFromLocalVillage")
// 		} else {
// 			services = VS.getAllSpecificServicesFromLocalVillage(c, serviceType, "h_stock", "showStockService", "ser := VS.getAllSpecificServicesFromLocalVillage")
// 		}
// 		servicesHTML := VS.servicesToHTML(c, services)
// 		render(c, gin.H{
// 			"services": servicesHTML,
// 			"role":     role,
// 			"position": string(serviceType),
// 		}, "inv-choose-service.html")
// 	}
// }

// // This handler shows the Deliveries of the service
// func (VS *Server) showStockService(c *gin.Context) {
// 	role := getRole(c)
// 	idService, b := VS.getIDParamFromURL(c, "id")
// 	if !b {
// 		return
// 	}
// 	service := VS.getServiceByID(c, idService, "h_stock", "showStockService", "ser := VS.getAllSpecificServicesFromLocalVillage")
// 	stocks := VS.getAllStockByServiceID(c, service.ID, "h_stock", "showStockService", "stocks := VS.getAllPositiveStockByServiceID")

// 	stocksMaterial, stocksProduct, stocksTool := VS.divideStocksByType(c, stocks)
// 	stocksMaterialHTML := VS.toViewstock(c, stocksMaterial)
// 	stocksProductHTML := VS.toViewstock(c, stocksProduct)
// 	stocksToolHTML := VS.toViewstock(c, stocksTool)

// 	render(c, gin.H{
// 		"role":               role,
// 		"stocksMaterialHTML": stocksMaterialHTML,
// 		"stocksProductHTML":  stocksProductHTML,
// 		"stocksToolHTML":     stocksToolHTML,
// 		"position":           string(service.Type),
// 	}, "inventory-service.html")
// }

// Gets the information of the ID, only takes until the barrier
// The second element return true if it was sucessfully and false in the opposite case
func (VS *Server) getArrayIDFromHTMLCleaningIt(c *gin.Context, barrier string, name string, element string) ([]bson.ObjectId, bool) {
	array := c.PostFormArray(name)
	var arrayID []bson.ObjectId
	for _, e := range array {
		// i := strings.Index(e, barrier)
		// e = e[0:i]
		if bson.IsObjectIdHex(e) {
			arrayID = append(arrayID, bson.ObjectIdHex(e))
		} else {
			VS.wrongFeedback(c, "Sorry, there is a problem with the "+element+"  ID :(")
			return arrayID, false
		}
	}
	return arrayID, true
}

// ServicePurchaseHTML is the Workshop Purchase with the Humand readable information
type ServicePurchaseHTML struct {
	ID           string `json:"id"`
	UserPhoto    string `json:"userphoto"`
	UserID       string `json:"userid"`
	Username     string `json:"username"`
	WorkshopName string `json:"workshopname"`
	ProductPhoto string `json:"productPhoto"`
	ProductName  string `json:"productName"`
	ProductID    string `json:"productid"`
	Date         string `json:"date"`
	DateUnix     int64  `json:"dateunix"`
	State        string `json:"state"`
	Price        int    `json:"price"`
	Photo        string `json:"photo"`
}

// // We mark the item as checked and we created or update the stock related
// func (VS *Server) makeStockDelivery(c *gin.Context, itemPackID DeliveryItem, idService bson.ObjectId) {
// 	// Make the stock
// 	stock, b := VS.getStockByItemAndService(c, itemPackID.IDItem, idService, "h_transactions", "makeStockDelivery", "stock := VS.getStockByItemAndService")
// 	if b {
// 		stock.Quantity += itemPackID.Quantity
// 		stock.DateLastUpdate = getTime()
// 		VS.addStockToDatabase(c, stock, "VS.addStockToDatabase", "makeStockDelivery", "stock := VS.getStockByItemAndService")
// 	} else {
// 		stock = Stock{
// 			ID:        bson.NewObjectId(),
// 			ItemID:    itemPackID.IDItem,
// 			ServiceID: idService,
// 			Quantity:  itemPackID.Quantity,
// 			// This changed with the last sync of the item
// 			DateLastUpdate: getTime(),
// 		}
// 	}
// 	VS.addStockToDatabase(c, stock, "VS.addStockToDatabase", "makeStockDelivery", "stock := VS.getStockByItemAndService")
// }
