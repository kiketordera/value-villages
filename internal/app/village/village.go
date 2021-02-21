package village

import (
	"context"
	"encoding/gob"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timshannon/bolthold"
	"gopkg.in/mgo.v2/bson"
)

// These are the global variables that we can use all accross the project
var (
	basePath = os.Getenv("GOPATH") + "/src/github.com/kiketordera/value-villages"
	// This are the rows in the tables in the UI
	titemsPerPageWorker = 10
	itemsPerTable       = 10
	// localVillage is the village introduced in the terminal when you initialice the app
	localVillage string
	// sessionTime is the seconds the session will be active
	sessionTime = 108000
	// This is the key to sign the cookies
	tokenSigningKey = []byte("SuperFancyToken:D")
	// IDVisualization is the dictionary with the ID of the service for the visualization of the data
	IDVisualization = make(map[bson.ObjectId]bson.ObjectId)
	Variables       = make(map[string]bson.ObjectId)

	// Dates
	FromDate int64
	ToDate   int64

	// ManagerID is the ID of the Manager who is rolling the WK in the moment (by cookie)
	ManagerID             bson.ObjectId
	DateLastManagerChange int64

	// InDevelopment is the flag to make some test to the app
	InDevelopment = true
	// FakeDate is the date we can set up in the app to check the times in the transactions and make tests
	FakeDate int64
	// TimeBeforeChanging is for calculate the time when we changed and fake that we are advancing in the time
	TimeBeforeChanging int64
	// LanguageAPP is the language of the app
	LanguageAPP  = English
	MainLanguage = English
	// URLbase is the URL used, in local will be localhost:8080, in the web will be
	URLbase = os.Getenv("URL")

	// TESTID || This information is only to make the test to the Duplicates in the DataBar
	TESTID = bson.NewObjectId()

	/** COLORS **/
	Blue       = color.RGBA{0, 163, 210, 0xff}
	Red        = color.RGBA{207, 31, 38, 0xff}
	LightBrown = color.RGBA{169, 124, 80, 0xff}
	DarkBrown  = color.RGBA{60, 36, 21, 0xff}
	Orange     = color.RGBA{247, 148, 29, 0xff}
	Pink       = color.RGBA{237, 3, 124, 0xff}
	Green      = color.RGBA{99, 177, 144, 0xff}
	Black      = color.RGBA{0, 0, 0, 0xff}
	White      = color.RGBA{255, 255, 255, 0xff}
)

// New DataBase object
func New() Server {
	localVillage = os.Getenv("VILLAGE")
	if localVillage == "" {
		log.Fatal("Village name not set")
	}

	// We open or create the DataBase and the Directory
	db, err := bolthold.Open(basePath+"/local-resources/"+localVillage+".db", 0666, nil)
	if err != nil {
		fmt.Println("We try to make the directory because we did not found it, here is the error: ")
		fmt.Println(err)
		fmt.Print("Here is the os.Getenv(GOPATH): ")
		fmt.Println(os.Getenv("GOPATH"))
		if os.Getenv("GOPATH") == "" {
			fmt.Println("Try writing the GOPATH with: ")
			fmt.Println("\033[31mexport GOPATH=$HOME/go\033[37m")
		}
		os.MkdirAll(basePath+"/local-resources/", os.ModePerm)
		db, err = bolthold.Open(basePath+"/local-resources/"+localVillage+".db", 0666, nil)
		if err != nil {
			// It should be log.Fatal instead of fmt.Print
			fmt.Print("We try to make the directory and did not work, here is the error: ")
			fmt.Print(err)
		}
	}
	DB := Server{
		DataBase: db,
	}

	// Admin user, central Village, central warehouse... for first log-in
	if localVillage == "Central" {
		DB.createAdminUser()
	}
	// Make sure the village exists in the Database
	// var village []Village
	// DB.DataBase.Find(&village, bolthold.Where("Name").Eq(localVillage))
	// if len(village) == 0 {
	// 	log.Fatal("Village name not found in the DataBase!")
	// }
	return DB
}

// Start the web server for a Server object
func (VS *Server) Start(wg *sync.WaitGroup) {
	// When the process is end, warning the group that is end
	defer wg.Done()
	FakeDate = time.Now().Unix()
	TimeBeforeChanging = time.Now().Unix()

	/*************** Here we initialice the gin framework with the custom parameters ***************/
	router := gin.Default()
	// Source of HTML
	// Dynamic routes for the images, uses the relative path whenever is in google Cloud, Raspberry or local development
	router.Static("/static", basePath+"/web/static")
	router.Static("/statichtml", basePath+"/web/html")
	router.Static("/local", basePath+"/local-resources")

	// Register the functions we use in all the HTML
	// // router.Delims("{[{", "}]}")
	// router.SetFuncMap(template.FuncMap{
	// 	"Language": getLanguage,
	// })
	// // They need to be in this order, this one is related with the functions in the templates
	// router.LoadHTMLFiles(basePath + "/web/html/*/*.html")
	// Then in the template, you call the function like: {{Language}}
	// This one is related with render the HTML
	router.LoadHTMLGlob(basePath + "/web/html/*/*.html")

	// VARIABLES
	FromDate = 0
	ToDate = getTime()

	// Register the elements that are available in different languages, and also all the elements to sync the events
	gob.Register(Server{})
	gob.Register(Village{})
	gob.Register(Access{})
	gob.Register(Sale{})
	gob.Register(User{})
	gob.Register(Synchronization{})
	gob.Register(Category{})
	gob.Register(CategoryType{})
	gob.Register(Item{})
	gob.Register(Stock{})
	gob.Register(Delivery{})
	gob.Register(WorkerOrder{})
	gob.Register(Payment{})
	gob.Register(VideoCourse{})
	gob.Register(VideoProblem{})
	gob.Register(Step{})
	gob.Register(Service{})
	gob.Register(ServiceType{})
	gob.Register(ServiceOrder{})
	gob.Register(VideoCheckList{})
	gob.Register(Conversation{})
	gob.Register(Assignment{})
	gob.Register(Message{})
	gob.Register(Report{})
	gob.Register(PDF{})
	gob.Register(Animal{})
	gob.Register(Translation{})
	gob.Register(ToDo{})
	gob.Register(ToDoChecked{})
	gob.Register(FieldDescription{})
	gob.Register(FieldCheckbox{})
	gob.Register(FieldPhoto{})
	gob.Register(FieldNumber{})
	gob.Register(CheckStock{})

	// Redirects when there is wrong route
	router.NoRoute(VS.redirect)

	/*************** GIN FRAMEWORKS ***************/
	// This implements a Favicon in our App
	// router.Use(favicon.New("./favicon.ico"))
	// app.Use(favicon.New("./favicon.ico"))
	// router.StaticFile("/favicon.ico", "./favicon.ico")

	// Uses the nice recovery for showing the HTML page for a nice user experience when an bug comes
	// router.Use(nice.Recovery(VS.recovery))

	// Log every request in the system
	// router.Use(logger.SetLogger())

	// Log every error in the system
	// router.Use(logger.SetLogger())

	// Checks user cookie for each request
	router.Use(checkToken())

	/*************** 	HANDLERS	 ***************/
	// Group article related routes from Workshop together
	comunication := router.Group("/comunication")
	{
		// Conversations
		comunication.GET("/conversations", ensureLoggedIn(), VS.showConversations)
		comunication.GET("/choose-users-conversation", ensureLoggedIn(), VS.showChooseUsersForConversation)
		comunication.POST("/choose-users-conversation", ensureLoggedIn(), VS.newConversation)
		comunication.GET("/chat/:id", ensureLoggedIn(), VS.showChat)
		comunication.POST("/chat/:id", ensureLoggedIn(), VS.addMessage)

		// Reports
		comunication.GET("/reports", ensureLoggedIn(), VS.showReports)
		comunication.GET("/new-report", ensureLoggedIn(), VS.showNewReport)
		comunication.POST("/new-report", ensureLoggedIn(), VS.newReport)
		comunication.GET("/see-report/:id", ensureLoggedIn(), VS.seeReport)
		comunication.POST("/see-report/:id", ensureLoggedIn(), VS.addMessageReport)
	}

	// Group article related routes from DATA together
	data := router.Group("/data")
	{
		// Users
		data.GET("/all-users", ensureIsManagerOrAdmin(), VS.showUsers)
		// New User
		data.GET("/new-user", ensureIsManagerOrAdmin(), VS.showNewUser)
		data.POST("/new-user", ensureIsManagerOrAdmin(), VS.newUser)
		// Edit User
		data.GET("/edit-user/:id", ensureIsManagerOrAdmin(), VS.showEditUser)
		data.POST("/edit-user/:id", ensureIsManagerOrAdmin(), VS.editUser)

		// Access
		data.GET("/access", ensureIsManagerOrAdmin(), VS.showAccess)
		data.GET("/edit-access/:id", ensureIsManagerOrAdmin(), VS.showEditAccess)
		data.POST("/edit-access/:id", ensureIsManagerOrAdmin(), VS.editAccess)

		// Villages & Services
		data.GET("/villages", ensureIsAdmin(), VS.showVillages)
		// New Village
		data.GET("/new-village", ensureIsAdmin(), VS.showNewVillage)
		data.POST("/new-village", ensureIsAdmin(), VS.newVillage)
		// Edit Village
		data.GET("/edit-village/:id", ensureIsAdmin(), VS.showEditVillage)
		data.POST("/edit-village/:id", ensureIsAdmin(), VS.editVillage)
		// New Service
		data.GET("/new-service", ensureIsAdmin(), VS.showNewService)
		data.POST("/new-service", ensureIsAdmin(), VS.newService)
		// Edit Service
		data.GET("/edit-service/:id", ensureIsAdmin(), VS.showEditService)
		data.POST("/edit-service/:id", ensureIsAdmin(), VS.editService)
		// Type of service
		data.GET("/types-service", ensureIsAdmin(), VS.showTypeServices)
		// New Type of service
		data.GET("/new-service-type", ensureIsAdmin(), VS.showNewServiceType)
		data.POST("/new-service-type", ensureIsAdmin(), VS.newServiceType)
		// Edit Type of service
		data.GET("/edit-service-type/:id", ensureIsAdmin(), VS.showEditServiceType)
		data.POST("/edit-service-type/:id", ensureIsAdmin(), VS.editServiceType)

		// Categories
		data.GET("/categories", ensureIsManagerOrAdmin(), VS.showCategories)
		data.GET("/new-category", ensureIsAdmin(), VS.showNewCategory)
		data.POST("/new-category", ensureIsAdmin(), VS.newCategory)
		data.GET("/edit-category/:id", ensureIsAdmin(), VS.showEditCategory)
		data.POST("/edit-category/:id", ensureIsAdmin(), VS.editCategory)
		data.GET("/see-items/:id", ensureIsManagerOrAdmin(), VS.showItems)
		data.GET("/new-item/:id", ensureIsAdmin(), VS.showNewItem)
		data.POST("/new-item/:id", ensureIsAdmin(), VS.newItem)
		data.GET("/edit-item/:id", ensureIsAdmin(), VS.showEditItem)
		data.POST("/edit-item/:id", ensureIsAdmin(), VS.editItem)
		data.GET("/new-type", ensureIsAdmin(), VS.showNewType)
		data.POST("/new-type", ensureIsAdmin(), VS.newType)
		data.GET("/edit-type/:id", ensureIsAdmin(), VS.showEditType)
		data.POST("/edit-type/:id", ensureIsAdmin(), VS.editType)
	}

	// Group article related routes from Workshop together
	deliveries := router.Group("/deliveries")
	{
		// Send from central
		deliveries.GET("/stock-created", ensureIsManagerOrAdmin(), VS.showDeliveriesFromCentral)
		deliveries.GET("/create-new-stock", ensureIsManagerOrAdmin(), VS.chooseServiceForDeliveryFromCentral)
		deliveries.GET("/new-stock", ensureIsManagerOrAdmin(), VS.chooseServiceForDeliveryFromCentral)
		deliveries.GET("/new-stock/:idservice/:iddelivery", ensureIsManagerOrAdmin(), VS.showItemsDeliveryFromCentral)
		deliveries.POST("/new-stock/:idservice/:iddelivery", ensureIsManagerOrAdmin(), VS.addItemsDeliveryGlobal)
		deliveries.GET("/delete-stock/:idstock", ensureIsManagerOrAdmin(), VS.deleteItemDeliveryGlobal)

		// Receive in the service
		deliveries.GET("/check-delivery/:iddelivery", ensureIsManagerOrAdmin(), VS.showConfirmDelivery)
		deliveries.POST("/check-delivery/:iddelivery", ensureIsAdmin(), VS.confirmDelivery)

		// Send from services
		deliveries.GET("/choose-service/:idservicetype", ensureIsManagerOrAdmin(), VS.chooseServiceForDeliveryFromService)
		deliveries.GET("/new-delivery/:idservice/:iddelivery/:idservicetype", ensureIsManagerOrAdmin(), VS.showItemsDeliveryFromService)
		deliveries.POST("/new-delivery/:idservice/:iddelivery/:idservicetype", ensureIsManagerOrAdmin(), VS.addItemsDeliveryFromService)
	}

	// Group article related routes from GENERAL together
	development := router.Group("/development")
	{
		// Time
		development.GET("/show-current-time", ensureIsAdmin(), VS.showCurrentTime(current))
		development.GET("/reset-time", ensureIsAdmin(), VS.showCurrentTime(reset))
		development.GET("/day-plus", ensureIsAdmin(), VS.showCurrentTime(dayPlus))
		development.GET("/day-minus", ensureIsAdmin(), VS.showCurrentTime(dayMinus))
		development.GET("/week-plus", ensureIsAdmin(), VS.showCurrentTime(weekPlus))
		development.GET("/week-minus", ensureIsAdmin(), VS.showCurrentTime(weekMinus))

		// Test Duplicated
		development.GET("/choose-test", ensureIsAdmin(), VS.chooseTest)
		development.GET("/make-test-duplicated", ensureIsAdmin(), VS.makeTestDuplicate)
		development.GET("/delete-test-duplicated", ensureIsAdmin(), VS.deleteInformationTestDuplicated)
		development.GET("/test-zip", ensureIsAdmin(), VS.testAudio)
		development.POST("/test-zip", ensureIsAdmin(), VS.newAudio)

	}

	// Group article related routes from GENERAL together
	general := router.Group("/")
	{
		// Presentation
		general.GET("/about", VS.about)
		general.GET("/design", VS.design)
		general.GET("/features", VS.features)
		general.GET("/manuals", VS.manuals)
		// Manuals
		general.GET("/manuals/:name", VS.showPDFDocument)

		// Log-in
		general.GET("/log-in", ensureNotLoggedIn(), VS.showLogin)
		general.POST("/log-in", ensureNotLoggedIn(), VS.login)
		general.GET("/log-out", ensureLoggedIn(), VS.logout)

		// Log in worker
		general.GET("log-in-worker", ensureNotLoggedIn(), VS.showLoginWorker)
		general.POST("/log-in-worker", ensureNotLoggedIn(), VS.login)

		// Dashboard
		general.GET("/dashboard", ensureLoggedIn(), VS.dashboard)

		general.GET("/good-feedback", ensureLoggedIn(), VS.goodFeedbackSimple)

		// Languages
		general.GET("/change-language/:language", ensureLoggedIn(), VS.changeLanguage)

	}
	// Group article related routes from GENERAL together
	performance := router.Group("/performance")
	{
		// Performance
		performance.GET("/self/:idservicetype", ensureIsWorkshopRelated(), VS.showPerformanceSelf)
		performance.GET("/local/:idservicetype", ensureIsManagerOrAdmin(), VS.showPerformanceLocalService)
		performance.GET("/all/:idservicetype", ensureIsAdmin(), VS.showPerformanceAllTypeService)
		performance.GET("/service/:id/:idservicetype", ensureIsManagerOrAdmin(), VS.showPerformanceSpecificService)
		performance.GET("/worker/:id/:idservicetype", ensureIsManagerOrAdmin(), VS.showPerformanceWorker)
	}

	// Group article related routes from GENERAL together
	records := router.Group("/records")
	{
		// Records
		records.GET("/see/actions-manager/:idservicetype", ensureIsManagerOrAdmin(), VS.showRecordsActionsManager)
		records.POST("/see/actions-manager/:idservicetype", ensureIsManagerOrAdmin(), VS.setUpService)
		records.GET("/see/actions-worker/:idservicetype", ensureIsWorker(), VS.showRecordsActionsWorker)
		records.GET("/choose-service/:idservicetype", ensureLoggedIn(), VS.showChooseServicePicker)
		records.POST("/choose-service/:idservicetype", ensureLoggedIn(), VS.setUpService)
		// Assignments
		records.GET("/assignments/:idservicetype", ensureIsManagerOrAdmin(), VS.showToolsAssignments)
		records.GET("/new-assignment/:idservicetype", ensureIsManagerOrAdmin(), VS.showNewAssignment)
		records.POST("/new-assignment/:idservicetype", ensureIsManagerOrAdmin(), VS.newAssignment)
		records.GET("/change-assignment/:id/:idservicetype", ensureIsManagerOrAdmin(), VS.showChangeAssignment)
		records.POST("/change-assignment/:id/:idservicetype", ensureIsManagerOrAdmin(), VS.changeAssignment)
		// Payments
		records.GET("/see/payments/:idservicetype", ensureIsManagerOrAdmin(), VS.showRecordsPayments)
		records.GET("/new/payment/:idservicetype", ensureIsManagerOrAdmin(), VS.showNewPayment)
		records.POST("/new/payment/:idservicetype", ensureIsManagerOrAdmin(), VS.newPayment)
		records.GET("/see/payment/:id/:idservicetype", ensureIsWorkshopRelated(), VS.seePaymentWorkerService)

		// Service Orders
		records.GET("/service-orders/:idservicetype", ensureIsManagerOrAdmin(), VS.showServiceOrders)
		records.GET("/new-service-order/:idservicetype", ensureIsAdmin(), VS.showNewServiceOrder)
		records.POST("/new-service-order/:idservicetype", ensureIsAdmin(), VS.newServiceOrder)
		records.GET("/edit-service-order/:id/:idservicetype", ensureIsAdmin(), VS.showEditServiceOrder)
		records.POST("/edit-service-order/:id/:idservicetype", ensureIsAdmin(), VS.editServiceOrder)

		// Worker Orders
		records.GET("/worker-orders/:idservicetype", ensureIsManagerOrAdmin(), VS.showOrdersWorker)
		records.GET("/new-worker-order/:idservicetype", ensureIsManagerOrAdmin(), VS.showNewWorkerOrder)
		records.POST("/new-worker-order/:idservicetype", ensureIsManagerOrAdmin(), VS.newWorkerOrder)

		// Stock generated
		records.GET("/stock-generated/:idservicetype", ensureIsManagerOrAdmin(), VS.showStockGenerated)
		records.GET("/workers-for-new-purchase-from-order/:idservicetype", ensureIsManagerOrAdmin(), VS.showWorkersForPurchase)
		records.GET("/new-purchase-from-order/:idworker/:idservicetype", ensureIsManagerOrAdmin(), VS.showNewPurchase)
		records.POST("/new-purchase-from-order/:idworker/:idservicetype", ensureIsManagerOrAdmin(), VS.newPurchase(false))
		records.GET("/new-free-purchase/:idservicetype", ensureIsManagerOrAdmin(), VS.showNewFreePurchase)
		records.POST("/new-free-purchase/:idservicetype", ensureIsManagerOrAdmin(), VS.newPurchase(true))

		// Deliveries
		records.GET("/send-or-receive-delivery/:idservicetype", ensureIsManagerOrAdmin(), VS.chooseDeliveriesSentReceived)
		records.GET("/deliveries-sent/:idservicetype", ensureIsManagerOrAdmin(), VS.showDeliveriesService(Sent))
		records.GET("/deliveries-received/:idservicetype", ensureIsManagerOrAdmin(), VS.showDeliveriesService(Received))
		records.GET("/check-delivery/:iddelivery/:idservicetype", ensureIsManagerOrAdmin(), VS.showConfirmDelivery)
		records.POST("/check-delivery/:iddelivery/:idservicetype", ensureIsAdmin(), VS.confirmDelivery)

		// Stock
		records.GET("/choose-type-stock/:idservicetype", ensureIsManagerOrAdmin(), VS.chooseTypeItem)
		records.GET("/see-stock/:itemtype/:idservicetype", ensureIsManagerOrAdmin(), VS.showStockService)
		records.GET("/choose-service-new-delivery/:idservicetype", ensureIsManagerOrAdmin(), VS.showStockService)

		// Checks
		records.GET("/check-service/:idservicetype", ensureIsManagerOrAdmin(), VS.showConfirmCheckApp)
		records.POST("/check-service/:idservicetype", ensureIsManagerOrAdmin(), VS.confirmCheckApp)

		// New Sale
		records.GET("/worker-sales/:idservicetype", ensureIsManagerOrAdmin(), VS.showWorkerSales)
		records.GET("/new-worker-sale/:idservicetype", ensureIsManagerOrAdmin(), VS.showNewSale)
		records.POST("/new-worker-sale/:idservicetype", ensureIsManagerOrAdmin(), VS.newSale)

		// WORKER Services --  Service related
		records.GET("/wk/index/:idservicetype", ensureIsWorker(), VS.showIndex)
		records.GET("/wk-sales/:idservicetype", ensureIsWorker(), VS.showRecordsSalesWorkerService)
		records.GET("/wk/orders/:idservicetype", ensureIsWorker(), VS.showRecordsOrdersWorkerService)
		records.GET("/wk/payments/:idservicetype", ensureIsWorker(), VS.showRecordsPaymentsWorkerService)
		records.GET("/wk/see/payment/:id/:idservicetype", ensureIsWorkshopRelated(), VS.seePaymentWorkerService)

		// WORKER GENERAL -- Not service related, this is global information for the user
		records.GET("/wk/wallet", ensureIsWorker(), VS.showWalletWorker)
		records.GET("/wk/sales", ensureIsWorker(), VS.showRecordsSalesWorkerGlobal)
		records.GET("/wk/orders", ensureIsWorker(), VS.showRecordsOrdersWorkerGlobal)
		records.GET("/wk/items", ensureIsWorker(), VS.showRecordsItemsWorker)
		records.GET("/wk/payments", ensureIsWorker(), VS.showRecordsPaymentsWorkerGlobal)
		records.GET("/wk/see/payment/:id", ensureIsWorkshopRelated(), VS.seePaymentWorkerGlobal)
	}

	// Group article related routes from Reports together
	reports := router.Group("/reports")
	{
		// QR
		reports.GET("/show-qr-codes", ensureIsAdmin(), VS.showQRCodes)
		reports.GET("/new-qr", ensureIsAdmin(), VS.showNewQRCodesForm)
		reports.POST("/new-qr", ensureIsAdmin(), VS.newPDFqr)
		reports.GET("/show-pdf/:name", ensureIsAdmin(), VS.showPDF)
	}

	// Group article related routes together
	settings := router.Group("/settings")
	{
		// Activity
		settings.GET("/activity", ensureLoggedIn(), VS.activity)
		settings.POST("/activity", ensureIsManagerOrAdmin(), VS.setUpServiceForAudits)
		// Activity Changes
		settings.GET("/activity/:id", ensureLoggedIn(), VS.activityChanges)

		// Sync
		settings.GET("/choose-sync-option", ensureIsAdmin(), VS.chooseOptionToSync)
		settings.GET("/exports", ensureIsAdmin(), VS.showSynchronizations(exp))
		settings.GET("/choose-village-to-sync", ensureIsAdmin(), VS.chooseVillageToSync)
		settings.GET("/new-export/:id", ensureIsAdmin(), VS.newSyncExport)
		settings.GET("/imports", ensureIsAdmin(), VS.showSynchronizations(imp))
		settings.GET("/new-import", anyoneCanSee(), VS.newImportDataBase)
		settings.POST("/new-import", anyoneCanSee(), VS.importDataBase)

		// General
		settings.GET("/general", ensureLoggedIn(), VS.showGeneralSettings)
		// Reset password
		settings.GET("/set-password", ensureLoggedIn(), VS.showSetPassword)
		settings.POST("/set-password", ensureLoggedIn(), VS.setPassword)
		// Reset password
		settings.GET("/set-password-worker", ensureLoggedIn(), VS.showSetPasswordWorker)
		settings.POST("/set-password-worker", ensureLoggedIn(), VS.setPassword)
		settings.GET("/set-password/:id", ensureIsManagerOrAdmin(), VS.resetPassword)
		// Introduce new password
		settings.GET("/change-password", ensureLoggedIn(), VS.showSetPassword)
		settings.POST("/change-password", ensureLoggedIn(), VS.setPassword)
	}

	// Group article related routes from GENERAL together
	slaughterhouse := router.Group("/slaughterhouse")
	{
		slaughterhouse.GET("/animals", ensureIsManagerOrAdmin(), VS.showAnimals)
		slaughterhouse.GET("/parts-animals/:animal", ensureIsManagerOrAdmin(), VS.showPartsAnimals)
		slaughterhouse.POST("/parts-animals/:animal", ensureIsManagerOrAdmin(), VS.newAnimal)
	}

	// Group article related routes from Workshop together
	todo := router.Group("/to-do")
	{
		todo.GET("/what-to-check", ensureIsManagerOrAdmin(), VS.showItemsToCheck)
		todo.GET("/new-to-do", ensureIsManagerOrAdmin(), VS.showNewToDo)
		todo.POST("/new-to-do", ensureIsManagerOrAdmin(), VS.newToDo)
		todo.GET("/see-checks/:id", ensureIsManagerOrAdmin(), VS.seeChecks)

		todo.GET("/to-dos-user", ensureLoggedIn(), VS.showToDosUser)
		todo.GET("/make-to-do/:id", ensureLoggedIn(), VS.showMakeToDo)
		todo.POST("/make-to-do/:id", ensureLoggedIn(), VS.makeToDo)
	}

	// Group article related routes from GENERAL together
	videos := router.Group("/videos")
	{
		videos.GET("/see-list/:idservicetype", ensureIsManagerOrAdmin(), VS.showListVideoCurses)
		videos.GET("/see-videos/:idservicetype", ensureIsWorkshopRelated(), VS.showVideosVideoCurses)
		// See properties
		videos.GET("/description-video-course/:id/:idservicetype", ensureIsWorkshopRelated(), VS.showDescriptionVideoCourse)
		videos.GET("/video-course-problems/:id/:idservicetype", ensureIsWorkshopRelated(), VS.showVideoCourseProblem)
		videos.GET("/step-video-course/:id/:step-number/:idservicetype", ensureIsWorkshopRelated(), VS.showStepVideoCourse)

		// Create and edit
		videos.GET("/new-video-course/:idservicetype", ensureIsAdmin(), VS.showNewVideoCourse)
		videos.POST("/new-video-course/:idservicetype", ensureIsAdmin(), VS.newVideoCourse)
		videos.GET("/edit-video-course/:id/:idservicetype", ensureIsAdmin(), VS.showEditVideoCourse)
		videos.POST("/edit-video-course/:id/:idservicetype", ensureIsAdmin(), VS.editVideoCourse)

		// Checklist
		videos.GET("/new-check/:idproduct/:idservicetype", ensureIsAdmin(), VS.showNewCheck)
		videos.POST("/new-check/:idproduct/:idservicetype", ensureIsAdmin(), VS.newCheck)
		videos.GET("/edit-check/:idcheck/:idservicetype", ensureIsAdmin(), VS.showEditCheck)
		videos.POST("/edit-check/:idcheck/:idservicetype", ensureIsAdmin(), VS.editCheck)
		// Steps
		videos.GET("/new-step/:id/:idservicetype", ensureIsAdmin(), VS.showNewStep)
		videos.POST("/new-step/:id/:idservicetype", ensureIsAdmin(), VS.newStep)
		videos.GET("/edit-step/:id/:current-step/:idservicetype", ensureIsAdmin(), VS.showEditStep)
		videos.POST("/edit-step/:id/:current-step/:idservicetype", ensureIsAdmin(), VS.editStep)

		// Problems
		videos.GET("/video-course-new-problem/:id/:idservicetype", ensureIsWorkshopRelated(), VS.showNewProblem)
		videos.POST("/video-course-new-problem/:id/:idservicetype", ensureIsWorkshopRelated(), VS.newProblem)
		// The videos for the worker without access to the service
		videos.GET("/see-all-videos", ensureIsWorkshopRelated(), VS.videoDashboard)
	}

	/*************** 	OTHER PARAMETERS	 ***************/

	// Declare a web server in the 8080 port
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	// Initialice a web server already declare, in background due to GO instruction
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful process shutdown from interrupts
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	// We wait here and do nothing until the server is told to shut down
	<-quit

	// Here we shut down the server SRV and give 5 seconds for it to shut down or it would be killed
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

// This method renders all templates and export the HTML to the browser the HTML
func render(c *gin.Context, data gin.H, templateName string) {
	// loggedInInterface, e := c.Get("isLoggedIn")
	// if e {
	// 	data["isLoggedIn"] = loggedInInterface.(bool)
	// }
	// } else {
	// 	fmt.Println("Redirecting to /log-in from render 1")
	// 	c.Redirect(http.StatusFound, "/log-in")
	// 	return
	// }
	// roleInterface, e := c.Get("role")
	// if e {
	// 	data["role"] = roleInterface.(string)
	// }
	// } else {
	// 	fmt.Println("Redirecting to /log-in from render 2")
	// 	c.Redirect(http.StatusFound, "/log-in")
	// 	return
	// }

	data["GLOBAL"] = getGobalTemplateVariables(c)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	default:
		// Respond with HTML

		c.HTML(http.StatusOK, templateName, data)
	}
}
