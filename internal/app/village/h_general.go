package village

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// Shows the PDF in the browser
func (VS *Server) showPDFDocument(c *gin.Context) {

	pathToPDF := basePath + "/web/static/documents/"
	fileName, b := VS.getParamFromURL(c, "name")
	if !b {
		return
	}
	targetPath := filepath.Join(pathToPDF, fileName)
	// This check is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
	if !strings.HasPrefix(filepath.Clean(targetPath), pathToPDF) {
		c.String(403, "Look like you attacking me")
		return
	}
	// Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	// Add this two headers to download instead of see in browser
	// c.Header("Content-Disposition", "attachment; filename="+fileName )
	// c.Header("Content-Type", "application/octet-stream")
	c.File(targetPath)
}

// Redirects the User to the dashboard if is logged in or to the log-in screen if is not
func (VS *Server) distributor(c *gin.Context) {
	r, e := c.Get("role")
	if e {
		// If the User is logged-in
		role := r.(string)
		if role != "" {
			fmt.Println("Redirecting to /dashboard from distributor")
			c.Redirect(http.StatusFound, "/dashboard")
		}
		return
	}
	fmt.Println("Redirecting to /log-in from distributor")
	c.Redirect(http.StatusFound, "/about")
}

// This handler shows the login page for the Manager and the Admin
func (VS *Server) showLogin(c *gin.Context) {
	render(c, gin.H{}, "login.html")
}

// This handler shows the login page for the Manager and the Admin
func (VS *Server) loadingContent(c *gin.Context) {
	render(c, gin.H{}, "loading-content.html")
}

// Shows the log-in page for the Worker
func (VS *Server) showLoginWorker(c *gin.Context) {
	workers := VS.getAllWorkersFromLocalVillage(c, "h_general", "showLoginWorker", "workers := VS.getAllWorkersFromLocalVillage")
	render(c, gin.H{
		"workers": workers,
	}, "login-worker.html")
}

// This handler checks if the user and password is correct and sets the cookie for the user
func (VS *Server) login(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	workerID := c.Request.FormValue("workerID")

	// Check that the Village introduced in the command line exists
	// if !VS.existLocalVillage(c, "h_general", "login", "VS.existLocalVillage") {
	if !VS.existLocalVillage(c, "h_general", "login", "VS.existLocalVillage") {
		// Check that if is not Central, and if the DB is empty (localvillage does not exist), then we redirect to import database for start
		if localVillage != "Central" {
			c.Redirect(http.StatusFound, "/settings/new-import/")
		} else {
			VS.getLocalVillage(c, "h_general", "login", "VS.getLocalVillage")
		}
		return
	}

	var user User
	exist := true
	if workerID != "" {
		id := bson.ObjectIdHex(workerID)
		user = VS.getUserByID(c, id, "h_general", "login", "user = v.getUserByID")
	} else {
		user, exist = VS.getUserByUsername(c, username, "h_general", "login", "user, exist = VS.getUserByUsername")
		// Every time the manager logs in we put his ID in the global variable
		if user.Role == Manager {
			if ManagerID != user.ID {
				DateLastManagerChange = getTime()
				ManagerID = user.ID
			}
		}
		if !exist {
			render(c, gin.H{}, "wrong-password.html")
			return
		}
	}

	if user.Role != Worker {
		ManagerID = user.ID
	}

	if user.Password == "" {
		setSession(c, user.Username, string(user.Role))
		if user.Role == Worker {
			fmt.Println("Redirecting to /set-password-worker")
			c.Redirect(http.StatusFound, "/settings/set-password-worker")
		} else {
			fmt.Println("Redirecting to /set-password")
			c.Redirect(http.StatusFound, "/settings/set-password")
		}
		return
	}

	// Check if the username/password combination is valid
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		// If the username/password is valid set the session cookie
		setSession(c, user.Username, string(user.Role))
		fmt.Println("Redirecting to /dashboard from login")
		c.Redirect(http.StatusFound, "/dashboard")
	} else {
		// Clear the cookie just in case
		// BUG
		// c.SetCookie("valuevillages", "", -1, "", "", http.SameSiteDefaultMode, false, true)
		c.Set("isLoggedIn", false)
		render(c, gin.H{}, "wrong-password.html")
	}

}

// Redirects when a user introduces a route that does not exists
func (VS *Server) redirect(c *gin.Context) {
	fmt.Println("Redirecting to /dashboard from redirect")
	c.Redirect(http.StatusFound, "/dashboard")
}

// Shows the recovery page when there is a fail in the program
func (VS *Server) recovery(c *gin.Context) {
	render(c, gin.H{
		"text": "Sorry, something went wrong :(",
	}, "wrong-feedback.html")
}

// Shows the Dashboard
func (VS *Server) dashboard(c *gin.Context) {
	role := getRole(c)
	serTypes := VS.getAllServiceTypes(c, "h_general", "dashboard", "serTypes := VS.getAllServiceTypes")

	if role == Worker {
		var typesWithAccess []ServiceType
		accesses := VS.getAllActiveAccessByUser(c, VS.getUserFomCookie(c, "h_general", "dashboard", "accesses := VS.getAllActiveAccessByUser").ID, "h_general", "dashboard", "getAllActiveAccessByUser(c, VS.getUserFomCookie")
		for _, tpe := range serTypes {
			for _, a := range accesses {
				service := VS.getServiceByID(c, a.IDService, "h_general", "dashboard", "service := VS.getServiceByID", true)
				if service.Type == tpe.ID {
					typesWithAccess = append(typesWithAccess, tpe)
				}
			}
		}
		serTypes = typesWithAccess
	}
	serHTML := VS.servicesTypeToServicesTypeHTML(c, serTypes)
	render(c, gin.H{
		"services": serHTML,
	}, "dashboard.html")
}

// Returns the current language of the App
func getLanguage(c *gin.Context) Languages {
	return LanguageAPP
}

// Changes the language of the app
func (VS *Server) changeLanguage(c *gin.Context) {
	lan, _ := VS.getParamFromURL(c, "language")
	language := Languages(lan)
	LanguageAPP = language
	c.Redirect(http.StatusFound, "/dashboard")
}

// Give a good feedback to the User
func (VS *Server) goodFeedback(c *gin.Context, URL string) {
	render(c, gin.H{}, "good-feedback.html")
	// We wait 2 seconds to redirect to give the user time to read the good feedback
	// time.Sleep(2 * time.Second)
	// c.Redirect(http.StatusFound, URL)
}

// Give a good feedback to the User
func (VS *Server) goodFeedbackSimple(c *gin.Context) {
	render(c, gin.H{}, "good-feedback.html")
}

// Give a bad feedback to the User
func (VS *Server) wrongFeedback(c *gin.Context, t string) {
	fmt.Print("Wrong feedback called")
	render(c, gin.H{
		"text": t,
	}, "wrong-feedback.html")
}

// Returns a true if a string is empty or false if not and gives an error in the UI if not
func (VS *Server) isEmpty(c *gin.Context, toCheck string, err string) bool {
	if toCheck != "" {
		return false
	}
	render(c, gin.H{
		"text": err,
	}, "wrong-feedback.html")
	return true
}
