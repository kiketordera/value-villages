package village

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// This handler shows the Conversations of the user
func (VS *Server) showConversations(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_comunication", "showChooseUsersForConversation", "user := VS.getUserFomCookie")
	conv := VS.getConversationsByUser(c, user.ID, "h_comunication", "showChooseUsersForConversation", "conv := VS.getConversationsByUser")
	convHTML := VS.conversationToHTML(c, user.ID, conv)
	JSON, err := json.Marshal(convHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"JSON": string(JSON),
		"conv": conv,
	}, "conversations.html")
}

// This handler shows the Users for choose for starting a conversation
func (VS *Server) showChooseUsersForConversation(c *gin.Context) {
	centralUser, _ := VS.getUserByUsername(c, string(Admin), "h_comunication", "showChooseUsersForConversation", "centralUser := VS.getUserByUsername")
	managerUsers := VS.getAllManagersFromSystem(c, "h_comunication", "showChooseUsersForConversation", "managerUsers := VS.getAllUsersFromLocalVillage")
	workersUsers := VS.getAllWorkersDependingRole(c, "h_comunication", "showChooseUsersForConversation", "workersUsers := VS.getAllWorkersFromTheSystem")
	managers := VS.usersToHTML(c, managerUsers)
	workers := VS.usersToHTML(c, workersUsers)
	render(c, gin.H{
		"centralUser": centralUser,
		"managers":    managers,
		"workers":     workers,
	}, "choose-user.html")
}

// This handler creates a new conversation
func (VS *Server) newConversation(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_comunication", "newConversation", "user := VS.getUserFomCookie")
	users, b := VS.getArrayIDFromHTML(c, "userID", "Conversation")
	if !b {
		return
	}
	users = append(users, user.ID)
	conver := Conversation{
		ID:    bson.NewObjectId(),
		Users: users,
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, conver, conver.ID, "h_comunication", "newConversation", "VS.addConversationToDatabase")
	c.Redirect(http.StatusFound, "/comunication/chat/"+conver.ID.Hex())
}

// This handler shows the Chat of a conversation
func (VS *Server) showChat(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_comunication", "showChat", "user := VS.getUserFomCookie")
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	conv := VS.getConversationByID(c, id, "h_comunication", "showChat", "conv := VS.getConversationByID")
	conv.Users = removeDuplicates(conv.Users)
	var users []User
	for _, uID := range conv.Users {
		u := VS.getUserByID(c, uID, "h_comunication", "showChat", "user := VS.getUserByID")
		users = append(users, u)
	}
	uHTML := VS.usersToHTML(c, users)
	messages := VS.getMessagesByConversation(c, id, "h_comunication", "showChat", "messages := VS.getMessagesByConversation")
	// If we access before make any change to don't go out of range
	var messagesPerDay []MessageDayHTML
	if len(messages) > 0 {
		messagerOrder := VS.orderMessagesPerDay(messages)
		messagesPerDay = VS.messagePerDayToHTML(c, user.ID, messagerOrder)
	}
	render(c, gin.H{
		"users":    uHTML,
		"messages": messagesPerDay,
	}, "chat.html")
}

// This handler adds the message to the conversation
func (VS *Server) addMessage(c *gin.Context) {
	idConver, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	user := VS.getUserFomCookie(c, "h_comunication", "addMessage", "user := VS.getUserFomCookie")
	// We leave open the possibility that there is no audio or photo/video
	photo, audio := "", ""
	p, _ := c.FormFile("photo")
	a, _ := c.FormFile("audio")
	id := bson.NewObjectId()
	if p != nil {
		photo, b = VS.getFileFromHTML(c, "photo", id.Hex(), basePath+"/local-resources/conversations/"+idConver.Hex()+"/")
		if !b {
			return
		}
	}
	if a != nil {
		audio, b = VS.getFileFromHTML(c, "audio", id.Hex(), basePath+"/local-resources/conversations/"+idConver.Hex()+"/")
		if !b {
			return
		}
	}
	text, _ := VS.getStringFromHTML(c, "message", "Message")
	msg := Message{
		ID:             id,
		IDUser:         user.ID,
		IDConversation: idConver,
		Text:           text,
		Photo:          photo,
		Audio:          audio,
		Date:           getTime(),
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, msg, msg.ID, "h_comunication", "addMessage", "VS.addElementToDatabaseWithRegisterFromCentral")
	c.Redirect(http.StatusFound, "/comunication/chat/"+idConver.Hex())
}

// This handler shows the Reports of the user
func (VS *Server) showReports(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_comunication", "showReports", "user := VS.getUserFomCookie")
	var reports []Report
	if user.Role == Admin {
		reports = VS.getReportsFromSystem(c, "h_comunication", "showReports", "reports = VS.getReportsFromSystem")
	} else {
		reports = VS.getReportsByUser(c, user.ID, "h_comunication", "showReports", "reports = VS.getReportsFromSystem")
	}
	reportsHTML := VS.reportsToHTML(c, user, reports)
	JSON, err := json.Marshal(reportsHTML)
	VS.checkOperation(c, "h_date", "showUsers", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"reports": reportsHTML,
		"JSON":    string(JSON),
	}, "reports.html")
}

// Shows the page to introduce the information of the new Report
func (VS *Server) showNewReport(c *gin.Context) {
	render(c, gin.H{}, "new-report.html")
}

// Creates the Report
func (VS *Server) newReport(c *gin.Context) {
	id := bson.NewObjectId()
	user := VS.getUserFomCookie(c, "h_comunication", "newReport", "user := VS.getUserFomCookie")
	var IDVillage bson.ObjectId
	if user.Role != Worker {
		IDVillage = VS.getLocalVillage(c, "h_comunication", "newReport", "IDVillage = VS.getOurVillage").ID
	} else {
		IDVillage = user.IDVillage
	}
	// We leave open the possibility that there is no audio or photo/video
	photo, audio := "", ""
	p, _ := c.FormFile("photo")
	a, _ := c.FormFile("audio")
	var b bool
	if p != nil {
		photo, b = VS.getPhotoFromHTML(c, "photo", id.Hex(), basePath+"/local-resources/users/"+user.Username+"/reports/")
		if !b {
			return
		}
	}
	if a != nil {
		audio, b = VS.getFileFromHTML(c, "audio", id.Hex(), basePath+"/local-resources/users/"+user.Username+"/reports/")
		if !b {
			return
		}
	}
	text, _ := VS.getStringFromHTML(c, "text", "Report")
	typeReport, b := VS.getStringFromHTML(c, "type", "Report")
	if !b {
		return
	}

	report := Report{
		ID:        id,
		IDUser:    user.ID,
		Text:      text,
		Photo:     photo,
		Audio:     audio,
		Date:      getTime(),
		Type:      TypeOfReport(typeReport),
		IDVillage: IDVillage,
		IsClose:   false,
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, report, report.ID, "h_comunication", "newReport", "VS.addElementToDatabaseWithRegisterFromCentral")
	VS.goodFeedback(c, "comunication/reports")
}

// Shows the report
func (VS *Server) seeReport(c *gin.Context) {
	user := VS.getUserFomCookie(c, "h_comunication", "seeReport", "user := VS.getUserFomCookie")
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	report := VS.getReportByID(c, id, "h_comunication", "seeReport", "report := VS.getReportsByID")
	rHTML := VS.reportToHTML(c, user, report)

	messages := VS.getMessagesByConversation(c, id, "h_comunication", "seeReport", "messages := VS.getMessagesByConversation")
	// If we access before make any change to don't go out of range
	var messagesPerDay []MessageDayHTML
	if len(messages) > 0 {
		messagerOrder := VS.orderMessagesPerDay(messages)
		messagesPerDay = VS.messagePerDayToHTML(c, user.ID, messagerOrder)
	}
	render(c, gin.H{
		"report":   rHTML,
		"messages": messagesPerDay,
	}, "see-report.html")
}

// This handler adds the message to the Report
func (VS *Server) addMessageReport(c *gin.Context) {
	idReport, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	rep := VS.getReportByID(c, idReport, "h_comunication", "addMessageReport", "rep := VS.getReportByID")
	if rep.IsClose {
		return
	}
	text, _ := VS.getStringFromHTML(c, "message", "Message")

	// We leave open the possibility that there is no audio or video
	photo, audio := "", ""
	p, _ := c.FormFile("photo")
	a, _ := c.FormFile("audio")
	id := bson.NewObjectId()
	if p != nil {
		photo, b = VS.getFileFromHTML(c, "photo", id.Hex(), basePath+"/local-resources/conversations/"+id.Hex()+"/")
		if !b {
			return
		}
	}
	if a != nil {
		audio, b = VS.getFileFromHTML(c, "audio", id.Hex(), basePath+"/local-resources/conversations/"+id.Hex()+"/")
		if !b {
			return
		}
	}
	msg := Message{
		ID:             id,
		IDConversation: idReport,
		Text:           text,
		Photo:          photo,
		Audio:          audio,
		Date:           getTime(),
		IDUser:         VS.getUserFomCookie(c, "h_comunication", "addMessageReport", "IDUser: VS.getUserFomCookie(").ID,
	}
	close := VS.getCheckBoxFromHTML(c, "closereport")
	if close {
		rep.IsClose = true
		VS.addElementToDatabaseWithRegisterFromCentral(c, rep, rep.ID, "h_comunication", "addMessageReport", "VS.addElementToDatabaseWithRegisterFromCentral")
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, msg, msg.ID, "h_comunication", "addMessageReport", "VS.addElementToDatabaseWithRegisterFromCentral")
	c.Redirect(http.StatusFound, "/comunication/see-report/"+idReport.Hex())
}
