package organizationtypes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var service = NewService()

func ListOrganizationTypes(c *gin.Context) {
	types, err := service.GetOrganizationTypes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types)
}

func GetOrganizationType(c *gin.Context) {
	t, err := service.GetOrganizationType(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	c.JSON(http.StatusOK, t)
}
