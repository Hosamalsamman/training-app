package countries

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var service = NewService()

func ListCountries(c *gin.Context) {

	orgs, err := service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

func GetCountry(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	org, err := service.GetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, org)
}
