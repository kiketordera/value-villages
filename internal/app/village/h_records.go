package village

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Shows the screen with all the actions that the worker haveinside the Service records
func (VS *Server) showRecordsActionsWorker(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showRecordsActionsWorker", "service, _ := VS.getOurService")
	if !b {
		return
	}
	servicesHTML := VS.getServicesPickerByType(c, service.TypeID)
	render(c, gin.H{
		"services":    servicesHTML,
		"service":     service,
		"serviceType": service.TypeID,
	}, "records-actions-worker.html")
}

// Shows the screen with all the actions that the Manager and Admin haveinside the Service records
func (VS *Server) showRecordsActionsManager(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showRecordsActionsManager", "service, b := VS.getOurService(")
	if !b {
		return
	}
	servicesHTML := VS.getServicesPickerByType(c, service.TypeID)
	render(c, gin.H{
		"services":    servicesHTML,
		"service":     service,
		"serviceType": service.TypeID,
	}, "records-actions-manager.html")
}

// When there is not Service chosen for representing the information, shows the Picker to select a service to see the information
func (VS *Server) showChooseServicePicker(c *gin.Context) {
	sType, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	servicesHTML := VS.getServicesPickerByType(c, sType)
	render(c, gin.H{
		"services":    servicesHTML,
		"serviceType": sType,
	}, "choose-service.html")
}

// Sets up the service for represent his information
func (VS *Server) setUpService(c *gin.Context) {
	serviceID, b := VS.getIDFromHTML(c, "serviceID", "Service")
	if !b {
		return
	}
	serviceTypeID, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	IDVisualization[serviceTypeID] = serviceID
	c.Redirect(http.StatusFound, "/records/see/actions-manager/"+serviceTypeID.Hex())
}

// Shows the screen with all the payments in the database of the specific service
func (VS *Server) showRecordsPayments(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showRecordsPayments", "service, b := VS.getOurService")
	if !b {
		return
	}
	payments := VS.getAllPaymentsFromSpecificService(c, service.ID, "h_records", "showRecordsPayments", "payments := VS.getAllPaymentsFromSpecificService")
	paymentsHTML := VS.toViewPayments(c, revertPayments(payments))
	JSON, err := json.Marshal(paymentsHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"paymentsHTML": paymentsHTML,
		"service":      service,
		"serviceType":  service.TypeID,
		"JSON":         string(JSON),
	}, "service-payments.html")
}

// Shows the page to introduce the information of the new Payment
func (VS *Server) showNewPayment(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showNewPayment", "service, b := VS.getOurService(c")
	if !b {
		return
	}
	//Get all Workers from DataBase in the local workshop
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_records", "showNewPayment", "workers := VS.getWorkersFromSpecificService")
	render(c, gin.H{
		"workers":     workers,
		"service":     service,
		"serviceType": service.TypeID,
	}, "new-payment.html")
}

// Stores the new Payment into the DataBase
func (VS *Server) newPayment(c *gin.Context) {
	id := bson.NewObjectId()
	service, b := VS.getOurService(c, "h_records", "newPayment", "service, b := VS.getOurService")
	if !b {
		return
	}
	workerID, b := VS.getIDFromHTML(c, "workerID", "Payment")
	if !b {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	worker := VS.getUserByID(c, workerID, "h_records", "newPayment", "worker := VS.getUserByID")
	photo, b := VS.getFileFromHTML(c, "photo", worker.Name+"_"+id.Hex(), basePath+"/local-resources/users/"+worker.Username+"/payments/")
	if !b {
		return
	}
	quantity, b := VS.getFloatFromHTML(c, "quantity", "Payment")
	if !b {
		return
	}
	payment := Payment{
		ID:        id,
		IDService: service.ID,
		IDWorker:  workerID,
		Quantity:  quantity,
		Photo:     photo,
		Date:      getTime(),
		IDManager: ManagerID,
	}
	VS.addElementToDatabaseWithRegister(c, payment, payment.ID, service.ID, "handlers_manager", "addPayment", "v.addPaymentToDatabase")
	VS.makePayment(c, payment)
	VS.goodFeedback(c, "records/see/payments/"+idservicetype.Hex())
}

// Shows the Payment information from inside a Service
func (VS *Server) seePaymentWorkerService(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "seePaymentWorkerService", "service, b := VS.getOurService")
	if !b {
		return
	}
	pID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	payment := VS.getPaymentByID(c, pID, "h_records", "seePaymentWorkerService", "payment := VS.getPaymentByID")
	user := VS.getUserByID(c, payment.IDWorker, "h_records", "seePaymentWorkerService", "user := VS.getUserByID")
	render(c, gin.H{
		"user":        user,
		"payment":     payment,
		"fromService": true,
		"serviceType": service.TypeID,
	}, "see-payment.html")
}

// Shows the Payment information from inside Global
func (VS *Server) seePaymentWorkerGlobal(c *gin.Context) {
	// service, b := VS.getOurService(c, "h_records", "seePaymentWorkerGlobal", "service, b := VS.getOurService")
	// if !b {
	// 	return
	// }
	pID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	payment := VS.getPaymentByID(c, pID, "h_records", "seePaymentWorkerGlobal", "payment := VS.getPaymentByID")
	user := VS.getUserByID(c, payment.IDWorker, "h_records", "seePaymentWorkerGlobal", "user := VS.getUserByID")
	render(c, gin.H{
		"user":    user,
		"payment": payment,
		// "serviceType": service.TypeID,
	}, "see-payment.html")
}

/*************************  ORDERS  *************************/

// Shows the screen with all the ServicesOrders in the database
func (VS *Server) showServiceOrders(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showServiceOrders", "service, b := VS.getOurService")
	if !b {
		return
	}
	prodOrders := VS.getAllServiceOrdersFromService(c, service.ID, "h_records", "showServiceOrders", "prodOrders := VS.getAllServiceOrdersFromService")
	prodOrdersHTML := VS.toViewWorkshopOrder(c, prodOrders)
	JSON, err := json.Marshal(prodOrdersHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"service":     service,
		"wOrders":     prodOrdersHTML,
		"serviceType": service.TypeID,
		"JSON":        string(JSON),
	}, "service-orders.html")
}

// Shows the page to introduce the information of the new ServiceOrder
func (VS *Server) showNewServiceOrder(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showNewServiceOrder", "service, b := VS.getOurService")
	if !b {
		return
	}
	products := VS.getAllVideoCoursesByServiceType(c, service.TypeID, "h_records", "showNewServiceOrder", "products := VS.getAllVideoCoursesByServiceType")
	services := VS.getAllSpecificServiceFromSystem(c, service.TypeID, "h_records", "showNewServiceOrder", "services := VS.getAllSpecificServiceFromSystem")
	wHTML := VS.servicesToHTML(c, services)
	render(c, gin.H{
		"products":    products,
		"workshops":   wHTML,
		"serviceType": service.TypeID,
	}, "form-service-order.html")
}

// Stores the new ServiceOrder information into the DataBase
func (VS *Server) newServiceOrder(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "newServiceOrder", "service, b := VS.getOurService")
	if !b {
		return
	}
	serviceOrder, b := VS.getNewServiceOrder(c, bson.NewObjectId())
	if !b {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	VS.addElementToDatabaseWithRegister(c, serviceOrder, serviceOrder.ID, service.ID, "h_records", "newServiceOrder", "VS.addElementToDatabaseWithRegister(c, serviceOrder")
	VS.goodFeedback(c, "/records/service-orders/"+idservicetype.Hex())
}

// Shows the page with the information of the ServiceOrder to modify
func (VS *Server) showEditServiceOrder(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showEditServiceOrder", "service, b := VS.getOurService")
	if !b {
		return
	}
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	wkOrder := VS.getServiceOrderByID(c, id, "h_records", "showEditServiceOrder", "wkOrder := VS.getServiceOrderByID")
	products := VS.getAllVideoCourses(c, "h_records", "showEditServiceOrder", "products := VS.getAllVideoCourses")
	workshops := VS.getAllSpecificServiceFromSystem(c, service.TypeID, "h_records", "showEditServiceOrder", "workshops := VS.getAllSpecificServiceFromSystem")
	render(c, gin.H{
		"po":          wkOrder,
		"products":    products,
		"workshops":   workshops,
		"edit":        true,
		"serviceType": service.TypeID,
	}, "form-service-order.html")
}

// Modifies the ServiceOrder into the Database
func (VS *Server) editServiceOrder(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "editServiceOrder", "service, b := VS.getOurService")
	if !b {
		return
	}
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldServiceOrder := VS.getServiceOrderByID(c, id, "h_records", "editServiceOrder", "oldServiceOrder := VS.getServiceOrderByID")
	newServiceOrder, b := VS.getNewServiceOrder(c, id)
	if !b {
		return
	}
	oldServiceOrder.Assigned = newServiceOrder.Assigned
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}

	// Audit
	audit := VS.getAudit(c, newServiceOrder, id, service.ID, Modified)
	modified := VS.dynamicAuditChange(c, newServiceOrder, oldServiceOrder, audit.ID)
	// Changes
	if oldServiceOrder.IDService != newServiceOrder.IDService {
		wkA := VS.getServiceByID(c, oldServiceOrder.IDService, "h_records", "editServiceOrder", "wkA := VS.getServiceByID(c, oldServiceOrder.ServiceID", true)
		wkB := VS.getServiceByID(c, newServiceOrder.IDService, "h_records", "editServiceOrder", "wkB := VS.getServiceByID", true)
		VS.checkChange(c, audit.ID, "Workshop", wkA.Name, wkB.Name)
		modified = true
	}
	if oldServiceOrder.IDService != newServiceOrder.IDService {
		proA, _ := VS.getVideoCourseByID(c, oldServiceOrder.IDService, "h_records", "editServiceOrder", "proA, _ := VS.getVideoCourseByID", true)
		proB, _ := VS.getVideoCourseByID(c, newServiceOrder.IDService, "h_records", "editServiceOrder", "proB, _ := VS.getVideoCourseByID", true)
		VS.checkChange(c, audit.ID, "Name", proA.Name, proB.Name)
		modified = true
	}
	if modified {
		VS.editElementInDatabaseWithRegister(c, newServiceOrder, oldServiceOrder, newServiceOrder.ID, service.ID, "h_records", "editServiceOrder", "VS.addElementToDatabase")
	}
	VS.goodFeedback(c, "/records/service-orders/"+idservicetype.Hex())
}

// Shows the page to introduce the information of the new WorkerOrder
func (VS *Server) showNewWorkerOrder(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showNewWorkerOrder", "service, b := VS.getOurService")
	if !b {
		return
	}
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_records", "showNewWorkerOrder", "workers := VS.getWorkersFromSpecificService")
	orders := VS.getAllServiceOrdersFromService(c, service.ID, "h_records", "showNewWorkerOrder", "orders := VS.getAllServiceOrdersFromService")
	ordersClean := cleanServiceOrders(orders)
	ordersHTML := VS.toViewWorkshopOrder(c, ordersClean)
	render(c, gin.H{
		"orders":      ordersHTML,
		"workers":     workers,
		"serviceType": service.TypeID,
	}, "new-worker-order.html")
}

// Stores the new WorkerOrder into the DataBase
func (VS *Server) newWorkerOrder(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "newWorkerOrder", "service, b := VS.getOurService")
	if !b {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	quantity, b := VS.getIntFromHTML(c, "quantity", "Worker Order")
	if !b {
		return
	}
	wkOrderID, b := VS.getIDFromHTML(c, "orderID", "Worker Order")
	if !b {
		return
	}
	wkOrder := VS.getServiceOrderByID(c, wkOrderID, "h_records", "newWorkerOrder", "wkOrder := VS.getServiceOrderByID")
	workerID, b := VS.getIDFromHTML(c, "workerID", "Worker Order")
	if !b {
		return
	}
	// Put quantity as the maximun non-realice Products to avoid they make more items thant assigned
	if quantity > (wkOrder.Quantity - wkOrder.AlreadyMade) {
		quantity = wkOrder.Quantity - wkOrder.AlreadyMade
	}
	wOrder := WorkerOrder{
		ID:          bson.NewObjectId(),
		IDService:   service.ID,
		IDProduct:   wkOrder.IDVideoCourse,
		IDWorker:    workerID,
		Quantity:    quantity,
		Date:        getTime(),
		AlreadyMade: 0,
		IDManager:   ManagerID,
	}
	VS.makeWorkerOrder(c, wOrder)
	VS.addElementToDatabaseWithRegister(c, wOrder, wOrder.ID, service.ID, "h_records", "newWorkerOrder", "VS.addElementToDatabaseWithRegister(c, wOrder")
	VS.goodFeedback(c, "records/worker-orders/"+idservicetype.Hex())
}

// Shows the screen with all the WorkerOrders in the database of the given Service
func (VS *Server) showOrdersWorker(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showOrdersWorker", "service, b := VS.getOurService")
	if !b {
		return
	}
	wkOrders := VS.getAllWorkerOrdersFromService(c, service.ID, "h_records", "showOrdersWorker", "wkOrders := VS.getAllWorkerOrdersFromService")
	ordersHTML := VS.toViewWorkerOrders(c, revertWorkerOrder(wkOrders))
	JSON, err := json.Marshal(ordersHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"service":     service,
		"ordersHTML":  ordersHTML,
		"JSON":        string(JSON),
		"serviceType": service.TypeID,
	}, "orders-worker.html")
}

// Shows the Stock generated inside the service
func (VS *Server) showStockGenerated(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showStockGenerated", "service, b := VS.getOurService")
	if !b {
		return
	}
	wkPurchases := VS.getAllPurchasesFromSpecificService(c, service.ID, "h_records", "showStockGenerated", "wkPurchases := VS.getAllPurchasesFromSpecificService")
	stockHTML := VS.stocksToHTML(c, revertServicePurchase(wkPurchases))
	JSON, err := json.Marshal(stockHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"stock":       stockHTML,
		"JSON":        string(JSON),
		"serviceType": service.TypeID,
	}, "purchases-service.html")
}

// Stores the new Transaction information into the DataBase
// The isFreePurchase indicates if the purchase is related to a ServiceOrder or no
func (VS *Server) newPurchase(isFreePurchase bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		service, b := VS.getOurService(c, "h_records", "newPurchase", "service, b := VS.getOurService")
		if !b {
			return
		}
		var workerID bson.ObjectId
		if isFreePurchase {
			workerID, b = VS.getIDFromHTML(c, "workerID", "Item")
			if !b {
				return
			}
		} else {
			workerID, b = VS.getIDParamFromURL(c, "idworker")
			if !b {
				return
			}
		}

		var users []bson.ObjectId
		users = append(users, workerID)
		stock := Stock{
			IDUser:            workerID,
			ServiceCreated:    service.ID,
			IDServiceLocation: service.ID,
			Date:              getTime(),
			IDManager:         ManagerID,
			TypeOfItem:        ServiceProduct,
			UsersRelated:      users,
		}

		// We check if the original source if from a video course or from a Item
		videoID, _ := VS.getIDFromHTMLWithoutChecking(c, "videoID")
		itemID, isItem := VS.getIDFromHTMLWithoutChecking(c, "itemID")
		orderID, isOrder := VS.getIDFromHTMLWithoutChecking(c, "orderID")
		fmt.Print(itemID)
		fmt.Println(isItem)

		var i Item

		if isOrder {
			order := VS.getWorkerOrderByID(c, orderID, "h_records", "newPurchase", "order := VS.getServiceOrderByID")
			i, isItem = VS.getItemByID(c, order.IDProduct, "h_records", "newPurchase", "i, isItem = VS.getItemByID", true)
			if !isItem {
				videoID = order.IDProduct
			} else {
				itemID = i.ID
			}
		}

		if !isItem {
			video, _ := VS.getVideoCourseByID(c, videoID, "h_records", "newPurchase", "video, _ := VS.getVideoCourseByID", true)
			fmt.Print("videoID 2")
			fmt.Println(videoID)

			QR, b := VS.getQRFromHTMLNew(c, ServiceProduct)
			if !b {
				return
			}
			stock.ID = QR
			stock.IsTrackable = true
			stock.IDVideoCourseOrItem = video.ID
			stock.Photo, b = VS.getFileFromHTML(c, "photo", "Purchase", basePath+"/local-resources/video-courses/"+video.Name+"/purchases/")
			if !b {
				return
			}
			stock.UsersRelated = append(stock.UsersRelated, VS.getUsersRelatedFromStock(c, stock.IDVideoCourseOrItem)...)
		} else {
			item, _ := VS.getItemByID(c, itemID, "h_records", "newPurchase", "item, _ := VS.getItemByID", true)
			cat := VS.getCategoryByID(c, item.IDCategory, "h_records", "newPurchase", "cat := VS.getCategoryByID")
			stock.IsTrackable = cat.IsTrackable
			stock.IDVideoCourseOrItem = item.ID
			tpe := VS.getCategoryTypeByID(c, cat.Type, "h_records", "newPurchase", "tpe := VS.getCategoryTypeByID")
			stock.Photo, b = VS.getFileFromHTML(c, "photo", "Purchase", basePath+"/local-resources/categories/"+tpe.Name+"/"+cat.Name+"/purchases/")
			if !b {
				return
			}
			if cat.IsTrackable {
				QR, b := VS.getQRFromHTMLNew(c, ServiceProduct)
				if !b {
					return
				}
				stock.ID = QR
			} else {
				stock.ID = bson.NewObjectId()
			}
		}
		fmt.Print("This is stock.VideoCourseOrItemID")
		fmt.Print(stock.IDVideoCourseOrItem)

		if isOrder {
			VS.makeServicePurchaseFromOrder(c, stock, orderID)
		} else {
			VS.makeServiceFreePurchase(c, stock)
		}
		idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
		if !b {
			return
		}
		VS.addElementToDatabaseWithRegister(c, stock, stock.ID, service.ID, "h_records", "newPurchase", "VS.addElementToDatabaseWithRegister")
		VS.goodFeedback(c, "records/stock-generated/"+idservicetype.Hex())
	}
}

// Shows the page to introduce the information of the new Purchase form a Service to a Worker
func (VS *Server) showNewFreePurchase(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showNewFreePurchase", "service, b := VS.getOurService")
	if !b {
		return
	}
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_records", "showNewFreePurchase", "workers := VS.getWorkersFromSpecificService")
	workersHTML := VS.usersToHTML(c, workers)
	videos := VS.getAllVideoCoursesByServiceType(c, service.TypeID, "h_records", "showNewFreePurchase", "videos := VS.getAllVideoCoursesByServiceType")
	videosHTML := VS.toViewVideoCourses(c, videos)
	items := VS.getAllItemsForSellServicePurchase(c, service.TypeID)
	render(c, gin.H{
		"workers":     workersHTML,
		"videos":      videosHTML,
		"items":       items,
		"serviceType": service.TypeID,
	}, "new-purchase.html")
}

// Shows the page to introduce the information of the new Purchase
func (VS *Server) showNewPurchase(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "showNewPurchase", "service, b := VS.getOurService")
	if !b {
		return
	}
	workerID, _ := VS.getIDParamFromURL(c, "idworker")
	user := VS.getUserByID(c, workerID, "h_records", "showNewPurchase", "user := VS.getUserByID")
	orders := VS.getWorkerOrderByWorkerID(c, workerID, "h_records", "showNewPurchase", "orders := VS.getWorkerOrderByWorkerID")
	ordersClean := cleanWorkerOrders(orders)
	wkOrdersHTML := VS.toViewWorkerOrders(c, ordersClean)
	render(c, gin.H{
		"orders":      wkOrdersHTML,
		"user":        user,
		"serviceType": service.TypeID,
	}, "new-purchase.html")
}

// This handler gives the option to the user to choose between send or received Deliveries
func (VS *Server) chooseDeliveriesSentReceived(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "chooseDeliveriesSentReceived", "service, b := VS.getOurService")
	if !b {
		return
	}
	render(c, gin.H{
		"serviceType": service.TypeID,
	}, "delivery-send-receive.html")
}

// This handler shows the Deliveries send or received (depend on the deliveryType) from the service given
func (VS *Server) showDeliveriesService(deliveryType DeliveryType) gin.HandlerFunc {
	return func(c *gin.Context) {
		service, b := VS.getOurService(c, "h_records", "showDeliveriesService", "service, b := VS.getOurService")
		if !b {
			return
		}
		var dPack []Delivery
		if deliveryType == Received {
			dPack = VS.getAllDeliveriesByServiceReceiverIDSent(c, service.ID, true)
		} else {
			dPack = VS.getAllDeliveriesSentByServiceID(c, service.ID, false)
		}
		dPackHTML := VS.deliveryToHTML(c, dPack)
		render(c, gin.H{
			"deliveryType": deliveryType,
			"dPack":        dPackHTML,
			"serviceType":  service.TypeID,
		}, "deliveries-service.html")
	}
}

// This handler shows the current Stock of the service given
func (VS *Server) showStockService(c *gin.Context) {
	iType, b := VS.getParamFromURL(c, "itemtype")
	if !b {
		return
	}
	itemType := TypeOfItem(iType)
	service, b := VS.getOurService(c, "h_records", "showStockService", "service, b := VS.getOurService")
	if !b {
		return
	}
	stocks := VS.getAllStockByServiceIDandTypeOfItem(c, service.ID, itemType, "h_records", "showStockService", "stocks := VS.getAllStockByServiceIDandTypeOfItem")
	stocksHTML := VS.stocksToHTML(c, stocks)
	JSON, err := json.Marshal(stocksHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)

	render(c, gin.H{
		"stock":       stocksHTML,
		"itemType":    iType,
		"serviceType": service.TypeID,
		"JSON":        string(JSON),
	}, "purchases-service.html")
}

// This handler gives the option to the user to choose which kind of stock he wants to see: Tools, Materials or Products made in the services
func (VS *Server) chooseTypeItem(c *gin.Context) {
	service, b := VS.getOurService(c, "h_records", "chooseTypeItem", "service, b := VS.getOurService")
	if !b {
		return
	}
	render(c, gin.H{
		"serviceType": service.TypeID,
	}, "choose-type-item.html")
}
