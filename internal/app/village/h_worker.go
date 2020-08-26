package village

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Shows the index page to the worker
func (VS *Server) showIndex(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showIndex", "service, b := VS.getOurService")
	if !b {
		return
	}
	render(c, gin.H{
		"serviceType": service.TypeID,
	}, "index.html")
}

// Shows all the transactions related with money to the worker (Material, products and payments)
func (VS *Server) showWalletWorker(c *gin.Context) {
	// service, b := VS.getOurService(c, "h_worker", "showWalletWorker", "service, b := VS.getOurService")
	// if !b {
	// 	return
	// }
	user := VS.getUserFomCookie(c, "h_worker", "showWalletWorker", "user := VS.getUserFomCookie")
	wallet := VS.getWalletSpecificWorker(c, user.ID)
	JSON, err := json.Marshal(wallet)
	VS.checkOperation(c, "h_worker", "showRecordsOrdersWorker", "JSON, err := json.Marshal(wallet)", err)
	render(c, gin.H{
		"JSON":   string(JSON),
		"wallet": wallet,
		// "serviceType": service.TypeID,
	}, "w-wallet.html")
}

// Shows ALL the stock generated of the worker
func (VS *Server) showRecordsSalesWorkerGlobal(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_worker", "showRecordsSalesWorkerGlobal", "user := v.getUserFomCookie")
	purchases := VS.getAllServicePurchasesFromSpecificWorker(c, user.ID, "h_worker", "showRecordsSalesWorkerGlobal", "purchases := VS.getAllServicePurchasesFromSpecificWorker")
	fmt.Print("longitude of purchases: ")
	fmt.Print(len(purchases))

	purchasesHTML := VS.stocksToHTML(c, purchases)
	JSON, err := json.Marshal(purchasesHTML)
	VS.checkOperation(c, "h_worker", "showRecordsSalesWorkerGlobal", "JSON, err := json.Marshal(purchasesHTML", err)
	render(c, gin.H{
		"sales": purchasesHTML,
		"JSON":  string(JSON),
	}, "w-records-sales.html")
}

// Shows the stock generated of the worker in the service
func (VS *Server) showRecordsSalesWorkerService(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showRecordsSalesWorkerService", "service, _ := VS.getOurService")
	if !b {
		return
	}
	user := VS.getUserFomCookie(c, "h_worker", "showRecordsSalesWorkerService", "user := v.getUserFomCookie")
	purchases := VS.getAllStockGeneratedFromSpecificWorkerAndService(c, user.ID, service.ID, "handlers_workshop", "showRecordsSalesWorkerService", "purchases := VS.getAllStockGeneratedFromSpecificWorkerAndService")
	fmt.Print("longitude of purchases:")
	fmt.Print(len(purchases))

	purchasesHTML := VS.stocksToHTML(c, purchases)
	JSON, err := json.Marshal(purchasesHTML)
	VS.checkOperation(c, "h_worker", "showRecordsSalesWorkerService", "JSON, err := json.Marshal(purchasesHTML", err)
	render(c, gin.H{
		"sales":       purchasesHTML,
		"JSON":        string(JSON),
		"fromService": true,
		"serviceType": service.TypeID,
	}, "w-records-sales.html")
}

// Shows all the Payments from the specific worker
func (VS *Server) showRecordsPaymentsWorkerGlobal(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_worker", "showRecordsPaymentsWorkerGlobal", "user := v.getUserFomCookie")
	payments := VS.getAllPaymentsFromSpecificWorker(c, user.ID, "h_worker", "showRecordsPaymentsWorkerGlobal", "payments := v.getAllPaymentsFromSpecificWorker")
	paymentsHTML := VS.toViewPayments(c, payments)
	JSON, err := json.Marshal(paymentsHTML)
	VS.checkOperation(c, "h_worker", "showRecordsPaymentsWorkerGlobal", "JSON, err := json.Marshal(paymentsHTML)", err)
	render(c, gin.H{
		"payments": paymentsHTML,
		"JSON":     string(JSON),
	}, "w-records-payments.html")
}

// Shows the Payments from the specific worker from the service
func (VS *Server) showRecordsPaymentsWorkerService(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showRecordsPaymentsWorkerService", "service, _ := VS.getOurService")
	if !b {
		return
	}
	user := VS.getUserFomCookie(c, "h_worker", "showRecordsPaymentsWorkerService", "user := v.getUserFomCookie")
	payments := VS.getAllPaymentsFromSpecificWorkerAndService(c, user.ID, service.ID, "handlers_worker", "showRecordsPayments", "payments := v.getAllPaymentsFromSpecificWorker")
	paymentsHTML := VS.toViewPayments(c, payments)
	JSON, err := json.Marshal(paymentsHTML)
	VS.checkOperation(c, "h_worker", "showRecordsPaymentsWorkerService", "JSON, err := json.Marshal(paymentsHTML)", err)
	render(c, gin.H{
		"payments":    paymentsHTML,
		"JSON":        string(JSON),
		"fromService": true,
		"serviceType": service.TypeID,
	}, "w-records-payments.html")
}

// Shows all the Worker Orders for the specific worker
func (VS *Server) showRecordsOrdersWorkerGlobal(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_worker", "showRecordsOrdersWorkerGlobal", "user := VS.getUserFomCookie)")
	orders := VS.getAllWorkerOrdersFromSpecificWorker(c, user.ID, "h_worker", "showRecordsOrdersWorkerGlobal", "orders := VS.getAllWorkerOrdersFromSpecificWorker")
	orders = revertWorkerOrder(orders)
	ordersHTML := VS.toViewWorkerOrdersWorkersView(c, orders)
	JSON, err := json.Marshal(ordersHTML)
	VS.checkOperation(c, "h_worker", "showRecordsOrdersWorkerGlobal", "JSON, err := json.Marshal(ordersHTML)", err)
	render(c, gin.H{
		"JSON":   string(JSON),
		"orders": ordersHTML,
	}, "w-records-orders.html")
}

// Shows the Worker Orders for the specific worker from that service
func (VS *Server) showRecordsOrdersWorkerService(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showRecordsOrdersWorkerService", "service, _ := VS.getOurService")
	if !b {
		return
	}
	user := VS.getUserFomCookie(c, "h_worker", "showRecordsOrdersWorkerService", "user := v.getUserFomCookie")
	orders := VS.getAllWorkerOrdersFromSpecificWorkerAndService(c, user.ID, service.ID, "h_worker", "showRecordsOrdersWorkerService", "orders := VS.getAllWorkerOrdersFromSpecificWorkerAndService")
	orders = revertWorkerOrder(orders)
	ordersHTML := VS.toViewWorkerOrdersWorkersView(c, orders)
	JSON, err := json.Marshal(ordersHTML)
	VS.checkOperation(c, "h_worker", "showRecordsOrdersWorkerService", "JSON, err := json.Marshal(ordersHTML))", err)
	render(c, gin.H{
		"JSON":        string(JSON),
		"orders":      ordersHTML,
		"fromService": true,
		"serviceType": service.TypeID,
	}, "w-records-orders.html")
}

// Shows all the Items made or produce for the specific worker
func (VS *Server) showRecordsItemsWorker(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_worker", "showRecordsItemsWorker", "user := v.getUserFomCookie")
	itemsSeelToWorker := VS.getAllSalesFromSpecificWorker(c, user.ID, "h_worker", "showRecordsItemsWorker", "itemsSeelToWorker := VS.getAllSalesFromSpecificWorker")
	itemsSeelToWorker = revertSales(itemsSeelToWorker)
	materialsHTML := VS.salesToHTML(c, itemsSeelToWorker)
	JSON, err := json.Marshal(materialsHTML)
	VS.checkOperation(c, "h_worker", "showRecordsItemsWorker", "JSON, err := json.Marshal(materialsHTML)", err)
	render(c, gin.H{
		"JSON": string(JSON),
	}, "w-records-items.html")
}

// Shows the page with the photo of all workers in order to select the worker to make the Purchase
func (VS *Server) showWorkersForPurchase(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showWorkersForPurchase", "service, _ := VS.getOurService")
	if !b {
		return
	}
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_worker", "showWorkersForPurchase", "workers := v.getAllWorkers")
	render(c, gin.H{
		"workers":     workers,
		"serviceType": service.TypeID,
	}, "show-workers-for-purchase.html")
}

// Shows the page to introduce the information of the new Service Sell
func (VS *Server) showNewSale(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showNewSale", "service, _ := VS.getOurService")
	if !b {
		return
	}
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_worker", "showNewSale", "workers := v.getAllWorkersFromLocalWorkshop")
	itemsToChoose := VS.getAllStockForServicePickerForSell(c, service.ID, "h_worker", "showNewSale", "itemsToChoose := VS.getAllPositiveStockByServiceID")
	render(c, gin.H{
		"workers":     workers,
		"items":       itemsToChoose,
		"serviceType": service.TypeID,
	}, "new-worker-sale.html")
}

// Stores the new Sale into the DataBase
func (VS *Server) newSale(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "newSale", "service, _ := VS.getOurService")
	if !b {
		return
	}
	// We clean the extra information needed for the slider
	stockID, b := VS.getIDFromHTMLCleaningIt(c, "|", "stock", "Stock")
	if !b {
		return
	}

	fmt.Print("this is the stock id: ")
	fmt.Print(stockID)

	userID, b := VS.getIDFromHTML(c, "workerID", "Sale")
	if !b {
		return
	}
	price := 0.0
	quantity := 1.0
	var oldStock Stock
	// Comprobar la checkbox
	QR, isQR := VS.getIDFromHTMLWithoutChecking(c, "qr")
	if isQR {
		fmt.Print("IS QR")
		stockID = QR
		oldStock = VS.getStockByID(c, stockID, "h_worker", "newSale", "oldStock = VS.getStockByID(")
	} else {
		fmt.Print("NO QR")
		quantity, b = VS.getFloatFromHTMLWithoutChecking(c, "quantity", "Sale")
		if !b {
			quantity = 1.0
		}
		oldStock = VS.getStockByID(c, stockID, "h_worker", "newSale", "oldStock = VS.getStockByID(c")
	}
	newStock := oldStock
	newStock.IDServiceLocation = userID
	VS.editElementInDatabaseWithRegister(c, oldStock, newStock, stockID, service.ID, "h_worker", "newSale", "v.editElementInDatabaseWithRegister")

	sale := Sale{
		ID:        bson.NewObjectId(),
		IDService: service.ID,
		IDStock:   stockID,
		IDWorker:  userID,
		IDManager: VS.getUserFomCookie(c, "h_worker", "newSale", "VS.getUserFomCookie").ID,
		Date:      getTime(),
		Price:     price,
		Quantity:  quantity,
	}

	fmt.Print("This is sale")
	fmt.Print(sale)

	VS.makeSale(c, sale)
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	VS.addElementToDatabaseWithRegister(c, sale, sale.ID, service.ID, "h_worker", "newSale", "v.addElementToDatabaseWithRegister")
	VS.goodFeedback(c, "records/worker-sales/"+idservicetype.Hex())
	fmt.Print("This is sale")
	fmt.Print(sale)
}

// Shows the screen with all the Worker Sales of the service
func (VS *Server) showWorkerSales(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showWorkerSales", "service, _ := VS.getOurService")
	if !b {
		return
	}
	wkSales := VS.getAllSalesFromSpecificService(c, service.ID, "h_worker", "showWorkerSales", "wkSales := VS.getAllSalesFromSpecificService")

	wkSalesHTML := VS.salesToHTML(c, wkSales)
	JSON, err := json.Marshal(wkSalesHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"service":     service,
		"ordersHTML":  wkSalesHTML,
		"serviceType": service.TypeID,
		"JSON":        string(JSON),
	}, "worker-sales.html")
}

/*************************  TOOLS  *************************/

// Shows the screen with all assignments
func (VS *Server) showToolsAssignments(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showToolsAssignments", "service, _ := VS.getOurService")
	if !b {
		return
	}
	assignments := VS.getAssignmentsByServiceID(c, service.ID, "h_worker", "showToolsAssignments", "assignments := VS.getAssignmentsByServiceID")
	assignmentsHTML := VS.toViewAssignments(c, assignments)
	render(c, gin.H{
		"assignments": assignmentsHTML,
		"serviceType": service.TypeID,
	}, "assignments.html")
}

// Shows the screen to make a new assignment
func (VS *Server) showNewAssignment(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "showNewAssignment", "service, _ := VS.getOurService")
	if !b {
		return
	}
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_worker", "showNewAssignment", "workers := v.getAllWorkersFromLocalWorkshop")
	items := VS.getAllStockForServicePickerForRent(c, service.ID, "h_worker", "showNewAssignment", "items := VS.getAllStockForServicePickerForRent")
	render(c, gin.H{
		"workers":     workers,
		"items":       items,
		"serviceType": service.TypeID,
	}, "new-assignment.html")
}

// Stores the new assignment
func (VS *Server) newAssignment(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "newAssignment", "service, _ := VS.getOurService")
	if !b {
		return
	}
	workerID, b := VS.getIDFromHTML(c, "workerID", "Worker")
	if !b {
		return
	}
	stockID, b := VS.getQRFromHTML(c)
	if !b {
		return
	}
	stock := VS.getStockByID(c, stockID, "h_worker", "newAssignment", "stock := VS.getStockByID")
	if stock.IDServiceLocation != service.ID {
		VS.wrongFeedback(c, "This object does not belong to this service")
		return
	}
	if !VS.isStockAvailableForAssignment(c, stock.ID, "handlers_manager", "newAssignment", "VS.isStockAvailableForAssignment") {
		VS.wrongFeedback(c, "This item is already assigned to another person")
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}

	assignment := Assignment{
		ID:        bson.NewObjectId(),
		IDService: service.ID,
		IDWorker:  workerID,
		IDManager: VS.getUserFomCookie(c, "h_worker", "newAssignment", "v.getUserFomCookie").ID,
		Date:      getTime(),
		IDStock:   stockID,
		IsBack:    false,
	}
	VS.addElementToDatabaseWithRegister(c, assignment, assignment.ID, service.ID, "h_worker", "newAssignment", "v.addElementToDatabaseWithRegister")
	VS.goodFeedback(c, "records/assignments/"+idservicetype.Hex())
}

// Shows the screen with the Assignment information to change
func (VS *Server) showChangeAssignment(c *gin.Context) {
	sType, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	idAssignment, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	ass := VS.getAssignmentByID(c, idAssignment, "h_worker", "showChangeAssignment", "ass := VS.getAssignmentByID")
	assHTML := VS.toViewAssignment(c, ass)
	render(c, gin.H{
		"ass":         assHTML,
		"serviceType": sType,
	}, "change-assignment.html")
}

// Stores the new Assignment information into the DataBase
func (VS *Server) changeAssignment(c *gin.Context) {
	service, b := VS.getOurService(c, "h_worker", "changeAssignment", "service, _ := VS.getOurService")
	if !b {
		return
	}
	idAssignment, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	ass := VS.getAssignmentByID(c, idAssignment, "h_worker", "changeAssignment", "ass := VS.getAssignmentByID")
	check := VS.getCheckBoxFromHTML(c, "returned")
	if check {
		ass.IsBack = true
		VS.addElementToDatabaseWithRegister(c, ass, ass.ID, service.ID, "h_worker", "changeAssignment", "v.addElementToDatabaseWithRegister")

	}
	VS.goodFeedback(c, "records/assignments/"+idservicetype.Hex())
}
