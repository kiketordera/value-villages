package village

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Gets the Checkbox information from the HTML in a POST request
func (VS *Server) getCheckBoxFromHTML(c *gin.Context, name string) bool {
	checkbox := c.Request.FormValue(name)
	fmt.Print("checkbox is: " + checkbox)
	if checkbox == "" {
		return false
	}
	return true
}

// Gets the String information of the HTML in a POST request
func (VS *Server) getStringFromHTML(c *gin.Context, name string, element string) (string, bool) {
	data := c.PostForm(name)
	// We remove the spaces before the first letter and after the last one
	dataClean := strings.TrimSpace(data)
	// Render error if the parameter was empty, this way we have a server-side validation
	if data == "" {
		render(c, gin.H{
			"text": "Sorry, there was a problem trying to process the " + name + " of the " + element,
		}, "wrong-feedback.html")
		return "", false
	}
	return dataClean, true
}

// Gets the Int information of the HTML in a POST request
func (VS *Server) getIntFromHTML(c *gin.Context, name string, element string) (int, bool) {
	data, err := strconv.Atoi(c.PostForm(name))
	VS.checkOperation(c, "h_view", "getIntFromHTML", "data, err := strconv.Atoi(c.PostForm(name))", err)
	// Render error if the parameter was empty, this way we have a server-side validation
	if string(data) == "" {
		render(c, gin.H{
			"text": "Sorry, there was a problem trying to process the " + name + " of the " + element,
		}, "wrong-feedback.html")
		return 0, false
	}
	return data, true
}

// Gets the Int information of the HTML in a POST request
func (VS *Server) getFloatFromHTML(c *gin.Context, name string, element string) (float64, bool) {
	data, err := strconv.ParseFloat(c.PostForm(name), 64)
	VS.checkOperation(c, "h_view", "getIntFromHTML", "data, err := strconv.Atoi(c.PostForm(name))", err)
	// Render error if the parameter was empty, this way we have a server-side validation
	if strconv.FormatFloat(data, 'f', -1, 64) == "" {
		render(c, gin.H{
			"text": "Sorry, there was a problem trying to process the " + name + " of the " + element,
		}, "wrong-feedback.html")
		return 0, false
	}
	return data, true
}

// Gets the Int information of the HTML in a POST request without checking for a correct value
func (VS *Server) getIntFromHTMLWithoutChecking(c *gin.Context, name string, element string) (int, bool) {
	data, err := strconv.Atoi(c.PostForm(name))
	VS.checkOperation(c, "h_view", "getIntFromHTMLWithoutChecking", "data, err := strconv.Atoi(c.PostForm(name))", err)
	// Render error if the parameter was empty, this way we have a server-side validation
	if string(data) == "" {
		return 0, false
	}
	return data, true
}

// Gets the Int information of the HTML in a POST request without checking for a correct value
func (VS *Server) getFloatFromHTMLWithoutChecking(c *gin.Context, name string, element string) (float64, bool) {
	data, err := strconv.ParseFloat(c.PostForm(name), 64)
	VS.checkOperation(c, "h_view", "getIntFromHTMLWithoutChecking", "data, err := strconv.Atoi(c.PostForm(name))", err)
	// Render error if the parameter was empty, this way we have a server-side validation
	if strconv.FormatFloat(data, 'f', -1, 64) == "" {
		return 0, false
	}
	return data, true
}

// Gets the File information of the HTML in a POST request with server validation
// and stores the file in the path given
func (VS *Server) getFileFromHTML(c *gin.Context, name string, IDelement string, directory string) (string, bool) {
	file, err := c.FormFile(name)
	VS.checkOperation(c, "h_view", "getFileFromHTML", "file, err := c.FormFile(name)", err)
	if VS.isEmpty(c, string(file.Filename), "Sorry, there was a problem trying to process the "+name+"  of the Element :(") {
		return "", false
	}
	if file.Filename != "" {
		file.Filename = IDelement + path.Ext(file.Filename)
		if err := c.SaveUploadedFile(file, directory+file.Filename); err != nil {
			os.MkdirAll(directory, os.ModePerm)
			if err = c.SaveUploadedFile(file, directory+file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				VS.wrongFeedback(c, "Error while uploading the File")
				return "", false
			}
		}
	}
	return file.Filename, true
}

// Gets the File information of the HTML in a POST request with server validation
// and stores the file in the path given
// returns the name of the file.
// IT IS ONLY USED FOR IMPORTING THE ZIP OF THE DATABASE FOR SYNC
// It check that the Village receiver in the sync is itself
func (VS *Server) getFileFromHTMLDB(c *gin.Context, name string, IDelement string, directory string) (string, bson.ObjectId, string, bool) {
	file, err := c.FormFile(name)
	VS.checkOperation(c, "h_view", "getFileFromHTML", "file, err := c.FormFile(name)", err)
	if VS.isEmpty(c, string(file.Filename), "Sorry, there was a problem trying to process the "+name+"  of the Element :(") {
		return "", bson.NewObjectId(), "", false
	}
	objectID := ""
	nameVillage := ""
	if file.Filename != "" {
		// Remove the Village for human identifier from the name of the File
		last := strings.LastIndex(file.Filename, "-")
		nameVillage = file.Filename[0:last]
		secondLast := strings.LastIndex(nameVillage, "-")
		nameVillage = nameVillage[secondLast+1 : len(nameVillage)]
		fmt.Println("This is the name of the Village: " + nameVillage)
		file.Filename = file.Filename[last+1 : len(file.Filename)]
		objectID = file.Filename[0 : len(file.Filename)-4]
		if !bson.IsObjectIdHex(objectID) {
			VS.wrongFeedback(c, "Error, the file has no ObjectID Identifier: "+objectID)
			return "", bson.NewObjectId(), "", false
		}
		if err := c.SaveUploadedFile(file, directory+file.Filename); err != nil {
			os.MkdirAll(directory, os.ModePerm)
			if err = c.SaveUploadedFile(file, directory+file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				VS.wrongFeedback(c, "Error while uploading the File")
				return "", bson.NewObjectId(), "", false
			}
		}
	}
	return file.Filename, bson.ObjectIdHex(objectID), nameVillage, true
}

// Gets the Photo information of the HTML in a POST request with server validation
func (VS *Server) getPhotoFromHTML(c *gin.Context, name string, IDelement string, directory string) (string, bool) {
	file, err := c.FormFile(name)
	// 	buffer := bufio.NewReader(ReadImage(file.Filename))
	//   _, err = buffer.Read(bytes)
	VS.checkOperation(c, "h_view", "getPhotoFromHTML", "file, err := c.FormFile(name)", err)
	if VS.isEmpty(c, file.Filename, "Sorry, there was a problem trying to process the "+name+"  of the Element :(") {
		return "", false
	}
	if file.Filename != "" {
		file.Filename = IDelement + path.Ext(file.Filename)
		if err := c.SaveUploadedFile(file, directory+file.Filename); err != nil {
			os.MkdirAll(directory, os.ModePerm)
			if err = c.SaveUploadedFile(file, directory+file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload photo err: %s", err.Error()))
				VS.wrongFeedback(c, "Error while uploading the Photo")
				return "", false
			}
		}
	}
	return file.Filename, true
}

// Gets the ArrayPhotos information of the HTML in a POST request with server validation
func (VS *Server) getArrayPhotosFromHTML(c *gin.Context, name string, directory string) ([]string, bool) {
	var photoNames []string
	// Multipart form
	form, err := c.MultipartForm()
	files := form.File[name]

	for _, file := range files {
		VS.checkOperation(c, "h_view", "getArrayPhotosFromHTML", "file, err := c.FormFile(name)", err)
		if VS.isEmpty(c, file.Filename, "Sorry, there was a problem trying to process the "+name+"  of the Element :(") {
			return photoNames, false
		}
		if file.Filename != "" {
			file.Filename = bson.NewObjectId().Hex() + path.Ext(file.Filename)
			if err := c.SaveUploadedFile(file, directory+file.Filename); err != nil {
				os.MkdirAll(directory, os.ModePerm)
				if err = c.SaveUploadedFile(file, directory+file.Filename); err != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("upload photo err: %s", err.Error()))
					VS.wrongFeedback(c, "Error while uploading the Photo")
					return photoNames, false
				}
			}
		}
		photoNames = append(photoNames, file.Filename)
	}

	return photoNames, true
}

// Gets the ID of the HTML in a POST request
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getIDFromHTML(c *gin.Context, name string, element string) (bson.ObjectId, bool) {
	var ID bson.ObjectId
	isID := c.PostForm(name)
	fmt.Print("this is plain")
	fmt.Print(isID)
	if bson.IsObjectIdHex(isID) {
		ID = bson.ObjectIdHex(isID)
	} else {
		VS.wrongFeedback(c, "Sorry, there is a problem with the "+name+" of the "+element+" :(")
		return ID, false
	}
	return ID, true
}

// Gets the information of the ID, only takes until the barrier and cleans the extra information
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getIDFromHTMLCleaningIt(c *gin.Context, barrier string, name string, element string) (bson.ObjectId, bool) {
	var ID bson.ObjectId
	isID := c.PostForm(name)
	i := strings.Index(isID, barrier)
	isID = isID[0:i]
	if bson.IsObjectIdHex(isID) {
		ID = bson.ObjectIdHex(isID)
	} else {
		VS.wrongFeedback(c, "Sorry, there is a problem with the "+element+"  ID :(")
		return ID, false
	}
	return ID, true
}

// Gets the information of the HTML in a POST request
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getIDFromHTMLWithoutChecking(c *gin.Context, name string) (bson.ObjectId, bool) {
	var ID bson.ObjectId
	isID := c.PostForm(name)
	if bson.IsObjectIdHex(isID) {
		ID = bson.ObjectIdHex(isID)
		return ID, true
	}
	return ID, false
}

// Gets the array ID of the HTML in a POST request
func (VS *Server) getArrayIDFromHTML(c *gin.Context, name string, element string) ([]bson.ObjectId, bool) {
	array := c.PostFormArray(name)
	var arrayID []bson.ObjectId
	for _, e := range array {
		if bson.IsObjectIdHex(e) {
			arrayID = append(arrayID, bson.ObjectIdHex(e))
		} else {
			VS.wrongFeedback(c, "Sorry, there is a problem with the "+element+"  ID :(")
			return arrayID, false
		}
	}
	if len(array) > 0 {
		return arrayID, true
	}
	VS.wrongFeedback(c, "Sorry, you did not select any Service in the "+element+"  ID :(")
	return arrayID, false
}

// Gets the ID information of the HTML in a POST request
func (VS *Server) getArrayIDFromHTMLWithoutChecking(c *gin.Context, name string, element string) ([]bson.ObjectId, bool) {
	array := c.PostFormArray(name)
	var arrayID []bson.ObjectId
	for _, e := range array {
		if bson.IsObjectIdHex(e) {
			arrayID = append(arrayID, bson.ObjectIdHex(e))
		} else {
			return arrayID, false
		}
	}
	if len(array) > 0 {
		return arrayID, true
	}
	return arrayID, false
}

// Create the human readable information for the Users element
func (VS *Server) usersToHTML(c *gin.Context, u []User) []UserHTML {
	var usersHTML []UserHTML
	for _, user := range u {
		usersHTML = append(usersHTML, VS.userToHTML(c, user))
	}
	return usersHTML
}

// Create the human readable information for the Users element
func (VS *Server) userToHTML(c *gin.Context, u User) UserHTML {
	var user UserHTML
	if u.Role == Worker {
		user.Village = VS.getVillageByID(c, u.IDVillage, "h_view", "userToHTML", "user.Village = VS.getVillageByID").Name
	}
	user.ID = u.ID
	user.Age = u.Age
	user.Balance = u.Balance
	user.Photo = "/local/users/" + u.Username + "/" + u.Photo
	user.Role = string(u.Role)
	user.Surname = u.Surname
	user.Name = u.Name
	user.Gender = string(u.Gender)
	user.Tribe = string(u.Tribe)
	user.Story = u.Story
	user.Username = u.Username
	user.IDVillage = u.IDVillage
	return user
}

// Transform the Services information from DataBase to Human readable (in array)
func (VS *Server) servicesToHTML(c *gin.Context, s []Service) []ServiceHTML {
	var servicesHTML []ServiceHTML
	for _, service := range s {
		servicesHTML = append(servicesHTML, VS.serviceToHTML(c, service))
	}
	return servicesHTML
}

// Transform the Services information from DataBase to Human readable (in individual)
func (VS *Server) serviceToHTML(c *gin.Context, s Service) ServiceHTML {
	var serviceHTML ServiceHTML
	serviceHTML.ID = s.ID
	serviceHTML.Name = s.Name
	serviceHTML.Balance = s.Balance
	serviceHTML.Type = VS.getServiceTypeByID(c, s.Type, "h_view", "serviceToHTML", "serviceHTML.Type = VS.getServiceTypeByID", true).Name
	serviceHTML.TypeID = VS.getServiceTypeByID(c, s.Type, "h_view", "serviceToHTML", "serviceHTML.TypeID = VS.getServiceTypeByID", true).ID
	serviceHTML.Village = VS.getVillageByID(c, s.IDVillage, "h_view", "serviceToHTML", "serviceHTML.Village = VS.getVillageByID").Name
	serviceHTML.VillageID = s.IDVillage.Hex()
	serviceHTML.DeliveryID = bson.NewObjectId()
	tpe := VS.getServiceTypeByID(c, s.Type, "handlers_view", "serviceToHTML", "tpe := VS.getServiceTypeByID", true)
	serviceHTML.Photo = "/local/services/" + tpe.Icon

	return serviceHTML
}

// Transform the Services Type information from DataBase to Human readable (in array)
func (VS *Server) servicesTypeToServicesTypeHTML(c *gin.Context, s []ServiceType) []ServiceTypeHTML {
	var servicesHTML []ServiceTypeHTML
	for _, service := range s {
		servicesHTML = append(servicesHTML, VS.serviceTypeToServiceTypeHTML(c, service))
	}
	return servicesHTML
}

// Transform the Services Type information from DataBase to Human readable (in individual)
func (VS *Server) serviceTypeToServiceTypeHTML(c *gin.Context, s ServiceType) ServiceTypeHTML {
	var serviceHTML ServiceTypeHTML
	serviceHTML.ID = s.ID
	serviceHTML.Name = s.Name
	serviceHTML.Icon = "/local/services/" + s.Icon
	serviceHTML.AllowNoQRPurchase = s.AllowNoQRPurchase
	role := getRole(c)
	if role == Admin {
		serviceHTML.Link = "/performance/all/" + s.ID.Hex()
	} else if role == Manager {
		serviceHTML.Link = "/performance/local/" + s.ID.Hex()
	} else if role == Worker {
		serviceHTML.Link = "/performance/self/" + s.ID.Hex()
	}
	return serviceHTML
}

// This methos transform the array of audits into a array of audits divided by the day of creation
func (VS *Server) auditPerDayToHTML(c *gin.Context, a [][]Audit) []AuditDayHTML {
	var aDay AuditDayHTML
	var auditDay []AuditDayHTML
	// previousDayAudit := a[i][0].Date.Date()
	if a == nil {
		return auditDay
	}
	if len(a) == 0 {
		return auditDay
	}
	for i := 0; i < len(a); i++ {
		aDay = (AuditDayHTML{})
		if len(a[i]) > 0 {
			y, m, d := time.Unix(a[i][0].Date, 0).Date()
			aDay.Date = time.Unix(a[i][0].Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
		}
		for j := 0; j < len(a[i]); j++ {
			if a[i][j].Description == Created {
				aDay.Created = append(aDay.Created, VS.auditToAuditHTML(c, a[i][j]))
			} else if a[i][j].Description == Modified {
				aDay.Modified = append(aDay.Modified, VS.auditToAuditHTML(c, a[i][j]))
			} else if a[i][j].Description == Deleted {
				aDay.Deleted = append(aDay.Deleted, VS.auditToAuditHTML(c, a[i][j]))
			}
		}
		auditDay = append(auditDay, aDay)
	}
	return auditDay
}

// This methos transform the array of messages into a array of messages divided by the day of creation
func (VS *Server) messagePerDayToHTML(c *gin.Context, idUser bson.ObjectId, a [][]Message) []MessageDayHTML {
	var aDay MessageDayHTML
	var auditDay []MessageDayHTML
	// previousDayAudit := a[i][0].Date.Date()
	if a == nil {
		return auditDay
	}
	if len(a) == 0 {
		return auditDay
	}
	for i := 0; i < len(a); i++ {
		aDay = (MessageDayHTML{})
		if len(a[i]) > 0 {
			y, m, d := time.Unix(a[i][0].Date, 0).Date()
			aDay.Date = time.Unix(a[i][0].Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
		}
		aDay.Messages = VS.messagesToHTML(c, idUser, a[i])
		auditDay = append(auditDay, aDay)
	}
	return auditDay
}

// This methos transform ONE audit into audit human readable information
func (VS *Server) auditToAuditHTML(c *gin.Context, a Audit) AuditHTML {
	var aHTML AuditHTML
	element := VS.getAnythingWithID(c, a.IDItem, "h_view", "auditToAuditHTML", "element := v.getAnythingWithID").Instance

	switch element.(type) {

	case ServiceOrder:
		s := element.(ServiceOrder)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		prod, _ := VS.getVideoCourseByID(c, s.IDVideoCourse, "h_view", "auditToAuditHTML", "prod, _ := VS.getVideoCourseByID", true)
		aHTML.Icon = "/local/video-courses/" + prod.Name + "/" + prod.Photo
		aHTML.Class = "image-square-rounded-table"
		aHTML.Name = VS.getServiceByID(c, s.IDService, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getServiceByID(c, s.ServiceID", true).Name
		aHTML.SecondLine = strconv.Itoa(s.Quantity)
		aHTML.Type = "WORKSHOP ORDER"
		return aHTML

	case Assignment:
		as := element.(Assignment)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		u := VS.getUserByID(c, as.IDWorker, "h_view", "auditToAuditHTML", "u := VS.getUserByID")
		aHTML.Icon = "/local/users/" + u.Username + "/" + u.Photo
		aHTML.Class = "image-rounded-table"
		item, _ := VS.getItemByID(c, as.IDStock, "h_view", "auditToAuditHTML", "item, _ := VS.getItemByID", true)
		aHTML.Name = item.Name
		aHTML.SecondLine = VS.getServiceByID(c, as.IDService, "h_view", "auditToAuditHTML", "aHTML.SecondLine = VS.getServiceByID", true).Name
		aHTML.Type = "Assignment"
		return aHTML

	case VideoCourse:
		p := element.(VideoCourse)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Class = "image-square-rounded-table"
		aHTML.Icon = "/local/video-courses/" + p.Name + "/" + p.Photo
		aHTML.Name = p.Name
		aHTML.Type = "VIDEO COURSE"
		return aHTML

	case ToDo:
		p := element.(ToDo)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Class = "image-square-rounded-table"
		aHTML.Icon = "/static/svg/calendar-color.svg"
		aHTML.Name = p.TitleToDo
		aHTML.Type = "TO DO"
		return aHTML

	case ToDoChecked:
		p := element.(ToDoChecked)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Class = "image-square-rounded-table"
		aHTML.Icon = "/static/svg/calendar-color.svg"
		todo := VS.getTodoByID(c, p.IDToDo, "h_view", "auditToAuditHTML", "todo := VS.getTodoByID")
		aHTML.Name = todo.TitleToDo
		aHTML.Type = "TO DO CHECKED"
		return aHTML

	case User:
		u := element.(User)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/local/users/" + u.Username + "/" + u.Photo
		aHTML.Class = "image-rounded-table"
		aHTML.Name = u.Name
		aHTML.Type = "USER"
		return aHTML

	case Category:
		cat := element.(Category)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		tpe := VS.getCategoryTypeByID(c, cat.Type, "h_view", "auditToAuditHTML", "tpe := VS.getCategoryTypeByID").Name
		aHTML.Icon = "/local/categories/" + tpe + "/" + cat.Name + "/" + cat.Photo
		aHTML.Class = "image-square-rounded-table"
		aHTML.Name = cat.Name
		aHTML.SecondLine = "Category"
		aHTML.Type = "CATEGORY"
		return aHTML

	case Item:
		i := element.(Item)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		cat := VS.getCategoryByID(c, i.IDCategory, "h_view", "auditToAuditHTML", "cat := VS.getCategoryByID")
		tpe := VS.getCategoryTypeByID(c, cat.Type, "h_view", "auditToAuditHTML", "tpe := VS.getCategoryTypeByID").Name
		aHTML.Icon = "/local/categories/" + tpe + "/" + cat.Name + "/" + i.Photo
		aHTML.Class = "image-square-rounded-table"
		aHTML.Name = i.Name
		aHTML.SecondLine = cat.Name
		aHTML.Type = "ITEM"
		return aHTML

	case Village:
		m := element.(Village)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Class = "image-square-rounded-table"
		aHTML.Icon = "/static/svg/tent.svg"
		aHTML.Name = m.Name
		aHTML.Type = "VILLAGE"
		return aHTML

	case ServiceType:
		m := element.(ServiceType)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		mHTML := VS.serviceTypeToServiceTypeHTML(c, m)
		aHTML.Class = "image-square-rounded-table"
		aHTML.Icon = mHTML.Icon
		aHTML.Name = mHTML.Name
		aHTML.Type = "SERVICE TYPE"
		return aHTML

	case Service:
		s := element.(Service)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		sHTML := VS.serviceToHTML(c, s)
		aHTML.Icon = sHTML.Photo
		aHTML.Class = "image-square-rounded-table"
		aHTML.Name = s.Name
		aHTML.Type = "SERVICE"
		aHTML.SecondLine = VS.getVillageByID(c, s.IDVillage, "h_view", "auditToAuditHTML", "aHTML.SecondLine = VS.getVillageByID").Name
		return aHTML

	case Access:
		u := element.(Access)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/access.svg"
		aHTML.Class = "image-square-rounded-table"
		aHTML.Name = VS.getUserByID(c, u.IDUser, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getUserByID").Name
		aHTML.SecondLine = VS.getServiceByID(c, u.IDService, "h_view", "auditToAuditHTML", "aHTML.SecondLine = VS.getServiceByID", true).Name
		aHTML.Type = "ACCESS"
		return aHTML

	case Report:
		rep := element.(Report)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/report.svg"
		aHTML.Class = "image-square-rounded-table"
		aHTML.Name = VS.getUserByID(c, rep.IDUser, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getUserByID").Name
		aHTML.SecondLine = strconv.FormatBool(rep.IsClose)
		aHTML.Type = "REPORT"
		return aHTML

	case Delivery:
		dl := element.(Delivery)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/deliveries.png"
		ser := VS.getServiceByID(c, dl.IDServiceEmitter, "h_view", "auditToAuditHTML", "ser := VS.getServiceByID", true)
		aHTML.Class = "image-activity-van"
		aHTML.Name = ser.Name
		aHTML.SecondLine = VS.getVillageByID(c, ser.IDVillage, "h_view", "auditToAuditHTML", "aHTML.SecondLine = VS.getVillageByID").Name
		aHTML.Type = "DELIVERY PACK"
		return aHTML

	case Step:
		s := element.(Step)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/video.svg"
		vid, _ := VS.getVideoCourseByID(c, s.IDVideoCourse, "h_view", "auditToAuditHTML", "vid, _ := VS.getVideoCourseByID", true)
		aHTML.Name = vid.Name
		aHTML.SecondLine = strconv.Itoa(s.IndexOrder) + ". " + s.Title
		aHTML.Type = "STEP"
		return aHTML

	case WorkerOrder:
		s := element.(WorkerOrder)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		user := VS.getUserByID(c, s.IDWorker, "h_view", "auditToAuditHTML", "user := VS.getUserByID")
		aHTML.Icon = "/local/users/" + user.Username + "/" + user.Photo
		vid, _ := VS.getVideoCourseByID(c, s.IDProduct, "h_view", "auditToAuditHTML", "vid, _ := VS.getVideoCourseByID", true)
		aHTML.Name = vid.Name
		aHTML.SecondLine = strconv.Itoa(s.Quantity)
		aHTML.Class = "image-rounded-table"
		aHTML.Type = "WORKER ORDER"
		return aHTML

	case Stock:
		s := element.(Stock)
		stk := VS.stockToHTML(c, s)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = stk.Photo
		aHTML.Class = "image-rounded-table"
		aHTML.Name = VS.getUserByID(c, s.IDUser, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getUserByID").Name
		aHTML.SecondLine = VS.getServiceByID(c, s.IDServiceLocation, "h_view", "auditToAuditHTML", "aHTML.SecondLine = VS.getServiceByID", false).Name
		aHTML.Type = "WORKSHOP PURCHASE OR DELIVERY"
		return aHTML

	case Payment:
		s := element.(Payment)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/payment.svg"
		aHTML.Name = VS.getUserByID(c, s.IDWorker, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getUserByID").Name
		aHTML.SecondLine = strconv.FormatFloat(s.Quantity, 'f', -1, 64)
		aHTML.Type = "PAYMENT"
		return aHTML

	case VideoProblem:
		s := element.(VideoProblem)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/Problem.svg"
		aHTML.Name = VS.getUserByID(c, s.IDWorker, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getUserByID").Name
		vid, _ := VS.getVideoCourseByID(c, s.IDVideoCourse, "h_view", "auditToAuditHTML", "vid, _ := VS.getVideoCourseByID", true)
		aHTML.SecondLine = vid.Name
		aHTML.Type = "PRODUCT PROBLEM"
		return aHTML

	case Sale:
		s := element.(Sale)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		stock := VS.getStockByID(c, s.IDStock, "h_view", "auditToAuditHTML", "stock := VS.getStockByID")
		item, b := VS.getItemByID(c, stock.IDVideoCourseOrItem, "h_view", "v", "item, b := VS.getItemByID", false)
		if b {
			cat := VS.getCategoryByID(c, item.IDCategory, "h_view", "auditToAuditHTML", "cat := VS.getCategoryByID")
			tpe := VS.getCategoryTypeByID(c, cat.Type, "h_view", "auditToAuditHTML", "tpe := VS.getCategoryTypeByID")
			aHTML.Icon = "/local/categories/" + tpe.Name + "/" + cat.Name + "/" + item.Photo
			aHTML.Name = VS.getUserByID(c, s.IDWorker, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getUserByID").Name
		} else {
			video, _ := VS.getVideoCourseByID(c, s.IDStock, "h_view", "auditToAuditHTML", "video, _ := VS.getVideoCourseByID", true)
			aHTML.Icon = "/local/video-courses/" + video.Name
			aHTML.Name = video.Name
		}
		aHTML.Class = "image-rounded-table"
		aHTML.SecondLine = strconv.FormatFloat(float64(s.Quantity), 'f', -1, 64)
		aHTML.Type = "SALE"
		return aHTML

	case VideoCheckList:
		lC := element.(VideoCheckList)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		p, _ := VS.getVideoCourseByID(c, lC.IDVideoCourse, "h_view", "auditToAuditHTML", "p, _ := VS.getVideoCourseByID", true)
		aHTML.Icon = "/local/video-courses/" + p.Name + "/checks/" + lC.Photo
		aHTML.Class = "image-square-rounded-table"
		aHTML.Name = p.Name
		aHTML.SecondLine = ""
		aHTML.Type = "PRODUCT CHECK LIST"
		return aHTML

	case CategoryType:
		cT := element.(CategoryType)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/type.svg"
		aHTML.Name = cT.Name
		aHTML.SecondLine = ""
		aHTML.Type = "CATEGORY TYPE"
		return aHTML

	case Conversation:
		cT := element.(Conversation)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		names := ""
		conv := VS.getConversationByID(c, cT.ID, "h_view", "auditToAuditHTML", "conv := VS.getConversationByID")
		for _, u := range conv.Users {
			names += ", " + VS.getUserByID(c, u, "h_view", "auditToAuditHTML", "names += \", \" + VS.getUserByID").Name
		}
		aHTML.Icon = "/static/svg/conversation.svg"
		aHTML.Name = "Conversation"
		aHTML.SecondLine = ""
		aHTML.Type = "CONVERSATION"
		return aHTML

	case Message:
		ms := element.(Message)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/conversation.svg"
		aHTML.Name = VS.getMessageByID(c, ms.ID, "h_view", "auditToAuditHTML", "aHTML.Name = VS.getMessageByID").Text
		aHTML.SecondLine = ""
		aHTML.Type = "MESSAGE"
		return aHTML

	case Synchronization:
		sync := element.(Synchronization)
		aHTML.Link = "/settings/activity/" + string(a.ID.Hex())
		aHTML.Icon = "/static/svg/sync.svg"
		aHTML.Name = "From " + VS.getVillageByID(c, sync.IDVillageEmitter, "h_view", "auditToAuditHTML", "VS.getVillageByID").Name
		aHTML.SecondLine = "To " + VS.getVillageByID(c, sync.IDVillageReceiver, "h_view", "auditToAuditHTML", "VS.getVillageByID").Name
		aHTML.Type = "SYNCHRONIZATION"
		return aHTML
	}
	return aHTML
}

// This method gets the bson.ObjectID Param from URL
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getIDParamFromURL(c *gin.Context, paramName string) (bson.ObjectId, bool) {
	var id bson.ObjectId
	bol := false
	if bson.IsObjectIdHex(c.Param(paramName)) {
		id = bson.ObjectIdHex(c.Param(paramName))
		bol = true
	} else {
		VS.wrongFeedback(c, "Sorry, the parameter for ID for the URL introduced is not valid :(")
	}
	return id, bol
}

// This method gets the string Param from URL
// The second element return true if it was successfully and false in the opposite case
func (VS *Server) getParamFromURL(c *gin.Context, nameParam string) (string, bool) {
	pos := ""
	bol := false
	if c.Param(nameParam) != "" {
		pos = c.Param(nameParam)
		bol = true
	} else {
		VS.wrongFeedback(c, "Sorry, the parameter for "+nameParam+" introduced is not valid :(")
	}
	return pos, bol
}

// Gets an array of users and return the AccessesHTML
func (VS *Server) getAccessFromUsersHTML(c *gin.Context, u []User) []AccessHTML {
	var accessHTML AccessHTML
	var accessesHTML []AccessHTML
	for _, user := range u {
		accessHTML.ID = user.ID
		accessHTML.Username = user.Username
		accessHTML.NameUser = user.Name
		accessHTML.NameVillage = VS.getVillageByID(c, user.IDVillage, "h_view", "getAccessFromUsersHTML", "accessHTML.NameVillage = VS.getVillageByID").Name
		accessHTML.PhotoUser = "/local/users/" + user.Username + "/" + user.Photo
		allAccess := VS.getAllActiveAccessByUser(c, user.ID, "h_view", "getAccessFromUsersHTML", "allAccess := VS.getAllActiveAccessByUser")
		for _, acc := range allAccess {
			serHTML := VS.serviceToHTML(c, VS.getServiceByID(c, acc.IDService, "h_view", "getAccessFromUsersHTML", "serHTML := VS.serviceToHTML(c, VS.getServiceByID", true))
			accessHTML.Services = append(accessHTML.Services, serHTML)
		}
		accessesHTML = append(accessesHTML, accessHTML)
		accessHTML.Services = nil
	}
	return accessesHTML
}

// This methos transform the array of Delivery into human readable information
func (VS *Server) deliveryToHTML(c *gin.Context, dPack []Delivery) []DeliveryHTML {
	var dpsHTML []DeliveryHTML
	var dpHTML DeliveryHTML
	for _, dp := range dPack {
		dpHTML.ID = dp.ID
		dpHTML.Date = strconv.Itoa(time.Unix(dp.Date, 0).Day()) + " / " + time.Unix(dp.Date, 0).Month().String()
		dpHTML.IsSent = dp.IsSent
		dpHTML.IsComplete = dp.IsComplete
		if dpHTML.IsSent && dp.IDManager != "" {
			dpHTML.ManagerName = VS.getUserByID(c, dp.IDManager, "h_view", "deliveryToHTML", "dpHTML.ManagerName = VS.getUserByID").Name
		}
		serEmitter := VS.getServiceByID(c, dp.IDServiceEmitter, "h_view", "deliveryToHTML", "serEmitter := VS.getServiceByID", true)
		dpHTML.ServiceEmitterName = serEmitter.Name
		dpHTML.ServiceEmitterID = serEmitter.ID
		serReceiver := VS.getServiceByID(c, dp.IDServiceReceiver, "h_view", "deliveryToHTML", "serReceiver := VS.getServiceByID", true)
		dpHTML.ServiceReceiverName = serReceiver.Name
		dpHTML.ServiceReceiverID = serReceiver.ID

		dpHTML.VillageName = VS.getVillageByID(c, serReceiver.IDVillage, "h_view", "deliveryToHTML", "dpHTML.VillageName = VS.getVillageByID").Name
		dpsHTML = append(dpsHTML, dpHTML)
	}
	return dpsHTML
}

// Returns an array of Villages and his services to select the service we want to make the delivery
// If the cond is true, removes from the array the service to avoid to make a delivery to itself
func (VS *Server) selectServiceDeliveryHTML(c *gin.Context, villages []Village, idService bson.ObjectId, isService bool) []SelectServiceHTML {
	var selectVillages []SelectServiceHTML
	var selectVillage SelectServiceHTML
	var services []Service
	for _, village := range villages {
		if !isService {
			services = VS.getAllServicesFromSpecificVillage(c, village.ID, "h_view", "selectServiceDeliveryHTML", "services = VS.getAllServicesFromSpecificVillage")
		} else {
			services = VS.getAllServicesFromSpecificVillageExceptItself(c, village.ID, idService, "h_view", "selectServiceDeliveryHTML", "services = VS.getAllServicesFromSpecificVillageExceptItself")
		}
		selectVillage.VillageName = village.Name
		selectVillage.Services = VS.servicesToHTML(c, services)
		selectVillages = append(selectVillages, selectVillage)
		services = nil
	}
	return selectVillages
}

// This method takes an array of Payment and returns an array of PaymentHTML
//  so we can represent the information as human readable in the website.
func (VS *Server) toViewPayments(c *gin.Context, payments []Payment) []PaymentHTML {
	var paymentsHTML []PaymentHTML
	var pay PaymentHTML
	for i := 0; i < len(payments); i++ {
		pay.ID = payments[i].ID
		pay.DateUnix = time.Unix(payments[i].Date, 0).UnixNano()

		//Here adds the Date
		pay.Date = strconv.Itoa(time.Unix(payments[i].Date, 0).Day()) + "/" + time.Unix(payments[i].Date, 0).Month().String()

		//Here adds the Worker
		worker := VS.getUserByID(c, payments[i].IDWorker, "h_view", "toViewPayments", "worker := VS.getUserByID(c, payments[i].WorkerID")
		pay.Worker = worker.Name + " " + worker.Surname

		//Here adds the Money
		pay.Quantity = payments[i].Quantity
		paymentsHTML = append(paymentsHTML, pay)
	}
	return paymentsHTML
}

// This method takes an array of WorkerOrder and returns an array of WorkerOrderHTML
// so we can represent the information as human readable in the website.
func (VS *Server) toViewWorkerOrders(c *gin.Context, orders []WorkerOrder) []WorkerOrderHTML {
	var ordersHTML []WorkerOrderHTML
	var ord WorkerOrderHTML
	for i := 0; i < len(orders); i++ {
		ord.ID = orders[i].ID
		ord.ItemID = orders[i].IDProduct

		//Here adds the Date
		ord.Date = strconv.Itoa(time.Unix(orders[i].Date, 0).Day()) + "/" + time.Unix(orders[i].Date, 0).Month().String()

		//Here adds the Product
		product, _ := VS.getVideoCourseByID(c, orders[i].IDProduct, "h_view", "toViewWorkerOrders", "product, _ := VS.getVideoCourseByID", true)
		ord.ItemName = product.Name
		ord.Photo = "/local/video-courses/" + product.Name + "/" + product.Photo

		//Here adds the Worker
		worker := VS.getUserByID(c, orders[i].IDWorker, "h_view", "toViewWorkerOrders", "worker := VS.getUserByID")
		ord.WorkerPhoto = worker.Photo
		ord.Name = worker.Name + " " + worker.Surname

		//Here adds the Quantity
		ord.Quantity = orders[i].Quantity

		//Here adds the status
		ord.Status = strconv.Itoa((orders[i].AlreadyMade * 100) / orders[i].Quantity)
		ord.Status += "%"
		ordersHTML = append(ordersHTML, ord)
	}
	return ordersHTML
}

// This method takes a Synchronization and returns a of SynchronizationHTML
// so we can represent the information as human readable in the website.
func (VS *Server) syncToHTML(c *gin.Context, s Synchronization) SynchronizationHTML {
	var sync SynchronizationHTML
	sync.ID = s.ID
	y, m, d := time.Unix(s.Date, 0).Date()
	sync.Date = time.Unix(s.Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
	sync.VillageEmitterName = VS.getVillageByID(c, s.IDVillageEmitter, "h_view", "syncToHTML", "sync.Village = VS.getVillageByID").Name
	sync.VillageReceiverName = VS.getVillageByID(c, s.IDVillageReceiver, "h_view", "syncToHTML", "sync.Village = VS.getVillageByID").Name
	sync.Admin = VS.getUserByID(c, s.IDAdmin, "h_view", "syncToHTML", "sync.Admin = VS.getUserByID").Name
	sync.Path = "/local/syncs/exports/" + sync.VillageReceiverName + "/" + "To-" + sync.VillageReceiverName + "-" + sync.ID.Hex() + ".zip"
	sync.IDVillageEmitter = s.IDVillageEmitter
	sync.IDVillageReceiver = s.IDVillageReceiver
	sync.IsDone = s.IsDone
	return sync
}

// This method takes an array of Synchronizations and returns an array of SynchronizationsHTML
// so we can represent properly the information in the website.
func (VS *Server) syncsToHTML(c *gin.Context, s []Synchronization) []SynchronizationHTML {
	var syncHTML []SynchronizationHTML
	for _, n := range s {
		syncHTML = append(syncHTML, VS.syncToHTML(c, n))
	}
	return syncHTML
}

// This method takes a Stock and returns a StockHTML
// so we can represent properly the information in the website.
func (VS *Server) stockToHTML(c *gin.Context, s Stock) StockHTML {
	var stock StockHTML
	stock.ID = s.ID
	y, m, d := time.Unix(s.Date, 0).Date()
	stock.Date = time.Unix(s.Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
	stock.DateUnix = time.Unix(s.Date, 0).UnixNano()

	stock.VideoOrItemID = s.IDVideoCourseOrItem
	stock.IsTrackable = s.IsTrackable
	prod, b := VS.getVideoCourseByID(c, s.IDVideoCourseOrItem, "h_view", "stockToHTML", "prod, b := VS.getVideoCourseByID", true)
	if b {
		stock.Photo = "/local/video-courses/" + prod.Name + "/purchases/" + s.Photo
		stock.VideoOrItemName = prod.Name
		stock.UnitType = Unit
		stock.Price = stock.Quantity * prod.Price
	} else {
		item, _ := VS.getItemByID(c, s.IDVideoCourseOrItem, "h_view", "stockToHTML", "item, _ := VS.getItemByID(c, purchase.VideoCourseOrItemID", true)
		itemHTML := VS.itemToHTML(c, item)
		stock.UnitType = item.UnitType
		stock.Photo = itemHTML.Photo
		stock.VideoOrItemName = itemHTML.Name
		stock.Price = stock.Quantity * item.Price
	}
	user := VS.getUserByID(c, s.IDUser, "h_view", "stockToHTML", "user := VS.getUserByID(c, purchase.UserID")
	stock.UserPhoto = "/local/users/" + user.Username + "/" + user.Photo
	stock.ItemOrVideoCourseID = prod.ID
	service := VS.getServiceByID(c, s.ServiceCreated, "h_view", "stockToHTML", "service := VS.getServiceByID", false)
	var ser Service
	if service == ser {
		stock.ServiceName = "In delivery"
	} else {
		stock.ServiceName = service.Name
	}
	stock.UserID = s.IDUser
	return stock
}

func (VS *Server) stocksToHTML(c *gin.Context, stocks []Stock) []StockHTML {
	var stocksHTML []StockHTML
	for _, s := range stocks {
		stocksHTML = append(stocksHTML, VS.stockToHTML(c, s))
	}
	return stocksHTML
}

func (VS *Server) todosToHTML(c *gin.Context, stocks []ToDo) []ToDoHTML {
	var todosHTML []ToDoHTML
	for _, s := range stocks {
		todosHTML = append(todosHTML, VS.todoToHTML(c, s))
	}
	return todosHTML
}

// This method takes an Todo and returns a TodoHTML
func (VS *Server) todoToHTML(c *gin.Context, todo ToDo) ToDoHTML {
	var todoHTML ToDoHTML
	todoHTML.ID = todo.ID
	todoHTML.IsTrackable = todo.IsTrackable
	todoHTML.Checkboxes = todo.Checkboxes
	todoHTML.Descriptions = todo.Descriptions
	todoHTML.IDuser = todo.IDUser
	todoHTML.IDvillage = todo.IDVillage
	todoHTML.NameUser = VS.getUserByID(c, todo.IDUser, "h_view", "todoToHTML", "todoHTML.NameUser = VS.getUserByID").Name
	todoHTML.NameVillage = VS.getVillageByID(c, todo.IDVillage, "h_view", "todoToHTML", "todoHTML.NameVillage = VS.getVillageByID").Name
	todoHTML.Numbers = todo.Numbers
	for _, p := range todo.Photos {
		user := VS.getUserByID(c, todo.IDUser, "h_view", "todoToHTML", "user := VS.getUserByID")
		newPhoto := FieldPhoto{
			ID:         p.ID,
			IDToDo:     p.IDToDo,
			TitlePhoto: p.TitlePhoto,
			Photo:      "/local/users/" + user.Username + "/to-dos/" + p.Photo,
		}
		todoHTML.Photos = append(todoHTML.Photos, newPhoto)
	}
	todoHTML.TimeChecking = todo.TimeChecking
	todoHTML.TitleToDo = todo.TitleToDo
	todoHTML.DescriptionToDo = todo.DescriptionToDo
	return todoHTML
}

// This method takes an array of TodoChecked and returns an array of TodoCheckedHTML
func (VS *Server) todoCheckedSToHTML(c *gin.Context, checkeds []ToDoChecked) []ToDoChecked {
	var sol []ToDoChecked
	for _, ch := range checkeds {
		sol = append(sol, VS.todoCheckedToHTML(c, ch))
	}
	return sol
}

// This method takes a TodoChecked and returns a TodoCheckedHTML
func (VS *Server) todoCheckedToHTML(c *gin.Context, checked ToDoChecked) ToDoChecked {
	var ch ToDoChecked
	ch.ID = checked.ID
	ch.IDToDo = checked.IDToDo
	ch.Numbers = checked.Numbers
	ch.Descriptions = checked.Descriptions
	ch.Date = checked.Date
	ch.Checkboxes = checked.Checkboxes
	todo := VS.getTodoByID(c, checked.IDToDo, "h_view", "todoCheckedToHTML", "todo := VS.getTodoByID(c, checked")
	for index, p := range checked.Photos {
		user := VS.getUserByID(c, todo.IDUser, "h_view", "todoCheckedToHTML", "user := VS.getUserByID(c, todo")
		newPhoto := FieldPhoto{
			ID:         p.ID,
			IDToDo:     p.IDToDo,
			TitlePhoto: todo.Photos[index].TitlePhoto,
			Photo:      "/local/users/" + user.Username + "/to-dos/checked/" + p.Photo,
		}
		ch.Photos = append(ch.Photos, newPhoto)
	}
	return ch
}

// This method takes an array of Sale and returns an array of SalesHTML
func (VS *Server) salesToHTML(c *gin.Context, wkS []Sale) []SaleHTML {
	var wkSales []SaleHTML
	var wkSale SaleHTML
	for _, sale := range wkS {
		wkSale.DateUnix = time.Unix(sale.Date, 0).UnixNano()
		wkSale.ID = sale.ID
		y, m, d := time.Unix(sale.Date, 0).Date()
		wkSale.Date = time.Unix(sale.Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
		fmt.Print("sale.IDStock")
		fmt.Print(sale.IDStock)
		stock := VS.getStockByID(c, sale.IDStock, "h_view", "salesToHTML", "stock := VS.getStockByID")
		fmt.Print("stock")
		fmt.Print(stock)
		item, b := VS.getItemByID(c, stock.IDVideoCourseOrItem, "h_view", "salesToHTML", "item, b := VS.getItemByID", true)
		fmt.Print("item")
		fmt.Print(item)
		if b {
			fmt.Print("Is an item, not a video, item.IDCategory:")
			fmt.Print(item.IDCategory)
			cat := VS.getCategoryByID(c, item.IDCategory, "h_view", "salesToHTML", "cat := VS.getCategoryByID")
			fmt.Print("Category:")

			fmt.Print(cat)
			tpe := VS.getCategoryTypeByID(c, cat.Type, "h_view", "salesToHTML", "tpe := VS.getCategoryTypeByID(c")
			wkSale.ItemPhoto = "/local/categories/" + tpe.Name + "/" + cat.Name + "/" + item.Photo
			wkSale.ItemName = item.Name
		} else {
			fmt.Print("Acording to the method, is a video:")
			video, _ := VS.getVideoCourseByID(c, sale.IDStock, "h_view", "salesToHTML", "video, _ := VS.getVideoCourseByID", true)
			wkSale.ItemPhoto = "/local/video-courses/" + video.Name + "/" + video.Photo
			wkSale.ItemName = item.Name
		}
		// It is supposed that here will not enter Video Products
		material, _ := VS.getItemByID(c, stock.IDVideoCourseOrItem, "h_view", "salesToHTML", "material, _ := VS.getItemByID", true)
		category := VS.getCategoryByID(c, material.IDCategory, "h_view", "salesToHTML", "category := VS.getCategoryByID")
		wkSale.CategoryName = category.Name
		wkSale.CategoryType = string(category.Type)
		wkSale.Price = sale.Price
		wkSale.ServiceName = VS.getServiceByID(c, sale.IDService, "h_view", "salesToHTML", "wkSale.ServiceName = VS.getServiceByID", true).Name
		user := VS.getUserByID(c, sale.IDWorker, "h_view", "salesToHTML", "user := VS.getUserByID")
		wkSale.UserPhoto = "/local/users/" + "/" + user.Username + "/" + user.Photo

		wkSale.Quantity = sale.Quantity
		wkSale.UserName = VS.getUserByID(c, sale.IDWorker, "h_view", "salesToHTML", "wkSale.Username = VS.getUserByID").Name
		wkSales = append(wkSales, wkSale)
	}
	return wkSales
}

// This method takes an array of VideoCourses and returns an array of VideoCourseHTML
func (VS *Server) toViewVideoCourses(c *gin.Context, products []VideoCourse) []VideoCourseHTML {
	var productsHTML []VideoCourseHTML
	//Create the HTML data and reverting the order of the slice
	for i := len(products) - 1; 0 <= i; i-- {
		productsHTML = append(productsHTML, VS.toViewVideoCourse(c, products[i]))
	}
	return productsHTML
}

// This method takes a VideoCourses and returns a VideoCourseHTML
func (VS *Server) toViewVideoCourse(c *gin.Context, products VideoCourse) VideoCourseHTML {
	var prodHTML VideoCourseHTML
	prodHTML.ID = products.ID
	prodHTML.VideoCourseID = products.ID

	prodHTML.VideoCourseName = products.Name
	prodHTML.Description = products.Description
	prodHTML.Price = products.Price
	prodHTML.Photo = "/local/video-courses/" + products.Name + "/" + products.Photo

	//Get all Steps from DataBase
	steps := VS.getAllStepsFromSpecificVideoCourse(c, products.ID, "h_view", "toViewVideoCourse", "steps := VS.getAllStepsFromSpecificVideoCourse")
	prodHTML.Steps = len(steps)
	prodHTML.StepsID = nil
	for j := len(steps) - 1; 0 <= j; j-- {
		prodHTML.StepsID = append(prodHTML.StepsID, steps[j].ID.Hex())
	}
	return prodHTML
}

// This method takes an array of ServiceOrder and returns an array of ServiceOrdersHTML
func (VS *Server) toViewWorkshopOrder(c *gin.Context, workshopOrders []ServiceOrder) []ServiceOrderHTML {
	var prodOrdersHTML []ServiceOrderHTML
	var prodOHTML ServiceOrderHTML
	for i := 0; i < len(workshopOrders); i++ {
		prodOHTML.ID = workshopOrders[i].ID
		prodOHTML.Assigned = workshopOrders[i].Assigned
		// We add the Product Name, get the Product from DataBase
		temProd, _ := VS.getVideoCourseByID(c, workshopOrders[i].IDVideoCourse, "h_view", "toViewWorkshopOrder", "temProd, _ := VS.getVideoCourseByID", true)
		prodOHTML.ProductName = temProd.Name
		prodOHTML.Photo = temProd.Photo
		// We add the Workshop name
		wk := VS.getServiceByID(c, workshopOrders[i].IDService, "h_view", "toViewWorkshopOrder", "wk := VS.getServiceByID", true)
		prodOHTML.Village = VS.getVillageByID(c, wk.IDVillage, "h_view", "toViewWorkshopOrder", "prodOHTML.Village = VS.getVillageByID").Name
		prodOHTML.WorkshopName = wk.Name
		// We add the Date
		prodOHTML.Deadline = strconv.Itoa(time.Unix(workshopOrders[i].Deadline, 0).Day()) + " / " + time.Unix(workshopOrders[i].Deadline, 0).Month().String()
		// We add the Quantity
		prodOHTML.Quantity = workshopOrders[i].Quantity
		switch workshopOrders[i].WindowPeriod {
		case 3:
			prodOHTML.Window = "3 Days"
			break

		case 7:
			prodOHTML.Window = "1 Week"
			break

		case 10:
			prodOHTML.Window = "10 Days"
			break

		case 14:
			prodOHTML.Window = "2 Weeks"
			break

		case 30:
			prodOHTML.Window = "1 Month"
			break

		default:
			prodOHTML.Window = "Error in input Window Period"
		}

		// We add the Already Made
		prodOHTML.AlreadyMade = workshopOrders[i].AlreadyMade

		if workshopOrders[i].Quantity != 0 {
			prodOHTML.Status = strconv.Itoa((100*workshopOrders[i].AlreadyMade)/workshopOrders[i].Quantity) + " %"
		} else {
			prodOHTML.Status = "You assigned Quantity = 0 "
		}
		prodOrdersHTML = append(prodOrdersHTML, prodOHTML)
	}
	return prodOrdersHTML
}

// This method take an Item and returns an ItemHTML, taking care if is a Video product or not
func (VS *Server) itemToHTML(c *gin.Context, item Item) ItemHTML {
	var iHTML ItemHTML
	iHTML.ID = item.ID
	iHTML.Name = item.Name
	prod, isProd := VS.getVideoCourseByID(c, item.ID, "h_view", "itemToHTML", "prod, isProd := VS.getVideoCourseByID", false)
	if isProd {
		iHTML.Photo = "/local/video-courses/" + prod.Name + "/purchases/" + "/" + prod.Photo
		iHTML.IsTrackable = true
	} else {
		cat := VS.getCategoryByID(c, item.IDCategory, "h_view", "itemToHTML", "cat := VS.getCategoryByID")
		iHTML.CategoryName = cat.Name
		iHTML.IsTrackable = cat.IsTrackable
		tpe := VS.getCategoryTypeByID(c, cat.Type, "h_view", "itemToHTML", "tpe := VS.getCategoryTypeByID")
		iHTML.Photo = "/local/categories/" + tpe.Name + "/" + cat.Name + "/" + item.Photo
	}
	iHTML.UnitType = item.UnitType
	return iHTML
}

// This method take an Item and returns an ItemHTML, taking care if is a Video product or not
func (VS *Server) itemToHTMLWtihQuantity(c *gin.Context, item Item, quantity float64) ItemHTML {
	var iHTML ItemHTML
	iHTML.ID = item.ID
	iHTML.Name = item.Name
	prod, isProd := VS.getVideoCourseByID(c, item.ID, "h_view", "itemToHTML", "prod, isProd := VS.getVideoCourseByID", false)
	if isProd {
		iHTML.Photo = "/local/video-courses/" + prod.Name + "/purchases/" + "/" + prod.Photo
		iHTML.IsTrackable = true
	} else {
		cat := VS.getCategoryByID(c, item.IDCategory, "h_view", "itemToHTML", "cat := VS.getCategoryByID")
		iHTML.CategoryName = cat.Name
		iHTML.IsTrackable = cat.IsTrackable
		tpe := VS.getCategoryTypeByID(c, cat.Type, "h_view", "itemToHTML", "tpe := VS.getCategoryTypeByID")
		iHTML.Photo = "/local/categories/" + tpe.Name + "/" + cat.Name + "/" + item.Photo
	}
	iHTML.UnitType = item.UnitType
	iHTML.Quantity = quantity
	return iHTML
}

// This method take an Item and returns an ItemHTML, taking care if is a Video product or not. It is used for the real stock in the Service
func (VS *Server) itemStockToHTML(c *gin.Context, item Item, serviceID bson.ObjectId) ItemHTML {
	var iHTML ItemHTML
	iHTML.ID = item.ID
	iHTML.Name = item.Name
	prod, isProd := VS.getVideoCourseByID(c, item.ID, "h_view", "itemStockToHTML", "prod, isProd := VS.getVideoCourseByID", false)
	if isProd {
		iHTML.Photo = "/local/video-courses/" + prod.Name + "/purchases/" + "/" + prod.Photo
	} else {
		cat := VS.getCategoryByID(c, item.IDCategory, "h_view", "itemStockToHTML", "cat := VS.getCategoryByID")
		iHTML.CategoryName = cat.Name
		tpe := VS.getCategoryTypeByID(c, cat.Type, "h_view", "itemStockToHTML", "tpe := VS.getCategoryTypeByID")
		iHTML.Photo = "/local/categories/" + tpe.Name + "/" + cat.Name + "/" + item.Photo
	}
	// VS.getAllPositiveStockByServiceID
	stockFromServiceHTML := VS.stocksToHTML(c, VS.getAllPositiveStockByServiceIDandItemOrVideoID(c, serviceID, item.ID, "h_view", "itemStockToHTML", "stockFromServiceHTML := VS.stocksToHTML("))
	if len(stockFromServiceHTML) != 0 {
		iHTML.Stock = append(iHTML.Stock, stockFromServiceHTML...)
	}
	return iHTML
}

// This method take an VideoCourse and returns an ItemHTML. It is used for the real stock in the Service
func (VS *Server) videoCourseStockToHTML(c *gin.Context, video VideoCourse, serviceID bson.ObjectId) ItemHTML {
	videoHTML := VS.toViewVideoCourse(c, video)
	var iHTML ItemHTML
	iHTML.ID = video.ID
	iHTML.Name = video.Name
	// iHTML.CategoryName = "Products"
	iHTML.Photo = videoHTML.Photo
	stockFromServiceHTML := VS.stocksToHTML(c, VS.getAllPositiveStockByServiceIDandItemOrVideoID(c, serviceID, video.ID, "h_view", "videoCourseStockToHTML", "stockFromServiceHTML := VS.stocksToHTML"))
	if len(stockFromServiceHTML) != 0 {
		stockFromServiceHTML[0].Quantity = float64(len(stockFromServiceHTML))
		iHTML.Stock = append(iHTML.Stock, stockFromServiceHTML[0])
	}
	return iHTML
}

// This method take an array of Item returns an array of itemsToHTML
func (VS *Server) itemsToHTML(c *gin.Context, items []Item) []ItemHTML {
	var itemsHTML []ItemHTML
	for _, item := range items {
		itemsHTML = append(itemsHTML, VS.itemToHTML(c, item))
	}
	return itemsHTML
}

// This converts the array of the Worker Orders into an array of WorkerOrdersHTML
func (VS *Server) toViewWorkerOrdersWorkersView(c *gin.Context, sliceOrders []WorkerOrder) []WorkerOrderHTML {
	var sliceOrdersHTML []WorkerOrderHTML
	var prov WorkerOrderHTML
	for i := 0; i < len(sliceOrders); i++ {
		prov.ID = sliceOrders[i].ID

		//This adds the photo
		product, _ := VS.getVideoCourseByID(c, sliceOrders[i].IDProduct, "h_view", "toViewWorkerOrdersWorkersView", "product, _ := VS.getVideoCourseByID(c, sliceOrders", true)
		//err = v.DataBase.Find(&product, bolthold.Where("ID").Eq(sliceOrders[i].ProductID))
		prov.Photo = "/local/video-courses/" + product.Name + "/" + product.Photo
		prov.ItemName = product.Name

		//This adds the date
		prov.DateUnix = time.Unix(sliceOrders[i].Date, 0).UnixNano()
		prov.Date = strconv.Itoa(time.Unix(sliceOrders[i].Date, 0).Day()) + "/" + strconv.Itoa(int(time.Unix(sliceOrders[i].Date, 0).Month()))

		//This adds the Quantity and AlreadyMade
		prov.Quantity = sliceOrders[i].Quantity
		prov.AlreadyMade = sliceOrders[i].AlreadyMade

		//This adds the image to show
		percent := (prov.AlreadyMade * 100) / prov.Quantity
		if percent == 0 {
			prov.Status = "/static/svg/status0.svg"
		} else if percent < 25 {
			prov.Status = "/static/svg/status25.svg"
		} else if percent < 50 {
			prov.Status = "/static/svg/status40.svg"
		} else if percent == 50 {
			prov.Status = "/static/svg/status50.svg"
		} else if percent < 75 {
			prov.Status = "/static/svg/status75.svg"
		} else if percent > 86 && percent < 90 {
			prov.Status = "/static/svg/status86.svg"
		} else if prov.AlreadyMade == prov.Quantity {
			prov.Status = "/static/svg/status-complete.svg"
		}
		sliceOrdersHTML = append(sliceOrdersHTML, prov)
	}
	return sliceOrdersHTML
}

// Create the human readable information for the table for Messages
func (VS *Server) messagesToHTML(c *gin.Context, idUser bson.ObjectId, msgs []Message) []MessageHTML {
	var messages []MessageHTML
	var message MessageHTML
	for _, m := range msgs {
		message.ID = m.ID
		message.IDConversation = m.IDConversation
		if m.Audio != "" {
			message.Audio = "/local/conversations/" + m.IDConversation.Hex() + "/" + m.Audio
		}
		if m.Photo != "" {
			message.Photo = "/local/conversations/" + m.IDConversation.Hex() + "/" + m.Photo
		}
		if m.IDUser == idUser {
			message.Class = "emitter"
		} else {
			message.Class = "receiver"
		}

		timeStamp := time.Unix(m.Date, 0)
		hr, min, sec := timeStamp.Clock()
		hour, minute, second := "", "", ""
		if hr < 10 {
			hour = "0" + strconv.Itoa(hr)
		} else {
			hour = strconv.Itoa(hr)
		}
		if min < 10 {
			minute = "0" + strconv.Itoa(min)
		} else {
			minute = strconv.Itoa(min)
		}
		if sec < 10 {
			second = "0" + strconv.Itoa(sec)
		} else {
			second = strconv.Itoa(sec)
		}
		message.Time = hour + ":" + minute + ":" + second
		message.Text = m.Text
		if m.IDUser != "" {
			message.User = VS.getUserByID(c, m.IDUser, "h_view", "messagesToHTML", "message.User = VS.getUserByID").Name
		}
		y, mo, d := time.Unix(m.Date, 0).Date()
		message.Date = time.Unix(m.Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + mo.String() + " of " + strconv.Itoa(y)
		messages = append(messages, message)
	}
	return messages
}

// Create the human readable information for the table for Conversations
func (VS *Server) conversationToHTML(c *gin.Context, idUser bson.ObjectId, u []Conversation) []ConversationHTML {
	var conversations []ConversationHTML
	var conv ConversationHTML
	var users []User
	for _, m := range u {
		users = nil
		conv.ID = m.ID
		m.Users = removeDuplicates(m.Users)
		for _, user := range m.Users {
			users = append(users, VS.getUserByID(c, user, "h_view", "conversationToHTML", "users = append(users, VS.getUserByID"))
		}
		conv.Users = VS.usersToHTML(c, users)
		msages := VS.getMessagesByConversation(c, m.ID, "h_view", "conversationToHTML", "msages := VS.getMessagesByConversation")
		if len(msages) != 0 {
			conv.LastMessage = msages[len(msages)-1].Text
		}
		conversations = append(conversations, conv)
	}
	return conversations
}

// Create the human readable information for one Report
func (VS *Server) reportToHTML(c *gin.Context, user User, r Report) ReportHTML {
	var rep ReportHTML
	rep.ID = r.ID
	rep.User = VS.getUserByID(c, user.ID, "h_view", "reportToHTML", "rep.User = VS.getUserByID").Name
	rep.Text = r.Text
	if r.Photo != "" {
		rep.Photo = "/local/users/" + user.Username + "/reports/" + r.Photo
	}
	rep.Audio = "bug"
	y, m, d := time.Unix(r.Date, 0).Date()
	rep.Date = time.Unix(r.Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
	rep.Type = strings.Title(string(r.Type))
	rep.IsClose = r.IsClose
	rep.Village = VS.getVillageByID(c, r.IDVillage, "h_view", "reportToHTML", "rep.Village = VS.getVillageByID").Name
	return rep
}

// Create the human readable information for an array of Reports
func (VS *Server) reportsToHTML(c *gin.Context, User User, reports []Report) []ReportHTML {
	var reportsHTML []ReportHTML
	var rep ReportHTML
	for _, r := range reports {
		rep = VS.reportToHTML(c, User, r)
		reportsHTML = append(reportsHTML, rep)
	}
	return reportsHTML
}

// This method takes an array of Category and returns an array of CategoryHTML
// so we can represent the human readable information in the website as human readable
func (VS *Server) toViewCategories(c *gin.Context, categories []Category) []CategoryHTML {
	var categoryHTML []CategoryHTML
	var cat CategoryHTML
	for i := 0; i < len(categories); i++ {

		//Here adds the common fields
		cat.ID = categories[i].ID
		cat.Name = categories[i].Name
		cat.TimeChecking = string(categories[i].TimeChecking)
		cat.Type = VS.getCategoryTypeByID(c, categories[i].Type, "h_view", "toViewCategories", "cat.Type = VS.getCategoryTypeByID(").Name
		cat.Photo = "/local/categories/" + cat.Type + "/" + categories[i].Name + "/" + categories[i].Photo
		cat.TypeItem = string(categories[i].TypeOfItem)

		serviceTypes := VS.getAllServiceTypes(c, "h_view", "toViewCategories", "serviceTypes := VS.getAllServiceTypes")
		for _, s := range serviceTypes {
			if ContainsID(categories[i].ServicesAccess, s.ID) {
				cat.Services = append(cat.Services, "/local/services/"+s.Icon)
			}
		}

		// We add the object
		categoryHTML = append(categoryHTML, cat)
		// We empty the array for the next loop
		cat.Services = nil

	}
	return categoryHTML
}

// This method takes an array of Assignments and returns an array of AssignmentsHTML
func (VS *Server) toViewAssignments(c *gin.Context, assignments []Assignment) []AssignmentHTML {
	var assignmentsHTML []AssignmentHTML
	var a AssignmentHTML
	for _, n := range assignments {
		a = VS.toViewAssignment(c, n)
		assignmentsHTML = append(assignmentsHTML, a)
	}
	return assignmentsHTML
}

// This method takes an Assignments and returns an AssignmentsHTML
func (VS *Server) toViewAssignment(c *gin.Context, assignments Assignment) AssignmentHTML {
	var a AssignmentHTML
	a.ID = assignments.ID
	item := VS.getStockByID(c, assignments.IDStock, "h_view", "toViewAssignment", "a.Item = VS.getItemByID")
	itemHTML := VS.stockToHTML(c, item)
	a.Item = itemHTML.VideoOrItemName
	a.PhotoItem = itemHTML.Photo
	a.IsBack = assignments.IsBack
	a.Date = strconv.Itoa(time.Unix(assignments.Date, 0).Day()) + " / " + time.Unix(assignments.Date, 0).Month().String()
	a.Manager = VS.getUserByID(c, assignments.IDManager, "h_view", "toViewAssignment", "a.Manager = VS.getUserByID(c").Name
	a.Service = VS.getServiceByID(c, assignments.IDService, "h_view", "toViewAssignment", "a.Service = VS.getServiceByID", true).Name
	worker := VS.getUserByID(c, assignments.IDWorker, "h_view", "toViewAssignment", "aHTML.Name =  v.getVillageByID")
	a.WorkerPhoto = "/local/users/" + worker.Username + "/" + worker.Photo
	a.Worker = worker.Name
	return a
}

// Convert the PDF into PDFHTML
func fromPDFtoPDFHTML(pdf PDF) PDFHTML {
	var pdfHTML PDFHTML
	y, m, d := time.Unix(pdf.Date, 0).Date()
	pdfHTML.Date = time.Unix(pdf.Date, 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)
	pdfHTML.ID = pdf.ID
	pdfHTML.Pages = pdf.Pages
	pdfHTML.TypeOfItem = string(pdf.TypeOfItem)
	pdfHTML.TypeOfPDF = string(pdf.TypeOfPDF)
	return pdfHTML
}

// Convert the []PDF into []PDFHTML
func fromPDFstoPDFHTMLs(pdf []PDF) []PDFHTML {
	var pdfHTML []PDFHTML
	for _, p := range pdf {
		pdfHTML = append(pdfHTML, fromPDFtoPDFHTML(p))
	}
	return pdfHTML
}
