package village

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// VILLAGES & SERVICES

// This method test if we can create 2 Villages with the same name, returns true if yes
func (VS *Server) testDuplicateVillage(c *gin.Context) bool {
	village := Village{
		ID:   bson.NewObjectId(),
		Name: "test",
	}
	return VS.addVillageToDatabase(c, village, "test", "testDuplicateVillage", "VS.addVillageToDatabase(c, village")
}

// This method test if we can create 2 Services with the same name in the same village, returns true if yes
func (VS *Server) testDuplicateService(c *gin.Context) bool {
	service := Service{
		ID:        bson.NewObjectId(),
		Name:      "test",
		IDVillage: TESTID,
	}
	return VS.addServiceToDatabase(c, service, "test", "testDuplicateService", "return VS.addServiceToDatabase")
}

// USERS

// This method test if we can create 2 workers with the same name in the same village, returns true if yes
func (VS *Server) testDuplicateUser(c *gin.Context) bool {
	user := User{
		ID:        bson.NewObjectId(),
		Username:  "test",
		Name:      "test",
		Role:      Worker,
		IDVillage: TESTID,
	}
	return VS.addUserToDatabase(c, user, "test", "testDuplicateUser", "return VS.addUserToDatabase")
}

// CATEGORIES

// This method test if we can create 2 Categories with the same name, returns true if yes
func (VS *Server) testDuplicateCategories(c *gin.Context) bool {
	category := Category{
		ID:   bson.NewObjectId(),
		Name: "test",
	}
	return VS.addCategoryToDatabase(c, category, "test", "testDuplicateCategories", "return VS.addCategoryToDatabase")
}

// This method test if we can create 2 Items with the same name, returns true if yes
func (VS *Server) testDuplicateItem(c *gin.Context) bool {
	item := Item{
		ID:         bson.NewObjectId(),
		Name:       "test",
		IDCategory: TESTID,
	}
	return VS.addItemToDatabase(c, item, "test", "testDuplicateItem", "return VS.addItemToDatabase")
}

// This method test if we can create 2 Categories Type with the same name, returns true if yes
func (VS *Server) testDuplicateCategoryType(c *gin.Context) bool {
	catType := CategoryType{
		ID:   bson.NewObjectId(),
		Name: "test",
	}
	return VS.addCategoryTypeToDatabase(c, catType, "test", "testDuplicateCategoryType", "return VS.addCategoryTypeToDatabase")
}

// This method test if we can create 2 Video Courses with the same name, returns true if yes
func (VS *Server) testDuplicateVideoCourse(c *gin.Context) bool {
	video := VideoCourse{
		ID:   bson.NewObjectId(),
		Name: "test",
	}
	return VS.addVideoCourseToDatabase(c, video, "test", "testDuplicateVideoCourse", "return VS.addVideoCourseToDatabase")
}
