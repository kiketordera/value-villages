package village

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

/*************************  USERS  *************************/
// This handler shows the all the users of the system if the user is Admin, or all the user of the village if the user is Manager
func (VS *Server) showUsers(c *gin.Context) {
	users := VS.getAllUsersDependingRole(c, "h_data", "showUsers", "users = VS.getAllUsersFromLocalVillage")
	uHTML := VS.usersToHTML(c, users)
	JSON, err := json.Marshal(uHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"JSON":  string(JSON),
		"users": users,
	}, "all-users.html")
}

// Shows the page to introduce the information of the new User
func (VS *Server) showNewUser(c *gin.Context) {
	villages := VS.getAllVillagesFromTheSystem(c, "h_data", "showNewUser", "villages := VS.getAllVillagesFromTheSystem")
	render(c, gin.H{
		"villages": villages,
	}, "form-user.html")
}

// Creates the new User
func (VS *Server) newUser(c *gin.Context) {
	newUser, b := VS.getNewUser(c, bson.NewObjectId(), "", 0)
	if !b {
		return
	}
	if VS.addUserToDatabase(c, newUser, "h_data", "newUser", "VS.addUserToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, newUser, newUser.ID)
		VS.goodFeedback(c, "/data/all-users")
	} else {
		VS.wrongFeedback(c, "Username not available!")
	}
}

// Shows the page with the information of the User to modify the information
func (VS *Server) showEditUser(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	user := VS.getUserByID(c, id, "h_data", "showEditUser", "user := v.getUserByID")
	userhtml := VS.userToHTML(c, user)
	userhtml.Role = strings.ToLower(userhtml.Role)
	userhtml.Gender = strings.ToLower(userhtml.Gender)
	villages := VS.getAllVillagesFromTheSystem(c, "h_data", "showEditUser", "villages := v.getAllVillages")
	render(c, gin.H{
		"user":     userhtml,
		"edit":     true,
		"villages": villages,
	}, "form-user.html")
}

// Modifies the new User
func (VS *Server) editUser(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldUser := VS.getUserByID(c, id, "h_data", "editUser", "oldUser := VS.getUserByID")
	newUser, b := VS.getNewUser(c, id, oldUser.Photo, oldUser.Balance)
	if !b {
		return
	}

	// Here we save the changes
	audit := VS.getAuditFromCentral(c, newUser, id, Modified)
	modified := VS.dynamicAuditChange(c, newUser, oldUser, audit.ID)
	// Here manually the changes that can not be automatized
	var i bson.ObjectId
	var oldVillage, newVillage string
	if oldUser.IDVillage != i {
		oldVillage = VS.getVillageByID(c, oldUser.IDVillage, "h_data", "editUser", "oldVillage := VS.getVillageByID").Name
	}
	if newUser.IDVillage != i {
		newVillage = VS.getVillageByID(c, newUser.IDVillage, "h_data", "editUser", "newVillage := VS.getVillageByID").Name
	}
	if VS.checkChange(c, audit.ID, "VillageID", oldVillage, newVillage) {
		VS.deleteAllAccessByUser(c, oldUser.ID, "h_data", "editUser", "VS.deleteAllAccessByUser")
		modified = true
	}
	if VS.checkChange(c, audit.ID, "Username", oldUser.Username, newUser.Username) {
		modified = true
		VS.move(basePath+"/local-resources/users/"+oldUser.Username+"/", basePath+"/local-resources/users/"+newUser.Username+"/")
	}
	if modified {
		if VS.addUserToDatabase(c, newUser, "h_data", "editUser", "VS.addUserToDatabase") {
			VS.registerAuditAndEvent(c, audit, newUser, newUser.ID)
		} else {
			VS.wrongFeedback(c, "Username not available!")
		}
	}
	VS.goodFeedback(c, "data/all-users")
}

/*************************  ACCESS  *************************/
// Shows the information about every worker to see his accesses to all the services
func (VS *Server) showAccess(c *gin.Context) {
	users := VS.getAllWorkersDependingRole(c, "h_data", "showAccess", "users = VS.getAllUsersFromLocalVillage")
	accHTML := VS.getAccessFromUsersHTML(c, users)
	JSON, err := json.Marshal(accHTML)
	VS.checkOperation(c, "h_data", "showAccess", "JSON, err := json.Marshal(accHTML)", err)
	render(c, gin.H{
		"users":    users,
		"accesses": accHTML,
		"JSON":     string(JSON),
	}, "access.html")
}

// Shows the page for edit the access for an User
func (VS *Server) showEditAccess(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	user := VS.getUserByID(c, id, "h_data", "showEditAccess", "user := VS.getUserByID")
	services := VS.getAllServicesFromSpecificVillage(c, user.IDVillage, "h_data", "showEditAccess", "services := v.getAllServicesFromLocalVillage")
	servicesHTML := VS.servicesToHTML(c, services)
	render(c, gin.H{
		"services": servicesHTML,
		"user":     user,
	}, "form-access.html")
}

// Edits the access to the services of the worker
func (VS *Server) editAccess(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	user := VS.getUserByID(c, id, "h_data", "editAccess", "user := VS.getUserByID")

	// Here for taking the different values of the checkboxes with the same name
	servicesToActivate := c.PostFormArray("service")
	VS.deleteAllAccessByUser(c, user.ID, "h_data", "editAccess", "servicesToActivate := c.PostFormArray")
	for _, serviceID := range servicesToActivate {
		a := Access{
			ID:        bson.NewObjectId(),
			IDUser:    user.ID,
			IDService: bson.ObjectIdHex(serviceID),
			IsActive:  true,
		}
		VS.addElementToDatabaseWithRegisterFromCentral(c, a, a.ID, "h_data", "editAccess", "VS.addElementToDatabaseWithRegisterFromCentral")
	}
	VS.goodFeedback(c, "data/access")
}

/*************************  VILLAGES & SERVICES  *************************/

// Shows the screen with all Villages and Services included in the sistem if the user is Admin, or from the local village if the user is Manager
func (VS *Server) showTypeServices(c *gin.Context) {
	servicesType := VS.getAllServiceTypes(c, "h_data", "showTypeServices", "servicesType := VS.getAllServiceTypes")
	servicesTypeHTML := VS.servicesTypeToServicesTypeHTML(c, servicesType)
	JSON, err := json.Marshal(servicesTypeHTML)
	VS.checkOperation(c, "h_date", "showTypeServices", "JSON, err := json.Marshal(servicesTypeHTML", err)
	render(c, gin.H{
		"services": servicesTypeHTML,
		"JSON":     string(JSON),
	}, "types-service.html")
}

// Shows the page to introduce the information of the new Service Type
func (VS *Server) showNewServiceType(c *gin.Context) {
	render(c, gin.H{}, "form-service-type.html")
}

// Creates a new Service Type
func (VS *Server) newServiceType(c *gin.Context) {
	serviceType, b := VS.getNewServiceType(c, bson.NewObjectId())
	if !b {
		return
	}
	if VS.addServiceTypeToDatabase(c, serviceType, "h_data", "newServiceType", "v.addServiceTypeToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, serviceType, serviceType.ID)
		VS.goodFeedback(c, "/data/types-service")
	} else {
		VS.wrongFeedback(c, "Sorry, there is already one Item with that name in the same category. Please choose another name or category")
	}
}

// Shows the page with the information of the Service Type to modify
func (VS *Server) showEditServiceType(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	serType := VS.getServiceTypeByID(c, id, "h_data", "showEditServiceType", "serType := VS.getServiceTypeByID", true)
	render(c, gin.H{
		"serType": serType,
		"edit":    true,
	}, "form-service-type.html")
}

// Edit the Service Type into the DataBase
func (VS *Server) editServiceType(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldserType := VS.getServiceTypeByID(c, id, "h_data", "editServiceType", "oldserType := VS.getServiceTypeByID", true)
	newserType, b := VS.getNewServiceType(c, id)
	if !b {
		return
	}
	VS.editElementInDatabaseWithRegisterFromCentral(c, newserType, oldserType, newserType.ID, "h_data", "editVillage", "VS.editElementInDatabaseWithRegisterFromCentral")
	VS.goodFeedback(c, "/data/types-service")
}

// Shows the screen with all Villages and Services included in the sistem
func (VS *Server) showVillages(c *gin.Context) {
	villages := VS.getAllVillagesFromTheSystem(c, "h_data", "showVillages", "villages := v.getAllVillagesFromTheSystem")
	services := VS.getAllServices(c, "h_data", "showVillages", "services := v.getAllServices")
	servicesHTML := VS.servicesToHTML(c, services)
	JSON, err := json.Marshal(servicesHTML)
	VS.checkOperation(c, "h_date", "showVillages", "JSON, err := json.Marshal(servicesHTML)", err)
	render(c, gin.H{
		"services": servicesHTML,
		"villages": villages,
		"JSON":     string(JSON),
	}, "villages.html")
}

// Shows the page to introduce the information of the new Village
func (VS *Server) showNewVillage(c *gin.Context) {
	render(c, gin.H{}, "form-village.html")
}

// Stores the new Village information into the DataBase
func (VS *Server) newVillage(c *gin.Context) {
	prefix, b := VS.getStringFromHTML(c, "prefix", "Village")
	if !b {
		return
	}
	village, b := VS.getNewVillage(c, bson.NewObjectId(), prefix)
	if !b {
		return
	}
	if VS.addVillageToDatabase(c, village, "h_data", "newVillage", "VS.addVillageToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, village, village.ID)
		VS.goodFeedback(c, "data/villages")
	} else {
		VS.wrongFeedback(c, "Sorry, there is already one Village with that Name or Prefix. Please choose another name ")
	}
}

// Shows the page with the information of the Village to modify
func (VS *Server) showEditVillage(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	vill := VS.getVillageByID(c, id, "h_data", "showEditVillage", "vill := v.getVillageByID")
	render(c, gin.H{
		"village": vill,
		"edit":    true,
	}, "form-village.html")
}

// Modifies Village information into the DataBase
func (VS *Server) editVillage(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldVillage := VS.getVillageByID(c, id, "h_data", "editVillage", "oldVillage := VS.getVillageByID")
	newVillage, b := VS.getNewVillage(c, id, oldVillage.Prefix)
	if !b {
		return
	}

	// Here we save the changes
	audit := VS.getAuditFromCentral(c, newVillage, id, Modified)
	modified := VS.dynamicAuditChange(c, newVillage, oldVillage, audit.ID)
	if modified {
		if VS.addVillageToDatabase(c, newVillage, "h_data", "editVillage", "VS.addVillageToDatabase") {
			VS.registerAuditAndEvent(c, audit, newVillage, newVillage.ID)
		} else {
			VS.wrongFeedback(c, "Sorry, there is already one Village with that Name or Prefix. Please choose another name ")
		}
		VS.goodFeedback(c, "/data/villages")
	}
}

// Shows the page to introduce the information of the new Service
func (VS *Server) showNewService(c *gin.Context) {
	villages := VS.getAllVillagesFromTheSystem(c, "h_data", "showNewService", "villages := v.getAllVillages")
	sType := VS.getAllServiceTypes(c, "h_data", "showNewService", "sType := VS.getAllServiceTypes")
	servicesTypeHTML := VS.servicesTypeToServicesTypeHTML(c, sType)
	render(c, gin.H{
		"villages": villages,
		"sType":    servicesTypeHTML,
	}, "form-service.html")
}

// Stores the new Service information into the DataBase
func (VS *Server) newService(c *gin.Context) {
	service, b := VS.getNewService(c, bson.NewObjectId())
	if !b {
		return
	}
	if VS.addServiceToDatabase(c, service, "h_data", "newService", "VS.addServiceToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, service, service.ID)
		VS.goodFeedback(c, "/data/villages")
	} else {
		VS.wrongFeedback(c, "Sorry, there is already one Service with that name in that Village. Please choose another name ")
	}
}

// Shows the page with the information of the Service to modify
func (VS *Server) showEditService(c *gin.Context) {
	sType := VS.getAllServiceTypes(c, "h_data", "showEditService", "sType := VS.getAllServiceTypes")
	servicesTypeHTML := VS.servicesTypeToServicesTypeHTML(c, sType)
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	service := VS.getServiceByID(c, id, "h_data", "showEditService", "service := VS.getServiceByID", true)
	// servHTML := VS.serviceToHTML(c, service)
	villages := VS.getAllVillagesFromTheSystem(c, "h_data", "showEditService", "villages := VS.getAllVillagesFromTheSystem")
	render(c, gin.H{
		"villages": villages,
		// "service": servHTML,
		"service": service,
		"edit":    true,
		"sType":   servicesTypeHTML,
	}, "form-service.html")
}

// Modified the Service into the Database
func (VS *Server) editService(c *gin.Context) {
	//Get the id of the Product from the Link
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldService := VS.getServiceByID(c, id, "h_data", "editService", "oldService := VS.getServiceByID", true)
	newService, b := VS.getNewService(c, id)
	if !b {
		return
	}
	fmt.Print("This is service complete: ")
	fmt.Println(newService)

	// Audit
	audit := VS.getAuditFromCentral(c, newService, id, Modified)
	modified := VS.dynamicAuditChange(c, newService, oldService, audit.ID)
	oldVillage := VS.getVillageByID(c, oldService.IDVillage, "h_data", "editService", "oldVillage := VS.getVillageByID").Name
	newVillage := VS.getVillageByID(c, newService.IDVillage, "h_data", "editService", "newVillage := VS.getVillageByID").Name
	if VS.checkChange(c, audit.ID, "VillageID", oldVillage, newVillage) {
		VS.deleteAllAccessByService(c, id, "h_data", "editService", "VS.deleteAllAccessByService")
		modified = true
	}
	if modified {
		if VS.addServiceToDatabase(c, newService, "h_data", "editService", "VS.addServiceToDatabase") {
			VS.registerAuditAndEvent(c, audit, newService, newService.ID)
		} else {
			VS.wrongFeedback(c, "Sorry, there is already one Service with that name in that Village. Please choose another name ")
		}
	}
	VS.goodFeedback(c, "/data/villages")
}

/*************************  CATEGORIES  *************************/

// Shows the screen with all the Categories in the database
func (VS *Server) showCategories(c *gin.Context) {
	cat := VS.getAllCategories(c, "h_data", "showCategories", "cat := v.getAllCategories")
	catHTML := VS.toViewCategories(c, cat)
	types := VS.getAllCategoryType(c, "h_data", "showCategories", "types := VS.getAllCategoryType")
	JSON, err := json.Marshal(catHTML)
	VS.checkOperation(c, "h_date", "showCategories", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"cat":   catHTML,
		"types": types,
		"JSON":  string(JSON),
	}, "categories.html")
}

// Shows the page to introduce the information of the new Category
func (VS *Server) showNewCategory(c *gin.Context) {
	types := VS.getAllCategoryType(c, "h_data", "showNewCategory", "types := v.getAllCategoryType")
	sType := VS.getAllServiceTypes(c, "h_data", "showNewCategory", "sType := VS.getAllServiceTypes")
	servicesTypeHTML := VS.servicesTypeToServicesTypeHTML(c, sType)
	render(c, gin.H{
		"types": types,
		"sType": servicesTypeHTML,
	}, "form-category.html")
}

// Stores the new Category information into the DataBase
func (VS *Server) newCategory(c *gin.Context) {
	// We only can introduce new information if the Language is the main Language
	if getLanguage(c) != MainLanguage {
		VS.wrongFeedback(c, "To add a Category in a Language that is not the Main language of the app, you need to first add the step in the main language, change to the language desired for the modification, and modify the step that you want to translate")
		return
	}
	newCategory, b := VS.getNewCategory(c, bson.NewObjectId())
	if !b {
		return
	}
	if VS.addCategoryToDatabase(c, newCategory, "h_data", "newCategory", "VS.addCategoryToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, newCategory, newCategory.ID)
		VS.goodFeedback(c, "/data/categories")
	} else {
		VS.wrongFeedback(c, "Sorry, there is already one Category with that name related with this Type, .Please choose another name or Type")
	}
}

// Shows the page to modify the information of the Category
func (VS *Server) showEditCategory(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	cat := VS.getCategoryByID(c, id, "h_data", "showEditCategory", "cat := VS.getCategoryByID")
	trans, b := VS.mayLookForTranslation(c, cat.ID)
	if b {
		cat = trans.(Category)
	}
	types := VS.getAllCategoryType(c, "h_data", "showEditCategory", "types := v.getAllCategoryType")
	sType := VS.getAllServiceTypes(c, "h_data", "showEditCategory", "sType := VS.getAllServiceTypes")
	servicesTypeHTML := VS.servicesTypeToServicesTypeHTML(c, sType)
	render(c, gin.H{
		"category": cat,
		"types":    types,
		"edit":     true,
		"sType":    servicesTypeHTML,
	}, "form-category.html")
}

// Modifies the Category information into the DataBase
func (VS *Server) editCategory(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldCat := VS.getCategoryByID(c, id, "h_data", "editCategory", "oldCat := VS.getCategoryByID")
	newCategory, b := VS.getNewCategory(c, id)
	if !b {
		return
	}
	// Check if we are editing another language than the main one
	if VS.isTranslation(c, newCategory, newCategory.ID) {
		return
	}
	newCategory.Items = oldCat.Items

	// Audit
	audit := VS.getAuditFromCentral(c, newCategory, id, Modified)
	modified := VS.dynamicAuditChange(c, newCategory, oldCat, audit.ID)
	oldType := VS.getCategoryTypeByID(c, oldCat.Type, "h_data", "editCategory", "oldType := VS.getCategoryTypeByID").Name
	newType := VS.getCategoryTypeByID(c, newCategory.Type, "h_data", "editCategory", "newType := VS.getCategoryTypeByID").Name
	if VS.checkChange(c, audit.ID, "OnlyRent", strconv.FormatBool(newCategory.IsTrackable), strconv.FormatBool(newCategory.IsTrackable)) {
		modified = true
	}
	if VS.checkChange(c, audit.ID, "Type", oldType, newType) {
		VS.move(basePath+"/local-resources/categories/"+oldType+"/"+oldCat.Name+"/", basePath+"/local-resources/categories/"+newType+"/"+newCategory.Name+"/")
		modified = true
	}
	if VS.checkChange(c, audit.ID, "Name", oldCat.Name, newCategory.Name) {
		VS.move(basePath+"/local-resources/categories/"+oldType+"/"+oldCat.Name+"/", basePath+"/local-resources/categories/"+newType+"/"+newCategory.Name+"/")
		modified = true
	}
	if modified {
		if VS.addCategoryToDatabase(c, newCategory, "h_data", "editCategory", "VS.addCategoryToDatabase") {
			VS.registerAuditAndEvent(c, audit, newCategory, newCategory.ID)
		} else {
			VS.wrongFeedback(c, "Sorry, there is already one Category with that name related with this Type, .Please choose another name or Type")
		}
	}
	VS.goodFeedback(c, "/data/categories")
}

// Shows the page to introduce the information of the new CategoryType
func (VS *Server) showNewType(c *gin.Context) {
	render(c, gin.H{
		"edit": false,
	}, "form-type.html")
}

// Stores the new CategoryType information into the DataBase
func (VS *Server) newType(c *gin.Context) {
	categoryType, b := VS.getNewCategoryType(c, bson.NewObjectId())
	if !b {
		return
	}
	if VS.addCategoryTypeToDatabase(c, categoryType, "h_data", "newType", "VS.addCategoryTypeToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, categoryType, categoryType.ID)
		VS.goodFeedback(c, "/data/categories")
	} else {
		VS.wrongFeedback(c, "Sorry, there is already one Category Type with that name. Please choose another name ")
	}
}

// Shows the page to modify the information of the new Category Type
func (VS *Server) showEditType(c *gin.Context) {
	//Get the id of the Type from the Link
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	catType := VS.getCategoryTypeByID(c, id, "h_data", "showEditType", "catType := VS.getCategoryTypeByID")
	render(c, gin.H{
		"edit":    true,
		"catType": catType,
	}, "form-type.html")
}

// Modifies the Category Type information into the DataBase
func (VS *Server) editType(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldCategoryType := VS.getCategoryTypeByID(c, id, "h_data", "editType", "oldCategoryType := VS.getCategoryTypeByID")
	newCategoryType, b := VS.getNewCategoryType(c, id)
	if !b {
		return
	}
	audit := VS.getAuditFromCentral(c, newCategoryType, id, Modified)
	modified := VS.checkChange(c, audit.ID, "Name", oldCategoryType.Name, newCategoryType.Name)
	if modified {
		// If we changed the name, we move the images of folder
		VS.move(basePath+"/local-resources/categories/"+oldCategoryType.Name+"/", basePath+"/local-resources/categories/"+newCategoryType.Name+"/")
		if VS.addCategoryTypeToDatabase(c, newCategoryType, "h_data", "editType", "VS.addCategoryTypeToDatabase") {
			VS.registerAuditAndEvent(c, audit, newCategoryType, newCategoryType.ID)
		} else {
			VS.wrongFeedback(c, "Sorry, there is already one Category Type with that name. Please choose another name ")
		}
	}
	VS.goodFeedback(c, "/data/categories")
}

// Shows all the items of the category selected
func (VS *Server) showItems(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	cat := VS.getCategoryByID(c, id, "h_data", "showItems", "cat := VS.getCategoryByID")
	items := VS.getAllItemsByCategory(c, id, "h_data", "showItems", "items := VS.getAllItemsByCategory")
	JSON, err := json.Marshal(items)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"items": items,
		"cat":   cat,
		"JSON":  string(JSON),
	}, "items-category.html")
}

// Shows the page to introduce the information of the new Item
func (VS *Server) showNewItem(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	cat := VS.getCategoryByID(c, id, "h_data", "showNewItem", "cat := v.getCategoryByID")
	render(c, gin.H{
		"cat": cat,
	}, "form-item.html")
}

// Stores the new Item into the DataBase
func (VS *Server) newItem(c *gin.Context) {
	catID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	newItem, b := VS.getNewItem(c, bson.NewObjectId(), catID)
	if !b {
		return
	}
	if VS.addItemToDatabase(c, newItem, "h_data", "newItem", "v.addItemToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, newItem, newItem.ID)
		VS.updateCategoryItems(c, catID, newItem)
		VS.goodFeedback(c, "/data/see-items/"+catID.Hex())
	} else {
		VS.wrongFeedback(c, "Sorry, there is already one Item with that name in the same category. Please choose another name or category")
	}
}

// Shows the page to modify the information of the Item
func (VS *Server) showEditItem(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	item, _ := VS.getItemByID(c, id, "h_data", "showEditItem", "item, _ := VS.getItemByID", true)
	trans, b := VS.mayLookForTranslation(c, item.ID)
	if b {
		item = trans.(Item)
	}
	cat := VS.getCategoryByID(c, item.IDCategory, "h_data", "showEditItem", "cat := v.getCategoryByID")
	categories := VS.getAllCategories(c, "h_data", "showEditItem", "categories := VS.getAllCategories")
	render(c, gin.H{
		"item":       item,
		"edit":       true,
		"cat":        cat,
		"categories": categories,
	}, "form-item.html")
}

// Modifies the Item into the DataBase
func (VS *Server) editItem(c *gin.Context) {
	itemID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldItem, _ := VS.getItemByID(c, itemID, "h_data", "editItem", "oldItem := VS.getItemByID", true)
	newItem, b := VS.getNewItem(c, itemID, oldItem.IDCategory)
	if !b {
		return
	}
	// Check if we are editing another language than the main one
	if VS.isTranslation(c, newItem, newItem.ID) {
		return
	}
	// Audit
	audit := VS.getAuditFromCentral(c, newItem, itemID, Modified)
	modified := VS.dynamicAuditChange(c, newItem, oldItem, audit.ID)

	// Here manually the changes that can not be automatized
	var oldCategory, newCategory string
	oldCategory = VS.getCategoryByID(c, oldItem.IDCategory, "h_data", "editItem", "oldCategory = VS.getCategoryByID").Name
	newCategory = VS.getCategoryByID(c, newItem.IDCategory, "h_data", "editItem", "newCategory = VS.getCategoryByID").Name
	if VS.checkChange(c, audit.ID, "CategoryID", oldCategory, newCategory) {
		modified = true
	}
	if modified {
		if VS.addItemToDatabase(c, newItem, "h_data", "editItem", "VS.addItemToDatabase") {
			VS.updateCategoryItems(c, oldItem.IDCategory, newItem)
			VS.registerAuditAndEvent(c, audit, newItem, newItem.ID)
		} else {
			VS.wrongFeedback(c, "Sorry, there is already one Item with that name in the same category. Please choose another name or category")
		}
	}
	VS.goodFeedback(c, "/data/see-items")
}
