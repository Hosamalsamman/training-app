package organizations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var service = NewService()

func ListOrganizations(c *gin.Context) {

	orgs, err := service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

func GetOrganization(c *gin.Context) {

	org, err := service.GetByID(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, org)
}
