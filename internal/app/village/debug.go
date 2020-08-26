package village

import (
	"github.com/gin-gonic/gin"
)

// This method checks if the operation was wrong, and if wrong, prints the error in the screen, the UI
func (VS *Server) checkOperation(c *gin.Context, class string, method string, operation string, err error) {
	if err != nil {
		render(c, gin.H{
			"Class":     class,
			"Method":    method,
			"Operation": operation,
			"ResultDB":  err,
		}, "ResultDB.html")
	}
}
