package village

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// This handler shows the List of Video Courses with access in the Service Type
func (VS *Server) showListVideoCurses(c *gin.Context) {
	service, b := VS.getOurService(c, "h_videos", "showListVideoCurses", "service, b := VS.getOurService")
	if !b {
		return
	}
	videos := VS.getAllVideoCoursesByServiceType(c, service.TypeID, "h_videos", "showListVideoCurses", "videos := VS.getAllVideoCoursesByServiceType")
	videosHTML := VS.toViewVideoCourses(c, videos)

	JSON, err := json.Marshal(videosHTML)
	VS.checkOperation(c, "h_videos", "showListVideoCurses", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"JSON":        string(JSON),
		"videos":      videosHTML,
		"serviceType": service.TypeID,
	}, "tutorial-list.html")
}

// Shows the page to introduce the information of the new VideoCourse
func (VS *Server) showNewVideoCourse(c *gin.Context) {
	service, b := VS.getOurService(c, "h_videos", "showNewVideoCourse", "service, b := VS.getOurService")
	if !b {
		return
	}
	sType := VS.getAllServiceTypes(c, "h_videos", "showNewVideoCourse", "sType := VS.getAllServiceTypes")
	servicesTypeHTML := VS.servicesTypeToServicesTypeHTML(c, sType)
	render(c, gin.H{
		"steps":       0,
		"serviceType": service.TypeID,
		"sType":       servicesTypeHTML,
	}, "form-video-course.html")
}

// Stores the new VideoCourse information into the DataBase
func (VS *Server) newVideoCourse(c *gin.Context) {
	// We only can introduce new information if the Language is the main Language
	if getLanguage(c) != MainLanguage {
		VS.wrongFeedback(c, "To add a Video Course in a Language that is not the Main language of the app, you need to first add the step in the main language, change to the language desired for the modification, and modify the step that you want to translate")
		return
	}
	video, b := VS.getNewVideoCourse(c, bson.NewObjectId())
	if !b {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}

	if VS.addVideoCourseToDatabase(c, video, "h_videos", "newVideoCourse", "v.addProductToDatabase") {
		VS.registerAddDatabaseMovementFromCentral(c, video, video.ID)
		VS.goodFeedback(c, "videos/see-list/"+idservicetype.Hex())
	} else {
		VS.wrongFeedback(c, "Sorry, there is already one Video with that name.Please choose another name ")
	}
}

// Shows the page with the information of the Video course to modify
func (VS *Server) showEditVideoCourse(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	video, _ := VS.getVideoCourseByID(c, id, "h_videos", "showEditVideoCourse", "video, _ := VS.getVideoCourseByID", true)
	// Check for a translation in another language
	trans, b := VS.mayLookForTranslation(c, video.ID)
	if b {
		video = trans.(VideoCourse)
	}
	//Get the steps of the product
	steps := VS.getAllStepsFromSpecificVideoCourse(c, video.ID, "h_videos", "showEditVideoCourse", "steps := VS.getAllStepsFromSpecificVideoCourse")
	service, b := VS.getOurService(c, "h_videos", "showEditVideoCourse", "service, b := VS.getOurService")
	if !b {
		return
	}
	sType := VS.getAllServiceTypes(c, "h_videos", "showEditVideoCourse", "sType := VS.getAllServiceTypes")
	servicesTypeHTML := VS.servicesTypeToServicesTypeHTML(c, sType)
	render(c, gin.H{
		"product":     video,
		"edit":        true,
		"steps":       len(steps),
		"serviceType": service.TypeID,
		"sType":       servicesTypeHTML,
	}, "form-video-course.html")
}

// Saves the new information of the Video course to the Database
func (VS *Server) editVideoCourse(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	oldVideo, _ := VS.getVideoCourseByID(c, id, "h_videos", "editVideoCourse", "oldVideo, _ := VS.getVideoCourseByID", true)
	newVideo, b := VS.getNewVideoCourse(c, id)
	if !b {
		return
	}
	// Check if we are editing another language than the main one
	if VS.isTranslation(c, newVideo, newVideo.ID) {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	// Audit
	audit := VS.getAuditFromCentral(c, newVideo, id, Modified)
	modified := VS.dynamicAuditChange(c, newVideo, oldVideo, audit.ID)
	if VS.checkChange(c, audit.ID, "Name", oldVideo.Name, newVideo.Name) {
		modified = true
		VS.move(basePath+"/local-resources/video-courses/"+oldVideo.Name+"/", basePath+"/local-resources/video-courses/"+newVideo.Name+"/")
	}
	if modified {
		if VS.addVideoCourseToDatabase(c, newVideo, "h_videos", "editVideoCourse", "if VS.addVideoCourseToDatabase") {
			VS.registerAuditAndEvent(c, audit, newVideo, newVideo.ID)
		} else {
			VS.wrongFeedback(c, "Sorry, there is already one Video with that name.Please choose another name ")
		}
	}
	VS.goodFeedback(c, "videos/see-list/"+idservicetype.Hex())
}

// Shows the screen with all the Video course in the database
func (VS *Server) showVideosVideoCurses(c *gin.Context) {
	service, b := VS.getOurService(c, "h_videos", "showVideosVideoCurses", "service, b := VS.getOurService")
	if !b {
		return
	}
	products := VS.getAllVideoCoursesByServiceType(c, service.TypeID, "h_videos", "showVideosVideoCurses", "products := VS.getAllVideoCoursesByServiceType")
	render(c, gin.H{
		"products":    products,
		"serviceType": service.TypeID,
	}, "tutorial-video.html")
}

// Shows all the services to select the videos in the dashboard for the workers
func (VS *Server) videoDashboard(c *gin.Context) {
	serTypes := VS.getAllServiceTypes(c, "h_videos", "videoDashboard", "serTypes := VS.getAllServiceTypes")
	serHTML := VS.servicesTypeToServicesTypeHTML(c, serTypes)
	render(c, gin.H{
		"services": serHTML,
	}, "video-dashboard.html")
}

// Shows the screen with the description of the video course, Photo, checks...
func (VS *Server) showDescriptionVideoCourse(c *gin.Context) {
	service, b := VS.getOurService(c, "h_videos", "showDescriptionVideoCourse", "service, _ := VS.getOurService")
	if !b {
		return
	}
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	video, _ := VS.getVideoCourseByID(c, id, "h_videos", "showDescriptionVideoCourse", "video, _ := VS.getVideoCourseByID", true)
	// Check for a translation in another language
	trans, b := VS.mayLookForTranslation(c, video.ID)
	if b {
		video = trans.(VideoCourse)
	}
	// Get the Checklist from the database
	list := VS.getChecklistByVideo(c, video.ID, "h_videos", "showDescriptionVideoCourse", "list := VS.getChecklistByVideo")
	render(c, gin.H{
		"product":     video,
		"id":          id,
		"list":        list,
		"serviceType": service.TypeID,
	}, "description-video-course.html")
}

// Shows the screen with the form for New Check for the specific video course
func (VS *Server) showNewCheck(c *gin.Context) {
	service, b := VS.getOurService(c, "h_videos", "showNewCheck", "service, _ := VS.getOurService")
	if !b {
		return
	}
	id, b := VS.getIDParamFromURL(c, "idproduct")
	if !b {
		return
	}
	video, _ := VS.getVideoCourseByID(c, id, "h_videos", "showNewCheck", "video, _ := VS.getVideoCourseByID", true)
	render(c, gin.H{
		"product":     video,
		"id":          id,
		"serviceType": service.TypeID,
	}, "form-check.html")
}

// Creates a new check for the video course
func (VS *Server) newCheck(c *gin.Context) {
	idProduct, b := VS.getIDParamFromURL(c, "idproduct")
	if !b {
		return
	}
	// We only can introduce new information if the Language is the main Language
	if getLanguage(c) != MainLanguage {
		VS.wrongFeedback(c, "To add a Video Check in a Language that is not the Main language of the app, you need to first add the step in the main language, change to the language desired for the modification, and modify the step that you want to translate")
		return
	}
	newCheck, b := VS.getNewProductCheckList(c, bson.NewObjectId(), idProduct)
	if !b {
		return
	}
	if VS.isTranslation(c, newCheck, newCheck.ID) {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, newCheck, newCheck.ID, "h_videos", "newCheck", "v.addCheckToDatabase")
	VS.goodFeedback(c, "new-check/:idproduct/"+idservicetype.Hex())
}

// Shows the screen with the information of the Check to modify
func (VS *Server) showEditCheck(c *gin.Context) {
	idCheck, b := VS.getIDParamFromURL(c, "idcheck")
	if !b {
		return
	}
	check := VS.getCheckByID(c, idCheck, "h_videos", "showEditCheck", "check := VS.getCheckByID", true)
	service, b := VS.getOurService(c, "h_videos", "showEditCheck", "service, b := VS.getOurService")
	if !b {
		return
	}
	render(c, gin.H{
		"check":       check,
		"edit":        true,
		"serviceType": service.TypeID,
	}, "form-check.html")
}

// Updates the new information of the Check in the database
func (VS *Server) editCheck(c *gin.Context) {
	idCheck, b := VS.getIDParamFromURL(c, "idcheck")
	if !b {
		return
	}
	oldCheck := VS.getCheckByID(c, idCheck, "h_videos", "editCheck", "oldCheck := VS.getCheckByID", true)
	newCheck, b := VS.getNewProductCheckList(c, idCheck, oldCheck.IDVideoCourse)
	if !b {
		return
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	VS.editElementInDatabaseWithRegisterFromCentral(c, newCheck, oldCheck, newCheck.ID, "h_videos", "editCheck", "VS.editElementInDatabaseWithRegisterFromCentral")
	VS.goodFeedback(c, "new-check/:idproduct/"+idservicetype.Hex())
}

// Shows the page to introduce the information of the new Step
func (VS *Server) showNewStep(c *gin.Context) {
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	video, _ := VS.getVideoCourseByID(c, id, "h_videos", "showNewStep", "video, _ := VS.getVideoCourseByID", true)
	toolsTotal, materialsTotal := VS.getToolsAndMaterials(c, video)
	service, b := VS.getOurService(c, "h_videos", "showNewStep", "service, b := VS.getOurService")
	if !b {
		return
	}

	// Calculate the num of steps
	steps := VS.getAllStepsFromSpecificVideoCourse(c, id, "h_videos", "showNewStep", "steps := VS.getAllStepsFromSpecificVideoCourse")
	render(c, gin.H{
		"product":     video,
		"numStep":     len(steps) + 1,
		"tools":       toolsTotal,
		"materials":   materialsTotal,
		"serviceType": service.TypeID,
	}, "form-step.html")
}

// Stores the new Step information into the DataBase
func (VS *Server) newStep(c *gin.Context) {
	// We only can introduce new information if the Language is the main Language
	if getLanguage(c) != MainLanguage {
		VS.wrongFeedback(c, "To add a Step in a Language that is not the Main language of the app, you need to first add the step in the main language, change to the language desired for the modification, and modify the step that you want to translate")
		return
	}
	step, b := VS.getNewStep(c, bson.NewObjectId())
	if !b {
		return
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, step, step.ID, "h_videos", "newStep", "VS.addElementToDatabaseWithRegisterFromCentral")
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	VS.goodFeedback(c, "videos/see-list/"+idservicetype.Hex())
}

// Shows the screen with the Step of the specific Video course from the database
func (VS *Server) showStepVideoCourse(c *gin.Context) {
	service, b := VS.getOurService(c, "h_videos", "showStepVideoCourse", "service, _ := VS.getOurService")
	if !b {
		return
	}
	prodID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	product, _ := VS.getVideoCourseByID(c, prodID, "h_videos", "showStepVideoCourse", "product := v.getProductByID", true)
	steps := VS.getAllStepsFromSpecificVideoCourse(c, prodID, "h_videos", "showStepVideoCourse", "steps := v.getAllStepsFromSpecificProduct")
	if len(steps) == 0 {
		c.Redirect(http.StatusFound, "/videos/new-step/"+prodID.Hex()+"/"+service.TypeID.Hex())
		return
	}

	// Get the /number param of the URL
	nURL, b := VS.getParamFromURL(c, "step-number")
	if !b {
		return
	}
	numbernoURL, err := strconv.Atoi(nURL)
	VS.checkOperation(c, "h_videos", "showStepVideoCourse", "product := v.getProductByID", err)

	step := steps[numbernoURL]
	// Check for a translation in another language
	trans, b := VS.mayLookForTranslation(c, step.ID)
	if b {
		step = trans.(Step)
	}

	var pagination []Pagination
	var p Pagination
	for i := 0; i < len(steps); i++ {
		p.Active = false
		p.Fakenumb = i + 1
		pagination = append(pagination, p)
	}
	pagination[numbernoURL].Active = true
	toolsHTML, materialsHTML := VS.getToolsAndMaterialsHTMLFromStep(c, step)
	render(c, gin.H{
		"product":     product,
		"step":        step,
		"id":          prodID,
		"pagination":  pagination,
		"tools":       toolsHTML,
		"materials":   materialsHTML,
		"serviceType": service.TypeID,
	}, "step-video-course.html")
}

// Shows the problems related with the video course
func (VS *Server) showVideoCourseProblem(c *gin.Context) {
	idProd, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	product, _ := VS.getVideoCourseByID(c, idProd, "h_videos", "showVideoCourseProblem", "product := v.getProductByID", true)
	user := VS.getUserFomCookie(c, "h_videos", "showVideoCourseProblem", "user := VS.getUserFomCookie")

	// Get all Problems from the Product from the user
	problemsUser := VS.getProblemsFromUserByProduct(c, idProd, user.ID, "handlers_manager", "showVideoCourseProblem", "problemsUser := VS.getProblemsFromUserByProduct")
	// Get all Problems from the Product NOT from the user
	problemsOthers := VS.getProblemsNOTFromUserByProduct(c, idProd, user.ID, "handlers_manager", "showVideoCourseProblem", "problemsOthers := VS.getProblemsNOTFromUserByProduct")
	service, b := VS.getOurService(c, "h_videos", "showVideoCourseProblem", "service, _ := VS.getOurService")
	if !b {
		return
	}

	render(c, gin.H{
		"problemsUser":   problemsUser,
		"problemsOthers": problemsOthers,
		"product":        product,
		"serviceType":    service.TypeID,
	}, "video-course-problems.html")
}

// Shows the page with the information of the Step to modify
func (VS *Server) showEditStep(c *gin.Context) {
	service, b := VS.getOurService(c, "h_videos", "showEditStep", "service, _ := VS.getOurService")
	if !b {
		return
	}
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}

	// Get the /number param of the URL
	numberURL, b := VS.getParamFromURL(c, "current-step")
	if !b {
		return
	}
	n, err := strconv.Atoi(numberURL)
	VS.checkOperation(c, "h_videos", "showEditStep", "n, err := strconv.Atoi(numberURL)", err)
	numberNoURL := n + 1

	// Get the Product from DataBase
	product, _ := VS.getVideoCourseByID(c, id, "h_videos", "showEditStep", "product, _ := VS.getVideoCourseByID", true)
	steps := VS.getAllStepsFromSpecificVideoCourse(c, id, "h_videos", "showEditStep", "steps := VS.getAllStepsFromSpecificVideoCourse")
	tools := VS.getAllItemsForRent(c, service.TypeID)
	materials := VS.getAllItemsForSell(c, service.TypeID)
	// step := VS.getStepByProductIDAndIndexOrder(c, product.ID, n, "h_videos", "showEditStep", "step := VS.getStepByProductIDAndIndexOrder")
	step := steps[n]

	fmt.Print("All ssteps i DB")
	fmt.Print(len(steps))
	// Check for a translation in another language
	trans, b := VS.mayLookForTranslation(c, step.ID)
	if b {
		step = trans.(Step)
	}

	var pagination []Pagination
	var p Pagination
	for i := 0; i < len(steps); i++ {
		p.Active = false
		p.Fakenumb = i + 1
		pagination = append(pagination, p)
	}
	pagination[n].Active = true

	render(c, gin.H{
		"product":     product,
		"numStep":     numberNoURL,
		"tools":       tools,
		"step":        step,
		"edit":        true,
		"materials":   materials,
		"id":          id.Hex(),
		"pagination":  pagination,
		"serviceType": service.TypeID,
	}, "form-step.html")
}

// Saves the new information of the Step to the Database
func (VS *Server) editStep(c *gin.Context) {
	prodID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	indexOrder, b := VS.getIntFromHTML(c, "indexOrder", "Step")
	if !b {
		return
	}
	oldStep := VS.getStepByProductIDAndIndexOrder(c, prodID, indexOrder, "h_videos", "editStep", "oldStep := VS.getStepByProductIDAndIndexOrder")
	newStep, b := VS.getNewStep(c, oldStep.ID)
	if !b {
		return
	}

	// Here for taking the different values of the checkboxes with the same name
	tools := c.Request.MultipartForm.Value["tool"]
	materials := c.Request.MultipartForm.Value["material"]

	// Warnings
	var warnings []string
	var w string
	for i := 1; i <= 6; i++ {
		w = c.PostForm("warning" + strconv.Itoa(i))
		if w != "" {
			warnings = append(warnings, w)
		}
	}
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
	if VS.isTranslation(c, newStep, newStep.ID) {
		return
	}
	VS.editElementInDatabaseWithRegisterFromCentral(c, newStep, oldStep, newStep.ID, "h_videos", "editStep", "VS.editElementInDatabaseWithRegisterFromCentral")
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	VS.goodFeedback(c, "videos/see-list/"+idservicetype.Hex())
}

// Shows the page to introduce the information of the new Problem
func (VS *Server) showNewProblem(c *gin.Context) {
	prodID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	video, _ := VS.getVideoCourseByID(c, prodID, "h_videos", "showNewProblem", "video, _ := VS.getVideoCourseByID", true)
	service, b := VS.getOurService(c, "h_videos", "showNewProblem", "service, b := VS.getOurService")
	if !b {
		return
	}
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_videos", "showNewProblem", "workers := VS.getWorkersFromSpecificService")
	render(c, gin.H{
		"workers":     workers,
		"video":       video,
		"serviceType": service.TypeID,
	}, "new-problem.html")
}

// Stores the new Problem information into the DataBase
func (VS *Server) newProblem(c *gin.Context) {
	videoID, _ := VS.getIDParamFromURL(c, "id")
	product, _ := VS.getVideoCourseByID(c, videoID, "h_videos", "newProblem", "product, _ := VS.getVideoCourseByID", true)
	workerID, b := VS.getIDFromHTML(c, "workerID", "Purchase")
	if !b {
		return
	}
	photo, b := VS.getFileFromHTML(c, "photo", "Purchase", basePath+"/local-resources/video-courses/"+product.Name+"/problems/photos/")
	if !b {
		return
	}
	audio, b := VS.getFileFromHTML(c, "audio", "Purchase", basePath+"/local-resources/video-courses/"+product.Name+"/problems/audios/")
	if !b {
		return
	}

	// Get the text of the problem
	problemText := c.Request.FormValue("problemText")
	problem := VideoProblem{
		ID:            bson.NewObjectId(),
		IDVideoCourse: videoID,
		TextProblem:   problemText,
		PhotoProblem:  photo,
		IDWorker:      workerID,
		Audio:         audio,
	}
	idservicetype, b := VS.getIDParamFromURL(c, "idservicetype")
	if !b {
		return
	}
	VS.addElementToDatabaseWithRegisterFromCentral(c, problem, problem.ID, "h_videos", "newProblem", "VS.addElementToDatabaseWithRegisterFromCentral")
	VS.goodFeedback(c, "/video-course-problems/"+videoID.Hex()+"/"+idservicetype.Hex())
}
