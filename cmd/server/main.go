package main

import (
	"training-app/db"
	"training-app/internal/clients"
	"training-app/internal/countries"
	"training-app/internal/governorates"
	"training-app/internal/jobTypeGroups"
	"training-app/internal/organizations"
	"training-app/internal/organizationtypes"
	"training-app/internal/workCenters"
	"training-app/internal/workGroups"
	"training-app/internal/workSites"
	"training-app/internal/qualificationTypes"
	"training-app/internal/qualifications"

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

	// clients
	clientRepo := clients.NewRepository(db.DB)
	clientService := clients.NewService(clientRepo)
	clientHandler := clients.NewHandler(clientService)

	r.GET("/clients", clientHandler.ListClients)
	r.GET("/client/:id", clientHandler.GetClient)

	// governorates
	govRepo := governorates.NewRepository(db.DB)
	govService := governorates.NewService(govRepo)
	govHandler := governorates.NewHandler(govService)

	r.GET("/governorates", govHandler.ListGovernorates)
	r.GET("/governorate/:id", govHandler.GetGovernorate)

	// gob type groups
	jobTypeGroupRepo := jobTypeGroups.NewRepository(db.DB)
	jobTypeGroupService := jobTypeGroups.NewService(jobTypeGroupRepo)
	jobTypeGroupHandler := jobTypeGroups.NewHandler(jobTypeGroupService)

	r.GET("/job-type-groups", jobTypeGroupHandler.ListJobTypeGroups)
	r.GET("/job-type-group/:id", jobTypeGroupHandler.GetJobTypeGroup)

	// work centers
	workCenterRepo := workCenters.NewRepository(db.DB)
	workCenterService := workCenters.NewService(workCenterRepo)
	workCenterHandler := workCenters.NewHandler(workCenterService)

	r.GET("/work-centers", workCenterHandler.ListWorkCenters)
	r.GET("/work-center/:id", workCenterHandler.GetWorkCenter)

	// work sites
	workSiteRepo := workSites.NewRepository(db.DB)
	workSiteService := workSites.NewService(workSiteRepo)
	workSiteHandler := workSites.NewHandler(workSiteService)

	r.GET("/work-sites", workSiteHandler.ListWorkSites)
	r.GET("/work-site/:id", workSiteHandler.GetWorkSite)

	// work groups
	workGroupRepo := workGroups.NewRepository(db.DB)
	workGroupService := workGroups.NewService(workGroupRepo)
	workGroupHandler := workGroups.NewHandler(workGroupService)

	// qualification types
	qualificationTypeRepo := qualificationTypes.NewRepository(db.DB)
	qualificationTypeService := qualificationTypes.NewService(qualificationTypeRepo)
	qualificationTypeHandler := qualificationTypes.NewHandler(qualificationTypeService)

	r.GET("/qualification-types", qualificationTypeHandler.ListQualificationTypes)
	r.GET("/qualification-type/:id", qualificationTypeHandler.GetQualificationType)

	//qualifications
	qualificationRepo := qualifications.NewRepository(db.DB)
	qualificationService := qualifications.NewService(qualificationRepo)
	qualificationHandler := qualifications.NewHandler(qualificationService)

	r.GET("/qualifications", qualificationHandler.ListQualifications)
	r.GET("/qualification/:id", qualificationHandler.GetQualification)

	r.GET("/work-groups", workGroupHandler.ListWorkGroups)
	r.GET("/work-group/:id", workGroupHandler.GetWorkGroup)

	r.Run(":8000")
}
