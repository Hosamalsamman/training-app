package main

import (
	"training-app/db"
	"training-app/internal/countries"
	"training-app/internal/organizations"
	"training-app/internal/organizationtypes"

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
	// func TenantMiddleware(c *gin.Context) {
	// 	clientID := getClientIDFromJWT(c)
	
	// 	scopedDB := db.DB.Where("client_id = ?", clientID)
	
	// 	c.Set("db", scopedDB)
	
	// 	c.Next()
	// }

	// r.GET("/", handlers.Home)
	// r.GET("/db-check", handlers.DbCheck)

	// organization types
	organizationTypeRepo := organizationtypes.NewRepository(db.DB)
	organizationTypeService := organizationtypes.NewService(organizationTypeRepo)
	organizationTypeHandler := organizationtypes.NewHandler(organizationTypeService)

	r.GET("/organization-types", organizationTypeHandler.ListOrganizationTypes)
	r.GET("/organization-type/:id", organizationTypeHandler.GetOrganizationType)

	// organizations
	organizationRepo := organizations.NewRepository(db.DB)
	organizationService := organizations.NewService(organizationRepo)
	organizationHandler := organizations.NewHandler(organizationService)

	r.GET("/organizations", organizationHandler.ListOrganizations)
	r.GET("/organization/:id", organizationHandler.GetOrganization)

	// countries
	countryRepo := countries.NewRepository(db.DB)
	countryService := countries.NewService(countryRepo)
	countryHandler := countries.NewHandler(countryService)

	r.GET("/countries", countryHandler.ListCountries)
	r.GET("/country/:id", countryHandler.GetCountry)

	r.Run(":8000")
}
