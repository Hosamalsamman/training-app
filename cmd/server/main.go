package main

import (
	"training-app/db"
	"training-app/internal/clients"
	"training-app/internal/countries"
	"training-app/internal/governorates"
	"training-app/internal/jobTypeGroups"
	"training-app/internal/organizations"
	"training-app/internal/organizationtypes"
	"training-app/internal/qualificationTypes"
	"training-app/internal/qualifications"
	"training-app/internal/workCenters"
	"training-app/internal/workGroups"
	"training-app/internal/workSites"
	"training-app/internal/jobs"
	"training-app/internal/grades"
	"training-app/internal/learningPaths"
	"training-app/internal/learningSubjects"
	"training-app/internal/pathGradeSubjects"
	"training-app/internal/pathGradeSubjectTerms"
	"training-app/internal/persons"
	"training-app/internal/departments"
	"training-app/internal/trainerSubjects"
	

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
	countryHandler := countries.New(db.DB)

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

	r.GET("/work-groups", workGroupHandler.ListWorkGroups)
	r.GET("/work-group/:id", workGroupHandler.GetWorkGroup)

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

	// jobs
	jobRepo := jobs.NewRepository(db.DB)
	jobService := jobs.NewService(jobRepo)
	jobHandler := jobs.NewHandler(jobService)

	r.GET("/jobs", jobHandler.ListJobs)
	r.GET("/job/:id", jobHandler.GetJob)

	// grades
	gradeRepo := grades.NewRepository(db.DB)
	gradeService := grades.NewService(gradeRepo)
	gradeHandler := grades.NewHandler(gradeService)

	r.GET("/grades", gradeHandler.ListGrades)
	r.GET("/grade/:id", gradeHandler.GetGrade)

	// learning paths
	learningPathRepo := learningPaths.NewRepository(db.DB)
	learningPathService := learningPaths.NewService(learningPathRepo)
	learningPathHandler := learningPaths.NewHandler(learningPathService)

	r.GET("/learning-paths", learningPathHandler.ListLearningPaths)
	r.GET("/learning-path/:id", learningPathHandler.GetLearningPath)

	// learning subjects
	learningSubjectRepo := learningSubjects.NewRepository(db.DB)
	learningSubjectService := learningSubjects.NewService(learningSubjectRepo)
	learningSubjectHandler := learningSubjects.NewHandler(learningSubjectService)

	r.GET("/learning-subjects", learningSubjectHandler.ListLearningSubjects)
	r.GET("/learning-subject/:id", learningSubjectHandler.GetLearningSubject)

	// path grade subjects
	pathGradeSubjectRepo := pathGradeSubjects.NewRepository(db.DB)
	pathGradeSubjectService := pathGradeSubjects.NewService(pathGradeSubjectRepo)
	pathGradeSubjectHandler := pathGradeSubjects.NewHandler(pathGradeSubjectService)

	r.GET("/path-grade-subjects", pathGradeSubjectHandler.ListPathGradeSubjects)
	r.GET("/path-grade-subject/:id", pathGradeSubjectHandler.GetPathGradeSubject)

	// path grge subject terms
	pathGradeSubjectTermRepo := pathGradeSubjectTerm.NewRepository(db.DB)
	pathGradeSubjectTermService := pathGradeSubjectTerm.NewService(pathGradeSubjectTermRepo)
	pathGradeSubjectTermHandler := pathGradeSubjectTerm.NewHandler(pathGradeSubjectTermService)

	r.GET("/path-grade-subject-terms", pathGradeSubjectTermHandler.ListPathGradeSubjectTerms)
	r.GET("/path-grade-subject-term/:id", pathGradeSubjectTermHandler.GetPathGradeSubjectTerm)

	// persons
	personRepo := persons.NewRepository(db.DB)
	personService := persons.NewService(personRepo)
	personHandler := persons.NewHandler(personService)

	r.GET("/persons", personHandler.ListPersons)
	r.GET("/person/:id", personHandler.GetPerson)

	// departments
	departmentRepo := departments.NewRepository(db.DB)
	departmentService := departments.NewService(departmentRepo)
	departmentHandler := departments.NewHandler(departmentService)

	r.GET("/departments", departmentHandler.ListDepartments)
	r.GET("/department/:id", departmentHandler.GetDepartment)

	// trainer subjects
	trainerSubjectRepo := trainerSubjects.NewRepository(db.DB)
	trainerSubjectService := trainerSubjects.NewService(trainerSubjectRepo)
	trainerSubjectHandler := trainerSubjects.NewHandler(trainerSubjectService)

	r.GET("/trainer-subjects", trainerSubjectHandler.ListTrainerSubjects)
	r.GET("/trainer-subject/:id", trainerSubjectHandler.GetTrainerSubject)

	r.Run(":8000")
}
