package main

import (
	"training-app/db"
	country "training-app/internal/countries"
	organization "training-app/internal/organizations"
	organizationtype "training-app/internal/organizationtypes"

	"github.com/gin-gonic/gin"

	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}
	db.Connect()

	r := gin.Default()

	// r.GET("/", handlers.Home)
	// r.GET("/db-check", handlers.DbCheck)

	// organization types
	r.GET("/organization-types", organizationtype.ListOrganizationTypes)
	r.GET("/organization-type/:id", organizationtype.GetOrganizationType)

	// organizations
	r.GET("/organizations", organization.ListOrganizations)
	r.GET("/organization/:id", organization.GetOrganization)

	// countries
	r.GET("/countries", country.ListCountries)
	r.GET("/country/:id", country.GetCountry)

	r.Run(":8000")
}
