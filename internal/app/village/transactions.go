package village

import (
	"encoding/gob"
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Make all the changes in the Service Order, Worker Order, inventory, Wallet and Service balance
func (VS *Server) makeServicePurchaseFromOrder(c *gin.Context, servicePurchase Stock, wkOrderID bson.ObjectId) {
	var cost float64
	item, b := VS.getItemByID(c, servicePurchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
	if b {
		cost = item.Price
	} else {
		video, _ := VS.getVideoCourseByID(c, servicePurchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
		cost = video.Price
	}

	/***********************  Products ***********************/
	// Add the purchase to the Worker Order
	// We take the first order from the Worker that is related with the Video Course
	order := VS.getWorkerOrderByID(c, wkOrderID, "h_transactions", "makeWorkshopPurchase", "VS.addWorkshopOrderToDatabase")
	order.AlreadyMade++
	VS.addElementToDatabaseWithRegister(c, order, order.ID, order.IDService, "h_transactions", "makeWorkshopPurchase", "VS.addWorkerOrderToDatabase")

	// Add the purchase to the Workshop Order
	wkOrder, b := VS.getServiceOrderFromServiceAndProductInDate(c, order.IDService, servicePurchase.IDVideoCourseOrItem, "handlers_workshop", "makeWorkshopPurchase", "wkOrder := VS.getWorkshopOrderByProductID")
	if b {
		wkOrder.AlreadyMade++
		VS.addElementToDatabaseWithRegister(c, wkOrder, wkOrder.ID, order.IDService, "h_transactions", "makeWorkshopPurchase", "VS.addWorkshopOrderToDatabase")
	}

	/***********************  Payment ***********************/
	// The service Pays
	wk := VS.getServiceByID(c, servicePurchase.ServiceCreated, "h_transactions", "makeWorkshopPurchase", "wk := VS.getServiceByID", true)
	wk.Balance -= cost
	VS.addServiceToDatabase(c, wk, "h_transactions", "makeWorkshopPurchase", "VS.addServiceToDatabase")

	// The Worker receives the money
	worker := VS.getUserByID(c, servicePurchase.IDUser, "h_transactions", "makeWorkshopPurchase", "worker := VS.getUserByID")
	worker.Balance += cost
	VS.addUserToDatabase(c, worker, "h_transactions", "makeWorkshopPurchase", "VS.addUserToDatabase")

	/***********************  STOCK ***********************/
	stock, b := VS.getStockByItemAndService(c, servicePurchase.IDVideoCourseOrItem, wk.ID, "h_transactions", "makeStockDelivery", "stock := VS.getStockByItemAndService")
	if b {
		VS.addElementToDatabaseWithRegister(c, stock, stock.ID, wk.ID, "VS.addStockToDatabase", "makeStockDelivery", "stock := VS.getStockByItemAndService")
	} else {
		newStock := servicePurchase
		newStock.ID = bson.NewObjectId()
		newStock.Date = getTime()
		newStock.UsersRelated = nil
		newStock.UsersRelated = append(newStock.UsersRelated, servicePurchase.IDUser)
		// newStock := Stock{
		// 	ID: bson.NewObjectId(),
		// 	IDVideoCourseOrItem: servicePurchase.IDVideoCourseOrItem,
		// 	Photo: servicePurchase.Photo,
		// 	IDServiceLocation: servicePurchase.IDServiceLocation,
		// 	TypeOfItem: servicePurchase.TypeOfItem,
		// 	Quantity: 1,
		// 	IDUser: servicePurchase.IDUser,
		// }
		VS.addElementToDatabaseWithRegister(c, newStock, newStock.ID, wk.ID, "VS.addStockToDatabase", "makeStockDelivery", "stock := VS.getStockByItemAndService")
	}
}

// Make all the changes in the inventory, Wallet and Service balance
func (VS *Server) makeServiceFreePurchase(c *gin.Context, servicePurchase Stock) {
	var cost float64
	item, b := VS.getItemByID(c, servicePurchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
	if b {
		cost = item.Price
	} else {
		video, _ := VS.getVideoCourseByID(c, servicePurchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
		cost = video.Price
	}
	/***********************  Payment ***********************/
	// The service Pays
	wk := VS.getServiceByID(c, servicePurchase.ServiceCreated, "h_transactions", "makeWorkshopPurchase", "wk := VS.getServiceByID", true)
	wk.Balance -= cost
	VS.addServiceToDatabase(c, wk, "h_transactions", "makeWorkshopPurchase", "VS.addServiceToDatabase")
	// The Worker receives the money
	worker := VS.getUserByID(c, servicePurchase.IDUser, "h_transactions", "makeWorkshopPurchase", "worker := VS.getUserByID")
	worker.Balance += cost
	VS.addUserToDatabase(c, worker, "h_transactions", "makeWorkshopPurchase", "VS.addUserToDatabase")
}

// Make all the changes in the Service Order, the assignments
func (VS *Server) makeWorkerOrder(c *gin.Context, wOrder WorkerOrder) {
	serviceOrder, b := VS.getServiceOrderFromServiceAndProductInDate(c, wOrder.IDService, wOrder.IDProduct, "h_transactions", "makeWorkerOrder", "wkOrder := VS.getWorkerOrderByID")
	if b {
		serviceOrder.Assigned += wOrder.Quantity
		VS.addElementToDatabaseWithRegister(c, serviceOrder, serviceOrder.ID, serviceOrder.IDService, "h_transactions", "makeWorkerOrder", "wkOrder := VS.getWorkerOrderByID")
	}
}

// Make all the changes in the Wallet, User balance and Service balance
func (VS *Server) makePayment(c *gin.Context, pay Payment) {
	// Workshop Pays
	wk := VS.getServiceByID(c, pay.IDService, "h_transactions", "makePayment", "wk := VS.getServiceByID", true)
	wk.Balance -= pay.Quantity
	VS.addServiceToDatabase(c, wk, "h_transactions", "makePayment", "VS.addServiceToDatabase")

	// Worker receives the payment
	worker := VS.getUserByID(c, pay.IDWorker, "h_transactions", "makePayment", "worker := VS.getUserByID")
	worker.Balance -= pay.Quantity
	VS.addUserToDatabase(c, worker, "h_transactions", "makePayment", "VS.addUserToDatabase")
}

// Make all the changes in the Wallet, Inventory, User balance and Service balance
func (VS *Server) makeSale(c *gin.Context, sale Sale) {
	// The worker pays for the materials
	worker := VS.getUserByID(c, sale.IDWorker, "h_transactions", "makeSale", "worker := VS.getUserByID")
	worker.Balance -= sale.Price
	VS.addUserToDatabase(c, worker, "h_transactions", "makeSale", "VS.addUserToDatabase")

	// The Workshop gets money for the materials
	wk := VS.getServiceByID(c, sale.IDService, "h_transactions", "makeSale", "wk := VS.getOurWorkshop", true)
	wk.Balance += sale.Price
	VS.addServiceToDatabase(c, wk, "h_transactions", "makeSale", "VS.addServiceToDatabase")

	// Make the stock (if you are selling it there is already a stock created)
	stock, b := VS.getStockByItemAndService(c, sale.IDStock, sale.IDService, "h_transactions", "makeStockDelivery", "stock := VS.getStockByItemAndService")
	if !b {
		fmt.Println("Error, stock not found")
		fmt.Print("this is sale.IDStock")
		fmt.Println(sale.IDStock)
		fmt.Print("this is sale.IDService")
		fmt.Println(sale.IDStock)
		return
	}
	stockOriginal := stock
	VS.editElementInDatabaseWithRegister(c, stock, stockOriginal, stock.ID, sale.IDService, "VS.addStockToDatabase", "makeStockDelivery", "stock := VS.getStockByItemAndService")
}

// registerAddDatabaseMovement adds the element to DB and registers the audit and the event related with the object created or modified
func (VS *Server) registerAddDatabaseMovement(c *gin.Context, object interface{}, objectID bson.ObjectId, serviceID bson.ObjectId) {
	// Get the Audit, EventDB and save the object
	// Audit
	audit := VS.getAudit(c, object, objectID, serviceID, Created)
	err := VS.DataBase.Upsert(audit.ID, audit)
	VS.checkOperation(c, "transactions", "registerAddDatabaseMovement", "Insert audit in DB-Upsert", err)

	// EventDB
	// We need to register the struct as an interface to save it on the EventDB
	gob.Register(object)
	event := VS.getEvent(c, nil, object, audit.ID)
	err = VS.DataBase.Upsert(event.ID, event)
	VS.checkOperation(c, "transactions", "registerAddDatabaseMovement", "Insert event in DB-Upsert", err)
}

// registerAddDatabaseMovement adds the element to DB and registers the audit and the event related with the object created or modified
func (VS *Server) registerAddDatabaseMovementForStartinggTheAPP(object interface{}, objectID bson.ObjectId, serviceID bson.ObjectId, IDAdmin bson.ObjectId, IDVillage bson.ObjectId) {
	// Get the Audit, EventDB and save the object
	// Audit
	audit := VS.getAuditForStartingTheAPP(object, objectID, IDAdmin, IDVillage, serviceID, Created)
	err := VS.DataBase.Upsert(audit.ID, audit)
	if err != nil {
		fmt.Print(err)
		return
	}
	// EventDB
	// We need to register the struct as an interface to save it on the EventDB
	gob.Register(object)
	event := VS.getEventForStartingAPP(nil, object, audit.ID, IDVillage)
	err = VS.DataBase.Upsert(event.ID, event)
}

// registerAddDatabaseMovementFromCentral adds the element to DB and registers the audit and the event related with the object created or modified. Sets the Central as ServiceID in the audit
func (VS *Server) registerAddDatabaseMovementFromCentral(c *gin.Context, object interface{}, objectID bson.ObjectId) {
	// Get the Audit, EventDB and save the object
	// Audit
	audit := VS.getAuditFromCentral(c, object, objectID, Created)
	err := VS.DataBase.Upsert(audit.ID, audit)
	VS.checkOperation(c, "transactions", "registerAddDatabaseMovement", "Insert audit in DB-Upsert", err)
	// EventDB
	// We need to register the struct as an interface to save it on the EventDB
	gob.Register(object)
	event := VS.getEvent(c, nil, object, audit.ID)
	err = VS.DataBase.Upsert(event.ID, event)
	VS.checkOperation(c, "transactions", "registerAddDatabaseMovement", "Insert event in DB-Upsert", err)
}

// This methos registers the audit and event from the data given
func (VS *Server) registerAuditAndEvent(c *gin.Context, audit Audit, object interface{}, id bson.ObjectId) {
	// Audit
	VS.addElementToDatabaseWithoutRegister(c, audit, audit.ID, "h_data", "newVillage", "addAuditToDatabase")
	// EventDB
	// We need to register the struct as an interface to save it on the EventDB
	gob.Register(object)
	event := VS.getEvent(c, nil, object, audit.ID)
	VS.addElementToDatabaseWithoutRegister(c, event, event.ID, "h_data", "newVillage", "VS.addEventToDatabase")
}
