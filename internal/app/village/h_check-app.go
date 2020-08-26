package village

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// This handler shows the form for receiving the delivery for the Manager in the service
func (VS *Server) showConfirmCheckApp(c *gin.Context) {
	sType, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	serviceHTML, _ := VS.getOurService(c, "h_check-app", "showConfirmCheckApp", "serviceHTML, _ := VS.getOurService")
	stocks := VS.getAllPositiveStockByServiceIDForMakeChecking(c, serviceHTML.ID, "h_check-app", "showConfirmCheckApp", "stocks := VS.getAllPositiveStockByServiceIDForMakeChecking")
	stockHTML := VS.stocksToHTML(c, stocks)
	render(c, gin.H{
		"service":      serviceHTML,
		"serviceType":  sType,
		"itemsChecked": stockHTML,
	}, "check-app-manager.html")
}

// This handler adds the items form the delivery to the service
func (VS *Server) confirmCheckApp(c *gin.Context) {
	serviceHTML, _ := VS.getOurService(c, "h_check-app", "confirmCheckApp", "serviceHTML, _ := VS.getOurService")
	newCheck := CheckStock{
		ID:        bson.NewObjectId(),
		Date:      getTime(),
		IDService: serviceHTML.ID,
		LastPhoto: 0,
	}
	// Comprobar la checkbox
	QR, isQR := VS.getIDFromHTMLWithoutChecking(c, "qr")
	if isQR {
		oldCheck, b := VS.getChecksFromItemAndServiceID(c, QR, serviceHTML.ID, "h_check-app", "confirmCheckApp", "oldCheck, b := VS.getChecksFromItemAndServiceID")
		if b {
			newCheck.ID = oldCheck.ID
			VS.editElementInDatabaseWithRegister(c, newCheck, oldCheck, oldCheck.ID, serviceHTML.ID, "h_check-app", "confirmCheckApp", "VS.editElementInDatabaseWithRegister")
		} else {
			newCheck.IDInstance = QR
			VS.addElementToDatabaseWithRegister(c, newCheck, newCheck.ID, serviceHTML.ID, "h_check-app", "confirmCheckApp", "VS.addElementToDatabaseWithRegister")
		}
	}
	itemsReceived, _ := VS.getArrayIDFromHTMLWithoutChecking(c, "item", "Delivery")
	for _, i := range itemsReceived {
		oldCheck, isOldCheck := VS.getChecksFromItemAndServiceID(c, i, serviceHTML.ID, "h_check-app", "confirmCheckApp", "oldCheck, isOldCheck := VS.getChecksFromItemAndServiceID")
		lastPhoto := oldCheck.LastPhoto
		if !isOldCheck {
			lastPhoto = newCheck.LastPhoto
			newCheck.ID = bson.NewObjectId()
		}
		newCheck.IDInstance = i
		switch lastPhoto {
		case 0:
			newCheck.Photo1, _ = VS.getPhotoFromHTML(c, i.Hex(), i.Hex()+"1", basePath+"/local/services/"+serviceHTML.Name+"/checks/")
		case 1:
			newCheck.Photo2, _ = VS.getPhotoFromHTML(c, i.Hex(), i.Hex()+"2", basePath+"/local/services/"+serviceHTML.Name+"/checks/")
		case 2:
			newCheck.Photo3, _ = VS.getPhotoFromHTML(c, i.Hex(), i.Hex()+"3", basePath+"/local/services/"+serviceHTML.Name+"/checks/")
		case 3:
			newCheck.Photo4, _ = VS.getPhotoFromHTML(c, i.Hex(), i.Hex()+"4", basePath+"/local/services/"+serviceHTML.Name+"/checks/")
		}
		newCheck.LastPhoto = lastPhoto + 1
		if isQR {
			newCheck.ID = oldCheck.ID
			VS.editElementInDatabaseWithRegister(c, newCheck, oldCheck, oldCheck.ID, serviceHTML.ID, "h_check-app", "confirmCheckApp", "VS.addElementToDatabaseWithRegister")
		} else {
			newCheck.IDInstance = QR
			VS.addElementToDatabaseWithRegister(c, newCheck, newCheck.ID, serviceHTML.ID, "h_check-app", "confirmCheckApp", "VS.addElementToDatabaseWithRegister")
		}
	}
	// VS.goodFeedback(c)
}

// Dashboard
// This handler shows the items that we can make for the user or manager to check, is from the DASHBOARD
func (VS *Server) showItemsToCheck(c *gin.Context) {
	todos := VS.getAllToDo(c, "h_check-app", "showItemsToCheck", "todos := VS.getAllToDo")
	todosHTML := VS.todosToHTML(c, todos)
	JSON, err := json.Marshal(todosHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"todosHTML": todosHTML,
		"JSON":      string(JSON),
	}, "to-dos.html")
}

// This handler shows the form to introduce the new to-do information
func (VS *Server) showNewToDo(c *gin.Context) {
	users := VS.getAllUsersFromTheSystemNotAdmin(c, "h_check-app", "showNewToDo", "users := VS.getAllUsersFromTheSystemNotAdmin")
	render(c, gin.H{
		"users": users,
	}, "form-to-do.html")
}

// This handler creates a new todo
func (VS *Server) newToDo(c *gin.Context) {
	userID, b := VS.getIDFromHTML(c, "workerID", "To do")
	if !b {
		return
	}
	title, b := VS.getStringFromHTML(c, "title", "To do")
	if !b {
		return
	}

	description, b := VS.getStringFromHTML(c, "description", "To do")
	if !b {
		return
	}

	timechecking, b := VS.getStringFromHTML(c, "timechecking", "To do")
	if !b {
		return
	}
	istrackable := VS.getCheckBoxFromHTML(c, "istrackable")
	id := bson.NewObjectId()
	if istrackable {
		id, b = VS.getQRFromHTMLNew(c, Task)
	}
	if !b {
		return
	}

	descriptions := c.PostFormArray("descriptionfield")
	numbers := c.PostFormArray("numbers")
	titlePhotos := c.PostFormArray("phototitle")
	user := VS.getUserByID(c, userID, "h_check-app", "newToDo", "user := VS.getUserByID")
	photos, _ := VS.getArrayPhotosFromHTML(c, "photo", basePath+"/local-resources/users/"+user.Username+"/to-dos/")
	checkboxes := c.PostFormArray("checkboxes")
	villageID := VS.getUserByID(c, userID, "h_check-app", "newToDo", "villageID := VS.getUserByID").IDVillage
	todo := ToDo{
		ID:              id,
		TitleToDo:       title,
		DescriptionToDo: description,
		IDUser:          userID,
		IDVillage:       villageID,
		TimeChecking:    TimecheckingType(timechecking),
		IsTrackable:     istrackable,
	}
	for _, ck := range checkboxes {
		newCheckbox := FieldCheckbox{
			ID:            bson.NewObjectId(),
			IDToDo:        todo.ID,
			TitleCheckbox: ck,
		}
		todo.Checkboxes = append(todo.Checkboxes, newCheckbox)
	}
	for index, p := range titlePhotos {
		newPhoto := FieldPhoto{
			ID:         bson.NewObjectId(),
			IDToDo:     todo.ID,
			TitlePhoto: p,
			Photo:      photos[index],
		}
		todo.Photos = append(todo.Photos, newPhoto)
	}
	for _, n := range numbers {
		newNumber := FieldNumber{
			ID:          bson.NewObjectId(),
			IDToDo:      todo.ID,
			TitleNumber: n,
		}
		todo.Numbers = append(todo.Numbers, newNumber)
	}
	for _, d := range descriptions {
		newDescription := FieldDescription{
			ID:          bson.NewObjectId(),
			IDToDo:      todo.ID,
			Description: d,
		}
		todo.Descriptions = append(todo.Descriptions, newDescription)
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, todo, todo.ID, "h_check-app", "newToDo", "VS.addElementToDatabaseWithRegisterFromCentral")
	VS.goodFeedback(c, "/to-do/what-to-check")
}

// This handler shows the items that we can make for the user or manager to check, is from the DASHBOARD
func (VS *Server) showToDosUser(c *gin.Context) {
	u := VS.getUserFomCookie(c, "h_check-app", "showToDosUser", "u := VS.getUserFomCookie")
	todos := VS.getToDoByUserIDForMakeChecking(c, u.ID, "h_check-app", "showToDosUser", "todos := VS.getToDoByUserIDForMakeChecking")
	todosHTML := VS.todosToHTML(c, todos)
	render(c, gin.H{
		"todosHTML": todosHTML,
	}, "to-dos-user.html")
}

// This handler shows the items that we can make for the user or manager to check, is from the DASHBOARD
func (VS *Server) showMakeToDo(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	todo := VS.getTodoByID(c, id, "h_check-app", "showMakeToDo", "todo := VS.getTodoByID(")
	todoHTML := VS.todoToHTML(c, todo)
	render(c, gin.H{
		"todo": todoHTML,
	}, "make-to-do.html")
}

// This handler shows the items that we can make for the user or manager to check, is from the DASHBOARD
func (VS *Server) makeToDo(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	todo := VS.getTodoByID(c, id, "h_check-app", "makeToDo", "todo := VS.getTodoByID")
	user := VS.getUserByID(c, todo.IDUser, "h_check-app", "makeToDo", "user := VS.getUserByID")

	descriptions := c.PostFormArray("descriptions")
	numbers := c.PostFormArray("numbers")
	photos, _ := VS.getArrayPhotosFromHTML(c, "photo", basePath+"/local-resources/users/"+user.Username+"/to-dos/checked/")
	// checkboxes := VS.getArrayCheckBoxesFromHTML(c, "checkboxes")
	checked := ToDoChecked{
		ID:     bson.NewObjectId(),
		Date:   getTime(),
		IDToDo: todo.ID,
	}

	if todo.IsTrackable {
		qr, b := VS.getQRFromHTML(c)
		if !b {
			VS.wrongFeedback(c, "It has been a problem with the QR code")
			return
		}
		if qr != todo.ID {
			VS.wrongFeedback(c, "The QR code scan is not the QR of this task")
			return
		}
	}

	for index, ck := range todo.Checkboxes {
		newCheckbox := FieldCheckbox{
			ID:            bson.NewObjectId(),
			IDToDo:        checked.ID,
			TitleCheckbox: todo.Checkboxes[index].TitleCheckbox,
			Checkbox:      VS.getCheckBoxFromHTML(c, hex.EncodeToString([]byte(ck.ID))),
		}
		checked.Checkboxes = append(checked.Checkboxes, newCheckbox)
	}
	for index, p := range photos {
		fmt.Print("photo made: " + strconv.Itoa(index))

		newPhoto := FieldPhoto{
			ID:         bson.NewObjectId(),
			IDToDo:     checked.ID,
			TitlePhoto: p,
			Photo:      photos[index],
		}
		checked.Photos = append(checked.Photos, newPhoto)
	}
	for index, n := range numbers {
		fmt.Print("number made: " + strconv.Itoa(index))

		numb, _ := strconv.Atoi(n)
		newNumber := FieldNumber{
			ID:          bson.NewObjectId(),
			IDToDo:      checked.ID,
			TitleNumber: todo.Numbers[index].TitleNumber,
			Number:      numb,
		}
		checked.Numbers = append(checked.Numbers, newNumber)
	}
	for index, d := range descriptions {
		fmt.Print("description made: " + strconv.Itoa(index))
		newDescription := FieldDescription{
			ID:          bson.NewObjectId(),
			IDToDo:      checked.ID,
			Description: d,
		}
		checked.Descriptions = append(checked.Descriptions, newDescription)
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, checked, checked.ID, "h_check-app", "makeToDo", "VS.addElementToDatabaseWithRegisterFromCentral")
	VS.goodFeedback(c, "to-do/to-dos-user")
}

// This handler shows the items that we can make for the user or manager to check, is from the DASHBOARD
func (VS *Server) seeChecks(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	checks := VS.getTodoCheckedByToDoID(c, id, "h_check-app", "seeChecks", "checks := VS.getTodoCheckedByToDoID")

	checksHTML := VS.todoCheckedSToHTML(c, checks)
	render(c, gin.H{
		"checks": checksHTML,
	}, "see-checks.html")
}
