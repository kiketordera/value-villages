package village

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// This handler changes the time from the developer selection and shows the new time that is set into the app
func (VS *Server) showCurrentTime(periodtime TimeChanges) gin.HandlerFunc {
	return func(c *gin.Context) {
		var t string
		switch periodtime {

		case current:
			FakeDate = getTime()
			y, m, d := time.Unix(getTime(), 0).Date()
			t = time.Unix(getTime(), 0).Weekday().String() + ", " + strconv.Itoa(d) + " of " + m.String() + " of " + strconv.Itoa(y)

		case reset:
			FakeDate = time.Now().Unix()
			TimeBeforeChanging = time.Now().Unix()
			c.Redirect(http.StatusFound, "/development/show-current-time")

		case dayPlus:
			t := time.Unix(FakeDate, 0)
			t = t.AddDate(0, 0, 1)
			FakeDate = t.Unix()
			TimeBeforeChanging = time.Now().Unix()
			c.Redirect(http.StatusFound, "/development/show-current-time")

		case dayMinus:
			t := time.Unix(FakeDate, 0)
			t = t.AddDate(0, 0, -1)
			FakeDate = t.Unix()
			TimeBeforeChanging = time.Now().Unix()
			c.Redirect(http.StatusFound, "/development/show-current-time")

		case weekPlus:
			t := time.Unix(FakeDate, 0)
			t = t.AddDate(0, 0, 7)
			FakeDate = t.Unix()
			TimeBeforeChanging = time.Now().Unix()
			c.Redirect(http.StatusFound, "/development/show-current-time")

		case weekMinus:
			t := time.Unix(FakeDate, 0)
			t = t.AddDate(0, 0, -7)
			FakeDate = t.Unix()
			TimeBeforeChanging = time.Now().Unix()
			c.Redirect(http.StatusFound, "/development/show-current-time")
		}
		render(c, gin.H{
			"time": t,
		}, "edit-time.html")
	}
}

// This method shows all the test options that we have for the development
func (VS *Server) chooseTest(c *gin.Context) {
	render(c, gin.H{
		"role": getRole(c),
	}, "choose-test.html")
}

// This method deletes all the information written by the duplicated tests
func (VS *Server) deleteInformationTestDuplicated(c *gin.Context) {
	VS.deleteInformationTest(c, "h_development", "deleteInformationTestDuplicated", "VS.deleteInformationTest")
	VS.goodFeedback(c, "/development/choose-test")
}

// This method makes all the duplicated test and shows the results
func (VS *Server) makeTestDuplicate(c *gin.Context) {
	render(c, gin.H{
		"dVillage":    VS.testDuplicateVillage(c),
		"dService":    VS.testDuplicateService(c),
		"dUser":       VS.testDuplicateUser(c),
		"dCatType":    VS.testDuplicateCategoryType(c),
		"dCategories": VS.testDuplicateCategories(c),
		"dItem":       VS.testDuplicateItem(c),
		"dVideo":      VS.testDuplicateVideoCourse(c),
		"role":        getRole(c),
	}, "test-duplicate.html")
}
