package village

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

/*************************  GLOBAL  *************************/

// This handler shows the different deliveries from central that was delivered to the services
func (VS *Server) showDeliveriesFromCentral(c *gin.Context) {
	service := VS.getCentralWarehouse(c, "h_deliveries", "showDeliveriesFromCentral", "service := VS.getCentralWarehouse")
	deliveries := VS.getAllDeliveriesSentByServiceID(c, service.ID, false)
	deliveriesHTML := VS.deliveryToHTML(c, deliveries)
	deliveriesHTML = revertDeliveries(deliveriesHTML)
	render(c, gin.H{
		"dPack": deliveriesHTML,
	}, "deliveries-from-central.html")
}

// This handler shows the different services available for a delivery from Central
func (VS *Server) chooseServiceForDeliveryFromCentral(c *gin.Context) {
	villages := VS.getAllVillagesFromTheSystem(c, "h_deliveries", "chooseServiceForDeliveryFromCentral", "villages := VS.getAllVillagesFromTheSystem")
	servicesToChoose := VS.selectServiceDeliveryHTML(c, villages, bson.NewObjectId(), false)
	render(c, gin.H{
		"services": servicesToChoose,
	}, "choose-service-for-delivery-from-central.html")
}

// This handler shows the form to include more stocks in the delivery, and shows the items that are already in the delivery from Central
func (VS *Server) showItemsDeliveryFromCentral(c *gin.Context) {
	serviceEmitter := VS.getCentralWarehouse(c, "h_deliveries", "showItemsDeliveryFromCentral", "serviceEmitter := VS.getCentralWarehouse")

	idDelivery, b := VS.getIDParamFromURL(c, "iddelivery")
	if !b {
		return
	}
	idServiceReceiver, b := VS.getIDParamFromURL(c, "idservice")
	if !b {
		return
	}
	serviceReceiver := VS.getServiceByID(c, idServiceReceiver, "h_deliveries", "showItemsDeliveryFromCentral", "serviceReceiver := VS.getServiceByID", true)
	serviceEmitterHTML := VS.serviceToHTML(c, serviceEmitter)
	serviceReceiverHTML := VS.serviceToHTML(c, serviceReceiver)
	itemsToChoose := VS.getAllITemsForServicePickerGlobal(c, "h_deliveries", "showItemsDeliveryFromCentral", "itemsToChoose := VS.getAllITemsForServicePickerGlobal")
	del, _ := VS.getDeliveryByID(c, idDelivery, "h_deliveries", "showItemsDeliveryFromCentral", "del, _ := VS.getDeliveryByID", true)
	sent := del.IsSent
	itemsChosen := VS.getAllStockByDelivery(c, idDelivery, "h_deliveries", "showItemsDeliveryFromCentral", "itemsChosen := VS.getAllStockByDelivery")
	itemsChosenHTML := VS.stocksToHTML(c, itemsChosen)

	render(c, gin.H{
		"serviceReceiver": serviceReceiverHTML,
		"serviceEmitter":  serviceEmitterHTML,
		"sent":            sent,
		"itemstochoose":   itemsToChoose,
		"itemsChosen":     itemsChosenHTML,
	}, "delivery-form-to-send.html")
}

// This handler adds the items of the global delivery
func (VS *Server) addItemsDeliveryGlobal(c *gin.Context) {
	serviceEmitter := VS.getCentralWarehouse(c, "h_deliveries", "addItemsDeliveryGlobal", "serviceEmitter := VS.getCentralWarehouse")
	idServiceReceiver, b := VS.getIDParamFromURL(c, "idservice")
	if !b {
		return
	}
	serviceReceiver := VS.getServiceByID(c, idServiceReceiver, "h_deliveries", "addItemsDeliveryGlobal", "serviceReceiver := VS.getServiceByID", true)

	// Create the delivery if does not exist
	idDelivery, b := VS.getIDParamFromURL(c, "iddelivery")
	if !b {
		return
	}
	delivery, isDelivery := VS.getDeliveryByID(c, idDelivery, "h_deliveries", "addItemsDeliveryGlobal", "delivery, isDelivery := VS.getDeliveryByID", true)
	if !isDelivery {
		delivery = Delivery{
			ID:                idDelivery,
			IDServiceEmitter:  serviceEmitter.ID,
			IDServiceReceiver: serviceReceiver.ID,
			Date:              getTime(),
			IsSent:            false,
			IDManager:         VS.getUserFomCookie(c, "h_deliveries", "addItemsDeliveryGlobal", "VS.getUserFomCookie").ID,
		}
	}

	// CHECKING about that the action is not to deliver
	// If is already send, we redirect to see the delivery, not to add or edit
	if delivery.IsSent {
		fmt.Print("Line 127")
		c.Redirect(http.StatusFound, "/deliveries/stock-created")
		return
	}

	// Check if the action is to deliver instead of add an item
	sent := VS.getCheckBoxFromHTML(c, "sent")
	if sent {
		delivery.IsSent = true
		VS.addElementToDatabaseWithRegisterFromCentral(c, delivery, delivery.ID, "h_deliveries", "addItemsDeliveryGlobal", "VS.addElementToDatabaseWithRegisterFromCentral")
		// Audit
		audit := VS.getAuditFromCentral(c, delivery, idDelivery, Modified)
		VS.registerAuditAndEvent(c, audit, delivery, delivery.ID)
		VS.goodFeedback(c, "/data/new-stock/"+serviceReceiver.ID.Hex()+"/"+delivery.ID.Hex())
		return
	}

	// We clean the extra information needed for the slider
	itemToDeliverID, b := VS.getIDFromHTMLCleaningIt(c, "|", "stock", "Stock")
	if !b {
		return
	}

	// We check the QR if is trackable
	item, _ := VS.getItemByID(c, itemToDeliverID, "h_deliveries", "addItemsDeliveryGlobal", "item, _ := VS.getItemByID", false)
	cat := VS.getCategoryByID(c, item.IDCategory, "h_deliveries", "addItemsDeliveryGlobal", "cat := VS.getCategoryByID")

	stockInDelivery := Stock{
		ID:                  bson.NewObjectId(),
		Photo:               item.Photo,
		IDServiceLocation:   delivery.ID,
		TypeOfItem:          cat.TypeOfItem,
		IDVideoCourseOrItem: item.ID,
		IDUser:              VS.getUserFomCookie(c, "h_deliveries", "addItemsDeliveryGlobal", "UserID: VS.getUserFomCookie").ID,
		Date:                getTime(),
		ServiceCreated:      serviceEmitter.ID,
		IsTrackable:         cat.IsTrackable,
	}
	if cat.IsTrackable {
		QR, b := VS.getQRFromHTMLNew(c, Tool)
		if !b {
			return
		}
		stockInDelivery.ID = QR
	}
	if !isDelivery {
		VS.addElementToDatabaseWithRegisterFromCentral(c, delivery, delivery.ID, "h_deliveries", "addItemsDeliveryGlobal", "VS.addElementToDatabaseWithRegisterFromCentral")
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, stockInDelivery, stockInDelivery.ID, "h_deliveries", "addItemsDeliveryGlobal", "VS.addElementToDatabaseWithRegisterFromCentral")
	// Delivery
	delivery.Stocks = append(delivery.Stocks, stockInDelivery.ID)
	VS.addElementToDatabaseWithRegisterFromCentral(c, delivery, delivery.ID, "h_deliveries", "addItemsDeliveryGlobal", "VS.addElementToDatabaseWithRegisterFromCentral")
	c.Redirect(http.StatusFound, "/deliveries/new-stock/"+serviceReceiver.ID.Hex()+"/"+delivery.ID.Hex())
}

// this handler deletes a item from the global delivery from central
func (VS *Server) deleteItemDeliveryGlobal(c *gin.Context) {
	idStock, b := VS.getIDParamFromURL(c, "idstock")
	if !b {
		return
	}
	stock := VS.getStockByID(c, idStock, "h_deliveries", "deleteItemDeliveryGlobal", "stock := VS.getStockByID")
	delivery, _ := VS.getDeliveryByID(c, stock.IDServiceLocation, "h_deliveries", "deleteItemDeliveryGlobal", "delivery, _ := VS.getDeliveryByID", true)
	if delivery.IsSent {
		return
	}
	VS.deleteAnything(c, stock, stock.ID, "h_deliveries", "deleteItemDeliveryGlobal", "VS.deleteAnything")
	// Recycle the QR
	if stock.IsTrackable {
		QR := QR{
			ID:         idStock,
			TypeOfItem: stock.TypeOfItem,
		}
		VS.addElementToDatabaseWithRegisterFromCentral(c, QR, QR.ID, "h_reports", "deleteItemDeliveryGlobal", "VS.addElementToDatabaseWithRegisterFromCentral")
	}
	delivery.Stocks = RemovesID(delivery.Stocks, idStock)
	VS.addElementToDatabaseWithRegisterFromCentral(c, delivery, delivery.ID, "h_deliveries", "deleteItemDeliveryGlobal", "VS.addElementToDatabaseWithRegisterFromCentral")
	VS.goodFeedback(c, "/data/new-stock/"+delivery.IDServiceReceiver.Hex()+"/"+delivery.ID.Hex())
}

// This handler shows the receive form for the delivery for the Manager in the service, to take photos or scan QR of the items when the delivery is received
func (VS *Server) showConfirmDelivery(c *gin.Context) {
	idDelivery, b := VS.getIDParamFromURL(c, "iddelivery")
	if !b {
		return
	}
	delivery, _ := VS.getDeliveryByID(c, idDelivery, "h_deliveries", "showConfirmDelivery", "delivery, _ := VS.getDeliveryByID", true)
	idService := delivery.IDServiceReceiver

	serviceReceiver := VS.getServiceByID(c, idService, "h_deliveries", "showConfirmDelivery", "serviceReceiver := VS.getServiceByID", true)
	serviceReceiverHTML := VS.serviceToHTML(c, serviceReceiver)

	itemsChecked := VS.getAllStockCheckedByDelivery(c, idDelivery, "h_deliveries", "showConfirmDelivery", "itemsChecked := VS.getAllStockCheckedByDelivery")
	itemsCheckedHTML := VS.stocksToHTML(c, itemsChecked)
	itemsNOTChecked := VS.getAllStockNOTCheckedByDelivery(c, idDelivery, "h_deliveries", "showConfirmDelivery", "itemsNOTChecked := VS.getAllStockNOTCheckedByDelivery")
	itemsNOTCheckedHTML := VS.stocksToHTML(c, itemsNOTChecked)

	render(c, gin.H{
		"service":         serviceReceiverHTML,
		"close":           delivery.IsComplete,
		"role":            getRole(c),
		"itemsChecked":    itemsCheckedHTML,
		"itemsNOTChecked": itemsNOTCheckedHTML,
	}, "delivery-form-receive.html")
}

// Mark as received the stock scanned by the QR or by the photo if it does not have QR
func (VS *Server) confirmDelivery(c *gin.Context) {
	service, b := VS.getOurService(c, "h_deliveries", "confirmDelivery", "service, b := VS.getOurService")
	if !b {
		return
	}

	// Comprobar la checkbox
	QR, b := VS.getIDFromHTMLWithoutChecking(c, "qr")
	if b {
		stock := VS.getStockByID(c, QR, "h_deliveries", "confirmDelivery", "stock := VS.getStockByID")
		stock.IDServiceLocation = service.ID
		audit := VS.getAudit(c, stock, QR, service.ID, Modified)
		VS.checkChange(c, audit.ID, "ServiceLocation", "Delivery", string(service.Type)+": "+service.Name)
		VS.addElementToDatabaseWithRegister(c, stock, stock.ID, service.ID, "h_deliveries", "confirmDelivery", "VS.addElementToDatabaseWithRegister")
	}
	idDelivery, b := VS.getIDParamFromURL(c, "iddelivery")
	if !b {
		return
	}
	delivery, _ := VS.getDeliveryByID(c, idDelivery, "h_deliveries", "confirmDelivery", "delivery, _ := VS.getDeliveryByID", true)
	if delivery.IsComplete {
		return
	}
	// Check if the action is to close the delivery
	if VS.getCheckBoxFromHTML(c, "confirmed") {
		delivery.IsComplete = true
		VS.addElementToDatabaseWithRegister(c, delivery, delivery.ID, service.ID, "h_deliveries", "confirmDelivery", "VS.addElementToDatabaseWithRegister")
	}
	itemsReceived, _ := VS.getArrayIDFromHTMLWithoutChecking(c, "item", "Delivery")
	for _, i := range itemsReceived {
		stock := VS.getStockByID(c, i, "h_deliveries", "confirmDelivery", "stock := VS.getStockByID")
		stock.IDServiceLocation = service.ID
		audit := VS.getAudit(c, delivery, QR, service.ID, Modified)
		VS.checkChange(c, audit.ID, "ServiceLocation", "Delivery", string(service.Type)+": "+service.Name)
		VS.addElementToDatabaseWithRegister(c, stock, stock.ID, service.ID, "h_deliveries", "confirmDelivery", "VS.addElementToDatabaseWithRegister")
	}
	VS.goodFeedback(c, "/delivery/check-delivery/"+idDelivery.Hex())
}

/*************************  Services  *************************/

// This handler shows the different services that we can select to make a delivery from a Service
func (VS *Server) chooseServiceForDeliveryFromService(c *gin.Context) {
	serviceHTML, _ := VS.getOurService(c, "h_deliveries", "chooseServiceForDeliveryFromService", "serviceHTML, _ := VS.getOurService")

	role := getRole(c)
	var villages []Village
	villages = append(villages, VS.getCentralVillage(c, "h_deliveries", "chooseServiceForDeliveryFromService", "villages = append(villages, VS.getCentralVillage"))
	if role == Admin {
		villages = VS.getAllVillagesFromTheSystem(c, "h_deliveries", "chooseServiceForDeliveryFromService", "villages = VS.getAllVillagesFromTheSystem")
	} else {
		villages = append(villages, VS.getLocalVillage(c, "h_deliveries", "chooseServiceForDeliveryFromService", "villages = append(villages, VS.getLocalVillage"))
	}

	servicesToChoose := VS.selectServiceDeliveryHTML(c, villages, serviceHTML.ID, true)
	render(c, gin.H{
		"services":    servicesToChoose,
		"serviceType": serviceHTML.TypeID,
	}, "choose-service-deliveries.html")
}

// This handler shows the form to include more stocks in the delivery, and shows the items that are already in the delivery from a Service
func (VS *Server) showItemsDeliveryFromService(c *gin.Context) {
	serviceEmitterHTML, _ := VS.getOurService(c, "h_deliveries", "showItemsDeliveryFromService", "serviceEmitterHTML, _ := VS.getOurService")

	idServiceReceiver, b := VS.getIDParamFromURL(c, "idservice")
	if !b {
		return
	}
	serviceReceiver := VS.getServiceByID(c, idServiceReceiver, "h_deliveries", "showItemsDeliveryFromService", "serviceReceiver := VS.getServiceByID", true)
	serviceReceiverHTML := VS.serviceToHTML(c, serviceReceiver)

	idDelivery, b := VS.getIDParamFromURL(c, "iddelivery")
	if !b {
		return
	}
	itemsToChoose := VS.getAllStockForServicePicker(c, serviceEmitterHTML.ID, "h_deliveries", "showItemsDeliveryFromService", "itemsToChoose := VS.getAllStockForServicePicker")

	del, _ := VS.getDeliveryByID(c, idDelivery, "h_deliveries", "showItemsDeliveryFromService", "del, _ := VS.getDeliveryByID", true)
	sent := del.IsSent
	itemsChosen := VS.getAllStockByDelivery(c, idDelivery, "h_deliveries", "showItemsDeliveryFromService", "itemsChosen := VS.getAllStockByDelivery")
	itemsChosenHTML := VS.stocksToHTML(c, itemsChosen)

	render(c, gin.H{
		"serviceType":     serviceEmitterHTML.TypeID,
		"serviceReceiver": serviceReceiverHTML,
		"serviceEmitter":  serviceEmitterHTML,
		"sent":            sent,
		"itemstochoose":   itemsToChoose,
		"itemsChosen":     itemsChosenHTML,
		"isFromService":   true,
	}, "delivery-form-to-send.html")
}

// This handler adds the item of the Service to the delivery
func (VS *Server) addItemsDeliveryFromService(c *gin.Context) {
	serviceEmitterHTML, _ := VS.getOurService(c, "h_deliveries", "addItemsDeliveryFromService", "serviceEmitterHTML, _ := VS.getOurService")

	idServiceReceiver, b := VS.getIDParamFromURL(c, "idservice")
	if !b {
		return
	}
	serviceReceiver := VS.getServiceByID(c, idServiceReceiver, "h_deliveries", "addItemsDeliveryFromService", "serviceReceiver := VS.getServiceByID", true)

	idDelivery, b := VS.getIDParamFromURL(c, "iddelivery")
	if !b {
		return
	}
	idServiceType, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	delivery, isDelivery := VS.getDeliveryByID(c, idDelivery, "h_deliveries", "addItemsDeliveryFromService", "delivery, isDelivery := VS.getDeliveryByID", true)
	if !isDelivery {
		delivery = Delivery{
			ID:                idDelivery,
			IDServiceEmitter:  serviceEmitterHTML.ID,
			IDServiceReceiver: serviceReceiver.ID,
			Date:              getTime(),
			IsSent:            false,
		}
	}

	// CHECKING about that the action is not to deliver
	// If is already send, we redirect to see the delivery, not to add or edit
	if delivery.IsSent {
		c.Redirect(http.StatusFound, "/deliveries/stock-created")
		return
	}

	// Check if the action is to deliver instead of add an item
	sent := VS.getCheckBoxFromHTML(c, "sent")
	if sent {
		delivery.IsSent = true
		VS.addElementToDatabaseWithRegister(c, delivery, delivery.ID, serviceEmitterHTML.ID, "h_deliveries", "addItemsDeliveryFromService", "VS.addElementToDatabaseWithRegister(c, delivery")
		VS.goodFeedback(c, "deliveries/new-delivery/"+idServiceReceiver.Hex()+"/"+idDelivery.Hex()+"/"+idServiceType.Hex())
		return
	}

	// Nothing to deliver
	quantityToDeliver, _ := VS.getIntFromHTML(c, "slider-range", "Item")
	if quantityToDeliver == 0 {
		VS.wrongFeedback(c, "You can not deliver 0 items")
		return
	}

	// We clean the extra information needed for the slider
	itemToDeliverID, b := VS.getIDFromHTMLCleaningIt(c, "|", "stock", "Stock")
	if !b {
		return
	}

	// We check the QR if is trackable
	item, b := VS.getItemByID(c, itemToDeliverID, "h_deliveries", "addItemsDeliveryFromService", "item, b := VS.getItemByID", false)
	var cat Category
	if b {
		cat = VS.getCategoryByID(c, item.IDCategory, "h_deliveries", "addItemsDeliveryFromService", "cat = VS.getCategoryByID")
	} else {
		// Is product from a video
		cat.IsTrackable = true
		video, _ := VS.getVideoCourseByID(c, item.IDCategory, "h_deliveries", "addItemsDeliveryFromService", "video, _ := VS.getVideoCourseByID", false)
		item.ID = video.ID
	}
	stock := VS.getStockByItemOrVideoID(c, itemToDeliverID, "h_deliveries", "addItemsDeliveryFromService", "stock := VS.getStockByItemOrVideoID")
	if cat.IsTrackable {
		QR, b := VS.getQRFromHTML(c)
		if !b {
			return
		}
		stock = VS.getStockByID(c, QR, "h_deliveries", "addItemsDeliveryFromService", "stock = VS.getStockByID(c, QR,")
	}
	if stock.IDServiceLocation != serviceEmitterHTML.ID {
		VS.wrongFeedback(c, "This object does not belong to this service")
		return
	}
	if stock.IDVideoCourseOrItem != item.ID {
		VS.wrongFeedback(c, "The QR does not belong to this object")
		return
	}
	stock.IDServiceLocation = delivery.ID
	VS.addElementToDatabaseWithRegister(c, stock, stock.ID, serviceEmitterHTML.ID, "h_deliveries", "addItemsDeliveryFromService", "addElementToDatabaseWithRegister(c, stock, stock")
	// Delivery
	delivery.Stocks = append(delivery.Stocks, stock.ID)
	VS.addElementToDatabaseWithoutRegister(c, delivery, delivery.ID, "h_deliveries", "addItemsDeliveryFromService", "VS.addElementToDatabaseWithoutRegister(c, delivery")
	// Refresh the page
	c.Redirect(http.StatusFound, "/deliveries/new-delivery/"+serviceReceiver.ID.Hex()+"/"+delivery.ID.Hex()+"/"+serviceEmitterHTML.TypeID.Hex())
}
