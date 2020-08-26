package village

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Shows the animals to slaughter
func (VS *Server) showAnimals(c *gin.Context) {
	render(c, gin.H{}, "animals.html")
}

// Shows the animals to slaughter
func (VS *Server) showPartsAnimals(c *gin.Context) {
	a, b := VS.getParamFromURL(c, "animal")
	if !b {
		return
	}
	animalType := AnimalType(a)
	render(c, gin.H{
		"animal": animalType,
	}, "parts-animal.html")
}

// Shows the animals to slaughter
func (VS *Server) newAnimal(c *gin.Context) {
	service, b := VS.getOurService(c, "h_slaughterhouse", "newAnimal", "service, b := VS.getOurService")
	if !b {
		return
	}
	a, b := VS.getParamFromURL(c, "animal")
	if !b {
		return
	}
	animalType := AnimalType(a)
	id := bson.NewObjectId()
	front, b := VS.getPhotoFromHTML(c, "FRONT", "Animal", "/local-resources/slaughterhouse/animals/"+id.Hex()+"/")
	if !b {
		return
	}
	back, b := VS.getPhotoFromHTML(c, "BACK", "Animal", "/local-resources/slaughterhouse/animals/"+id.Hex()+"/")
	if !b {
		return
	}
	right, b := VS.getPhotoFromHTML(c, "RIGHT", "Animal", "/local-resources/slaughterhouse/animals/"+id.Hex()+"/")
	if !b {
		return
	}
	left, b := VS.getPhotoFromHTML(c, "LEFT", "Animal", "/local-resources/slaughterhouse/animals/"+id.Hex()+"/")
	if !b {
		return
	}

	animal := Animal{
		ID:         id,
		Front:      front,
		Back:       back,
		Right:      right,
		Left:       left,
		IsClosed:   false,
		AnimalType: animalType,
	}
	VS.addElementToDatabaseWithRegister(c, animal, animal.ID, service.ID, "h_slaughterhouse", "newAnimal", "VS.addAnimalToDatabase")
	VS.goodFeedback(c, "slaughterhouse/animals")
}
