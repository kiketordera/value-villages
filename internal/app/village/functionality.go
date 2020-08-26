package village

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/timshannon/bolthold"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// This creates the admin user for the first login if it not created
// This method ONLY executes the VERY FIRST TIME you initiallice the app, with no Database
func (VS *Server) createAdminUser() error {

	username, _ := string(Admin), string(Admin)
	if !VS.isUsernameAvailable(username) {
		return errors.New("Admin already exists")
	}
	VS.initialData()
	return nil
}

// This write the hardcoded data THE VERY FIRST TIME the app is initialized
func (VS *Server) initialData() {
	// Create the Admin user
	_, password := string(Admin), string(Admin)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("All the INITIAL DATA has been created")
	copyDirectory(basePath+"/data-for-new-start/", basePath+"/local-resources/")

	// Create the Central Village
	village := Village{
		ID:     bson.NewObjectId(),
		Name:   "Central",
		Prefix: "CENT",
	}

	// Central warehouse
	sType := ServiceType{
		ID:                bson.NewObjectId(),
		Name:              "Warehouse",
		Icon:              "warehouse.svg",
		AllowNoQRPurchase: false,
	}

	// Create the Central Warehouse
	wh := Service{
		ID:        bson.NewObjectId(),
		Name:      "Central Warehouse",
		IDVillage: village.ID,
		Balance:   0,
		Type:      sType.ID,
	}

	// Create the first Admin User
	user := User{
		ID:        bson.NewObjectId(),
		Username:  string(Admin),
		Role:      Admin,
		Name:      string(Admin),
		Password:  string(passwordHash),
		Photo:     "vv-logo.png",
		IDVillage: village.ID,
	}

	VS.addElementToDatabaseWithRegisterForStartingAPP(village, village.ID, wh.ID, village.ID, user.ID, "h_functionality", "initialData()", "// Village  err = VS.addElementToDatabaseWithRegister")
	VS.addElementToDatabaseWithRegisterForStartingAPP(sType, sType.ID, wh.ID, village.ID, user.ID, "h_functionality", "initialData()", "// Village  err = VS.addElementToDatabaseWithRegister")
	VS.addElementToDatabaseWithRegisterForStartingAPP(wh, wh.ID, wh.ID, village.ID, user.ID, "h_functionality", "initialData()", "// Village  err = VS.addElementToDatabaseWithRegister")
	VS.addElementToDatabaseWithRegisterForStartingAPP(user, user.ID, wh.ID, village.ID, user.ID, "h_functionality", "initialData()", "// Village  err = VS.addElementToDatabaseWithRegister")
}

// Gets the Role of the User from the cookie
func getRole(c *gin.Context) Role {
	u, _ := c.Get("role")
	var role Role
	if u != nil {
		role = Role(u.(string))
		return role
	}
	return role
}

// Gets the User from the cookie
func (VS *Server) getUserFomCookie(c *gin.Context, class string, method string, operation string) User {
	// We get here the username from the Cookie
	u, _ := c.Get("username")
	user, _ := VS.getUserByUsername(c, u.(string), class, method, operation)
	return user
}

// Creates an Audit to track the activity of the database
func (VS *Server) getAudit(c *gin.Context, object interface{}, IDitem bson.ObjectId, IDservice bson.ObjectId, description AuditType) Audit {
	var audit Audit
	audit.ID = bson.NewObjectId()
	user := VS.getUserFomCookie(c, "h_functionality", "getAudit", "user := VS.getUserFomCookie")
	audit.IDUser = user.ID
	audit.IDItem = IDitem
	audit.Description = description
	audit.IDService = IDservice
	audit.Date = getTime()
	audit.IDVillage = VS.getLocalVillage(c, "h_functionality", "getAudit", "audit.IDVillage = VS.getLocalVillage").ID
	audit = VS.addInformationAuditItself(audit, object)
	return audit
}

// Creates an Audit to track the activity of the database
// ONLY EXECUTER THE VERY FIRST TIME THE APP INITIALICES (when there is no DataBase)
func (VS *Server) getAuditForStartingTheAPP(object interface{}, IDitem bson.ObjectId, IDAdmin bson.ObjectId, IDVillage bson.ObjectId, IDservice bson.ObjectId, description AuditType) Audit {
	var audit Audit
	audit.ID = bson.NewObjectId()
	audit.IDUser = IDAdmin
	audit.IDItem = IDitem
	audit.Description = description
	audit.IDService = IDservice
	audit.Date = getTime()
	audit.IDVillage = IDVillage
	audit = VS.addInformationAuditItself(audit, object)
	return audit
}

// Creates an Audit to track the activity of the database
func (VS *Server) getAuditFromCentral(c *gin.Context, object interface{}, IDitem bson.ObjectId, description AuditType) Audit {
	var audit Audit
	audit.ID = bson.NewObjectId()
	user := VS.getUserFomCookie(c, "h_functionality", "getAuditFromCentral", "user := VS.getUserFomCookie")
	audit.IDUser = user.ID
	audit.IDItem = IDitem
	audit.Description = description
	audit.IDService = VS.getCentralWarehouse(c, "h_functionality", "getAuditFromCentral", "audit.IDservice = VS.getCentralWarehouse").ID
	audit.Date = getTime()
	audit.IDVillage = VS.getLocalVillage(c, "h_functionality", "getAuditFromCentral", "audit.IDVillage = VS.getLocalVillage").ID
	audit = VS.addInformationAuditItself(audit, object)
	return audit
}

// Creates an EventDB to track the events of the database
func (VS *Server) getEvent(c *gin.Context, previous interface{}, after interface{}, auditID bson.ObjectId) EventDB {
	var event EventDB
	event.ID = bson.NewObjectId()
	event.Previous = previous
	event.After = after
	event.Date = getTime()
	event.IDAudit = auditID
	event.IDVillage = VS.getLocalVillage(c, "h_functionality", "getAuditFromCentral", "audit.IDVillage = VS.getLocalVillage").ID
	return event
}

// Creates an EventDB to track the events of the database
// Creates the Events of the initialData when start without DataBase
func (VS *Server) getEventForStartingAPP(previous interface{}, after interface{}, auditID bson.ObjectId, villageID bson.ObjectId) EventDB {
	var event EventDB
	event.ID = bson.NewObjectId()
	event.Previous = previous
	event.After = after
	event.Date = getTime()
	event.IDAudit = auditID
	event.IDVillage = villageID
	return event
}

// Divides all the audits per day
func (VS *Server) orderAuditPerDay(a []Audit) [][]Audit {
	var auditsPerDay [][]Audit
	previousDayAudit := time.Unix(a[0].Date, 0).Day()
	var aux []Audit

	for i := 0; i < len(a); i++ {
		// If the date is different than the previous one, we change the index in the slice
		if previousDayAudit != time.Unix(a[i].Date, 0).Day() {
			previousDayAudit = time.Unix(a[i].Date, 0).Day()
			auditsPerDay = append(auditsPerDay, aux)
			aux = nil
		}
		aux = append(aux, a[i])
	}
	// This one is for adding the last loop
	auditsPerDay = append(auditsPerDay, aux)
	return auditsPerDay
}

// Divides all the messages per day
func (VS *Server) orderMessagesPerDay(a []Message) [][]Message {
	var messages [][]Message
	previousDayAudit := time.Unix(a[0].Date, 0).Day()
	var aux []Message

	for i := 0; i < len(a); i++ {
		// If the date is different than the previous one, we change the index in the slice
		if previousDayAudit != time.Unix(a[i].Date, 0).Day() {
			previousDayAudit = time.Unix(a[i].Date, 0).Day()
			messages = append(messages, aux)
			aux = nil
		}
		aux = append(aux, a[i])
	}
	// This one is for adding the last loop
	messages = append(messages, aux)
	return messages
}

// Check is two values given are different, if yes it creates a change and returns true
func (VS *Server) checkChange(c *gin.Context, IDAudit bson.ObjectId, nameAttribute, before, after string) bool {
	if !strings.EqualFold(before, after) {
		var ch Change
		ch.ID = bson.NewObjectId()
		ch.IDAudit = IDAudit
		ch.NameAttribute = nameAttribute
		ch.Before = before
		ch.After = after
		VS.addElementToDatabaseWithoutRegister(c, ch, ch.ID, "h_functionality", "checkChange", "VS.addElementToDatabaseWithoutRegister ")
		return true
	}
	return false
}

// Rename the folder to show the images when we change a value related with the directory
func (VS *Server) move(oldPath string, newPath string) {
	// copy.Copy(oldPath, newPath)
	copyDirectory(oldPath, newPath)
	os.RemoveAll(oldPath)
}

// ContainsAccess tells whether []Access contains an Access
func ContainsAccess(a []Access, x bson.ObjectId) bool {
	for _, n := range a {
		if x == n.IDService {
			return true
		}
	}
	return false
}

// ContainsItemsPicker tells whether []ItemsPicker contains a ItemsPicker
func ContainsItemsPicker(a []ItemsPicker, x ItemsPicker) bool {
	for _, n := range a {
		if x.CategoryName == n.CategoryName {
			return true
		}
	}
	return false
}

// ContainsString tells whether arrayString a contains string x
func ContainsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// ContainsID tells whether arrayObjectId a contains ObjectId x
func ContainsID(a []bson.ObjectId, x bson.ObjectId) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// RemovesID give us the array without the element given
func RemovesID(array []bson.ObjectId, element bson.ObjectId) []bson.ObjectId {
	if !ContainsID(array, element) {
		return array
	}
	var solution []bson.ObjectId
	for _, e := range array {
		if e != element {
			solution = append(solution, e)
		}
	}
	return solution
}

// This method orders in alphabetic order the array of Wallet
func walletSort(animals []Wallet) {
	for i := 0; i < len(animals); i++ {
		sweepWallet(animals)
	}
}

// This method orders in alphabetic order the array of Wallet
func sweepWallet(animals []Wallet) {
	N := len(animals)
	firstIndex := 0
	secondIndex := 1
	for secondIndex < N {
		firstVal := animals[firstIndex]
		secondVal := animals[secondIndex]

		var a int64
		a = firstVal.DateUnix
		b := secondVal.DateUnix

		if a > b {
			animals[firstIndex] = secondVal
			animals[secondIndex] = firstVal
		}
		firstIndex++
		secondIndex++
	}
}

// Updates the items in the array in the category by alphabetical order
func (VS *Server) updateOrderToCategory(c *gin.Context, items []Item, cat Category) {
	var id []bson.ObjectId
	for i := 0; i < len(items); i++ {
		id = append(id, items[i].getID())
	}
	cat.Items = id
	VS.addCategoryToDatabase(c, cat, "h_functionality", "updateOrderToCategory", "VS.addCategoryToDatabase")
}

// This method give us all the items that his category haveaccess to the service provided
func (VS *Server) getAllItemsByService(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []ItemsPicker {
	service := VS.getServiceByID(c, serviceID, "h_functionality", "getAllItemsByService", "service := VS.getServiceByID", true)
	allCategories := VS.getAllCategories(c, "h_functionality", "getAllItemsByService", "allCategories := VS.getAllCategories")
	// We clean and take only the categories with access to the Item
	var categoriesWithItem []Category
	for _, category := range allCategories {
		if ContainsID(category.ServicesAccess, service.Type) {
			categoriesWithItem = append(categoriesWithItem, category)
		}
	}
	// Now we take all the items from the categories
	var itemsHTML []ItemsPicker
	var item ItemsPicker
	for _, category := range categoriesWithItem {
		item.CategoryName = category.Name
		item.Items = VS.itemsToHTML(c, VS.getAllItemsByCategory(c, category.ID, "h_functionality", "getAllItemsByService", "item.Stocks = VS.itemsToHTML(c, VS.getAllItemsByCategory"))
		itemsHTML = append(itemsHTML, item)
	}
	return itemsHTML
}

// This method give us all the items in the system for the picker to delivery from Central to a Village
func (VS *Server) getAllITemsForServicePickerGlobal(c *gin.Context, class string, method string, operation string) []ItemsPicker {
	// Categories
	allCategories := VS.getAllCategories(c, "h_functionality", "getAllITemsForServicePickerGlobal", "allCategories := VS.getAllCategories")
	// Now we take all the items from the categories
	var itemsHTML []ItemsPicker
	var item ItemsPicker
	for _, category := range allCategories {
		item.CategoryName = category.Name
		items := VS.getAllItemsByCategory(c, category.ID, "functionality", "getAllITemsForServicePickerGlobal", "item.Items = VS.getAllItemsByCategory")
		item.Items = append(item.Items, VS.itemsToHTML(c, items)...)
		// If the category haveitems
		if len(item.Items) != 0 {
			itemsHTML = append(itemsHTML, item)
		}
		item.Items = nil
	}
	// Videos no because is a labour from the services
	return itemsHTML
}

// This method give us all the items with stock to the picker for delivery
func (VS *Server) getAllStockForServicePicker(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []ItemsPicker {
	stock := VS.getAllPositiveStockByServiceID(c, serviceID, "functionality", "getAllStockForServicePicker", "stock := VS.getAllPositiveStockByServiceID")

	// Now we take all the items from the stocks
	var itemsHTML []ItemsPicker
	var item ItemsPicker
	for _, s := range stock {

		video, b := VS.getVideoCourseByID(c, s.IDVideoCourseOrItem, "functionality", "getAllStockForServicePicker", "video, b := VS.getVideoCourseByID(c, s.VideoCourseOrItemID", false)
		if !b {
			itemReal, _ := VS.getItemByID(c, s.IDVideoCourseOrItem, "functionality", "getAllStockForServicePicker", "itemReal, _ := VS.getItemByID(c, s.VideoCourseOrItemID", false)
			cat := VS.getCategoryByID(c, itemReal.IDCategory, "functionality", "getAllStockForServicePicker", "cat := VS.getCategoryByID(c, itemReal")
			item.CategoryName = cat.Name
			// itemReal.ID = s.ID
			itemRealHTML := VS.itemToHTML(c, itemReal)
			item.Items = append(item.Items, itemRealHTML)
		} else {
			item.CategoryName = "Products from the service"
			i := VS.videoCourseStockToHTML(c, video, serviceID)
			// i.ID = s.ID
			item.Items = append(item.Items, i)
		}
		itemsHTML = append(itemsHTML, item)
		item.Items = nil
	}

	// // Videos
	// 	allvideos := VS.getAllVideoCoursesByServiceType(c, service.Type, "functionality","getAllItemsByService","allCategories := VS.getAllCategories")
	// 	for _, video := range allvideos {
	// 		item.CategoryName = "Products from the service"
	// 		positiveStock := VS.getAllPositiveStockByServiceIDandItemOrVideoID(c, serviceID, video.ID, "h_deliveries", "showItemsDelivery", "VS.getAllPositiveStockByServiceIDandItemOrVideoID")
	// 		if len(positiveStock) != 0 {
	// 			item.Items = append(item.Items, VS.videoCourseStockToHTML(c, video, serviceID))
	// 		}
	// 		itemsHTML = append(itemsHTML, item)
	// 		item.Items = nil
	// 	}
	return itemsHTML
}

// This method give us all the items with stock to the picker for rent
func (VS *Server) getAllStockForServicePickerForRent(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []ItemsPicker {
	stock := VS.getAllPositiveStockByServiceID(c, serviceID, "functionality", "getAllStockForServicePickerForRent", "stock := VS.getAllPositiveStockByServiceID")
	fmt.Print("This is positive stock: ")
	fmt.Println(stock)
	// Now we take all the items from the stocks
	var itemsHTML []ItemsPicker
	var item ItemsPicker
	for _, s := range stock {
		itemReal, b := VS.getItemByID(c, s.IDVideoCourseOrItem, "functionality", "getAllStockForServicePickerForRent", "itemReal, b := VS.getItemByID", false)
		if b {
			fmt.Print("inside if: ")
			cat := VS.getCategoryByID(c, itemReal.IDCategory, "functionality", "getAllStockForServicePickerForRent", "cat := VS.getCategoryByID(c, itemReal")
			if cat.TypeOfItem != Tool {
				break
			}
			item.CategoryName = cat.Name
			// itemReal.ID = s.ID
			itemRealHTML := VS.itemToHTML(c, itemReal)
			item.Items = append(item.Items, itemRealHTML)
			itemsHTML = append(itemsHTML, item)
		} else {
			fmt.Print("outside if: ")
		}
		item.Items = nil
	}
	fmt.Print("itemsHTML: ")
	fmt.Print(itemsHTML)
	return itemsHTML
}

// This method give us all the items with stock to the picker for sell
func (VS *Server) getAllStockForServicePickerForSell(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []StockPicker {
	stocks := VS.getAllPositiveStockByServiceID(c, serviceID, "functionality", "getAllStockForServicePickerForSell", "stock := VS.getAllPositiveStockByServiceID")
	// Now we take all the items from the stocks
	var stocksHTML []StockPicker
	var stock StockPicker
	for _, s := range stocks {
		itemReal, b := VS.getItemByID(c, s.IDVideoCourseOrItem, "functionality", "getAllStockForServicePickerForSell", "itemReal, b := VS.getItemByID", false)
		if b {
			cat := VS.getCategoryByID(c, itemReal.IDCategory, "functionality", "getAllStockForServicePickerForSell", "cat := VS.getCategoryByID(c, itemReal")
			if cat.TypeOfItem == Tool {
				break
			}
			stock.CategoryName = cat.Name
			sHTML := VS.stockToHTML(c, s)
			stock.Stocks = append(stock.Stocks, sHTML)
			stocksHTML = append(stocksHTML, stock)
		}
		stock.Stocks = nil
	}
	return stocksHTML
}

// getTotal gives us the total money (spent, earn and total) to represent it in the performance page from the whole system
func (VS *Server) getTotal(c *gin.Context, sales []Sale, wkPurchases []Stock, payments []Payment) (int, int, int) {
	var increase, decrease, cost float64

	// Increase
	for _, sale := range sales {
		increase += sale.Price
	}
	// Decrease
	for _, purchase := range wkPurchases {
		item, b := VS.getItemByID(c, purchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
		if b {
			cost = item.Price
		} else {
			video, _ := VS.getVideoCourseByID(c, purchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
			cost = video.Price
		}
		decrease -= cost
	}
	for _, payment := range payments {
		decrease -= payment.Quantity
	}
	return int(increase), int(decrease), int(increase + decrease)
}

// getTotal gives us the total money (spent, earn and total) to represent it in the performance page for a specific user
func (VS *Server) getTotalWorker(c *gin.Context, sales []Sale, wkPurchases []Stock, payments []Payment) (int, int, int) {
	var increase, decrease, cost float64
	// Increase
	for _, sell := range sales {
		decrease -= sell.Price
	}
	for _, payment := range payments {
		decrease -= payment.Quantity
	}

	// Decrease
	for _, purchase := range wkPurchases {
		item, b := VS.getItemByID(c, purchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
		if b {
			cost = item.Price
		} else {
			video, _ := VS.getVideoCourseByID(c, purchase.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
			cost = video.Price
		}
		decrease += cost
	}
	return int(increase), int(decrease), int(increase + decrease)
}

// This method cleans the array of WorkerOrders and leave only the orders that are not complete
func cleanWorkerOrders(array []WorkerOrder) []WorkerOrder {
	array = revertWorkerOrder(array)
	var ordersClean []WorkerOrder
	for i := 0; i < len(array); i++ {
		if array[i].Quantity > array[i].AlreadyMade {
			ordersClean = append(ordersClean, array[i])
		}
	}
	ordersClean = removeDuplicatesWKOrders(ordersClean)
	return ordersClean
}

// This method cleans the array of ServiceOrder and leave only the orders that are not complete and the window period time is not started
func cleanServiceOrders(array []ServiceOrder) []ServiceOrder {
	array = revertServiceOrder(array)
	var ordersClean []ServiceOrder
	for i := 0; i < len(array); i++ {
		if array[i].Quantity > array[i].AlreadyMade && array[i].Deadline > (int64(array[i].WindowPeriod)*86400+getTime()) {
			ordersClean = append(ordersClean, array[i])
		}
	}
	return ordersClean
}

/*************************  REVERT ORDER  *************************/

// Reverts the order of the array Payment
func revertPayments(input []Payment) []Payment {
	if len(input) == 0 {
		return input
	}
	return append(revertPayments(input[1:]), input[0])
}

// Reverts the order of the array WorkerOrder
func revertWorkerOrder(input []WorkerOrder) []WorkerOrder {
	if len(input) == 0 {
		return input
	}
	return append(revertWorkerOrder(input[1:]), input[0])
}

// Reverts the order of the array ServiceOrder
func revertServiceOrder(input []ServiceOrder) []ServiceOrder {
	if len(input) == 0 {
		return input
	}
	return append(revertServiceOrder(input[1:]), input[0])
}

// Reverts the order of the array ServicePurchase
func revertServicePurchase(input []Stock) []Stock {
	if len(input) == 0 {
		return input
	}
	return append(revertServicePurchase(input[1:]), input[0])
}

// RevertQR the order of the array SALES
func revertPDF(input []PDF) []PDF {
	if len(input) == 0 {
		return input
	}
	return append(revertPDF(input[1:]), input[0])
}

// Reverts the order of the array WorkerOrder
func revertDeliveries(input []DeliveryHTML) []DeliveryHTML {
	if len(input) == 0 {
		return input
	}
	return append(revertDeliveries(input[1:]), input[0])
}

// This method give us all the items that his category is only for rent and haveaccess to the service provided
func (VS *Server) getAllItemsForRent(c *gin.Context, serviceTypeID bson.ObjectId) []ItemsPicker {
	allCategories := VS.getAllCategoriesOnlyForRent(c, "functionality", "getAllItemsForRent", "allCategories := VS.getAllCategoriesOnlyForRent")
	// We clean and take only the categories with access to the Item
	var categoriesWithItem []Category
	for _, category := range allCategories {
		if ContainsID(category.ServicesAccess, serviceTypeID) {
			categoriesWithItem = append(categoriesWithItem, category)
		}
	}
	// Now we take all the items from the categories
	var itemsHTML []ItemsPicker
	var item ItemsPicker
	for _, category := range categoriesWithItem {
		item.CategoryName = category.Name
		item.Items = VS.itemsToHTML(c, VS.getAllItemsByCategory(c, category.ID, "functionality", "getAllItemsForRent", "item.Stocks = VS.itemsToHTML(c, VS.getAllItemsByCategory"))
		itemsHTML = append(itemsHTML, item)
	}
	return itemsHTML
}

// This method give us all the items that his category is only for sell haveaccess to the service provided
func (VS *Server) getAllItemsForSell(c *gin.Context, serviceTypeID bson.ObjectId) []ItemsPicker {
	allCategories := VS.getAllCategoriesOnlyForSell(c, "functionality", "getAllItemsForSell", "allCategories := VS.getAllCategoriesOnlyForSell")
	// We clean and take only the categories with access to the Service
	var categoriesWithItem []Category
	for _, category := range allCategories {
		if ContainsID(category.ServicesAccess, serviceTypeID) {
			categoriesWithItem = append(categoriesWithItem, category)
		}
	}
	// Now we take all the items from the categories
	var itemsHTML []ItemsPicker
	var item ItemsPicker
	for _, category := range categoriesWithItem {
		item.CategoryName = category.Name
		item.Items = VS.itemsToHTML(c, VS.getAllItemsByCategory(c, category.ID, "functionality", "getAllItemsForSell", "item.Stocks = VS.itemsToHTML(c, VS.getAllItemsByCategory"))
		itemsHTML = append(itemsHTML, item)
	}
	return itemsHTML
}

// This method give us all the items that his category haveaccess to the service provided for a Service Purchase
func (VS *Server) getAllItemsForSellServicePurchase(c *gin.Context, serviceTyperID bson.ObjectId) []ItemHTML {
	allCategories := VS.getAllCategoriesOnlyForSell(c, "functionality", "getAllItemsForSellServicePurchase", "allCategories := VS.getAllCategoriesOnlyForSell")
	// We clean and take only the categories with access to the Service
	var categoriesWithItem []Category
	for _, category := range allCategories {
		if ContainsID(category.ServicesAccess, serviceTyperID) {
			categoriesWithItem = append(categoriesWithItem, category)
		}
	}
	// Now we take all the items from the categories
	var itemsHTML []ItemHTML
	var totalItemsHTML []ItemHTML
	for _, category := range categoriesWithItem {
		itemsHTML = VS.itemsToHTML(c, VS.getAllItemsByCategory(c, category.ID, "functionality", "getAllItemsForSellServicePurchase", "itemsHTML = VS.itemsToHTML(c, VS.getAllItemsByCategory"))
		totalItemsHTML = append(totalItemsHTML, itemsHTML...)
	}
	return totalItemsHTML
}

// Transform the Purchases information from DataBase to Human readable (in individual)
func (VS *Server) purchasesToWalletWorker(c *gin.Context, p Stock) Wallet {
	var cost float64
	prod, b := VS.getVideoCourseByID(c, p.IDVideoCourseOrItem, "h_functionality", "purchasesToWalletWorker", "prod, b := VS.getVideoCourseByID", true)
	var w Wallet
	w.ID = p.ID.Hex()
	w.Date = strconv.Itoa(time.Unix(p.Date, 0).Day()) + "/" + time.Unix(p.Date, 0).Month().String()
	if b {
		w.ItemPhoto = "/local/video-courses/" + prod.Name + "/purchases/" + "/" + p.Photo
		cost = prod.Price
	} else {
		item, _ := VS.getItemByID(c, p.IDVideoCourseOrItem, "h_functionality", "purchasesToWalletWorker", "item, _ := VS.getItemByID(c, p.", true)
		cat := VS.getCategoryByID(c, item.IDCategory, "h_functionality", "purchasesToWalletWorker", "cat := VS.getCategoryByID")
		tpe := VS.getCategoryTypeByID(c, cat.Type, "h_functionality", "purchasesToWalletWorker", "tpe := VS.getCategoryTypeByID")
		w.ItemPhoto = "/local/categories/" + tpe.Name + "/" + cat.Name + "/" + item.Photo
		cost = item.Price
	}
	w.Price = cost
	w.ItemName = prod.Name
	w.DateUnix = time.Unix(p.Date, 0).UnixNano()
	w.ServiceName = VS.getServiceByID(c, p.ServiceCreated, "h_functionality", "purchasesToWalletWorker", "w.ServiceName = VS.getServiceByID", true).Name
	return w
}

// Transform the Sales information from DataBase to Human readable (in individual)
func (VS *Server) saleToWalletWorker(c *gin.Context, s Sale) Wallet {
	mat, _ := VS.getItemByID(c, s.IDStock, "h_functionality", "saleToWalletWorker", "mat, _ := VS.getItemByID", true)
	var w Wallet
	w.ID = s.ID.Hex()
	w.Date = strconv.Itoa(time.Unix(s.Date, 0).Day()) + "/" + time.Unix(s.Date, 0).Month().String()
	cat := VS.getCategoryByID(c, mat.IDCategory, "h_functionality", "saleToWalletWorker", "cat := VS.getCategoryByID(c, mat")
	tpe := VS.getCategoryTypeByID(c, cat.Type, "h_functionality", "saleToWalletWorker", "tpe := VS.getCategoryTypeByID")
	w.ItemName = mat.Name
	w.ItemPhoto = "/local-resources/categories/" + tpe.Name + "/" + cat.Name + "/" + mat.Photo
	w.Price = -mat.Price
	w.ServiceName = VS.getServiceByID(c, s.IDService, "h_functionality", "saleToWalletWorker", "w.ServiceName = VS.getServiceByID", true).Name
	w.DateUnix = time.Unix(s.Date, 0).UnixNano()
	return w
}

// Transform the Payments information from DataBase to Human readable (in individual)
func (VS *Server) paymentToWalletWorker(c *gin.Context, p Payment) Wallet {
	var w Wallet
	w.ID = p.ID.Hex()
	w.Date = strconv.Itoa(time.Unix(p.Date, 0).Day()) + "/" + time.Unix(p.Date, 0).Month().String()
	w.ItemPhoto = "/static/svg/money-body.svg"
	w.Price = -p.Quantity
	w.DateUnix = time.Unix(p.Date, 0).UnixNano()
	w.ItemName = "Payment"
	w.ServiceName = VS.getServiceByID(c, p.IDService, "h_functionality", "paymentToWalletWorker", "w.ServiceName = VS.getServiceByID(c, p.", true).Name
	return w
}

// Give the Time as an int64, and depending if we are in development (return FakeTime) or in Service (return RealTime)
func getTime() int64 {
	if InDevelopment {
		// With the sum we fake than are advancing in the time
		return FakeDate + (time.Now().Unix() - TimeBeforeChanging)
	}
	return time.Now().Unix()
}

// Removes the elements duplicated in the array
func removeDuplicates(intSlice []bson.ObjectId) []bson.ObjectId {
	keys := make(map[bson.ObjectId]bool)
	list := []bson.ObjectId{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Removes the elements duplicated in the array
func removeDuplicatesWKOrders(intSlice []WorkerOrder) []WorkerOrder {
	keys := make(map[bson.ObjectId]bool)
	list := []WorkerOrder{}
	for _, entry := range intSlice {
		if _, value := keys[entry.IDProduct]; !value {
			keys[entry.IDProduct] = true
			list = append(list, entry)
		}
	}
	return list
}

//******* OLD METHOD, this is related with wkhtmltopdf but gives problems with docker, even tho the methods works in macOS perfectly with wkhtmltopdf installed  ************//
// Generates a PDF given the URL and stores it in the given PATH
func (VS *Server) generatePDFfromURL(c *gin.Context, URL string, directoryToSave string, id bson.ObjectId) {
	// Create new PDF generator
	pdfg, err := pdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(pdf.OrientationLandscape)

	// Create a new input page from an URL
	page := pdf.NewPage(URL)

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	name := id.Hex() + ".pdf"
	wholePath := basePath + directoryToSave + "/" + name
	err = pdfg.WriteFile(wholePath)
	if err != nil {
		os.MkdirAll(basePath+directoryToSave, os.ModePerm)
		err = pdfg.WriteFile(wholePath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Transform the PDF from the BD to make the PDF view for QR codes
func getPagePDFqr(QRstring []string, pdf PDF) PDFPageqrHTML {
	var page PDFPageqrHTML
	// Create and add the QR ID image
	qr, err := qrcode.New(pdf.ID.Hex(), qrcode.High)
	qr.ForegroundColor = Red
	qr.BackgroundColor = White
	img, err := qr.PNG(256)
	if err != nil {
		fmt.Print(err)
	}
	img2 := base64.StdEncoding.EncodeToString(img)
	page.ID = img2
	page.QRs = QRstring
	page.TypeOfItem = strings.ToUpper(string(pdf.TypeOfItem))
	y, m, d := time.Unix(pdf.Date, 0).Date()
	page.Date = time.Unix(pdf.Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
	return page
}

// Returns the list of services in HTML for the picker given the type
func (VS *Server) getServicesPickerByType(c *gin.Context, serTypeID bson.ObjectId) []ServiceHTML {
	role := getRole(c)
	var services []Service
	if role == Admin {
		services = VS.getAllSpecificServiceFromSystem(c, serTypeID, "h_functionality", "getServicesPickerByType", "services = VS.getAllSpecificServiceFromSystem")
	} else {
		services = VS.getAllSpecificServicesFromLocalVillage(c, serTypeID, "h_functionality", "getServicesPickerByType", "services = VS.getAllSpecificServicesFromLocalVillage")
	}
	servicesHTML := VS.servicesToHTML(c, services)

	return servicesHTML
}

// Returns the list of services in HTML for the picker
func (VS *Server) getAllServicesPicker(c *gin.Context) []ServiceHTML {
	role := getRole(c)
	var services []Service
	if role == Admin {
		services = VS.getAllServices(c, "h_functionality", "showRecordsWorkshop", "services = VS.getAllServices")
	} else {
		services = VS.getAllServicesFromLocalVillage(c, "h_functionality", "showRecordsWorkshop", "services = VS.getAllServicesFromLocalVillage")
	}
	servicesHTML := VS.servicesToHTML(c, services)
	return servicesHTML
}

// Returns the list of tools and materials for make the VideoCourse, taking in account if the service haveaccess to the category of items
// and avoiding the duplicates
func (VS *Server) getToolsAndMaterials(c *gin.Context, course VideoCourse) ([]ItemsPicker, []ItemsPicker) {
	var toolsTotal []ItemsPicker
	var materialsTotal []ItemsPicker
	for _, service := range course.ServicesAccess {
		tools := VS.getAllItemsForRent(c, service)
		for _, t := range tools {
			if !ContainsItemsPicker(toolsTotal, t) {
				toolsTotal = append(toolsTotal, t)
			}
		}
		materials := VS.getAllItemsForSell(c, service)
		for _, m := range materials {
			if !ContainsItemsPicker(materialsTotal, m) {
				materialsTotal = append(materialsTotal, m)
			}
		}
	}
	return toolsTotal, materialsTotal
}

// Stores the translation if the language is not the main one
func (VS *Server) isTranslation(c *gin.Context, i interface{}, id bson.ObjectId) bool {
	if getLanguage(c) != MainLanguage {
		trans, b := VS.getTranslationByIDInstanceAndLanguage(c, id, "h_functionality", "isTranslation", "trans := VS.getTranslationByIDInstanceAndLanguage")
		if !b {
			trans = Translation{
				ID:         bson.NewObjectId(),
				IDInstance: id,
				Instance:   i,
				Language:   getLanguage(c),
			}
		} else {
			trans.Instance = i
		}
		VS.addElementToDatabaseWithRegisterFromCentral(c, trans, trans.ID, "h_functionality", "isTranslation", "VS.addElementToDatabaseWithRegisterFromCentral")
		// VS.registerAuditAndEvent(c, VS.getAuditFromCentral(c, trans, id, Modified), trans, trans.ID)
		VS.goodFeedback(c, "/dashboard")
		return true
	}
	return false
}

// Returns the ID obtained from the QR, and takes care that is a new QR to use
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getQRFromHTMLNew(c *gin.Context, typeProduct TypeOfItem) (bson.ObjectId, bool) {
	QR, b := VS.getIDFromHTML(c, "qr", "QR")
	if !b {
		return bson.NewObjectId(), false
	}
	if !VS.isQRavailableForAssignment(c, typeProduct, QR, "h_functionality", "getQRFromHTMLNew", "if !VS.isQRavailableForAssignment") {
		VS.wrongFeedback(c, "Sorry but the QR is not for purchase or had been already used!")
		return bson.NewObjectId(), false
	}
	return QR, true
}

// Returns the ID obtained from the QR
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getQRFromHTML(c *gin.Context) (bson.ObjectId, bool) {
	QR, b := VS.getIDFromHTML(c, "qr", "QR")
	if !b {
		return bson.NewObjectId(), false
	}
	return QR, true
}

// Returns all the users related in the whole proccess of the stock to be added
func (VS *Server) getUsersRelatedFromStock(c *gin.Context, videoID bson.ObjectId) []bson.ObjectId {
	materials := VS.getAllMaterialsForCompleteVideoCourse(c, videoID)
	var users []bson.ObjectId
	for _, mat := range materials {
		users = append(users, VS.getStockByItemOrVideoID(c, mat, "h_functionality", "getUsersRelatedFromStock", "users = append(users, VS.getStockByItemOrVideoID").UsersRelated...)
	}
	return users
}

// Gives us the COMPLETE list of materials needed in the whole Videocourse, taking in account all the steps
// and avoiding the duplicates
func (VS *Server) getAllMaterialsForCompleteVideoCourse(c *gin.Context, videoID bson.ObjectId) []bson.ObjectId {
	steps := VS.getAllStepsFromSpecificVideoCourse(c, videoID, "h_functionality", "getAllMaterialsForCompleteVideoCourse", "steps := VS.getAllStepsFromSpecificVideoCourse")
	var materialsTotal []bson.ObjectId
	for _, step := range steps {
		for _, mat := range step.MaterialsNeeded {
			if !ContainsID(materialsTotal, mat) {
				materialsTotal = append(materialsTotal, mat)
			}
		}
	}
	return materialsTotal
}

// This method returns all the Users from the whole system if the user are admin and from the local village if the user are Manager
func (VS *Server) getAllUsersDependingRole(c *gin.Context, class string, method string, operation string) []User {
	role := getRole(c)
	var users []User
	if role == Admin {
		users = VS.getAllUsersFromTheSystemNotSelf(c, "h_functionality", "getAllUsersDependingRole", "users = VS.getAllUsersFromTheSystem")
	} else {
		users = VS.getAllUsersFromLocalVillage(c, "h_functionality", "getAllUsersDependingRole", "users = VS.getAllUsersFromLocalVillage")
	}
	return users
}

// This method returns all the Workers from the whole system if the user are admin and from the local village if the user are Manager
func (VS *Server) getAllWorkersDependingRole(c *gin.Context, class string, method string, operation string) []User {
	role := getRole(c)
	var users []User
	if role == Admin {
		users = VS.getAllWorkersFromTheSystem(c, "h_functionality", "getAllWorkersDependingRole", "users := v.getAllWorkersFromTheSystem")
	} else {
		users = VS.getAllWorkersFromLocalVillage(c, "h_functionality", "getAllWorkersDependingRole", "users = v.getAllWorkersFromLocalVillage")
	}
	return users
}

// This method iterates over the Attributes of an Struct, compares 2 items struct and saves the changes if there was changes
func (VS *Server) dynamicAuditChange(c *gin.Context, newObject interface{}, oldObject interface{}, auditID bson.ObjectId) bool {

	fmt.Print("For reflecting this is what I get")
	fmt.Print(newObject)

	modified := false
	var number int
	var date int64
	var boolean bool
	var arrayID []bson.ObjectId
	var arrayString []string

	// Check if both are items struct, and of the same type
	if reflect.TypeOf(newObject) != reflect.TypeOf(oldObject) {
		return false
	}
	if reflect.ValueOf(newObject).Kind() == reflect.Struct && reflect.ValueOf(oldObject).Kind() == reflect.Struct {
		for i := 0; i < reflect.ValueOf(newObject).NumField(); i++ {
			nameAtribute := reflect.TypeOf(newObject).Field(i).Name
			if reflect.TypeOf(newObject).Field(i).Type == reflect.TypeOf(bson.NewObjectId()) {
				fmt.Println("This is ID")
				// Is a BSON ID Object!
				valueNewAttribute, _ := reflect.ValueOf(newObject).Field(i).Interface().(bson.ObjectId)
				valueOldAttribute, _ := reflect.ValueOf(oldObject).Field(i).Interface().(bson.ObjectId)
				if VS.checkChange(c, auditID, nameAtribute, valueOldAttribute.Hex(), valueNewAttribute.Hex()) {
					modified = true
				}
			} else if reflect.TypeOf(newObject).Field(i).Type == reflect.TypeOf(number) {
				fmt.Println("This is FLOAT")
				// Is an float64
				valueNewAttribute, _ := reflect.ValueOf(newObject).Field(i).Interface().(float64)
				valueOldAttribute, _ := reflect.ValueOf(oldObject).Field(i).Interface().(float64)
				valueOld := strconv.FormatFloat(valueOldAttribute, 'f', -1, 64)
				valueNew := strconv.FormatFloat(valueNewAttribute, 'f', -1, 64)
				fmt.Println("This is new")
				fmt.Print(valueNew)
				fmt.Println("This is old")
				fmt.Print(valueOld)
				if VS.checkChange(c, auditID, nameAtribute, valueOld, valueNew) {
					modified = true
				}
			} else if reflect.TypeOf(newObject).Field(i).Type == reflect.TypeOf(number) {
				fmt.Println("This is INT")
				// Is an int
				valueNewAttribute, _ := reflect.ValueOf(newObject).Field(i).Interface().(int)
				valueOldAttribute, _ := reflect.ValueOf(oldObject).Field(i).Interface().(int)
				if VS.checkChange(c, auditID, nameAtribute, strconv.Itoa(valueOldAttribute), strconv.Itoa(valueNewAttribute)) {
					modified = true
				}
			} else if reflect.TypeOf(newObject).Field(i).Type == reflect.TypeOf(date) {
				fmt.Println("This is DATE")
				// Is an int64
				valueNewAttribute, _ := reflect.ValueOf(newObject).Field(i).Interface().(int64)
				valueOldAttribute, _ := reflect.ValueOf(oldObject).Field(i).Interface().(int64)
				if VS.checkChange(c, auditID, nameAtribute, strconv.FormatInt(valueOldAttribute, 10), strconv.FormatInt(valueNewAttribute, 10)) {
					modified = true
				}
			} else if reflect.TypeOf(newObject).Field(i).Type == reflect.TypeOf(boolean) {
				fmt.Println("This is bool")
				// Is a bool
				valueNewAttribute, _ := reflect.ValueOf(newObject).Field(i).Interface().(bool)
				valueOldAttribute, _ := reflect.ValueOf(oldObject).Field(i).Interface().(bool)
				if VS.checkChange(c, auditID, nameAtribute, strconv.FormatBool(valueOldAttribute), strconv.FormatBool(valueNewAttribute)) {
					modified = true
				}
			} else if reflect.TypeOf(newObject).Field(i).Type == reflect.TypeOf(arrayID) {
				fmt.Println("This is arrayID")
				// Is a IDArray, by default we update it
				valueAttribute, _ := reflect.ValueOf(newObject).Field(i).Interface().([]bson.ObjectId)
				for _, id := range valueAttribute {
					if VS.checkChange(c, auditID, nameAtribute, "ID Array", id.Hex()) {
						modified = true
					}
				}
			} else if reflect.TypeOf(newObject).Field(i).Type == reflect.TypeOf(arrayString) {
				// Is a ArrayString, by default we update it
				valueAttribute, _ := reflect.ValueOf(newObject).Field(i).Interface().([]string)
				for _, data := range valueAttribute {
					if VS.checkChange(c, auditID, nameAtribute, "String Array", data) {
						modified = true
					}
				}
			} else {
				valueNewAttribute := reflect.ValueOf(newObject).Field(i).String()
				valueOldAttribute := reflect.ValueOf(oldObject).Field(i).String()
				if VS.checkChange(c, auditID, nameAtribute, valueOldAttribute, valueNewAttribute) {
					modified = true
				}
			}
		}
	}
	return modified
}

// This method iterates over the Attributes of an Struct, and add all the information in the Audit
func (VS *Server) addInformationAuditItself(audit Audit, object interface{}) Audit {
	var number int
	var date int64
	var boolean bool
	var arrayID []bson.ObjectId
	var arrayString []string
	audit.InformationObject = make(map[string]string)
	for i := 0; i < reflect.ValueOf(object).NumField(); i++ {
		nameAtribute := reflect.TypeOf(object).Field(i).Name
		if reflect.TypeOf(object).Field(i).Type == reflect.TypeOf(bson.NewObjectId()) {
			// Is a BSON ID Object!
			valueAttribute, _ := reflect.ValueOf(object).Field(i).Interface().(bson.ObjectId)
			audit.InformationObject[nameAtribute] = valueAttribute.Hex()
		} else if reflect.TypeOf(object).Field(i).Type == reflect.TypeOf(number) {
			// Is an int
			valueAttribute, _ := reflect.ValueOf(object).Field(i).Interface().(int)
			audit.InformationObject[nameAtribute] = strconv.Itoa(valueAttribute)
		} else if reflect.TypeOf(object).Field(i).Type == reflect.TypeOf(date) {
			// Is an int
			valueAttribute, _ := reflect.ValueOf(object).Field(i).Interface().(int64)
			audit.InformationObject[nameAtribute] = strconv.FormatInt(valueAttribute, 10)
		} else if reflect.TypeOf(object).Field(i).Type == reflect.TypeOf(boolean) {
			// Is a bool
			valueAttribute, _ := reflect.ValueOf(object).Field(i).Interface().(bool)
			audit.InformationObject[nameAtribute] = strconv.FormatBool(valueAttribute)
		} else if reflect.TypeOf(object).Field(i).Type == reflect.TypeOf(arrayID) {
			// Is a IDArrays
			// valueAttribute, _ := reflect.ValueOf(object).Field(i).Interface().(arrayID)
			valueAttribute, _ := reflect.ValueOf(object).Field(i).Interface().([]bson.ObjectId)
			for _, id := range valueAttribute {
				audit.InformationObject[nameAtribute] = id.Hex()
			}

		} else if reflect.TypeOf(object).Field(i).Type == reflect.TypeOf(arrayString) {
			// Is a ArrayString
			valueAttribute, _ := reflect.ValueOf(object).Field(i).Interface().([]string)
			for _, data := range valueAttribute {
				audit.InformationObject[nameAtribute] = data
			}
		} else {
			valueAttribute := reflect.ValueOf(object).Field(i).String()
			audit.InformationObject[nameAtribute] = valueAttribute
		}
	}
	return audit
}

// This method create the global variables inside the HTML templates
func getGobalTemplateVariables(c *gin.Context) map[string]string {
	m := make(map[string]string)
	m["language"] = string(getLanguage(c))
	m["role"] = string(getRole(c))
	m["itemspertable"] = strconv.Itoa(itemsPerTable)
	m["URL"] = c.Request.URL.String()
	m["object"] = "Object empty"
	return m
}

// This method returns the translation if available
func (VS *Server) mayLookForTranslation(c *gin.Context, id bson.ObjectId) (interface{}, bool) {
	if getLanguage(c) != MainLanguage {
		trans, b := VS.getTranslationByIDInstanceAndLanguage(c, id, "h_functionality", "mayLookForTranslation", "trans, b := VS.getTranslationByIDInstanceAndLanguage")
		if b {
			return trans.Instance, true
		}
		return nil, false
	}
	return nil, false
}

// Returns the Materials and Tools for HTML representation from a step given
func (VS *Server) getToolsAndMaterialsHTMLFromStep(c *gin.Context, step Step) ([]ItemHTML, []ItemHTML) {
	//We collect the tools and materials from the array of the step
	var tools []Item
	for i := 0; i < len(step.ToolsNeeded); i++ {
		toolID := step.ToolsNeeded[i]
		tool, _ := VS.getItemByID(c, toolID, "h_functionality", "getToolsAndMaterialsHTMLFromStep", "tool := v.getToolByID", true)
		tools = append(tools, tool)
	}
	toolsHTML := VS.itemsToHTML(c, tools)
	var materials []Item
	for i := 0; i < len(step.MaterialsNeeded); i++ {
		materialID := step.MaterialsNeeded[i]
		material, _ := VS.getItemByID(c, materialID, "h_functionality", "getToolsAndMaterialsHTMLFromStep", "material := v.getMaterialByID", true)
		materials = append(materials, material)
	}
	materialsHTML := VS.itemsToHTML(c, materials)
	return toolsHTML, materialsHTML
}

// Returns and array of Wallet from the Specific worker
func (VS *Server) getWalletSpecificWorker(c *gin.Context, userID bson.ObjectId) []Wallet {
	sales := VS.getAllSalesFromSpecificWorker(c, userID, "h_functionality", "getWalletSpecificWorker", "sales := v.getAllSalesFromSpecificWorker")
	payments := VS.getAllPaymentsFromSpecificWorker(c, userID, "h_functionality", "getWalletSpecificWorker", "payments := v.getAllPaymentsFromSpecificWorker")
	stocks := VS.getAllServicePurchasesFromSpecificWorker(c, userID, "h_functionality", "getWalletSpecificWorker", "purchases := v.getAllWorkshopPurchasesFromSpecificWorker")
	fmt.Print("sales is: ")
	fmt.Print(len(sales))
	fmt.Print("payments is: ")
	fmt.Print(len(payments))
	fmt.Print("stocks is: ")
	fmt.Print(len(stocks))
	return VS.getWallet(c, sales, payments, stocks)
}

// Produces the Wallet information
func (VS *Server) getWallet(c *gin.Context, sales []Sale, payments []Payment, stocks []Stock) []Wallet {
	var wallet []Wallet
	for _, sale := range sales {
		wallet = append(wallet, VS.saleToWalletWorker(c, sale))
	}
	for _, pay := range payments {
		wallet = append(wallet, VS.paymentToWalletWorker(c, pay))
	}
	for _, stock := range stocks {
		wallet = append(wallet, VS.purchasesToWalletWorker(c, stock))
	}
	walletSort(wallet)
	// We put the balance
	if len(wallet) != 0 {
		// Putting the balance
		if (len(wallet) - 1) >= 0 {
			wallet[len(wallet)-1].Balance = wallet[len(wallet)-1].Price
		}
		for i := len(wallet) - 1; 0 < i; i-- {
			wallet[i-1].Balance = wallet[i].Balance + wallet[i-1].Price
		}
	}
	return wallet
}

// Takes the Date from the DatePicker in the HTML
func (VS *Server) getDateFromHTML(c *gin.Context, name string, element string) int64 {
	deadline, b := VS.getStringFromHTML(c, name, element)
	if !b {
		return -1
	}
	day, _ := strconv.Atoi(deadline[:2])
	month, _ := strconv.Atoi(deadline[3:5])
	year, _ := strconv.Atoi(deadline[6:])
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Unix()
	return date
}

// This method returns all the Audits that needs to be synchronized with the village
func (VS *Server) getAllObjectsFromAudits(c *gin.Context, audits []Audit) []Anything {
	var objects []Anything
	for _, a := range audits {
		objects = append(objects, VS.getAnythingWithID(c, a.IDItem, "h_functionality", "getAllObjetsFromAuditsLastSynchronization", "objects = append(objects, VS.getAnythingWithID"))
	}
	return objects
}

// This method returns the DataBase that the user will download to synchronize the villages with the central DataBase
// newSynchronizationID is the Id of the new synchronization we are making, we pass the parameter to avoid to have the lastestSync date from the one we are just creating
// Will only create a new DataBase if there is something new to add
func (VS *Server) createDataBaseForSync(c *gin.Context, villageReceiver Village, newSynchronizationID bson.ObjectId) (Server, string, bool) {
	audits := VS.getAllAuditsFromLastSynchronization(c, villageReceiver.ID, newSynchronizationID, "h_functionality", "createDataBaseForSync", "audits := VS.getAllAuditsFromLastSynchronization")
	if len(audits) == 0 {
		VS.wrongFeedback(c, "There is no new audits to sync")
		return Server{}, "", false
	}
	events := VS.getAllEventsFromLastSynchronization(c, villageReceiver.ID, newSynchronizationID, "h_functionality", "createDataBaseForSync", "events := VS.getAllEventsFromLastSynchronization")
	// We need to add the audit of the sync we are making
	auditSync := VS.getAuditByIDItemID(c, newSynchronizationID, "h_functionality", "createDataBaseForSync", "auditSync := VS.getAuditByIDItemID")
	audits = append(audits, auditSync)
	objects := VS.getAllObjectsFromAudits(c, audits)

	// We open or create the DataBase and the Directory
	db, err := bolthold.Open(basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/"+newSynchronizationID.Hex()+".db", 0666, nil)
	if err != nil {
		os.MkdirAll(basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/", os.ModePerm)
		db, err = bolthold.Open(basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/"+newSynchronizationID.Hex()+".db", 0666, nil)
	}
	DB := Server{
		DataBase: db,
	}
	// Copy audits of the elements
	for _, a := range audits {
		DB.addElementToDatabaseWithoutRegister(c, a, a.ID, "h_functionality", "createDataBaseForSync", "DB.addElementToDatabaseWithoutRegister(c, a, a.ID")
	}
	// Copy events of the elements
	for _, e := range events {
		DB.addElementToDatabaseWithoutRegister(c, e, e.ID, "h_functionality", "createDataBaseForSync", "DB.addElementToDatabaseWithoutRegister(c, e, e.ID")
	}
	fmt.Print("Audits copied: ")
	fmt.Println(len(audits))
	fmt.Print("Events copied: ")
	fmt.Println(len(events))
	fmt.Print("Objects copied: ")
	fmt.Println(len(objects))

	// Copy elements
	for _, o := range objects {
		DB.addElementToDatabaseWithoutRegister(c, o.Instance, o.ID, "h_functionality", "createDataBaseForSync", "DB.addElementToDatabaseWithoutRegister(c, e, e.ID")

		// Here we move all the media: Photos, audios and videos
		switch o.Instance.(type) {

		case ServiceType:
			s := o.Instance.(ServiceType)
			copyFile(basePath+"/local-resources/services/"+s.Icon, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/services/"+s.Icon)

		case User:
			u := o.Instance.(User)
			fmt.Println("User executed")
			copyFile(basePath+"/local-resources/users/"+u.Username+"/"+u.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/users/"+u.Username+"/"+u.Photo)

		case Category:
			cat := o.Instance.(Category)
			fmt.Println("Category executed")
			catType := VS.getCategoryTypeByID(c, cat.Type, "h_functionality", "createDataBaseForSync", "cat.Type = VS.getCategoryTypeByID(").Name
			copyFile(basePath+"/local-resources/categories/"+catType+"/"+cat.Name+"/"+cat.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/categories/"+catType+"/"+cat.Name+"/"+cat.Photo)

		case Item:
			i := o.Instance.(Item)
			fmt.Println("Item executed")
			prod, isProd := VS.getVideoCourseByID(c, i.ID, "h_functionality", "createDataBaseForSync", "prod, isProd := VS.getVideoCourseByID", false)
			if isProd {
				copyFile(basePath+"/local-resources/video-courses/"+prod.Name+"/purchases/"+"/"+prod.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+prod.Name+"/purchases/"+prod.Photo)
			} else {
				cat := VS.getCategoryByID(c, i.IDCategory, "h_functionality", "createDataBaseForSync", "cat := VS.getCategoryByID")
				tpe := VS.getCategoryTypeByID(c, cat.Type, "h_functionality", "createDataBaseForSync", "tpe := VS.getCategoryTypeByID")
				copyFile(basePath+"/local-resources/categories/"+tpe.Name+"/"+cat.Name+"/"+i.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/categories/"+tpe.Name+"/"+cat.Name+"/"+i.Photo)
			}

		case Stock:
			stock := o.Instance.(Stock)
			fmt.Println("Stock executed")
			prod, b := VS.getVideoCourseByID(c, stock.IDVideoCourseOrItem, "h_functionality", "createDataBaseForSync", "prod, b := VS.getVideoCourseByID", true)
			if b {
				stock.Photo = "/local-resources/video-courses/" + prod.Name + "/purchases/" + stock.Photo
				// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+prod.Name+"/purchases/", os.ModePerm)
				copyFile(basePath+"/local-resources/video-courses/"+prod.Name+"/purchases/"+stock.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+prod.Name+"/purchases/"+stock.Photo)
			}

		case Payment:
			p := o.Instance.(Payment)
			fmt.Println("Payment executed")
			user := VS.getUserByID(c, p.IDWorker, "h_functionality", "createDataBaseForSync", "user := VS.getUserByID")
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/users/"+user.Username+"/payments/", os.ModePerm)
			copyFile(basePath+"/local-resources/users/"+user.Username+"/payments/"+p.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/users/"+user.Username+"/payments/"+p.Photo)

		case VideoCourse:
			v := o.Instance.(VideoCourse)
			fmt.Println("VideoCourse executed")
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+v.Name+"/", os.ModePerm)
			copyFile(basePath+"/local-resources/video-courses/"+v.Name+"/"+v.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+v.Name+"/"+v.Photo)

		case VideoProblem:
			v := o.Instance.(VideoProblem)
			fmt.Println("VideoProblem executed")
			video, _ := VS.getVideoCourseByID(c, v.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video := VS.getVideoCourseByID", true)
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/problems/photos/", os.ModePerm)
			copyFile(basePath+"/local-resources/video-courses/"+video.Name+"/problems/photos/"+v.PhotoProblem, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/problems/photos/"+v.PhotoProblem)
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/problems/audios/", os.ModePerm)
			copyFile(basePath+"/local-resources/video-courses/"+video.Name+"/problems/audios/"+v.Audio, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/problems/audios/"+v.Audio)

		case Step:
			s := o.Instance.(Step)
			fmt.Println("Step executed")
			video, _ := VS.getVideoCourseByID(c, s.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video := VS.getVideoCourseByID", true)
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/steps/", os.ModePerm)
			copyFile(basePath+"/local-resources/video-courses/"+video.Name+"/steps/"+s.Video, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/steps/"+s.Video)
			copyFile(basePath+"/local-resources/video-courses/"+video.Name+"/steps/"+s.Audio, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/steps/"+s.Audio)

		case VideoCheckList:
			v := o.Instance.(VideoCheckList)
			fmt.Println("VideoCheckList executed")
			video, _ := VS.getVideoCourseByID(c, v.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video := VS.getVideoCourseByID", true)
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/checks/", os.ModePerm)
			copyFile(basePath+"/local-resources/video-courses/"+video.Name+"/checks/"+v.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/checks/"+v.Photo)
			copyFile(basePath+"/local-resources/video-courses/"+video.Name+"/checks/"+v.Audio, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/video-courses/"+video.Name+"/checks/"+v.Audio)

		case Message:
			m := o.Instance.(Message)
			fmt.Println("Message executed")
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/conversations/"+m.IDConversation.Hex()+"/", os.ModePerm)
			copyFile(basePath+"/local-resources/conversations/"+m.IDConversation.Hex()+"/"+m.Audio, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/conversations/"+m.IDConversation.Hex()+"/"+m.Audio)
			copyFile(basePath+"/local-resources/conversations/"+m.IDConversation.Hex()+"/"+m.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/conversations/"+m.IDConversation.Hex()+"/"+m.Photo)

		case Report:
			r := o.Instance.(Report)
			fmt.Println("Report executed")
			user := VS.getUserByID(c, r.IDUser, "h_functionality", "createDataBaseForSync", "user := VS.getUserByID")
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/users/"+user.Username+"/reports/", os.ModePerm)
			copyFile(basePath+"/local-resources/users/"+user.Username+"/reports/"+r.Audio, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/users/"+user.Username+"/reports/"+r.Audio)
			copyFile(basePath+"/local-resources/users/"+user.Username+"/reports/"+r.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/users/"+user.Username+"/reports/"+r.Photo)

		case Animal:
			a := o.Instance.(Animal)
			fmt.Println("Animal executed")
			os.MkdirAll(basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/slaughterhouse/animals/"+a.ID.Hex()+"/", os.ModePerm)
			copyDirectory(basePath+"/local-resources/slaughterhouse/animals/"+a.ID.Hex()+"/", basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/slaughterhouse/animals/"+a.ID.Hex()+"/")
			// copy.Copy(basePath+"/local-resources/slaughterhouse/animals/"+a.ID.Hex()+"/"+a.Front, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/slaughterhouse/animals/"+a.ID.Hex()+"/")
			// copy.Copy(basePath+"/local-resources/slaughterhouse/animals/"+a.ID.Hex()+"/"+a.Right, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/slaughterhouse/animals/"+a.ID.Hex()+"/")
			// copy.Copy(basePath+"/local-resources/slaughterhouse/animals/"+a.ID.Hex()+"/"+a.Left, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/slaughterhouse/animals/"+a.ID.Hex()+"/")
			// copy.Copy(basePath+"/local-resources/slaughterhouse/animals/"+a.ID.Hex()+"/"+a.Back, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/slaughterhouse/animals/"+a.ID.Hex()+"/")

		case FieldPhoto:
			fP := o.Instance.(FieldPhoto)
			fmt.Println("FieldPhoto executed")
			toDo := VS.getTodoByID(c, fP.IDToDo, "h_functionality", "createDataBaseForSync", "toDo := VS.getTodoByID")
			user := VS.getUserByID(c, toDo.IDUser, "h_functionality", "createDataBaseForSync", "user := VS.getUserByID")
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/users/"+user.Username+"/to-dos/", os.ModePerm)
			copyFile(basePath+"/local-resources/users/"+user.Username+"/to-dos/"+fP.Photo, basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/users/"+user.Username+"/to-dos/"+fP.Photo)

		case CheckStock:
			cS := o.Instance.(CheckStock)
			fmt.Println("CheckStock executed")
			service := VS.getServiceByID(c, cS.IDService, "h_functionality", "createDataBaseForSync", "service := VS.getServiceByID", true)
			// os.MkdirAll(basePath+"/local-resources/syncs/exports/"+village.Name+"/"+newSynchronizationID.Hex()+"/services/"+service.Name+"/checks/", os.ModePerm)
			copyDirectory(basePath+"/local-resources/services/"+service.Name+"/checks/", basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/services/"+service.Name+"/checks/")
			// copy.Copy(basePath+"/local-resources/services/"+service.Name+"/checks/"+cS.Photo1, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/services/"+service.Name+"/checks/")
			// copy.Copy(basePath+"/local-resources/services/"+service.Name+"/checks/"+cS.Photo2, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/services/"+service.Name+"/checks/")
			// copy.Copy(basePath+"/local-resources/services/"+service.Name+"/checks/"+cS.Photo3, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/services/"+service.Name+"/checks/")
			// copy.Copy(basePath+"/local-resources/services/"+service.Name+"/checks/"+cS.Photo4, basePath+"/local-resources/syncs/"+village.Name+"/"+newSynchronizationID.Hex()+"/services/"+service.Name+"/checks/")
		}
	}
	zipDirectory(basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/"+newSynchronizationID.Hex()+"/", basePath+"/local-resources/syncs/exports/"+villageReceiver.Name+"/", "To-"+villageReceiver.Name+"-"+newSynchronizationID.Hex(), "/"+newSynchronizationID.Hex()+"/")
	return DB, "/local/syncs/" + villageReceiver.Name + "/" + newSynchronizationID.Hex() + ".db", true
}

// This method imports all the objects, audits and events from the database imported to the local Database
func (VS *Server) importElementsFromDatabases(c *gin.Context, syncID bson.ObjectId) bool {

	// We open or create the DataBase and the Directory
	db, err := bolthold.Open(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/"+syncID.Hex()+".db", 0666, nil)
	if err != nil {
		VS.wrongFeedback(c, "There have been a problem trying to open the import DataBase")
	}
	importedDB := Server{
		DataBase: db,
	}
	// Check for the village receiver in the DB with the IDVillageEmitter of the Sync
	sync, _ := importedDB.getSynchronizationByID(c, syncID, "h_functionality.go", "importElementsFromDatabases", "events := importedDB.getAllEvents", true)
	vil := VS.getLocalVillage(c, "h_settings", "newSync", "VS.getLocalVillage")
	if vil.ID != sync.IDVillageReceiver {
		VS.wrongFeedback(c, "This Synchronization is not for this village on the IDVillageReceiver")
		fmt.Print("Local village id is: ")
		fmt.Println(vil.ID.Hex())
		fmt.Print("Sync is: ")
		fmt.Println(sync)
		return false
	}
	audits := importedDB.getAllAudits(c, "h_functionality.go", "importElementsFromDatabases", "audits := importedDB.getAllAudits")
	events := importedDB.getAllEvents(c, "h_functionality.go", "importElementsFromDatabases", "events := importedDB.getAllEvents")
	objects := importedDB.getAllObjectsFromAudits(c, audits)

	syncUpdated := sync
	syncUpdated.IsDone = true
	syncUpdated.DateSyncReceiver = getTime()
	var villageEmitter Village
	var villageEmpty Village
	VS.DataBase.Get(sync.IDVillageEmitter, &villageEmitter)
	if villageEmitter == villageEmpty {
		fmt.Print("Is looking for the village in the imported DataBase")
		villageEmitter = importedDB.getVillageByID(c, sync.IDVillageEmitter, "h_functionality.go", "importElementsFromDatabases", "events := importedDB.getAllEvents")
	}

	fmt.Print("Village name is: " + villageEmitter.Name)

	fmt.Print("Audits to copy: ")
	fmt.Println(len(audits))
	fmt.Print("Events to copy: ")
	fmt.Println(len(events))
	fmt.Print("Objects to copy: ")
	fmt.Println(len(objects))

	// Copy audits of the elements
	for _, a := range audits {
		VS.addElementToDatabaseWithoutRegister(c, a, a.ID, "h_functionality.go", "importElementsFromDatabases", "VS.addElementToDatabaseWithoutRegister")
	}
	// Copy events of the elements
	for _, e := range events {
		VS.addElementToDatabaseWithoutRegister(c, e, e.ID, "h_functionality.go", "importElementsFromDatabases", "VS.addElementToDatabaseWithoutRegister")
	}

	// Copy elements
	// HERE WE NEED TO USE METHODS OF THE DB IN ROW
	// Because will look for elements, that if it does not find, will look them in the another database, and we need
	// to avoid the error messages of not getting the elements in the database, that is why we use raw methods
	// and not the ones created with the error messages
	for _, o := range objects {
		VS.addElementToDatabaseWithoutRegister(c, o.Instance, o.ID, "h_functionality.go", "importElementsFromDatabases", "VS.addElementToDatabaseWithoutRegister")

		// Here we move all the media: Photos, audios and videos
		switch o.Instance.(type) {

		case ServiceType:
			s := o.Instance.(ServiceType)
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/services/"+s.Icon, basePath+"/local-resources/services/"+s.Icon)

		case User:
			u := o.Instance.(User)
			fmt.Println("User executed")
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/users/"+u.Username+"/"+u.Photo, basePath+"/local-resources/users/"+u.Username+"/"+u.Photo)

		case Category:
			cat := o.Instance.(Category)
			fmt.Println("Category executed")
			var catType CategoryType
			var catTypeEmpty CategoryType
			VS.DataBase.Get(cat.Type, &catType)
			if catType == catTypeEmpty {
				fmt.Print("Is looking for the CategoryType in the imported DataBase")
				importedDB.DataBase.Get(cat.Type, &catType)
			}
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/categories/"+catType.Name+"/"+cat.Name+"/"+cat.Photo, basePath+"/local-resources/categories/"+catType.Name+"/"+cat.Name+"/"+cat.Photo)

		case Item:
			i := o.Instance.(Item)
			fmt.Println("Item executed")
			var videoEmpty VideoCourse
			video, isProd := VS.getVideoCourseByID(c, i.ID, "h_functionality", "createDataBaseForSync", "video, isProd := VS.getVideoCourseByID", false)
			if video.ID == videoEmpty.ID {
				fmt.Print("Is looking for the VideoCourse in the imported DataBase")
				video, isProd = importedDB.getVideoCourseByID(c, i.ID, "h_functionality", "createDataBaseForSync", "video, isProd = importedDB.getVideoCourseByID", false)
			}

			if isProd {
				copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/purchases/"+video.Photo, basePath+"/local-resources/video-courses/"+video.Name+"/purchases/"+"/"+video.Photo)
			} else {
				var cat Category
				var catEmpty Category
				VS.DataBase.Get(i.IDCategory, &cat)
				if cat.ID == catEmpty.ID {
					fmt.Print("Is looking for the Category in the imported DataBase")
					cat = importedDB.getCategoryByID(c, i.IDCategory, "h_functionality", "createDataBaseForSync", "cat := VS.getCategoryByID")
				}
				var catType CategoryType
				var catTypeEmpty CategoryType
				VS.DataBase.Get(cat.Type, &catType)
				if catType.ID == catTypeEmpty.ID {
					fmt.Print("Is looking for the CategoryType Type in the imported DataBase")
					catType = importedDB.getCategoryTypeByID(c, cat.Type, "h_functionality", "createDataBaseForSync", "tpe := VS.getCategoryTypeByID")
				}
				copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/categories/"+catType.Name+"/"+cat.Name+"/"+i.Photo, basePath+"/local-resources/categories/"+catType.Name+"/"+cat.Name+"/"+i.Photo)
			}

		case Stock:
			stock := o.Instance.(Stock)
			fmt.Println("Stock executed")
			var videoEmpty VideoCourse
			video, b := VS.getVideoCourseByID(c, stock.IDVideoCourseOrItem, "h_functionality", "createDataBaseForSync", "video, b := VS.getVideoCourseByID", false)
			if video.ID == videoEmpty.ID {
				fmt.Print("Is looking for the VideoCourse in the imported DataBase")
				video, b = importedDB.getVideoCourseByID(c, stock.IDVideoCourseOrItem, "h_functionality", "createDataBaseForSync", "video, b = importedDB.getVideoCourseByID", false)
			}
			if b {
				stock.Photo = "/local-resources/video-courses/" + video.Name + "/purchases/" + stock.Photo
				copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/purchases/"+stock.Photo, basePath+"/local-resources/video-courses/"+video.Name+"/purchases/"+stock.Photo)
			}

		case Payment:
			p := o.Instance.(Payment)
			fmt.Println("Payment executed")
			var user User
			var userEmpty User
			VS.DataBase.Get(p.IDWorker, &user)
			if user.ID == userEmpty.ID {
				fmt.Print("Is looking for the VideoCourse in the imported DataBase")
				user = importedDB.getUserByID(c, p.IDWorker, "h_functionality", "createDataBaseForSync", "user = importedDB.getUserByID")
			}
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/users/"+user.Username+"/payments/"+p.Photo, basePath+"/local-resources/users/"+user.Username+"/payments/"+p.Photo)

		case VideoCourse:
			v := o.Instance.(VideoCourse)
			fmt.Println("VideoCourse executed")
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+v.Name+"/"+v.Photo, basePath+"/local-resources/video-courses/"+v.Name+"/"+v.Photo)

		case VideoProblem:
			v := o.Instance.(VideoProblem)
			fmt.Println("VideoProblem executed")
			video, _ := VS.getVideoCourseByID(c, v.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video := VS.getVideoCourseByID", false)
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/problems/photos/"+v.PhotoProblem, basePath+"/local-resources/video-courses/"+video.Name+"/problems/photos/"+v.PhotoProblem)
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/problems/audios/"+v.Audio, basePath+"/local-resources/video-courses/"+video.Name+"/problems/audios/"+v.Audio)

		case Step:
			s := o.Instance.(Step)
			fmt.Println("Step executed")
			video, _ := VS.getVideoCourseByID(c, s.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video := VS.getVideoCourseByID", true)
			var videoEmpty VideoCourse
			if video.ID == videoEmpty.ID {
				fmt.Print("Is looking for the VideoCourse in the imported DataBase")
				video, _ = importedDB.getVideoCourseByID(c, s.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video, b = importedDB.getVideoCourseByID", false)
			}
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/steps/"+s.Video, basePath+"/local-resources/video-courses/"+video.Name+"/steps/"+s.Video)
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/steps/"+s.Audio, basePath+"/local-resources/video-courses/"+video.Name+"/steps/"+s.Audio)

		case VideoCheckList:
			v := o.Instance.(VideoCheckList)
			fmt.Println("VideoCheckList executed")
			video, _ := VS.getVideoCourseByID(c, v.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video := VS.getVideoCourseByID", true)
			var videoEmpty VideoCourse
			if video.ID == videoEmpty.ID {
				fmt.Print("Is looking for the VideoCourse in the imported DataBase")
				video, _ = importedDB.getVideoCourseByID(c, v.IDVideoCourse, "h_functionality", "createDataBaseForSync", "video, b = importedDB.getVideoCourseByID", false)
			}
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/checks/"+v.Photo, basePath+"/local-resources/video-courses/"+video.Name+"/checks/"+v.Photo)
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/video-courses/"+video.Name+"/checks/"+v.Audio, basePath+"/local-resources/video-courses/"+video.Name+"/checks/"+v.Audio)

		case Message:
			m := o.Instance.(Message)
			fmt.Println("Message executed")
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/conversations/"+m.IDConversation.Hex()+"/"+m.Audio, basePath+"/local-resources/conversations/"+m.IDConversation.Hex()+"/"+m.Audio)
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/conversations/"+m.IDConversation.Hex()+"/"+m.Photo, basePath+"/local-resources/conversations/"+m.IDConversation.Hex()+"/"+m.Photo)

		case Report:
			r := o.Instance.(Report)
			fmt.Println("Report executed")
			var user User
			var userEmpty User
			VS.DataBase.Get(r.IDUser, &user)
			if user.ID == userEmpty.ID {
				fmt.Print("Is looking for the User in the imported DataBase")
				user = importedDB.getUserByID(c, r.IDUser, "h_functionality", "createDataBaseForSync", "user := VS.getUserByID")
			}
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/users/"+user.Username+"/reports/"+r.Audio, basePath+"/local-resources/users/"+user.Username+"/reports/"+r.Audio)
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/users/"+user.Username+"/reports/"+r.Photo, basePath+"/local-resources/users/"+user.Username+"/reports/"+r.Photo)

		case Animal:
			a := o.Instance.(Animal)
			fmt.Println("Animal executed")
			copyDirectory(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/slaughterhouse/animals/"+a.ID.Hex()+"/", basePath+"/local-resources/slaughterhouse/animals/"+a.ID.Hex()+"/")

		case FieldPhoto:
			fP := o.Instance.(FieldPhoto)
			fmt.Println("FieldPhoto executed")
			var toDo ToDo
			var toDoEmpty ToDo
			VS.DataBase.Get(fP.IDToDo, &toDo)
			if toDo.ID == toDoEmpty.ID {
				fmt.Print("Is looking for the toDo in the imported DataBase")
				toDo = importedDB.getTodoByID(c, fP.IDToDo, "h_functionality", "createDataBaseForSync", "toDo := VS.getTodoByID")
			}
			var user User
			var userEmpty User
			VS.DataBase.Get(toDo.IDUser, &user)
			if user.ID == userEmpty.ID {
				fmt.Print("Is looking for the User in the imported DataBase")
				user = importedDB.getUserByID(c, toDo.IDUser, "h_functionality", "createDataBaseForSync", "user := VS.getUserByID")
			}
			copyFile(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/users/"+user.Username+"/to-dos/"+fP.Photo, basePath+"/local-resources/users/"+user.Username+"/to-dos/"+fP.Photo)

		case CheckStock:
			cS := o.Instance.(CheckStock)
			fmt.Println("CheckStock executed")
			var service Service
			var serviceEmpty Service
			VS.DataBase.Get(cS.IDService, &service)
			if service.ID == serviceEmpty.ID {
				fmt.Print("Is looking for the CheckStock in the imported DataBase")
				service = importedDB.getServiceByID(c, cS.IDService, "h_functionality", "createDataBaseForSync", "service := VS.getServiceByID", true)
			}
			copyDirectory(basePath+"/local-resources/syncs/imports/"+syncID.Hex()+"/services/"+service.Name+"/checks/", basePath+"/local-resources/services/"+service.Name+"/checks/")
		}
	}

	audits = VS.getAllAudits(c, "h_functionality.go", "importElementsFromDatabases", "audits := importedDB.getAllAudits")
	events = VS.getAllEvents(c, "h_functionality.go", "importElementsFromDatabases", "events := importedDB.getAllEvents")
	objects = VS.getAllObjectsFromAudits(c, audits)
	for _, o := range objects {
		shortDescriptionInterface(o)
	}
	villages := VS.getAllVillagesFromTheSystem(c, "h_functionality.go", "importElementsFromDatabases", "audits := importedDB.getAllAudits")

	fmt.Print("Audits copied: ")
	fmt.Println(len(audits))
	fmt.Print("Events copied: ")
	fmt.Println(len(events))
	fmt.Print("Objects copied: ")
	fmt.Println(len(objects))
	for _, o := range objects {
		shortDescriptionInterface(o)
	}
	fmt.Print("Villages copied: ")
	fmt.Println(len(villages))

	// Once all the elements are imported, we make the update to the sync
	VS.editElementInDatabaseWithRegister(c, syncUpdated, sync, sync.ID, VS.getLocalVillage(c, "h_functionality.go", "importElementsFromDatabases", "events := importedDB.getAllEvents").ID, "h_functionality.go", "importElementsFromDatabases", "events := importedDB.getAllEvents")
	return true
}

// This method returns all the Audits that needs to be synchronized with the villageEmitter
func (VS *Server) testAudio(c *gin.Context) {
	render(c, gin.H{}, "audio2.html")
}

// This method returns all the Audits that needs to be synchronized with the villageEmitter
func (VS *Server) testPDF(c *gin.Context) {

	// var pathsIMG []string
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for l := 0; l < 5; l++ {
		pdf.AddPage()
		x, y := 20.0, 10.0
		for i, k := 0, 0; i < 77; i++ {
			// Create and add the QR ID image
			qrID := bson.NewObjectId().Hex()
			qr, err := qrcode.New(qrID, qrcode.High)
			if err != nil {
				fmt.Print(err)
			}
			qr.ForegroundColor = Red
			qr.BackgroundColor = White
			img, err := qr.PNG(140)
			if err != nil {
				fmt.Print(err)
			}
			// convert []byte to image for saving to file
			image, _, _ := image.Decode(bytes.NewReader(img))
			//save the imgByte to file
			os.Mkdir(basePath+"/local-resources/pdf/test/", os.ModePerm)
			out, err := os.Create(basePath + "/local-resources/pdf/test/" + qrID + ".png")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = png.Encode(out, image)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// pathsIMG = append(pathsIMG, basePath+"/local-resources/pdf/test"+qrID+".png")
			pdf.Image(basePath+"/local-resources/pdf/test/"+qrID+".png", x, y, nil) //print image
			x += 75
			k++
			if k >= 7 {
				y += 75
				x = 20
				k = 0
			}
		}
	}

	pdf.WritePdf(basePath + "/local-resources/pdf/qr/" + bson.NewObjectId().Hex() + ".pdf")
	fmt.Print("PDF done")
	// Delete the images of the QR that we created before
	os.RemoveAll(basePath + "/local-resources/pdf/test/")
	// This is for set specific typography in the PDF
	// err = pdf.AddTTFFont("Nunito", "../ttf/Loma.ttf")
	// err = pdf.SetFont("Nunito", "", 14)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	return
	// }
}

// This method returns all the Audits that needs to be synchronized with the villageEmitter
func (VS *Server) newAudio(c *gin.Context) {
}
