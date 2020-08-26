package village

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Creates the User with the information given and from the HTML and returns the object
// It takes care about setting up the username
func (VS *Server) getNewUser(c *gin.Context, id bson.ObjectId, oldPhoto string, balance float64) (User, bool) {
	username, b := VS.getStringFromHTML(c, "username", "User")
	if !b {
		return User{}, false
	}
	roleString, b := VS.getStringFromHTML(c, "role", "User")
	if !b {
		return User{}, false
	}
	role := Role(roleString)
	villageID, b := VS.getIDFromHTML(c, "village", "User")
	if !b {
		return User{}, false
	}
	name, b := VS.getStringFromHTML(c, "name", "User")
	if !b {
		return User{}, false
	}
	surname, b := VS.getStringFromHTML(c, "surname", "User")
	if !b {
		return User{}, false
	}
	age, b := VS.getIntFromHTML(c, "age", "User")
	if !b {
		return User{}, false
	}
	genderString, b := VS.getStringFromHTML(c, "gender", "User")
	gender := GenderType(genderString)
	tribeString, b := VS.getStringFromHTML(c, "tribe", "User")
	if !b {
		return User{}, false
	}
	story, b := VS.getStringFromHTML(c, "description", "User")
	if !b {
		return User{}, false
	}

	// Manager & Administrator are non-village related users, but the workers are village-related
	if role == Worker {
		newUsername := VS.getVillageByID(c, villageID, "get_HTML", "getNewUser", "newUsername = v.getVillageByID").Prefix
		username = newUsername + "_" + username
	} else {
		villageID = "" // VS.getCentral(c, "h_data", "newUser", "newUsername = v.getVillageByID").ID
	}
	// We need to put this step after the username checking
	var photo string
	if oldPhoto == "" {
		photo, b = VS.getPhotoFromHTML(c, "photo", username, basePath+"/local-resources/users/"+username+"/")
		if !b {
			return User{}, false
		}
	} else {
		photo = oldPhoto
	}

	newUser := User{
		ID:        id,
		Name:      name,
		Surname:   surname,
		Age:       age,
		Gender:    gender,
		Photo:     photo,
		Tribe:     TribeType(tribeString),
		Role:      role,
		Username:  username,
		Balance:   balance,
		IDVillage: villageID,
		Story:     story,
	}
	return newUser, true
}

// Creates the Village with the information given and from the HTML and returns the object
func (VS *Server) getNewVillage(c *gin.Context, id bson.ObjectId, prefix string) (Village, bool) {
	name, b := VS.getStringFromHTML(c, "name", "Village")
	if !b {
		return Village{}, false
	}
	newVillage := Village{
		ID:     id,
		Name:   name,
		Prefix: prefix,
	}
	return newVillage, true
}

// Creates the ServiceType with the information given and from the HTML and returns the object
func (VS *Server) getNewServiceType(c *gin.Context, id bson.ObjectId) (ServiceType, bool) {
	name, b := VS.getStringFromHTML(c, "name", "Village")
	if !b {
		return ServiceType{}, false
	}
	photo, b := VS.getPhotoFromHTML(c, "photo", name, basePath+"/local-resources/services/")
	if !b {
		return ServiceType{}, false
	}
	newServiceType := ServiceType{
		ID:                id,
		Name:              name,
		Icon:              photo,
		AllowNoQRPurchase: VS.getCheckBoxFromHTML(c, "allowfree"),
	}
	return newServiceType, true
}

// Creates the Service with the information given and from the HTML and returns the object
func (VS *Server) getNewService(c *gin.Context, id bson.ObjectId) (Service, bool) {
	villageID, b := VS.getIDFromHTML(c, "village", "Service")
	if !b {
		return Service{}, false
	}
	name, b := VS.getStringFromHTML(c, "name", "Service")
	if !b {
		return Service{}, false
	}
	balance, b := VS.getFloatFromHTML(c, "balance", "Service")
	if !b {
		return Service{}, false
	}
	serviceTypeID, b := VS.getIDFromHTML(c, "servicetype", "Service Type")
	if !b {
		return Service{}, false
	}
	fmt.Print("Balance in the HTML is: ")
	fmt.Println(balance)

	newService := Service{
		ID:        id,
		Name:      name,
		Balance:   balance,
		IDVillage: villageID,
		Type:      serviceTypeID,
	}
	return newService, true
}

// Creates the Category with the information given and from the HTML and returns the object
func (VS *Server) getNewCategory(c *gin.Context, id bson.ObjectId) (Category, bool) {
	name, b := VS.getStringFromHTML(c, "name", "Category")
	if !b {
		return Category{}, false
	}
	cTypeID, b := VS.getIDFromHTML(c, "type", "Category")
	if !b {
		return Category{}, false
	}
	cType := VS.getCategoryTypeByID(c, cTypeID, "get_HTML", "getNewCategory", "cType := VS.getCategoryTypeByID")
	photo, b := VS.getFileFromHTML(c, "photo", name, basePath+"/local-resources/categories/"+cType.Name+"/"+name+"/")
	if !b {
		return Category{}, false
	}
	// Here for taking the different values of the services checkboxes
	servicesWithAccess, b := VS.getArrayIDFromHTML(c, "service", "Service")
	if !b {
		return Category{}, false
	}
	description, b := VS.getStringFromHTML(c, "description", "Category")
	if !b {
		return Category{}, false
	}
	timecheckingString, b := VS.getStringFromHTML(c, "timechecking", "Category")
	if !b {
		return Category{}, false
	}
	typeofitemString, b := VS.getStringFromHTML(c, "itemtype", "Category")
	if !b {
		return Category{}, false
	}

	category := Category{
		ID:             id,
		Name:           name,
		Description:    description,
		Photo:          photo,
		TimeChecking:   TimecheckingType(timecheckingString),
		Type:           cTypeID,
		ServicesAccess: servicesWithAccess,
		IsTrackable:    VS.getCheckBoxFromHTML(c, "istrackable"),
		TypeOfItem:     TypeOfItem(typeofitemString),
	}

	return category, true
}

// Creates the CategoryType with the information given and from the HTML and returns the object
func (VS *Server) getNewCategoryType(c *gin.Context, id bson.ObjectId) (CategoryType, bool) {
	name, b := VS.getStringFromHTML(c, "name", "Category Type")
	if !b {
		return CategoryType{}, false
	}

	tpe := CategoryType{
		ID:   id,
		Name: name,
	}
	return tpe, true
}

// Creates the Item with the information given and from the HTML and returns the object
func (VS *Server) getNewItem(c *gin.Context, id bson.ObjectId, idCategory bson.ObjectId) (Item, bool) {

	catID, b := VS.getIDFromHTMLWithoutChecking(c, "category")
	var cat Category
	if !b {
		// If new
		cat = VS.getCategoryByID(c, idCategory, "get_HTML", "getNewItem", "cat := VS.getCategoryByID")
	} else {
		// If edit
		cat = VS.getCategoryByID(c, catID, "get_HTML", "getNewItem", "cat := VS.getCategoryByID")
	}
	catType := VS.getCategoryTypeByID(c, cat.Type, "get_HTML", "getNewItem", "catType := VS.getCategoryTypeByID")
	name, b := VS.getStringFromHTML(c, "name", "Item")
	if !b {
		return Item{}, false
	}
	photo, b := VS.getPhotoFromHTML(c, "photo", name, basePath+"/local-resources/categories/"+catType.Name+"/"+cat.Name+"/")
	if !b {
		return Item{}, false
	}
	price, b := VS.getFloatFromHTML(c, "cost", "Item")
	if !b {
		return Item{}, false
	}
	description, b := VS.getStringFromHTML(c, "description", "Item")
	if !b {
		return Item{}, false
	}
	unittype, b := VS.getStringFromHTML(c, "unittype", "Category")
	if !b {
		return Item{}, false
	}

	item := Item{
		ID:          id,
		Name:        name,
		Price:       price,
		Photo:       photo,
		Description: description,
		IDCategory:  cat.ID,
		UnitType:    UnitType(unittype),
	}
	return item, true
}

// Creates the VideoCourse with the information given and from the HTML and returns the object
func (VS *Server) getNewVideoCourse(c *gin.Context, id bson.ObjectId) (VideoCourse, bool) {
	name, b := VS.getStringFromHTML(c, "name", "VideoCourse")
	if !b {
		return VideoCourse{}, false
	}
	photo, b := VS.getFileFromHTML(c, "photo", name, basePath+"/local-resources/video-courses/"+name+"/")
	if !b {
		return VideoCourse{}, false
	}
	description, b := VS.getStringFromHTML(c, "description", "VideoCourse")
	if !b {
		return VideoCourse{}, false
	}
	price, b := VS.getFloatFromHTML(c, "price", "VideoCourse")
	if !b {
		return VideoCourse{}, false
	}
	servicesWithAccess, b := VS.getArrayIDFromHTML(c, "service", "Service")
	if !b {
		return VideoCourse{}, false
	}
	video := VideoCourse{
		ID:             id,
		Name:           name,
		Description:    description,
		Price:          price,
		Photo:          photo,
		ServicesAccess: servicesWithAccess,
	}
	return video, true
}

// Creates the ProductCheckList with the information given and from the HTML and returns the object
func (VS *Server) getNewProductCheckList(c *gin.Context, idCheck bson.ObjectId, idProduct bson.ObjectId) (VideoCheckList, bool) {
	product, b := VS.getVideoCourseByID(c, idProduct, "handlers_workshop", "showNewStep", "products := v.getAllProducts", true)
	if !b {
		return VideoCheckList{}, false
	}
	photo, b := VS.getFileFromHTML(c, "photo", idCheck.Hex(), basePath+"/local-resources/video-courses/"+product.Name+"/checks/")
	if !b {
		return VideoCheckList{}, false
	}
	audio, b := VS.getFileFromHTML(c, "audio", idCheck.Hex(), basePath+"/local-resources/video-courses/"+product.Name+"/checks/")
	if !b {
		return VideoCheckList{}, false
	}
	description, b := VS.getStringFromHTML(c, "description", "Product CheckList")
	if !b {
		return VideoCheckList{}, false
	}

	check := VideoCheckList{
		ID:            idCheck,
		IDVideoCourse: idProduct,
		Description:   description,
		Audio:         audio,
		Photo:         photo,
	}
	return check, true
}

// Creates the Step with the information given and from the HTML and returns the object
func (VS *Server) getNewStep(c *gin.Context, id bson.ObjectId) (Step, bool) {
	prodID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return Step{}, false
	}
	product, b := VS.getVideoCourseByID(c, prodID, "get_HTML", "getNewStep", "products := VS.getVideoCourseByID", true)
	if !b {
		return Step{}, false
	}
	indexOrder, b := VS.getIntFromHTML(c, "indexOrder", "Step")
	if !b {
		return Step{}, false
	}
	video, b := VS.getFileFromHTML(c, "video", strconv.Itoa(indexOrder), basePath+"/local-resources/video-courses/"+product.Name+"/steps/")
	if !b {
		return Step{}, false
	}
	// audio, b := VS.getFileFromHTML(c, "audio", strconv.Itoa(indexOrder), basePath+"/local-resources/video-courses/"+product.Name+"/steps/")
	// if !b {
	// 	return Step{}, false
	// }
	title, b := VS.getStringFromHTML(c, "title", "Step")
	if !b {
		return Step{}, false
	}
	description, b := VS.getStringFromHTML(c, "description", "Step")
	if !b {
		return Step{}, false
	}

	// Here for taking the different values of the checkboxes with the same name
	tools := c.Request.MultipartForm.Value["tool"]
	materials := c.Request.MultipartForm.Value["material"]
	warnings := c.Request.MultipartForm.Value["warning"]

	// // Warnings
	// var warnings []string
	// var w string
	// for i := 1; i <= 6; i++ {
	// 	w = c.PostForm("warning" + strconv.Itoa(i))
	// 	if w != "" {
	// 		warnings = append(warnings, w)
	// 	}
	// }
	// Now we transform into ID the Items
	var idTools []bson.ObjectId
	for _, tool := range tools {
		idTools = append(idTools, bson.ObjectIdHex(tool))
	}
	// Now we transform into ID the materials
	var idMaterials []bson.ObjectId
	for _, mat := range materials {
		idMaterials = append(idMaterials, bson.ObjectIdHex(mat))
	}
	newStep := Step{
		ID:              id,
		IDVideoCourse:   prodID,
		IndexOrder:      indexOrder,
		Title:           title,
		Description:     description,
		Warnings:        warnings,
		Video:           video,
		ToolsNeeded:     idTools,
		Audio:           video,
		MaterialsNeeded: idMaterials,
	}
	return newStep, true
}

// Creates the Service Order with the information given and from the HTML and returns the object
func (VS *Server) getNewServiceOrder(c *gin.Context, id bson.ObjectId) (ServiceOrder, bool) {
	date := VS.getDateFromHTML(c, "deadline", "Workshop Order")
	videoID, b := VS.getIDFromHTML(c, "productID", "Workshop Order")
	if !b {
		return ServiceOrder{}, false
	}
	serviceID, b := VS.getIDFromHTML(c, "workshopID", "Workshop Order")
	if !b {
		return ServiceOrder{}, false
	}
	quantity, b := VS.getIntFromHTML(c, "quantity", "Workshop Order")
	if !b {
		return ServiceOrder{}, false
	}
	window, b := VS.getIntFromHTML(c, "window", "Workshop Order")
	if !b {
		return ServiceOrder{}, false
	}
	sOrder := ServiceOrder{
		ID:            id,
		IDVideoCourse: videoID,
		IDService:     serviceID,
		Quantity:      quantity,
		AlreadyMade:   0,
		Assigned:      0,
		Deadline:      date,
		WindowPeriod:  WindowPeriodType(window),
	}
	return sOrder, true
}
