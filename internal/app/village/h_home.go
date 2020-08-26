package village

import (
	"github.com/gin-gonic/gin"
)

// Shows the about page of the app
func (VS *Server) about(c *gin.Context) {
	render(c, gin.H{}, "presentation-about.html")
}

// Shows the design page of the app
func (VS *Server) design(c *gin.Context) {
	render(c, gin.H{}, "presentation-design.html")
}

// Shows the features page of the app
func (VS *Server) features(c *gin.Context) {
	render(c, gin.H{}, "presentation-features.html")
}

// Shows the manuals page of the app
func (VS *Server) manuals(c *gin.Context) {
	render(c, gin.H{}, "presentation-manuals.html")
}
