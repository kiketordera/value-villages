package village

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timshannon/bolthold"
	"gopkg.in/mgo.v2/bson"
)

// "github.com/goinggo/tracelog"

/*************************    SYNCHRONIZATION    *************************/

// This method returns all the Synchronizations made with the DBs of the villages
func (VS *Server) getAllSynchronizations(c *gin.Context) []Synchronization {
	var syncs []Synchronization
	VS.DataBase.Find(&syncs, nil)
	return syncs
}

// This method returns all the Synchronizations made with the DBs as the Village gives as emitter
func (VS *Server) getAllSynchronizationsByVillageEmitter(c *gin.Context, IDVillageEmitter bson.ObjectId, class string, method string, operation string) []Synchronization {
	var syncs []Synchronization
	err := VS.DataBase.Find(&syncs, bolthold.Where("IDVillageEmitter").Eq(IDVillageEmitter))
	VS.checkOperation(c, class, method, operation, err)
	return syncs
}

// This method returns all the Synchronizations made with the DBs as the Village gives as receiver
func (VS *Server) getAllSynchronizationsByVillageReceiver(c *gin.Context, IDVillageReceiver bson.ObjectId, class string, method string, operation string) []Synchronization {
	var syncs []Synchronization
	err := VS.DataBase.Find(&syncs, bolthold.Where("IDVillageReceiver").Eq(IDVillageReceiver))
	VS.checkOperation(c, class, method, operation, err)
	return syncs
}

// This method finds the Village in database by his ID and returns it, if not found prints in the screen error
func (VS *Server) getSynchronizationByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string, check bool) (Synchronization, bool) {
	var sync []Synchronization
	err := VS.DataBase.Find(&sync, bolthold.Where("ID").Eq(id))
	if check {
		VS.checkOperation(c, class, method, operation, err)
	}
	if len(sync) == 1 {
		return sync[0], true
	}
	return Synchronization{}, false
}

// This method returns the date of the last Synchronization from the village
// newSyncID is the Id of the new synchronization we are making, we pass the parameter to avoid to have the lastestSync date from the one we are just creating
func (VS *Server) getLastSynchronizationDateByVillageReceiver(c *gin.Context, id bson.ObjectId, newSyncID bson.ObjectId, class string, method string, operation string) int64 {
	var syncs []Synchronization
	err := VS.DataBase.Find(&syncs, bolthold.Where("IDVillageReceiver").Eq(id).And("ID").Ne(newSyncID))
	VS.checkOperation(c, class, method, operation, err)
	if len(syncs) > 0 {
		var lastSync int64
		lastSync = syncs[0].Date
		for _, s := range syncs {
			if s.Date > lastSync {
				lastSync = s.Date
			}
		}
		fmt.Print("Last synchronization was in: ")
		fmt.Println(lastSync)
		return lastSync
	}
	fmt.Println("Last synchronization not found ")
	return 0
}

// This method returns all the Audits that needs to be synchronized with the village
// newSyncID is the Id of the new synchronization we are making, we pass the parameter to avoid to have the lastestSync date from the one we are just creating
func (VS *Server) getAllAuditsFromLastSynchronization(c *gin.Context, IDVillage bson.ObjectId, newSyncID bson.ObjectId, class string, method string, operation string) []Audit {
	lastSync := VS.getLastSynchronizationDateByVillageReceiver(c, IDVillage, newSyncID, "h_database", "getAllAuditsFromLastSynchronization", "lastSync := VS.getLastSynchronizationDateByVillageReceiver")
	var audits []Audit
	var err error
	// // Get audits from the village
	// if lastSync == 0 {
	// 	err = VS.DataBase.Find(&audits, bolthold.Where("IDVillage").Eq(IDVillage))
	// } else {
	// 	err = VS.DataBase.Find(&audits, bolthold.Where("Date").Gt(lastSync).And("IDVillage").Eq(IDVillage))
	// }
	// VS.checkOperation(c, class, method, operation, err)

	// Get audits from the Central
	if lastSync == 0 {
		if VS.checkWeAreCentral(c, "h_database", "getAllAuditsFromLastSynchronization", "if VS.checkWeAreCentral") {
			fmt.Println("line 93 executed database.go")
			err = VS.DataBase.Find(&audits, bolthold.Where("IDService").Eq(VS.getCentralWarehouse(c, "h_database", "getAllAuditsFromLastSynchronization", "VS.DataBase.Find(&audits, bolthold.Where(IDService).Eq(VS.getCentralWarehouse").ID))
		} else {
			fmt.Println("line 96 executed database.go")
			err = VS.DataBase.Find(&audits, bolthold.Where("IDVillage").Eq(VS.getLocalVillage(c, "h_database", "getAllAuditsFromLastSynchronization", "audits, bolthold.Where(IDVillage).Eq(VS.getLocalVillage").ID))
			fmt.Println("These are the audits: ")
			fmt.Println(audits)

		}
	} else {
		fmt.Print(lastSync)
		if VS.checkWeAreCentral(c, "h_database", "getAllAuditsFromLastSynchronization", "if VS.checkWeAreCentral") {
			fmt.Println("line 102 executed database.go, lastsync: ")
			err = VS.DataBase.Find(&audits, bolthold.Where("Date").Gt(lastSync).And("IDService").Eq(VS.getCentralWarehouse(c, "h_database", "getAllAuditsFromLastSynchronization", "VS.DataBase.Find(&audits, bolthold.Where(Date).Ge(lastSync).And(IDService).Eq(VS.getCentralWarehouse").ID))
		} else {
			fmt.Println("line 105 executed database.go, lastsync: ")
			err = VS.DataBase.Find(&audits, bolthold.Where("Date").Gt(lastSync).And("IDVillage").Eq(VS.getLocalVillage(c, "h_database", "getAllAuditsFromLastSynchronization", "VS.DataBase.Find(&audits, bolthold.Where(Date).Ge(lastSync).And(IDVillage).Eq(VS.getLocalVillage").ID))
		}
	}
	VS.checkOperation(c, class, method, operation, err)
	fmt.Print("the audits are: ")
	fmt.Println(len(audits))
	return audits
}

func (VS *Server) checkWeAreCentral(c *gin.Context, class string, method string, operation string) bool {
	vil := VS.getLocalVillage(c, "h_database", "checkWeAreCentral", "VS.vil := VS.getLocalVillage")
	if vil.Name == "Central" {
		return true
	}
	return false
}

// This method returns all the Events from the Database
func (VS *Server) getAllEvents(c *gin.Context, class string, method string, operation string) []EventDB {
	var events []EventDB
	err := VS.DataBase.Find(&events, nil)
	VS.checkOperation(c, class, method, operation, err)
	return events
}

// This method returns all the Events that needs to be synchronized with the village
// newSyncID is the Id of the new synchronization we are making, we pass the parameter to avoid to have the lastestSync date from the one we are just creating
func (VS *Server) getAllEventsFromLastSynchronization(c *gin.Context, IDVillage bson.ObjectId, newSyncID bson.ObjectId, class string, method string, operation string) []EventDB {
	lastSync := VS.getLastSynchronizationDateByVillageReceiver(c, IDVillage, newSyncID, "h_database", "getAllEventsFromLastSynchronization", "lastSync := VS.getLastSynchronizationDateByVillageReceiver")
	var events []EventDB
	var err error
	// if lastSync == 0 {
	// 	err = VS.DataBase.Find(&events, bolthold.Where("IDVillage").Eq(IDVillage))
	// } else {
	// 	err = VS.DataBase.Find(&events, bolthold.Where("Date").Gt(lastSync).And("IDVillage").Eq(IDVillage))
	// }
	// VS.checkOperation(c, class, method, operation, err)
	// fmt.Print("The number of events after first block of code are: ")
	// fmt.Println(len(events))
	// Get events from the Central
	if lastSync == 0 {
		fmt.Println("line 95 executed in getAllEventsFromLastSynchronization in database.go")
		if VS.checkWeAreCentral(c, "h_database", "getAllAuditsFromLastSynchronization", "if VS.checkWeAreCentral") {
			fmt.Print("Line 145 executed")
			err = VS.DataBase.Find(&events, bolthold.Where("IDVillage").Eq(VS.getCentralVillage(c, "h_database", "getAllEventsFromLastSynchronization", "VS.DataBase.Find(&events, bolthold.Where(IDVillage).Eq(VS.getCentralVillage").ID))
		} else {
			fmt.Print("Line 148 executed")
			err = VS.DataBase.Find(&events, bolthold.Where("IDVillage").Eq(VS.getLocalVillage(c, "h_database", "getAllEventsFromLastSynchronization", "VS.DataBase.Find(&events, bolthold.Where(IDVillage).Eq(VS.getLocalVillage").ID))
		}
	} else {
		fmt.Print("line 98 executed in getAllEventsFromLastSynchronization in database.go, lastsync: ")
		fmt.Println(lastSync)
		if VS.checkWeAreCentral(c, "h_database", "getAllAuditsFromLastSynchronization", "if VS.checkWeAreCentral") {
			fmt.Print("Line 155 executed")
			err = VS.DataBase.Find(&events, bolthold.Where("Date").Gt(lastSync).And("IDVillage").Eq(VS.getCentralVillage(c, "h_database", "getAllEventsFromLastSynchronization", "VS.DataBase.Find(&events, bolthold.Where(Date).Ge(lastSync).And(IDVillage).Eq(VS.getCentralVillage").ID))
		} else {
			fmt.Print("Line 158 executed")
			err = VS.DataBase.Find(&events, bolthold.Where("Date").Gt(lastSync).And("IDVillage").Eq(VS.getLocalVillage(c, "h_database", "getAllEventsFromLastSynchronization", "VS.DataBase.Find(&events, bolthold.Where(Date).Ge(lastSync).And(IDVillage).Eq(VS.getLocalVillage").ID))
		}
	}
	VS.checkOperation(c, class, method, operation, err)
	fmt.Print("The number of events in the end are: ")
	fmt.Println(len(events))
	return events
}

/*************************    ADD / EDIT    *************************/

// This method adds any element to the DB and takes care about all related , putting all the changes and the audits
func (VS *Server) addElementToDatabaseWithRegister(c *gin.Context, i interface{}, ID bson.ObjectId, IDservice bson.ObjectId, class string, method string, operation string) {
	err := VS.DataBase.Upsert(ID, i)
	VS.checkOperation(c, class, method, operation, err)
	VS.registerAddDatabaseMovement(c, i, ID, IDservice)
}

// This method adds any element to the DB and takes care about all related , putting all the changes and the audits
// Does not have the c variable, so does not check for corredt operations in the Databade, but only executes THE VERY FIRST TIME the app initiates (when there is no DB)
func (VS *Server) addElementToDatabaseWithRegisterForStartingAPP(i interface{}, ID bson.ObjectId, IDservice bson.ObjectId, IDVillage bson.ObjectId, IDAdmin bson.ObjectId, class string, method string, operation string) {
	VS.DataBase.Upsert(ID, i)
	VS.registerAddDatabaseMovementForStartinggTheAPP(i, ID, IDservice, IDAdmin, IDVillage)
}

// This method adds any element to the DB and takes care about all related , putting all the changes and the audits
func (VS *Server) addElementToDatabaseWithRegisterFromCentral(c *gin.Context, i interface{}, ID bson.ObjectId, class string, method string, operation string) {
	err := VS.DataBase.Upsert(ID, i)
	VS.checkOperation(c, class, method, operation, err)
	VS.registerAddDatabaseMovementFromCentral(c, i, ID)
}

// This method adds any element to the DB and does not make the changes and the audits
func (VS *Server) addElementToDatabaseWithoutRegister(c *gin.Context, i interface{}, ID bson.ObjectId, class string, method string, operation string) {
	fmt.Println("Element added with ID: " + ID.Hex())
	err := VS.DataBase.Upsert(ID, i)
	VS.checkOperation(c, class, method, operation, err)
}

// This method edit the element into the DB and takes care about all the changes and the audits
// IDservice is the ID where occur the transaction
func (VS *Server) editElementInDatabaseWithRegister(c *gin.Context, newObject interface{}, oldObject interface{}, IDboth bson.ObjectId, IDservice bson.ObjectId, class string, method string, operation string) {
	err := VS.DataBase.Upsert(IDboth, newObject)
	VS.checkOperation(c, class, method, operation, err)
	// Here we save the changes
	audit := VS.getAudit(c, newObject, IDboth, IDservice, Modified)
	modified := VS.dynamicAuditChange(c, newObject, oldObject, audit.ID)
	if modified {
		VS.registerAuditAndEvent(c, audit, newObject, IDboth)
	}
}

// This method edit the element into the DB and takes care about all the changes and the audits
// ID service is the ID where occur the transaction
func (VS *Server) editElementInDatabaseWithRegisterFromCentral(c *gin.Context, newObject interface{}, oldObject interface{}, IDboth bson.ObjectId, class string, method string, operation string) {
	err := VS.DataBase.Upsert(IDboth, newObject)
	VS.checkOperation(c, class, method, operation, err)
	// Here we save the changes
	audit := VS.getAuditFromCentral(c, newObject, IDboth, Modified)
	modified := VS.dynamicAuditChange(c, newObject, oldObject, audit.ID)
	if modified {
		VS.registerAuditAndEvent(c, audit, newObject, IDboth)
	}
}

/***********************  VILLAGE ***********************/

// This method adds (or updates) the Village created to the Database,
// we don't allow to have 2 Villages with the same name or the same prefix
// Returns true if the operatios was successfully and false if it fails
func (VS *Server) addVillageToDatabase(c *gin.Context, p Village, class string, method string, operation string) bool {
	var villageName []Village
	VS.DataBase.Find(&villageName, bolthold.Where("Name").Eq(p.Name))
	var villagePref []Village
	VS.DataBase.Find(&villagePref, bolthold.Where("Prefix").Eq(p.Prefix))
	if len(villageName)+len(villageName) != 0 {
		// Here we check if we are editing the Village
		if villageName[0].ID == p.ID {
			err := VS.DataBase.Upsert(p.ID, p)
			VS.checkOperation(c, class, method, operation, err)
			return true
		}
		return false
	}
	err := VS.DataBase.Upsert(p.ID, p)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// This method finds the Village in database by his ID and returns it, if not found prints in the screen error
func (VS *Server) getVillageByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Village {
	var village Village
	err := VS.DataBase.Get(id, &village)
	VS.checkOperation(c, class, method, operation, err)
	return village
}

// This method returns our Local Village, that is the village of the server from the comand line to initiate the app
func (VS *Server) getLocalVillage(c *gin.Context, class string, method string, operation string) Village {
	// localvillage is the variable we write on the terminal to execute and run the app
	var village []Village
	err := VS.DataBase.Find(&village, bolthold.Where("Name").Eq(localVillage))
	VS.checkOperation(c, class, method, operation, err)
	if len(village) != 1 {
		VS.wrongFeedback(c, "Sorry, but the village name put in the comand line, does not exist in the DB!")
		return Village{}
	}
	return village[0]
}

// This method returns our Local Village, that is the village of the server from the comand line to initiate the app
func (VS *Server) existLocalVillage(c *gin.Context, class string, method string, operation string) bool {
	// localvillage is the variable we write on the terminal to execute and run the app
	var village []Village
	err := VS.DataBase.Find(&village, bolthold.Where("Name").Eq(localVillage))
	VS.checkOperation(c, class, method, operation, err)
	if len(village) != 1 {
		return false
	}
	return true
}

// This method returns the Central village
func (VS *Server) getCentralVillage(c *gin.Context, class string, method string, operation string) Village {
	// localvillage is the variable we write on the terminal to execute and run the app
	var village []Village
	err := VS.DataBase.Find(&village, bolthold.Where("Name").Eq("Central"))
	VS.checkOperation(c, class, method, operation, err)
	if len(village) != 1 {
		VS.wrongFeedback(c, "Sorry, was a problem accessing the Central Village!")
		return Village{}
	}
	return village[0]
}

// This method returns the Central Warehouse
func (VS *Server) getCentralWarehouse(c *gin.Context, class string, method string, operation string) Service {
	//localvillage is the variable we write on the terminal to execute and run the app
	central := VS.getCentralVillage(c, class, method, operation)
	var service []Service
	err := VS.DataBase.Find(&service, bolthold.Where("Name").Eq("Central Warehouse").And("IDVillage").Eq(central.ID))
	VS.checkOperation(c, class, method, operation, err)
	if len(service) != 1 {
		VS.wrongFeedback(c, "Sorry, was a problem accessing the Central Warehouse!")
		return Service{}
	}
	return service[0]
}

// This method returns all the Villages of our System
func (VS *Server) getAllVillagesFromTheSystem(c *gin.Context, class string, method string, operation string) []Village {
	var villages []Village
	err := VS.DataBase.Find(&villages, nil)
	VS.checkOperation(c, class, method, operation, err)
	return villages
}

// This method returns all the Villages of our System except the local Village
func (VS *Server) getAllVillagesFromTheSystemButLocal(c *gin.Context, class string, method string, operation string) []Village {
	localVillage := VS.getLocalVillage(c, "database", "getAllVillagesFromTheSystemButLocal", "localVillage := VS.getLocalVillage")
	var villages []Village
	err := VS.DataBase.Find(&villages, bolthold.Where("ID").Ne(localVillage.ID))
	VS.checkOperation(c, class, method, operation, err)
	return villages
}

/***********************  SERVICES ***********************/

// This method returns all the ServicesTypes of our system
func (VS *Server) getAllServiceTypes(c *gin.Context, class string, method string, operation string) []ServiceType {
	var services []ServiceType
	err := VS.DataBase.Find(&services, nil)
	VS.checkOperation(c, class, method, operation, err)
	return services
}

// This method finds the ServiceType in database by his ID and returns it.
func (VS *Server) getServiceTypeByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string, check bool) ServiceType {
	var ser ServiceType
	err := VS.DataBase.Get(id, &ser)
	if check {
		VS.checkOperation(c, class, method, operation, err)
	}
	return ser
}

// This method returns all the Services of our system
func (VS *Server) getAllServices(c *gin.Context, class string, method string, operation string) []Service {
	var services []Service
	err := VS.DataBase.Find(&services, nil)
	VS.checkOperation(c, class, method, operation, err)
	return services
}

// This method finds the Service in database by his ID and returns it.
func (VS *Server) getServiceByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string, check bool) Service {
	var ser Service
	err := VS.DataBase.Get(id, &ser)
	if check {
		VS.checkOperation(c, class, method, operation, err)
	}
	return ser
}

// This method adds (or updates) the Service created to the Database
// Returns true if the operatios was successfully and false if it fails
// We do not allow to have 2 services with he same name
func (VS *Server) addServiceToDatabase(c *gin.Context, s Service, class string, method string, operation string) bool {
	var service []Service
	err := VS.DataBase.Find(&service, bolthold.Where("Name").Eq(s.Name).And("IDVillage").Eq(s.IDVillage))
	VS.checkOperation(c, class, method, operation, err)
	if len(service) != 0 {
		if service[0].ID == s.ID {
			err := VS.DataBase.Upsert(s.ID, s)
			VS.checkOperation(c, class, method, operation, err)
			return true
		}
		return false
	}
	err = VS.DataBase.Upsert(s.ID, s)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// Returns true if the User given have access to the service given and false otherwise
func (VS *Server) hasAccess(c *gin.Context, userID bson.ObjectId, serviceID bson.ObjectId, class string, method string, operation string) bool {
	var access []Access
	err := VS.DataBase.Find(&access, bolthold.Where("IDUser").Eq(userID).And("IDService").Eq(serviceID).And("IsActive").Eq(true))
	VS.checkOperation(c, class, method, operation, err)
	if len(access) > 0 {
		return true
	}
	return false
}

// This method DELETE all the access by User id
func (VS *Server) deleteAllAccessByUser(c *gin.Context, userID bson.ObjectId, class string, method string, operation string) []Access {
	var acc []Access
	err := VS.DataBase.Find(&acc, bolthold.Where("IDUser").Eq(userID))
	VS.checkOperation(c, class, method, operation, err)
	for _, a := range acc {
		a.IsActive = false
		VS.addElementToDatabaseWithoutRegister(c, a, a.ID, "handlers_dataBase", "deleteAllAccessByUser", "v.addAuditToDatabase(c")
		VS.checkOperation(c, class, method, operation, err)
	}
	return acc
}

// This method DELETE all the access by Service id
func (VS *Server) deleteAllAccessByService(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []Access {
	var acc []Access
	err := VS.DataBase.Find(&acc, bolthold.Where("IDService").Eq(serviceID))
	VS.checkOperation(c, class, method, operation, err)
	for _, a := range acc {
		VS.addElementToDatabaseWithoutRegister(c, a, a.ID, "handlers_dataBase", "deleteAllAccessByUser", "v.addAuditToDatabase(c")
		err := VS.DataBase.Delete(a.ID, a)
		VS.checkOperation(c, class, method, operation, err)
	}
	return acc
}

// Returns the service of the local village, or the service representation (global variables in village.go) if there is more then 1 or if Admin user
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getOurService(c *gin.Context, class string, method string, operation string) (ServiceHTML, bool) {
	sType, _ := VS.getIDParamFromURL(c, "idservicetype")
	role := getRole(c)

	var service []Service
	if role == Admin {
		VS.DataBase.Find(&service, bolthold.Where("Type").Eq(sType))
	} else {
		vil := VS.getLocalVillage(c, "database", "getOurService", "vil := VS.getOurVillage")
		VS.DataBase.Find(&service, bolthold.Where("IDVillage").Eq(vil.ID).And("Type").Eq(sType))
	}

	if len(service) == 0 {
		c.Redirect(http.StatusFound, "/records/choose-service/"+sType.Hex())
		return ServiceHTML{}, false
	}

	if len(service) == 1 {
		IDVisualization[sType] = service[0].ID
	}

	var IDServiceEmpty bson.ObjectId
	serviceToView := IDVisualization[sType]
	if IDServiceEmpty == serviceToView {
		c.Redirect(http.StatusFound, "/records/choose-service/"+sType.Hex())
		return ServiceHTML{}, false
	}
	ser := VS.getServiceByID(c, IDVisualization[sType], "h_database", "getOurService", "ser := VS.getServiceByID(c, IDVisualization", true)
	serviceHTML := VS.serviceToHTML(c, ser)
	return serviceHTML, true
}

// This method Returns all the services that are from the Specific type given
func (VS *Server) getAllSpecificServiceFromSystem(c *gin.Context, serTypeID bson.ObjectId, class string, method string, operation string) []Service {
	var ser []Service
	err := VS.DataBase.Find(&ser, bolthold.Where("Type").Eq(serTypeID))
	VS.checkOperation(c, class, method, operation, err)
	return ser
}

// This method Returns all the services from the local village with a specific a type
func (VS *Server) getAllSpecificServicesFromLocalVillage(c *gin.Context, serTypeID bson.ObjectId, class string, method string, operation string) []Service {
	vil := VS.getLocalVillage(c, class, method, operation)
	var ser []Service
	err := VS.DataBase.Find(&ser, bolthold.Where("Type").Eq(serTypeID).And("IDVillage").Eq(vil.ID))
	VS.checkOperation(c, class, method, operation, err)
	return ser
}

// This method Returns all the services from the local village
func (VS *Server) getAllServicesFromLocalVillage(c *gin.Context, class string, method string, operation string) []Service {
	vil := VS.getLocalVillage(c, class, method, operation)
	var ser []Service
	err := VS.DataBase.Find(&ser, bolthold.Where("IDVillage").Eq(vil.ID))
	VS.checkOperation(c, class, method, operation, err)
	return ser
}

// This method Returns all Sales of the system related with the specific type of service
func (VS *Server) getAllSalesFromTheSystemByServiceType(c *gin.Context, serTypeID bson.ObjectId, class string, method string, operation string) []Sale {
	services := VS.getAllSpecificServiceFromSystem(c, serTypeID, class, method, operation)
	var salesService []Sale
	var totalSales []Sale
	for _, s := range services {
		err := VS.DataBase.Find(&salesService, bolthold.Where("IDService").Eq(s.ID))
		VS.checkOperation(c, class, method, operation, err)
		totalSales = append(totalSales, salesService...)
		salesService = nil
	}
	return totalSales
}

// This method Returns all Payments of the system related with the specific type of service
func (VS *Server) getAllPaymentsFromTheSystemByServiceType(c *gin.Context, serTypeID bson.ObjectId, class string, method string, operation string) []Payment {
	services := VS.getAllSpecificServiceFromSystem(c, serTypeID, class, method, operation)
	var paymentsService []Payment
	var totalPayments []Payment
	for _, s := range services {
		err := VS.DataBase.Find(&paymentsService, bolthold.Where("IDService").Eq(s.ID))
		VS.checkOperation(c, class, method, operation, err)
		totalPayments = append(totalPayments, paymentsService...)
		paymentsService = nil
	}
	return totalPayments
}

// This method Returns all StockGenerated of the system related with the specific type of service
func (VS *Server) getAllStockGeneratedFromTheSystemByServiceType(c *gin.Context, serTypeID bson.ObjectId, class string, method string, operation string) []Stock {
	services := VS.getAllSpecificServiceFromSystem(c, serTypeID, class, method, operation)
	var stockService []Stock
	var totalStock []Stock
	for _, s := range services {
		err := VS.DataBase.Find(&stockService, bolthold.Where("ServiceCreated").Eq(s.ID))
		VS.checkOperation(c, class, method, operation, err)
		totalStock = append(totalStock, stockService...)
		stockService = nil
	}
	return totalStock
}

// This method returns all the WorkerOrders from the services given
func (VS *Server) getAllWorkerOrdersFromService(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []WorkerOrder {
	var orders []WorkerOrder
	err := VS.DataBase.Find(&orders, bolthold.Where("IDService").Eq(serviceID))
	VS.checkOperation(c, class, method, operation, err)
	return orders
}

// This method returns all the Worker Order from a specific worker and service
func (VS *Server) getAllWorkerOrdersFromSpecificWorkerAndService(c *gin.Context, idWorker bson.ObjectId, idService bson.ObjectId, class string, method string, operation string) []WorkerOrder {
	var orders []WorkerOrder
	err := VS.DataBase.Find(&orders, bolthold.Where("IDService").Eq(idService).And("IDWorker").Eq(idWorker))
	VS.checkOperation(c, class, method, operation, err)
	return orders
}

// This method returns all the Worker Order of a specific worker
func (VS *Server) getAllWorkerOrdersFromSpecificWorker(c *gin.Context, idWorker bson.ObjectId, class string, method string, operation string) []WorkerOrder {
	var orders []WorkerOrder
	err := VS.DataBase.Find(&orders, bolthold.Where("IDWorker").Eq(idWorker))
	VS.checkOperation(c, class, method, operation, err)
	return orders
}

// This method finds the Video Course in database by his ID and returns it.
func (VS *Server) getVideoCourseByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string, check bool) (VideoCourse, bool) {
	var products []VideoCourse
	err := VS.DataBase.Find(&products, bolthold.Where("ID").Eq(id))
	if check {
		VS.checkOperation(c, class, method, operation, err)
		if len(products) == 0 {
			VS.checkOperation(c, "No Video Course found", method, operation, err)
			return VideoCourse{}, false
		}
	}
	if len(products) == 0 {
		return VideoCourse{}, false
	}
	return products[0], true
}

// This method returns all the Workers with Access to our workshop
func (VS *Server) getWorkersFromSpecificService(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []User {
	var workers []User
	var workersWithAccess []User
	err := VS.DataBase.Find(&workers, bolthold.Where("Role").Eq(Worker))
	for _, worker := range workers {
		if VS.hasAccess(c, worker.ID, serviceID, class, method, operation) {
			workersWithAccess = append(workersWithAccess, worker)
		}
	}
	VS.checkOperation(c, class, method, operation, err)
	return workersWithAccess
}

// This method finds all the Worker Orders related with a specific worker
func (VS *Server) getWorkerOrderByWorkerID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) []WorkerOrder {
	var orders []WorkerOrder
	err := VS.DataBase.Find(&orders, bolthold.Where("IDWorker").Eq(id))
	VS.checkOperation(c, class, method, operation, err)
	// if len(orders)==0 {VS.sendError(c, class, method, operation, err)}
	return orders
}

// This method finds the Worker Order in database by his ID and returns it
func (VS *Server) getWorkerOrderByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) WorkerOrder {
	var order WorkerOrder
	err := VS.DataBase.Get(id, &order)
	VS.checkOperation(c, class, method, operation, err)
	// if order == (WorkerOrder{}) {VS.sendError(c, class, method, operation, err)}
	return order
}

// Returns all the Sales from the specific Worker
func (VS *Server) getAllSalesFromSpecificWorker(c *gin.Context, workerID bson.ObjectId, class string, method string, operation string) []Sale {
	var sales []Sale
	err := VS.DataBase.Find(&sales, bolthold.Where("IDWorker").Eq(workerID))
	VS.checkOperation(c, class, method, operation, err)
	return sales
}

// This method returns all the Payments from the specific Worker
func (VS *Server) getAllPaymentsFromSpecificWorker(c *gin.Context, id bson.ObjectId, class string, method string, operation string) []Payment {
	var payments []Payment
	err := VS.DataBase.Find(&payments, bolthold.Where("IDWorker").Eq(id))
	VS.checkOperation(c, class, method, operation, err)
	return payments
}

// This method returns all the Payments of the specific Worker and Service
func (VS *Server) getAllPaymentsFromSpecificWorkerAndService(c *gin.Context, wkid bson.ObjectId, serviceid bson.ObjectId, class string, method string, operation string) []Payment {
	var payments []Payment
	err := VS.DataBase.Find(&payments, bolthold.Where("IDWorker").Eq(wkid).And("IDService").Eq(serviceid))
	VS.checkOperation(c, class, method, operation, err)
	return payments
}

// This method returns all the ServicePurchases of the specific Worker
func (VS *Server) getAllServicePurchasesFromSpecificWorker(c *gin.Context, workerID bson.ObjectId, class string, method string, operation string) []Stock {
	var workshopsP []Stock
	err := VS.DataBase.Find(&workshopsP, bolthold.Where("IDUser").Eq(workerID))
	VS.checkOperation(c, class, method, operation, err)
	return workshopsP
}

// This method returns all the Stock Generated from the specific Worker and Service
func (VS *Server) getAllStockGeneratedFromSpecificWorkerAndService(c *gin.Context, workerID bson.ObjectId, serviceID bson.ObjectId, class string, method string, operation string) []Stock {
	var workshopsP []Stock
	err := VS.DataBase.Find(&workshopsP, bolthold.Where("IDUser").Eq(workerID).And("ServiceCreated").Eq(serviceID))
	VS.checkOperation(c, class, method, operation, err)
	return workshopsP
}

// This method returns all video courses from the Database
func (VS *Server) getAllVideoCourses(c *gin.Context, class string, method string, operation string) []VideoCourse {
	var videos []VideoCourse
	err := VS.DataBase.Find(&videos, nil)
	VS.checkOperation(c, class, method, operation, err)
	return videos
}

// This method adds (or updates) the Product created to the Database,
// we don't allow to have 2 VideoCourse with the same name, so returns true if it can store it
func (VS *Server) addVideoCourseToDatabase(c *gin.Context, p VideoCourse, class string, method string, operation string) bool {
	var video []VideoCourse
	VS.DataBase.Find(&video, bolthold.Where("Name").Eq(p.Name))
	if len(video) != 0 {
		if video[0].ID == p.ID {
			err := VS.DataBase.Upsert(p.ID, p)
			VS.checkOperation(c, class, method, operation, err)
			return true
		}
		return false
	}
	err := VS.DataBase.Upsert(p.ID, p)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// This method returns all the steps from a specific VideoCourseID
func (VS *Server) getAllStepsFromSpecificVideoCourse(c *gin.Context, id bson.ObjectId, class string, method string, operation string) []Step {
	var steps []Step
	err := VS.DataBase.Find(&steps, bolthold.Where("IDVideoCourse").Eq(id))
	VS.checkOperation(c, class, method, operation, err)
	return steps
}

// This method returns all the ServiceOrders of our whole system
func (VS *Server) getAllServiceOrdersFromService(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []ServiceOrder {
	var serviceOrders []ServiceOrder
	err := VS.DataBase.Find(&serviceOrders, bolthold.Where("IDService").Eq(serviceID))
	VS.checkOperation(c, class, method, operation, err)
	return serviceOrders
}

// This method returns a ServiceOrder that is not all assigned,
// to discount the quantity assigned in the service order
func (VS *Server) getServiceOrderFromServiceAndProductInDate(c *gin.Context, serviceID bson.ObjectId, videoID bson.ObjectId, class string, method string, operation string) (ServiceOrder, bool) {
	var serviceOrders []ServiceOrder
	err := VS.DataBase.Find(&serviceOrders, bolthold.Where("IDService").Eq(serviceID).And("IDVideoCourse").Eq(videoID))
	VS.checkOperation(c, class, method, operation, err)
	fmt.Print("serviceOrders")
	fmt.Print(serviceOrders)

	serviceOrders = cleanServiceOrders(serviceOrders)
	fmt.Print("service Clean")
	fmt.Print(serviceOrders)
	if len(serviceOrders) > 0 {
		return serviceOrders[0], true
	}
	return ServiceOrder{}, false
}

// This method finds the ServiceOrder in database by his ID and returns it.
func (VS *Server) getServiceOrderByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) ServiceOrder {
	var pOrder ServiceOrder
	err := VS.DataBase.Get(id, &pOrder)
	VS.checkOperation(c, class, method, operation, err)
	// if pOrder == (WorkshopOrder{}) {VS.sendError(c, class, method, operation, err)}
	return pOrder
}

// This method finds the Payment in database by his ID and returns it.
func (VS *Server) getPaymentByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Payment {
	var pay Payment
	err := VS.DataBase.Get(id, &pay)
	VS.checkOperation(c, class, method, operation, err)
	// if pay == (Payment{}) {v.sendError(c, class, method, operation, err)}
	return pay
}

// This method finds the Checklist in DataBase from the video
func (VS *Server) getChecklistByVideo(c *gin.Context, videoID bson.ObjectId, class string, method string, operation string) []VideoCheckList {
	var list []VideoCheckList
	err := VS.DataBase.Find(&list, bolthold.Where("IDVideoCourse").Eq(videoID))
	VS.checkOperation(c, class, method, operation, err)
	// We look for translations just in case the language is another
	var categoriesTranslate []VideoCheckList
	if getLanguage(c) != MainLanguage {
		// We need to register the struct as an interface to save it on the EventDB
		gob.Register(Category{})
		gob.Register(Item{})
		for _, cat := range list {
			trans, b := VS.getTranslationByIDInstanceAndLanguage(c, cat.ID, "h_data", "showEditCategory", "trans, b := VS.getTranslationByIDInstanceAndLanguage(c, cat.ID")
			if b {
				// Needs type assertion
				categoriesTranslate = append(categoriesTranslate, trans.Instance.(VideoCheckList))
			} else {
				categoriesTranslate = append(categoriesTranslate, cat)
			}
		}
		return categoriesTranslate
	}
	return list
}

// This method returns all the problem from the user by product
func (VS *Server) getProblemsFromUserByProduct(c *gin.Context, productID bson.ObjectId, userID bson.ObjectId, class string, method string, operation string) []VideoProblem {
	var problems []VideoProblem
	err := VS.DataBase.Find(&problems, bolthold.Where("IDVideoCourse").Eq(productID).And("IDWorker").Eq(userID))
	VS.checkOperation(c, class, method, operation, err)
	return problems
}

// This method returns all the problem NOT from the user by product
func (VS *Server) getProblemsNOTFromUserByProduct(c *gin.Context, productID bson.ObjectId, userID bson.ObjectId, class string, method string, operation string) []VideoProblem {
	var problems []VideoProblem
	err := VS.DataBase.Find(&problems, bolthold.Where("IDVideoCourse").Eq(productID).And("IDWorker").Ne(userID))
	VS.checkOperation(c, class, method, operation, err)
	return problems
}

// This method finds the Check in database by his ID and returns it
func (VS *Server) getCheckByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string, check bool) VideoCheckList {
	var checkL []VideoCheckList
	err := VS.DataBase.Find(&checkL, bolthold.Where("ID").Eq(id))
	if check {
		VS.checkOperation(c, class, method, operation, err)
		if len(checkL) == 0 {
			// VS.sendError(c, "No Check found", method, operation, err)
			return VideoCheckList{}
		}
	}
	return checkL[0]
}

// This method returns all the Sales from specific service
func (VS *Server) getAllSalesFromSpecificService(c *gin.Context, idWK bson.ObjectId, class string, method string, operation string) []Sale {
	var sales []Sale
	err := VS.DataBase.Find(&sales, bolthold.Where("IDService").Eq(idWK))
	VS.checkOperation(c, class, method, operation, err)
	return sales
}

// This method returns all the Payments of our workshop
func (VS *Server) getAllPaymentsFromSpecificService(c *gin.Context, idWK bson.ObjectId, class string, method string, operation string) []Payment {
	var pay []Payment
	err := VS.DataBase.Find(&pay, bolthold.Where("IDService").Eq(idWK))
	VS.checkOperation(c, class, method, operation, err)
	return pay
}

// This method returns all the Purchases from the service given
func (VS *Server) getAllPurchasesFromSpecificService(c *gin.Context, idWK bson.ObjectId, class string, method string, operation string) []Stock {
	var s []Stock
	err := VS.DataBase.Find(&s, bolthold.Where("ServiceCreated").Eq(idWK))
	VS.checkOperation(c, class, method, operation, err)
	return s
}

/***********************  USERS ***********************/

// This method adds (or updates) the User created to the Database
// Returns true if the operatios was sucessfully and false if it fails
func (VS *Server) addUserToDatabase(c *gin.Context, u User, class string, method string, operation string) bool {
	// We check if the username is available, if not send error mesage and break
	var users []User
	if err := VS.DataBase.Find(&users, bolthold.Where("Username").Eq(u.Username)); err != nil {
		return false
	}
	if len(users) > 0 {
		if users[0].ID != u.ID {
			VS.wrongFeedback(c, "The Username is already taken")
			return false
		}
	}
	err := VS.DataBase.Upsert(u.ID, u)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// This method finds and return the User in database by his ID
func (VS *Server) getUserByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) User {
	var user User
	err := VS.DataBase.Get(id, &user)
	VS.checkOperation(c, class, method, operation, err)
	return user
}

// This method finds and return the Audit in database by his ID and returns it
func (VS *Server) getAuditByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Audit {
	var a Audit
	err := VS.DataBase.Get(id, &a)
	VS.checkOperation(c, class, method, operation, err)
	return a
}

// This method finds and return the Audit with the IDItem related.
// It assures that can only retrieve audits with an IDItem that can not be repeated in another audit (non modificable objects, like Synchronizations)
func (VS *Server) getAuditByIDItemID(c *gin.Context, IDItem bson.ObjectId, class string, method string, operation string) Audit {
	var audit []Audit
	err := VS.DataBase.Find(&audit, bolthold.Where("IDItem").Eq(IDItem))
	VS.checkOperation(c, class, method, operation, err)
	if len(audit) == 0 {
		VS.wrongFeedback(c, "Audit for object not found")
		return Audit{}
	} else if len(audit) > 1 {
		VS.wrongFeedback(c, "Audit for object repeated")
		return Audit{}
	}
	return audit[0]
}

// This method finds the User in database by his username and returns it.
func (VS *Server) getUserByUsername(c *gin.Context, username string, class string, method string, operation string) (User, bool) {
	var users []User
	err := VS.DataBase.Find(&users, bolthold.Where("Username").Eq(username))
	VS.checkOperation(c, class, method, operation, err)
	if users == nil {
		nickname := VS.getLocalVillage(c, class, method, operation).Prefix
		nickname += nickname
		err = VS.DataBase.Find(&users, bolthold.Where("Username").Eq(nickname))
		if users == nil {
			return User{}, false
		}
	}
	var user = users[0]
	return user, true
}

// This method returns all the Users from the whole system that are NOT ADMIN
func (VS *Server) getAllUsersFromTheSystemNotAdmin(c *gin.Context, class string, method string, operation string) []User {
	var users []User
	err := VS.DataBase.Find(&users, bolthold.Where("Role").Ne(Admin))
	VS.checkOperation(c, class, method, operation, err)
	return users
}

// This method returns all the Users from the whole system that are NOT ITSELF
func (VS *Server) getAllUsersFromTheSystemNotSelf(c *gin.Context, class string, method string, operation string) []User {
	u := VS.getUserFomCookie(c, "database", "getAllUsersFromTheSystemNotSelf", "VS.getUserFomCookie")
	var users []User
	err := VS.DataBase.Find(&users, bolthold.Where("ID").Ne(u.ID))
	VS.checkOperation(c, class, method, operation, err)
	return users
}

// This method returns all the Workers from the whole system
func (VS *Server) getAllWorkersFromTheSystem(c *gin.Context, class string, method string, operation string) []User {
	var users []User
	err := VS.DataBase.Find(&users, bolthold.Where("Role").Eq(Worker))
	VS.checkOperation(c, class, method, operation, err)
	return users
}

// This method returns all the Users from the local village
func (VS *Server) getAllUsersFromLocalVillage(c *gin.Context, class string, method string, operation string) []User {
	vill := VS.getLocalVillage(c, class, method, operation)
	var users []User
	err := VS.DataBase.Find(&users, bolthold.Where("IDVillage").Eq(vill.ID))
	VS.checkOperation(c, class, method, operation, err)
	return users
}

// This method returns all the Managers of our system
func (VS *Server) getAllManagersFromSystem(c *gin.Context, class string, method string, operation string) []User {
	var users []User
	err := VS.DataBase.Find(&users, bolthold.Where("Role").Eq(Manager))
	VS.checkOperation(c, class, method, operation, err)
	return users
}

// This method returns all the Workers of our local village
func (VS *Server) getAllWorkersFromLocalVillage(c *gin.Context, class string, method string, operation string) []User {
	vill := VS.getLocalVillage(c, class, method, operation)
	var users []User
	err := VS.DataBase.Find(&users, bolthold.Where("IDVillage").Eq(vill.ID).And("Role").Eq(Worker))
	VS.checkOperation(c, class, method, operation, err)
	return users
}

// This method returns all the Active accesses by User
func (VS *Server) getAllActiveAccessByUser(c *gin.Context, userID bson.ObjectId, class string, method string, operation string) []Access {
	var acc []Access
	err := VS.DataBase.Find(&acc, bolthold.Where("IDUser").Eq(userID).And("IsActive").Eq(true))
	VS.checkOperation(c, class, method, operation, err)
	return acc
}

// This method returns all the Services of the local village
func (VS *Server) getAllServicesFromSpecificVillage(c *gin.Context, id bson.ObjectId, class string, method string, operation string) []Service {
	var services []Service
	err := VS.DataBase.Find(&services, bolthold.Where("IDVillage").Eq(id))
	VS.checkOperation(c, class, method, operation, err)
	return services
}

// This method returns all the Services of the local village
func (VS *Server) getAllServicesFromSpecificVillageExceptItself(c *gin.Context, id bson.ObjectId, serviceID bson.ObjectId, class string, method string, operation string) []Service {
	var services []Service
	err := VS.DataBase.Find(&services, bolthold.Where("IDVillage").Eq(id).And("ID").Ne(serviceID))
	VS.checkOperation(c, class, method, operation, err)
	return services
}

/***********************  INTERNAL DATABASE ***********************/

// This method returns all the Audits of the User
func (VS *Server) getAuditsByUser(c *gin.Context, userID bson.ObjectId, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, bolthold.Where("IDUser").Eq(userID))
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns all the Audits of the Village
func (VS *Server) getAuditsByVillage(c *gin.Context, IDVillage bson.ObjectId, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, bolthold.Where("IDVillage").Eq(IDVillage))
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns all the Audits happend in the service
func (VS *Server) getAuditsByService(c *gin.Context, IDservice bson.ObjectId, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, bolthold.Where("IDService").Eq(IDservice))
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns all the Audits of the User
func (VS *Server) getAllAudits(c *gin.Context, class string, method string, operation string) []Audit {
	var aud []Audit
	err := VS.DataBase.Find(&aud, nil)
	VS.checkOperation(c, class, method, operation, err)
	return aud
}

// This method returns anything that have the ID as ID
func (VS *Server) getAnythingWithID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Anything {
	var object Anything

	// Check if is an Item
	var it []Item
	VS.DataBase.Find(&it, bolthold.Where("ID").Eq(id))
	if len(it) > 0 {
		object.ID = it[0].ID
		object.Instance = it[0]
		return object
	}

	// Check if is a To Do
	var todo []ToDo
	VS.DataBase.Find(&todo, bolthold.Where("ID").Eq(id))
	if len(todo) > 0 {
		object.ID = todo[0].ID
		object.Instance = todo[0]
		return object
	}

	// Check if is a To Do
	var todoChecked []ToDoChecked
	VS.DataBase.Find(&todoChecked, bolthold.Where("ID").Eq(id))
	if len(todoChecked) > 0 {
		object.ID = todoChecked[0].ID
		object.Instance = todoChecked[0]
		return object
	}

	// Check a is a Assignment
	var a []Assignment
	VS.DataBase.Find(&a, bolthold.Where("ID").Eq(id))
	if len(a) > 0 {
		object.ID = a[0].ID
		object.Instance = a[0]
		return object
	}

	// Check if is a DeliveryPack
	var dp []Delivery
	VS.DataBase.Find(&dp, bolthold.Where("ID").Eq(id))
	if len(dp) > 0 {
		object.ID = dp[0].ID
		object.Instance = dp[0]
		return object
	}

	// Check if is a Category
	var category []Category
	VS.DataBase.Find(&category, bolthold.Where("ID").Eq(id))
	if len(category) > 0 {
		object.ID = category[0].ID
		object.Instance = category[0]
		return object
	}

	// Check if is a Product
	var product []VideoCourse
	VS.DataBase.Find(&product, bolthold.Where("ID").Eq(id))
	if len(product) > 0 {
		object.ID = product[0].ID
		object.Instance = product[0]
		return object
	}

	// Check if is a Product
	var serType []ServiceType
	VS.DataBase.Find(&serType, bolthold.Where("ID").Eq(id))
	if len(serType) > 0 {
		object.ID = serType[0].ID
		object.Instance = serType[0]
		return object
	}

	// Check if is a Step
	var step []Step
	VS.DataBase.Find(&step, bolthold.Where("ID").Eq(id))
	if len(step) > 0 {
		object.ID = step[0].ID
		object.Instance = step[0]
		return object
	}

	// Check if is a Service
	var service []Service
	VS.DataBase.Find(&service, bolthold.Where("ID").Eq(id))
	if len(service) > 0 {
		object.ID = service[0].ID
		object.Instance = service[0]
		return object
	}

	// Check if is a ServiceOrder
	var workshopOrder []ServiceOrder
	VS.DataBase.Find(&workshopOrder, bolthold.Where("ID").Eq(id))
	if len(workshopOrder) > 0 {
		object.ID = workshopOrder[0].ID
		object.Instance = workshopOrder[0]
		return object
	}

	// Check if is an WorkerOrder
	var payment []Payment
	VS.DataBase.Find(&payment, bolthold.Where("ID").Eq(id))
	if len(payment) > 0 {
		object.ID = payment[0].ID
		object.Instance = payment[0]
		return object
	}

	// Check if is a WorkerOrder
	var workerOrder []WorkerOrder
	VS.DataBase.Find(&workerOrder, bolthold.Where("ID").Eq(id))
	if len(workerOrder) > 0 {
		object.ID = workerOrder[0].ID
		object.Instance = workerOrder[0]
		return object
	}

	// Check if is an Audit
	var audit []Audit
	VS.DataBase.Find(&audit, bolthold.Where("ID").Eq(id))
	if len(audit) > 0 {
		object.ID = audit[0].ID
		object.Instance = audit[0]
		return object
	}

	// Check a is a Report
	var r []Report
	VS.DataBase.Find(&r, bolthold.Where("ID").Eq(id))
	if len(r) > 0 {
		object.ID = r[0].ID
		object.Instance = r[0]
		return object
	}

	// Check if is a User
	var user []User
	VS.DataBase.Find(&user, bolthold.Where("ID").Eq(id))
	if len(user) > 0 {
		object.ID = user[0].ID
		object.Instance = user[0]
		return object
	}

	// Check if is a Village
	var village []Village
	VS.DataBase.Find(&village, bolthold.Where("ID").Eq(id))
	if len(village) > 0 {
		object.ID = village[0].ID
		object.Instance = village[0]
		return object
	}

	// Check if is an Access
	var acc []Access
	VS.DataBase.Find(&acc, bolthold.Where("ID").Eq(id))
	if len(acc) > 0 {
		object.ID = acc[0].ID
		object.Instance = acc[0]
		return object
	}

	// Check if is an Stock
	var wkPurchase []Stock
	VS.DataBase.Find(&wkPurchase, bolthold.Where("ID").Eq(id))
	if len(wkPurchase) > 0 {
		object.ID = wkPurchase[0].ID
		object.Instance = wkPurchase[0]
		return object
	}

	// Check if is an VideoProblem
	var problem []VideoProblem
	VS.DataBase.Find(&problem, bolthold.Where("ID").Eq(id))
	if len(problem) > 0 {
		object.ID = problem[0].ID
		object.Instance = problem[0]
		return object
	}

	// Check if is an Sale
	var sale []Sale
	VS.DataBase.Find(&sale, bolthold.Where("ID").Eq(id))
	if len(sale) > 0 {
		object.ID = sale[0].ID
		object.Instance = sale[0]
		return object
	}

	// Check if is an VideoCheckList
	var check []VideoCheckList
	VS.DataBase.Find(&check, bolthold.Where("ID").Eq(id))
	if len(check) > 0 {
		object.ID = check[0].ID
		object.Instance = check[0]
		return object
	}

	// Check if is a CategoryType
	var catType []CategoryType
	VS.DataBase.Find(&catType, bolthold.Where("ID").Eq(id))
	if len(catType) > 0 {
		object.ID = catType[0].ID
		object.Instance = catType[0]
		return object
	}

	// Check if is a Message
	var ms []Message
	VS.DataBase.Find(&ms, bolthold.Where("ID").Eq(id))
	if len(ms) > 0 {
		object.ID = ms[0].ID
		object.Instance = ms[0]
		return object
	}

	// Check if is a Conversation
	var conv []Conversation
	VS.DataBase.Find(&conv, bolthold.Where("ID").Eq(id))
	if len(conv) > 0 {
		object.ID = conv[0].ID
		object.Instance = conv[0]
		return object
	}

	// Check if is a Sync
	var sync []Synchronization
	VS.DataBase.Find(&sync, bolthold.Where("ID").Eq(id))
	if len(sync) > 0 {
		object.ID = sync[0].ID
		object.Instance = sync[0]
		return object
	}

	fmt.Print("An audit without case!!")

	return Anything{}
}

// This method returns all the Changes of the given Audit ID
func (VS *Server) getChangesByAudit(c *gin.Context, auditID bson.ObjectId, class string, method string, operation string) []Change {
	var ch []Change
	err := VS.DataBase.Find(&ch, bolthold.Where("IDAudit").Eq(auditID))
	VS.checkOperation(c, class, method, operation, err)
	return ch
}

/***********************  CATEGORIES ***********************/

// This method adds (or updates) the Category created to the Database, and checks if there is a Category with the same name before saving it
// return true if it can savew it and false otherwise
func (VS *Server) addCategoryToDatabase(c *gin.Context, ca Category, class string, method string, operation string) bool {
	var cat []Category
	VS.DataBase.Find(&cat, bolthold.Where("Name").Eq(ca.Name).And("Type").Eq(ca.Type))
	if len(cat) != 0 {
		if cat[0].ID == ca.ID {
			err := VS.DataBase.Upsert(ca.ID, ca)
			VS.checkOperation(c, class, method, operation, err)
			return true
		}
		return false
	}
	err := VS.DataBase.Upsert(ca.ID, ca)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// This method finds the Cartegory in database by his ID and returns it
func (VS *Server) getCategoryByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Category {
	var category Category
	err := VS.DataBase.Get(id, &category)
	VS.checkOperation(c, class, method, operation, err)
	return category
}

// This method finds the n in database by his ID and language and returns it
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getTranslationByIDInstanceAndLanguage(c *gin.Context, id bson.ObjectId, class string, method string, operation string) (Translation, bool) {
	var translation []Translation
	err := VS.DataBase.Find(&translation, bolthold.Where("IDInstance").Eq(id).And("Language").Eq(getLanguage(c)))
	VS.checkOperation(c, class, method, operation, err)
	if len(translation) == 0 {
		return Translation{}, false
	}
	return translation[0], true
}

// This method returns all the Categories of our system
func (VS *Server) getAllCategories(c *gin.Context, class string, method string, operation string) []Category {
	var categories []Category
	err := VS.DataBase.Find(&categories, nil)
	VS.checkOperation(c, class, method, operation, err)
	// We look for translations just in case the language is another
	var categoriesTranslate []Category
	if getLanguage(c) != MainLanguage {
		// We need to register the struct as an interface to save it on the EventDB
		gob.Register(Category{})
		gob.Register(Item{})
		for _, cat := range categories {
			trans, b := VS.getTranslationByIDInstanceAndLanguage(c, cat.ID, "h_data", "showEditCategory", "trans, b := VS.getTranslationByIDInstanceAndLanguage(c, cat.ID")
			if b {
				// Needs type assertion
				categoriesTranslate = append(categoriesTranslate, trans.Instance.(Category))
			} else {
				categoriesTranslate = append(categoriesTranslate, cat)
			}
		}
		return categoriesTranslate
	}
	return categories
}

// This method returns all the Categories of our system
func (VS *Server) getAllCategoriesOnlyForRent(c *gin.Context, class string, method string, operation string) []Category {
	var categories []Category
	err := VS.DataBase.Find(&categories, nil)
	VS.checkOperation(c, class, method, operation, err)
	var catOnlyRent []Category
	for _, cat := range categories {
		if cat.TypeOfItem == Tool {
			catOnlyRent = append(catOnlyRent, cat)
		}
	}
	return catOnlyRent
}

// This method returns all the Categories that are only for Sell to the workers
func (VS *Server) getAllCategoriesOnlyForSell(c *gin.Context, class string, method string, operation string) []Category {
	var categories []Category
	err := VS.DataBase.Find(&categories, nil)
	VS.checkOperation(c, class, method, operation, err)
	var catOnlySell []Category
	for _, cat := range categories {
		if cat.TypeOfItem != Tool {
			catOnlySell = append(catOnlySell, cat)
		}
	}
	return catOnlySell
}

// This method returns all Category Type
func (VS *Server) getAllCategoryType(c *gin.Context, class string, method string, operation string) []CategoryType {
	var cType []CategoryType
	err := VS.DataBase.Find(&cType, nil)
	VS.checkOperation(c, class, method, operation, err)
	return cType
}

// This method finds the CategoryType in database by his ID and returns it
func (VS *Server) getCategoryTypeByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) CategoryType {
	var cType CategoryType
	err := VS.DataBase.Get(id, &cType)
	VS.checkOperation(c, class, method, operation, err)
	return cType
}

// This method finds the CategoryType in database by his name and returns it
func (VS *Server) getCategoryTypeByName(c *gin.Context, name string, class string, method string, operation string) CategoryType {
	var cType []CategoryType
	err := VS.DataBase.Find(&cType, bolthold.Where("Name").Eq(name))
	VS.checkOperation(c, class, method, operation, err)
	return cType[0]
}

// This method returns true is we found a CategoryType with the name provided
func (VS *Server) existCategoryTypeName(name string) bool {
	var ct []CategoryType
	VS.DataBase.Find(&ct, bolthold.Where("Name").Eq(name))
	if len(ct) > 0 {
		return true
	}
	return false
}

// This method adds the Category type to DB
func (VS *Server) addCategoryTypeToDatabase(c *gin.Context, ct CategoryType, class string, method string, operation string) bool {
	if VS.existCategoryTypeName(ct.Name) {
		return false
	}
	err := VS.DataBase.Upsert(ct.ID, ct)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// This method finds the Item in the database by his ID and returns it
func (VS *Server) getItemByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string, check bool) (Item, bool) {
	var material []Item
	err := VS.DataBase.Find(&material, bolthold.Where("ID").Eq(id))
	if check {
		VS.checkOperation(c, class, method, operation, err)
		if len(material) == 0 {
			VS.checkOperation(c, "No Material found", method, operation, err)
			return Item{}, false
		}
	}
	if len(material) != 0 {
		return material[0], true
	}
	return Item{}, false
}

// This method returns all the Items from the Category given
func (VS *Server) getAllItemsByCategory(c *gin.Context, categoryID bson.ObjectId, class string, method string, operation string) []Item {
	var items []Item
	err := VS.DataBase.Find(&items, bolthold.Where("IDCategory").Eq(categoryID))
	VS.checkOperation(c, "No Material found", method, operation, err)
	// We look for translations just in case the language is another
	var itemsTranslate []Item
	if getLanguage(c) != MainLanguage {
		// We need to register the struct as an interface to save it on the EventDB
		gob.Register(Item{})
		for _, i := range items {
			trans, b := VS.getTranslationByIDInstanceAndLanguage(c, i.ID, "h_data", "showEditCategory", "trans := VS.getTranslationByIDAndLanguage")
			if b {
				// Needs type assertion
				itemsTranslate = append(itemsTranslate, trans.Instance.(Item))
			} else {
				itemsTranslate = append(itemsTranslate, i)
			}
		}
		bubbleSort(itemsTranslate)
		return itemsTranslate
	}
	bubbleSort(items)
	return items
}

/***********************  DELIVERIES ***********************/

// This method returns all the Stock from the Delivery given
func (VS *Server) getAllStockByDelivery(c *gin.Context, deliveryID bson.ObjectId, class string, method string, operation string) []Stock {
	del, _ := VS.getDeliveryByID(c, deliveryID, "database", "getAllStockByDelivery", "del := VS.getDeliveryByID", true)
	fmt.Print("This is delivery")
	fmt.Print(del)
	var stocks []Stock
	for _, s := range del.Stocks {
		stocks = append(stocks, VS.getStockByID(c, s, "database", "getAllStockByDelivery", "stocks = append(stocks, VS.getStockByID"))
	}
	return stocks
}

// This method returns all the Stock checked from the Delivery given
func (VS *Server) getAllStockCheckedByDelivery(c *gin.Context, deliveryID bson.ObjectId, class string, method string, operation string) []Stock {
	delivery, _ := VS.getDeliveryByID(c, deliveryID, "database", "getAllStockCheckedByDelivery", "delivery, _ := VS.getDeliveryByID", true)
	var stocks []Stock
	for _, s := range delivery.Stocks {
		stock := VS.getStockByID(c, s, "database", "getAllStockCheckedByDelivery", "stock := VS.getStockByID")
		if stock.IDServiceLocation != deliveryID {
			stocks = append(stocks, stock)
		}
	}
	return stocks
}

// This method returns all the Stock NOT Checked from the Delivery given
func (VS *Server) getAllStockNOTCheckedByDelivery(c *gin.Context, deliveryID bson.ObjectId, class string, method string, operation string) []Stock {
	var stocks []Stock
	VS.DataBase.Find(&stocks, bolthold.Where("IDServiceLocation").Eq(deliveryID))
	return stocks
}

// This method returns all the deliveries that are sent from the Service given
func (VS *Server) getAllDeliveriesSentByServiceID(c *gin.Context, serviceID bson.ObjectId, isSent bool) []Delivery {
	var dPack []Delivery
	var err error
	if isSent {
		err = VS.DataBase.Find(&dPack, bolthold.Where("IDServiceEmitter").Eq(serviceID).And("IsSent").Eq(true))
	} else {
		err = VS.DataBase.Find(&dPack, bolthold.Where("IDServiceEmitter").Eq(serviceID))
	}
	VS.checkOperation(c, "", "", "ya todo", err)
	return dPack
}

// This method returns all the deliveries that are received from the Service given
func (VS *Server) getAllDeliveriesByServiceReceiverIDSent(c *gin.Context, serviceID bson.ObjectId, isSent bool) []Delivery {
	var dPack []Delivery
	var err error
	if isSent {
		err = VS.DataBase.Find(&dPack, bolthold.Where("IDServiceReceiver").Eq(serviceID).And("IsSent").Eq(true))
	} else {
		err = VS.DataBase.Find(&dPack, bolthold.Where("IDServiceReceiver").Eq(serviceID))
	}
	VS.checkOperation(c, "", "", "ya todo", err)
	return dPack
}

// This method finds the Delivery in the database by his ID and returns it
func (VS *Server) getDeliveryByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string, check bool) (Delivery, bool) {
	var dp []Delivery
	err := VS.DataBase.Find(&dp, bolthold.Where("ID").Eq(id))
	if check {
		VS.checkOperation(c, class, method, operation, err)
		if len(dp) == 0 {
			VS.checkOperation(c, "No Deliveries found", method, operation, err)
			return Delivery{}, false
		}
	}
	return dp[0], true
}

// This method adds (or updates) the Item created to the Database
// Returns true if the operatios was successfully and false if it fails
func (VS *Server) addItemToDatabase(c *gin.Context, i Item, class string, method string, operation string) bool {
	var item []Item
	VS.DataBase.Find(&item, bolthold.Where("Name").Eq(i.Name).And("IDCategory").Eq(i.IDCategory))
	if len(item) != 0 {
		if item[0].ID == i.ID {
			err := VS.DataBase.Upsert(i.ID, i)
			VS.checkOperation(c, class, method, operation, err)
			return true
		}
		return false
	}
	err := VS.DataBase.Upsert(i.ID, i)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// This method adds (or updates) the ServiceType created to the Database
// Returns true if the operatios was successfully and false if it fails
func (VS *Server) addServiceTypeToDatabase(c *gin.Context, i ServiceType, class string, method string, operation string) bool {
	var item []Item
	VS.DataBase.Find(&item, bolthold.Where("Name").Eq(i.Name).And("ID").Eq(i.ID))
	if len(item) != 0 {
		if item[0].ID == i.ID {
			err := VS.DataBase.Upsert(i.ID, i)
			VS.checkOperation(c, class, method, operation, err)
			return true
		}
		return false
	}
	err := VS.DataBase.Upsert(i.ID, i)
	VS.checkOperation(c, class, method, operation, err)
	return true
}

// This method update the items inside the category
func (VS *Server) updateCategoryItems(c *gin.Context, catID bson.ObjectId, i Item) {
	cat := VS.getCategoryByID(c, catID, "h_database", "updateCategoryItems", "cat := VS.getCategoryByID")
	cat.Items = append(cat.Items, i.ID)
	VS.addCategoryToDatabase(c, cat, "h_database", "updateCategoryItems", "VS.addCategoryToDatabase")
	items := VS.getAllItemsByCategory(c, cat.ID, "h_database", "updateCategoryItems", "items := VS.getAllItemsByCategory")
	cate := VS.getCategoryByID(c, cat.ID, "h_database", "updateCategoryItems", "cate := VS.getCategoryByID")
	cate.Items = nil
	bubbleSort(items)
	for _, i := range items {
		cate.Items = append(cate.Items, i.ID)
	}
	VS.addCategoryToDatabase(c, cate, "h_database", "updateCategoryItems", "items := VS.getAllItemssByCategory")
}

/***********************  COMUNICATION ***********************/

// This method finds the messages related to a conversation given
func (VS *Server) getMessagesByConversation(c *gin.Context, conversationID bson.ObjectId, class string, method string, operation string) []Message {
	var mesages []Message
	err := VS.DataBase.Find(&mesages, bolthold.Where("IDConversation").Eq(conversationID))
	VS.checkOperation(c, class, method, operation, err)
	return mesages
}

// This method finds the Conversation in database by his ID and returns it, if not found prints in the screen error
func (VS *Server) getConversationByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Conversation {
	var conv Conversation
	err := VS.DataBase.Get(id, &conv)
	VS.checkOperation(c, class, method, operation, err)
	return conv
}

// This method finds the Message in database by his ID and returns it, if not found prints in the screen error
func (VS *Server) getMessageByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Message {
	var conv Message
	err := VS.DataBase.Get(id, &conv)
	VS.checkOperation(c, class, method, operation, err)
	return conv
}

// This method returns all the Conversations by User
func (VS *Server) getConversationsByUser(c *gin.Context, userID bson.ObjectId, class string, method string, operation string) []Conversation {
	var conv []Conversation
	err := VS.DataBase.Find(&conv, nil)
	VS.checkOperation(c, class, method, operation, err)
	var userConv []Conversation
	for _, n := range conv {
		if ContainsID(n.Users, userID) {
			userConv = append(userConv, n)
		}
	}
	return userConv
}

// This method returns all the Reports from the User
func (VS *Server) getReportsByUser(c *gin.Context, userID bson.ObjectId, class string, method string, operation string) []Report {
	var reports []Report
	err := VS.DataBase.Find(&reports, bolthold.Where("IDUser").Eq(userID))
	VS.checkOperation(c, class, method, operation, err)
	return reports
}

// This method returns the report by the ID given
func (VS *Server) getReportByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Report {
	var rep Report
	err := VS.DataBase.Get(id, &rep)
	VS.checkOperation(c, class, method, operation, err)
	return rep
}

// This method finds and return the User in database by his ID
func (VS *Server) getTodoCheckedByToDoID(c *gin.Context, todoID bson.ObjectId, class string, method string, operation string) []ToDoChecked {
	var todo []ToDoChecked
	err := VS.DataBase.Find(&todo, bolthold.Where("IDToDo").Eq(todoID))
	VS.checkOperation(c, class, method, operation, err)
	return todo
}

// This method finds and return the Todo in database by his ID
func (VS *Server) getTodoByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) ToDo {
	var todo ToDo
	err := VS.DataBase.Get(id, &todo)
	VS.checkOperation(c, class, method, operation, err)
	return todo
}

// This method returns all the ToDo of our System
func (VS *Server) getAllToDo(c *gin.Context, class string, method string, operation string) []ToDo {
	var todo []ToDo
	err := VS.DataBase.Find(&todo, nil)
	VS.checkOperation(c, class, method, operation, err)
	return todo
}

// This method returns the stock that is from the video or item given
func (VS *Server) getStockByItemOrVideoID(c *gin.Context, idItemOrVideo bson.ObjectId, class string, method string, operation string) Stock {
	var stock []Stock
	VS.DataBase.Find(&stock, bolthold.Where("IDVideoCourseOrItem").Eq(idItemOrVideo))
	if len(stock) != 0 {
		return stock[len(stock)-1]
	}
	return Stock{}
}

// This method finds and return the Stock in database by his ID
func (VS *Server) getStockByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Stock {
	var s Stock
	err := VS.DataBase.Get(id, &s)
	VS.checkOperation(c, class, method, operation, err)
	return s
}

// This method returns all the Positive Stock given the serviceID
func (VS *Server) getAllPositiveStockByServiceID(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []Stock {
	var dPack []Stock
	err := VS.DataBase.Find(&dPack, bolthold.Where("IDServiceLocation").Eq(serviceID).And("Quantity").Gt(0.0))
	VS.checkOperation(c, class, method, operation, err)
	return dPack
}

// This method returns all the ToDo from the User given, for checking, so takes care to see if exist already a check and should appear or not
func (VS *Server) getToDoByUserIDForMakeChecking(c *gin.Context, userID bson.ObjectId, class string, method string, operation string) []ToDo {
	var todos []ToDo
	var todosToCheck []ToDo
	err := VS.DataBase.Find(&todos, bolthold.Where("IDUser").Eq(userID))
	VS.checkOperation(c, class, method, operation, err)

	for _, s := range todos {
		shouldLastCheck := getTime()
		t := time.Unix(shouldLastCheck, 0)
		switch s.TimeChecking {

		case OnlyOnce:
			var newChecks []ToDoChecked
			VS.DataBase.Find(&newChecks, bolthold.Where("IDToDo").Eq(s.ID))
			if len(newChecks) != 0 {
				// Don't check
				shouldLastCheck = 0
			}
		// Every month
		case EveryMonth:
			t = t.AddDate(0, -1, 0)
			shouldLastCheck = t.Unix()
		// Every week
		case EveryWeek:
			t = t.AddDate(0, 0, -7)
			shouldLastCheck = t.Unix()
		// Every day
		case EveryDay:
			t = t.AddDate(0, 0, -1)
			shouldLastCheck = t.Unix()
		}
		// First check
		var newChecks []ToDoChecked
		VS.DataBase.Find(&newChecks, bolthold.Where("IDToDo").Eq(s.ID))
		if len(newChecks) == 0 {
			todosToCheck = append(todosToCheck, s)
		}
		// Old checks
		var checks []ToDoChecked
		VS.DataBase.Find(&checks, bolthold.Where("IDToDo").Eq(s.ID).And("Date").Lt(shouldLastCheck))
		if len(checks) > 0 {
			todosToCheck = append(todosToCheck, s)
		}
	}
	return todosToCheck
}

// This method returns all the ToDo from the Service given, for checking, so takes care to see if exist already a check and should appear or not
func (VS *Server) getAllPositiveStockByServiceIDForMakeChecking(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []Stock {
	var stocks []Stock
	var stocksToCheck []Stock
	err := VS.DataBase.Find(&stocks, bolthold.Where("IDServiceLocation").Eq(serviceID).And("Quantity").Gt(0.0))
	VS.checkOperation(c, class, method, operation, err)

	for _, s := range stocks {
		var categ Category
		shouldLastCheck := getTime()
		t := time.Unix(shouldLastCheck, 0)
		item, b := VS.getItemByID(c, s.IDVideoCourseOrItem, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "item, b := VS.getItemByID", false)
		if !b {
			// Every which time that we should check the products made in the services, in this case 1 week
			t = t.AddDate(0, 0, 7)
			shouldLastCheck = t.Unix()
		} else {
			cat := VS.getCategoryByID(c, item.IDCategory, "h_database", "getAllPositiveStockByServiceIDForMakeChecking", "cat := VS.getCategoryByID")
			categ = cat
			fmt.Print("This is shouldLastCheck: ")
			fmt.Print(shouldLastCheck)
			switch cat.TimeChecking {
			// Don't check
			case Never:
				shouldLastCheck = 0
			// Every month
			case EveryMonth:
				t = t.AddDate(0, -1, 0)
				shouldLastCheck = t.Unix()
			// Every week
			case EveryWeek:
				t = t.AddDate(0, 0, -7)
				shouldLastCheck = t.Unix()
			// Every day
			case EveryDay:
				t = t.AddDate(0, 0, -1)
				shouldLastCheck = t.Unix()
			// Every time Manager change village
			case ChangeVillage:
				shouldLastCheck = 0
				stocksToCheck = append(stocksToCheck, s)
			}
		}
		// First check
		var newChecks []CheckStock
		VS.DataBase.Find(&newChecks, bolthold.Where("IDInstance").Eq(s.ID))
		if len(newChecks) == 0 && categ.TimeChecking != Never {
			stocksToCheck = append(stocksToCheck, s)
		}
		// Old checks
		var checks []CheckStock
		VS.DataBase.Find(&checks, bolthold.Where("IDInstance").Eq(s.ID).And("Date").Lt(shouldLastCheck))
		if len(checks) > 0 {
			stocksToCheck = append(stocksToCheck, s)
		}
	}
	return stocksToCheck
}

// This method returns all the Positive stock from the Service and item|video given
func (VS *Server) getAllPositiveStockByServiceIDandItemOrVideoID(c *gin.Context, serviceID bson.ObjectId, itemID bson.ObjectId, class string, method string, operation string) []Stock {
	var stock []Stock
	err := VS.DataBase.Find(&stock, bolthold.Where("IDServiceLocation").Eq(serviceID).And("Quantity").Gt(0.0).And("IDVideoCourseOrItem").Eq(itemID))
	VS.checkOperation(c, class, method, operation, err)
	return stock
}

// This method returns all the Stock from the Service and TypeOfItem given
func (VS *Server) getAllStockByServiceIDandTypeOfItem(c *gin.Context, serviceID bson.ObjectId, tpe TypeOfItem, class string, method string, operation string) []Stock {
	var stock []Stock
	err := VS.DataBase.Find(&stock, bolthold.Where("IDServiceLocation").Eq(serviceID).And("TypeOfItem").Eq(tpe))
	VS.checkOperation(c, class, method, operation, err)
	return stock
}

// This method returns all the Assignments from the Service given
func (VS *Server) getAssignmentsByServiceID(c *gin.Context, serviceID bson.ObjectId, class string, method string, operation string) []Assignment {
	var as []Assignment
	err := VS.DataBase.Find(&as, bolthold.Where("IDService").Eq(serviceID))
	VS.checkOperation(c, "No Material found", method, operation, err)
	return as
}

// This method returns true is the item is available for assignment and false if not
func (VS *Server) isStockAvailableForAssignment(c *gin.Context, stockID bson.ObjectId, class string, method string, operation string) bool {
	var as []Assignment
	VS.DataBase.Find(&as, bolthold.Where("IDStock").Eq(stockID).And("IsBack").Eq(false))
	if len(as) > 0 {
		return false
	}
	return true
}

// This method finds and returns the assignment by his ID
func (VS *Server) getAssignmentByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Assignment {
	var rep Assignment
	err := VS.DataBase.Get(id, &rep)
	VS.checkOperation(c, class, method, operation, err)
	return rep
}

// This method returns all the Reports of our System
func (VS *Server) getReportsFromSystem(c *gin.Context, class string, method string, operation string) []Report {
	var reports []Report
	err := VS.DataBase.Find(&reports, nil)
	VS.checkOperation(c, class, method, operation, err)
	return reports
}

/***********************  PDF ***********************/

// This method adds (or updates) the Access created to the Database
func (VS *Server) deleteQRfromDatabase(c *gin.Context, a QR, class string, method string, operation string) {
	err := VS.DataBase.Delete(a.ID, a)
	VS.checkOperation(c, class, method, operation, err)
}

// This method returns all the Villages of our System
func (VS *Server) getQRbyID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) QR {
	var rep QR
	err := VS.DataBase.Get(id, &rep)
	VS.checkOperation(c, class, method, operation, err)
	return rep
}

// This method returns all the Villages of our System
func (VS *Server) getAllPDFQR(c *gin.Context, class string, method string, operation string) []PDF {
	var rep []PDF
	err := VS.DataBase.Find(&rep, bolthold.Where("TypeOfPDF").Eq(PDFQR))
	VS.checkOperation(c, class, method, operation, err)
	return rep
}

// This method give us all the Video Courses by a ServiceType
// Returns the Videos in the correspondenting language
func (VS *Server) getAllVideoCoursesByServiceType(c *gin.Context, serviceTypeID bson.ObjectId, class string, method string, operation string) []VideoCourse {
	allVideos := VS.getAllVideoCourses(c, "database", "getAllVideoCoursesByServiceType", "allCategories := VS.getAllCategories")
	// We clean and take only the video courses with access to the Item
	var videosWithService []VideoCourse
	for _, video := range allVideos {
		if ContainsID(video.ServicesAccess, serviceTypeID) {
			videosWithService = append(videosWithService, video)
		}
	}
	// We look for translations just in case the language is another
	var videosTranslate []VideoCourse
	if getLanguage(c) != MainLanguage {
		// We need to register the struct as an interface to save it on the EventDB
		gob.Register(VideoCourse{})
		for _, video := range allVideos {
			trans, b := VS.getTranslationByIDInstanceAndLanguage(c, video.ID, "h_data", "showEditCategory", "trans, b := VS.getTranslationByIDInstanceAndLanguage(c, cat.ID")
			if b {
				// Needs type assertion

				videosTranslate = append(videosTranslate, trans.Instance.(VideoCourse))
			} else {
				videosTranslate = append(videosTranslate, video)
			}
		}
		return videosTranslate
	}
	return videosWithService
}

// This method true is the QR is available for a new use and false if is not, and if available, it deletes it from the available QR table
func (VS *Server) isQRavailableForAssignment(c *gin.Context, serType TypeOfItem, idQR bson.ObjectId, class string, method string, operation string) bool {
	var qrs []QR
	err := VS.DataBase.Find(&qrs, bolthold.Where("TypeOfItem").Eq(serType).And("ID").Eq(idQR))
	VS.checkOperation(c, class, method, operation, err)
	if len(qrs) == 1 {
		code := VS.getQRbyID(c, idQR, "handlers_manager", "showAddPaymentM", "workers := v.getAllWorkers")
		VS.deleteQRfromDatabase(c, code, "handlers_manager", "showAddPaymentM", "workers := v.getAllWorkers")
		return true
	}
	return false
}

// *****************************  DELETE *****************************

// This method Deletes anything given
func (VS *Server) deleteAnything(c *gin.Context, u interface{}, id bson.ObjectId, class string, method string, operation string) {
	err := VS.DataBase.Delete(id, u)
	VS.checkOperation(c, class, method, operation, err)
}

// This method deletes the information of the Duplicate tests
func (VS *Server) deleteInformationTest(c *gin.Context, class string, method string, operation string) {
	var village []Village
	err := VS.DataBase.Find(&village, bolthold.Where("Name").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range village {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}

	var service []Service
	err = VS.DataBase.Find(&service, bolthold.Where("Name").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range service {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}

	var manager []User
	err = VS.DataBase.Find(&manager, bolthold.Where("Username").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range manager {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}

	var worker []User
	err = VS.DataBase.Find(&worker, bolthold.Where("Username").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range worker {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}

	var catType []CategoryType
	err = VS.DataBase.Find(&catType, bolthold.Where("Name").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range catType {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}

	var cat []Category
	err = VS.DataBase.Find(&cat, bolthold.Where("Name").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range cat {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}

	var i []Item
	err = VS.DataBase.Find(&i, bolthold.Where("Name").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range i {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}

	var video []VideoCourse
	err = VS.DataBase.Find(&video, bolthold.Where("Name").Eq("test"))
	VS.checkOperation(c, class, method, operation, err)
	for _, v := range video {
		VS.deleteAnything(c, v, v.ID, class, method, operation)
	}
}

// This method deletes a Stock given
func (VS *Server) deleteStock(c *gin.Context, stock Stock, class string, method string, operation string) {
	err := VS.DataBase.Delete(stock.ID, stock)
	VS.checkOperation(c, class, method, operation, err)
}

// This method finds and return the Step in database by his ID
func (VS *Server) getStepByID(c *gin.Context, id bson.ObjectId, class string, method string, operation string) Step {
	var s Step
	err := VS.DataBase.Get(id, &s)
	VS.checkOperation(c, class, method, operation, err)
	return s
}

// This method finds the step by the product and the index
func (VS *Server) getStepByProductIDAndIndexOrder(c *gin.Context, videoCourseID bson.ObjectId, index int, class string, method string, operation string) Step {
	var s []Step
	err := VS.DataBase.Find(&s, bolthold.Where("IDVideoCourse").Eq(videoCourseID).And("IndexOrder").Eq(index))
	VS.checkOperation(c, class, method, operation, err)
	return s[0]
}

// This method returns all the Checks from the item and ServiceID
func (VS *Server) getChecksFromItemAndServiceID(c *gin.Context, itemID bson.ObjectId, serviceID bson.ObjectId, class string, method string, operation string) (CheckStock, bool) {
	var check []CheckStock
	err := VS.DataBase.Find(&check, bolthold.Where("IDInstance").Eq(itemID).And("IDService").Eq(serviceID))
	VS.checkOperation(c, class, method, operation, err)
	if len(check) > 0 {
		return check[0], true
	}
	return CheckStock{}, false
}
