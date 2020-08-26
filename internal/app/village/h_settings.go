package village

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// Shows the activity of the user
func (VS *Server) activity(c *gin.Context) {
	u := VS.getUserFomCookie(c, "h_settings", "activity", "u := VS.getUserFomCookie")
	IDServiceVisualization := Variables["audit"]

	var audits []Audit
	if u.Role == Admin {
		var a bson.ObjectId
		if a == IDServiceVisualization {
			fmt.Print("All audios")
			audits = VS.getAllAudits(c, "h_settings", "activity", "audits = VS.getAuditsByVillage")
		} else {
			fmt.Print("All audios from service")
			fmt.Print("Variables[audit]")
			fmt.Println(IDServiceVisualization)

			audits = VS.getAuditsByService(c, IDServiceVisualization, "h_settings", "activity", "audits = VS.getAuditsByVillage")
			fmt.Print("All audits by service are: ")
			fmt.Println(len(audits))

		}
	} else if u.Role == Manager {
		vill := VS.getLocalVillage(c, "h_settings", "activity", "vill := VS.getLocalVillage")
		audits = VS.getAuditsByVillage(c, vill.ID, "h_settings", "activity", "audits = VS.getAuditsByVillage")
	} else {
		audits = VS.getAuditsByUser(c, u.ID, "h_settings", "activity", "audits = VS.getAuditsByUser")
	}

	var audHTML []AuditDayHTML
	// If we access before make any change to don't go out of range
	if len(audits) > 0 {
		// Second we order them in an array of different days and convert it to HTML
		a := VS.orderAuditPerDay(audits)
		audHTML = VS.auditPerDayToHTML(c, a)
	}
	servicesHTML := VS.getAllServicesPicker(c)
	render(c, gin.H{
		"services":        servicesHTML,
		"audits":          audHTML,
		"role":            u.Role,
		"serviceVisualID": IDServiceVisualization,
	}, "activity.html")
}

// Sets up the workshop for represent his information
func (VS *Server) setUpServiceForAudits(c *gin.Context) {
	serviceID, b := VS.getIDFromHTMLWithoutChecking(c, "serviceID")
	if !b {
		delete(Variables, "audit")
	}
	// fromDate := VS.getDateFromHTML(c, "fromdate", "Activity")
	// toDate := VS.getDateFromHTML(c, "todate", "Activity")
	// isFrom := VS.getCheckBoxFromHTML(c, "isfrom", "Activity")
	// isTo := VS.getCheckBoxFromHTML(c, "isto", "Activity")
	// fmt.Print("fromDate: ")
	// fmt.Println(fromDate)
	// fmt.Print("toDate: ")
	// fmt.Println(toDate)
	// fmt.Print("isFrom: ")
	// fmt.Println(isFrom)
	// fmt.Print("isTo: ")
	// fmt.Println(isTo)

	// if !isFrom {
	// 	FromDate = fromDate
	// }
	// if !isTo {
	// 	ToDate = toDate
	// }
	Variables["audit"] = serviceID
	c.Redirect(http.StatusFound, "/settings/activity")
}

// Shows the log-in page for the normal user
func (VS *Server) activityChanges(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	audit := VS.getAuditByID(c, id, "h_settings", "activityChanges", "audit := VS.getAuditByID")
	// if audit.Description == Modified {
	// 	audit.InformationObject = nil
	// }
	changes := VS.getChangesByAudit(c, id, "h_settings", "activityChanges", "changes := VS.getChangesByAudit")

	render(c, gin.H{
		"changes":           changes,
		"informationobject": audit.InformationObject,
	}, "activity-changes.html")
}

// Shows the screen to introduce the new password
func (VS *Server) showSetPassword(c *gin.Context) {
	render(c, gin.H{}, "reset-password.html")
}

// Sets the new password to the User
func (VS *Server) setPassword(c *gin.Context) {
	// Password
	password, b := VS.getStringFromHTML(c, "password", "User")
	if !b {
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	VS.checkOperation(c, "h_settings", "setPassword", "passwordHash, err := bcrypt.GenerateFromPassword", err)

	// Set the password to the User
	u := VS.getUserFomCookie(c, "h_settings", "setPassword", "u := VS.getUserFomCookie")
	u.Password = string(passwordHash)
	b = VS.addUserToDatabase(c, u, "h_settings", "setPassword", "VS.addUserToDatabase(c, u")
	if !b {
		VS.wrongFeedback(c, "The User could not be updated")
	}

	// Audit
	audit := VS.getAuditFromCentral(c, u, u.ID, Modified)
	// Changes
	VS.checkChange(c, audit.ID, "Password", "Private  : )", "Private ; )")
	// VS.addElementToDatabaseWithRegisterFromCentral(c, audit, audit.ID, "h_settings", "setPassword", "VS.checkChange(c, audit.ID")
	VS.registerAuditAndEvent(c, audit, u, u.ID)
	c.Redirect(http.StatusFound, "/dashboard")
}

// Sets the new password of the User to ""
func (VS *Server) resetPassword(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	// Reset the password to the User
	u := VS.getUserByID(c, id, "h_settings", "setPassword", "u := VS.getUserFomCookie")
	u.Password = ""
	b = VS.addUserToDatabase(c, u, "h_settings", "setPassword", "VS.addUserToDatabase(c, u")
	if !b {
		VS.wrongFeedback(c, "The User could not be updated")
	}
	// Audit
	audit := VS.getAuditFromCentral(c, u, u.ID, Modified)
	// Changes
	VS.checkChange(c, audit.ID, "Change ", "Made by ", "Manager")
	VS.checkChange(c, audit.ID, "Password", "Private  : )", "Private ; )")
	VS.addElementToDatabaseWithRegisterFromCentral(c, audit, audit.ID, "h_settings", "setPassword", "VS.checkChange(c, audit.ID")
	VS.goodFeedback(c, "")
}

// Shows the screen to introduce the new password for the Worker
func (VS *Server) showSetPasswordWorker(c *gin.Context) {
	render(c, gin.H{}, "reset-password-worker.html")
}

// Shows the general settings the user can make to his account
func (VS *Server) showGeneralSettings(c *gin.Context) {
	render(c, gin.H{}, "general.html")
}

// This handler shows all the synchronizations of the system between the Central DataBase and the villages
func (VS *Server) showSynchronizations(opt SyncOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var syncs []Synchronization
		if opt == exp {
			syncs = VS.getAllSynchronizationsByVillageEmitter(c, VS.getLocalVillage(c, "h_settings", "showSynchronizations", "getAllSynchronizationsByVillageEmitter(c, VS.getLocalVillage").ID, "h_settings", "showSynchronizations", "syncs = VS.getAllSynchronizationsByVillageEmitter")
		} else if opt == imp {
			syncs = VS.getAllSynchronizationsByVillageReceiver(c, VS.getLocalVillage(c, "h_settings", "showSynchronizations", "getAllSynchronizationsByVillageReceiver(c, VS.getLocalVillage").ID, "h_settings", "showSynchronizations", "syncs = VS.getAllSynchronizationsByVillageEmitter")
		}
		sHTML := VS.syncsToHTML(c, syncs)
		JSON, err := json.Marshal(sHTML)
		VS.checkOperation(c, "h_settings", "showSynchronizations", "sHTML := VS.syncsToHTML(c, syncs)", err)
		render(c, gin.H{
			"JSON":  string(JSON),
			"opt":   opt,
			"syncs": syncs,
		}, "synchronizations.html")
	}
}

// This handler shows all the different villages to make a new synchronization
func (VS *Server) chooseOptionToSync(c *gin.Context) {
	render(c, gin.H{}, "choose-sync-options.html")
}

// This handler shows all the different villages to make a new synchronization
func (VS *Server) chooseVillageToSync(c *gin.Context) {
	var villages []Village
	if VS.checkWeAreCentral(c, "h_settings", "chooseVillageToSync", "if VS.checkWeAreCentral") {
		villages = VS.getAllVillagesFromTheSystemButLocal(c, "h_settings", "chooseVillageToSync", "JSON, err := json.Marshal(uHTML)")
	} else {
		villages = append(villages, VS.getCentralVillage(c, "h_settings", "chooseVillageToSync", "villages = append(villages, VS.getCentralVillage"))
	}
	render(c, gin.H{
		"villages": villages,
	}, "choose-village-sync.html")
}

// This handler exports a new DataBase to be Sync with the village chosen
func (VS *Server) newSyncExport(c *gin.Context) {
	VS.loadingContent(c)
	villageID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	village := VS.getVillageByID(c, villageID, "h_settings", "newSync", "village := VS.getVillageByID(c, villageID")
	sync := Synchronization{
		ID:                bson.NewObjectId(),
		Date:              getTime(),
		IDVillageReceiver: villageID,
		IDVillageEmitter:  VS.getLocalVillage(c, "h_settings", "newSync", "IDVillageEmitter: VS.getLocalVillage").ID,
		IDAdmin:           VS.getUserFomCookie(c, "h_settings", "newSync", "AdminID:   VS.getUserFomCookie").ID,
		IsDone:            false,
	}
	if len(VS.getAllAuditsFromLastSynchronization(c, village.ID, sync.ID, "h_settings", "newSync", "audits := VS.getAllAuditsFromLastSynchronization")) == 0 {
		VS.wrongFeedback(c, "There is no new information to sync")
		return
	}

	VS.addElementToDatabaseWithRegister(c, sync, sync.ID, VS.getLocalVillage(c, "h_settings", "newSync", "VS.getLocalVillage").ID, "h_settings", "newSync", "addElementToDatabaseWithRegister)")
	_, _, b = VS.createDataBaseForSync(c, village, sync.ID)
	if !b {
		return
	}
	VS.goodFeedback(c, "settings/exports")
}

// This handler shows all the different villages to make a new synchronization
func (VS *Server) newImportDataBase(c *gin.Context) {
	render(c, gin.H{}, "import-database.html")
}

// This handler shows all the different villages to make a new synchronization
func (VS *Server) importDataBase(c *gin.Context) {
	zip, ID, name, b := VS.getFileFromHTMLDB(c, "zip", "zip", basePath+"/local-resources/syncs/imports/")
	if !b {
		return
	}

	// Watch out not to sync two times the same ZIP file
	_, b = VS.getSynchronizationByID(c, ID, "h_settings", "newSync", "VS.getLocalVillage", false)
	if b {
		VS.wrongFeedback(c, "This ZIP file has been already Synchronized")
		return
	}
	// Watch out not to sync a ZIP file with another village Receiver
	if !VS.existLocalVillage(c, "h_settings", "newSync", "if VS.existLocalVillage") {
		// Check for the village receiver in the name of the file
		if localVillage != name {
			VS.wrongFeedback(c, "This Synchronization is not for this village on the name of the File")
			return
		}
	}

	VS.loadingContent(c)
	fmt.Print("File name is: " + zip)
	unzipDirectory(basePath+"/local-resources/syncs/imports/"+zip, basePath+"/local-resources/syncs/imports/")
	VS.importElementsFromDatabases(c, ID)
	VS.goodFeedback(c, "settings/imports")

	fmt.Print("Synchronization finished!")
}
