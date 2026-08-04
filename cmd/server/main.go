package main

import (
	"training-app/db"
	"training-app/internal/clients"
	"training-app/internal/countries"
	"training-app/internal/courseSessions"
	"training-app/internal/courses"
	"training-app/internal/departments"
	"training-app/internal/documentationType"
	"training-app/internal/documentations"
	"training-app/internal/governorates"
	"training-app/internal/grades"
	"training-app/internal/jobTypeGroups"
	"training-app/internal/jobs"
	"training-app/internal/learningPaths"
	"training-app/internal/learningSubjects"
	"training-app/internal/learningTerms"
	"training-app/internal/organizations"
	"training-app/internal/organizationtypes"
	"training-app/internal/pathGradeSubjectTerms"
	"training-app/internal/pathGradeSubjects"
	"training-app/internal/persons"
	"training-app/internal/qualificationTypes"
	"training-app/internal/qualifications"
	"training-app/internal/trainerSubjects"
	"training-app/internal/trainingRooms"
	"training-app/internal/workCenters"
	"training-app/internal/workGroups"
	"training-app/internal/workSites"
	"training-app/internal/evaluationType"

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
	organizationTypeHandler := organizationtypes.New(db.DB)

	r.GET("/organization-types", organizationTypeHandler.ListOrganizationTypes)
	r.GET("/organization-type/:id", organizationTypeHandler.GetOrganizationType)

	// organizations
	organizationHandler := organizations.New(db.DB)

	r.GET("/organizations", organizationHandler.ListOrganizations)
	r.GET("/organization/:id", organizationHandler.GetOrganization)

	// countries
	countryHandler := countries.New(db.DB)

	r.GET("/countries", countryHandler.ListCountries)
	r.GET("/country/:id", countryHandler.GetCountry)

	// clients
	clientHandler := clients.New(db.DB)

	r.GET("/clients", clientHandler.ListClients)
	r.GET("/client/:id", clientHandler.GetClient)

	// governorates
	govHandler := governorates.New(db.DB)

	r.GET("/governorates", govHandler.ListGovernorates)
	r.GET("/governorate/:id", govHandler.GetGovernorate)

	// gob type groups
	jobTypeGroupHandler := jobTypeGroups.New(db.DB)

	r.GET("/job-type-groups", jobTypeGroupHandler.ListJobTypeGroups)
	r.GET("/job-type-group/:id", jobTypeGroupHandler.GetJobTypeGroup)

	// work centers
	workCenterHandler := workCenters.New(db.DB)

	r.GET("/work-centers", workCenterHandler.ListWorkCenters)
	r.GET("/work-center/:id", workCenterHandler.GetWorkCenter)

	// work sites
	workSiteHandler := workSites.New(db.DB)

	r.GET("/work-sites", workSiteHandler.ListWorkSites)
	r.GET("/work-site/:id", workSiteHandler.GetWorkSite)

	// work groups
	workGroupHandler := workGroups.New(db.DB)

	r.GET("/work-groups", workGroupHandler.ListWorkGroups)
	r.GET("/work-group/:id", workGroupHandler.GetWorkGroup)

	// qualification types
	qualificationTypeHandler := qualificationTypes.New(db.DB)

	r.GET("/qualification-types", qualificationTypeHandler.ListQualificationTypes)
	r.GET("/qualification-type/:id", qualificationTypeHandler.GetQualificationType)

	//qualifications
	qualificationHandler := qualifications.New(db.DB)

	r.GET("/qualifications", qualificationHandler.ListQualifications)
	r.GET("/qualification/:id", qualificationHandler.GetQualification)

	// jobs
	jobHandler := jobs.New(db.DB)

	r.GET("/jobs", jobHandler.ListJobs)
	r.GET("/job/:id", jobHandler.GetJob)

	// grades
	gradeHandler := grades.New(db.DB)

	r.GET("/grades", gradeHandler.ListGrades)
	r.GET("/grade/:id", gradeHandler.GetGrade)

	// learning paths
	learningPathHandler := learningPaths.New(db.DB)

	r.GET("/learning-paths", learningPathHandler.ListLearningPaths)
	r.GET("/learning-path/:id", learningPathHandler.GetLearningPath)

	// learning subjects
	learningSubjectHandler := learningSubjects.New(db.DB)

	r.GET("/learning-subjects", learningSubjectHandler.ListLearningSubjects)
	r.GET("/learning-subject/:id", learningSubjectHandler.GetLearningSubject)

	// path grade subjects
	pathGradeSubjectHandler := pathGradeSubjects.New(db.DB)

	r.GET("/path-grade-subjects", pathGradeSubjectHandler.ListPathGradeSubjects)
	r.GET("/path-grade-subject/:id", pathGradeSubjectHandler.GetPathGradeSubject)

	// path grge subject terms
	pathGradeSubjectTermHandler := pathGradeSubjectTerm.New(db.DB)

	r.GET("/path-grade-subject-terms", pathGradeSubjectTermHandler.ListPathGradeSubjectTerms)
	r.GET("/path-grade-subject-term/:id", pathGradeSubjectTermHandler.GetPathGradeSubjectTerm)

	// persons
	personHandler := persons.New(db.DB)

	r.GET("/persons", personHandler.ListPersons)
	r.GET("/person/:id", personHandler.GetPerson)

	// departments
	departmentHandler := departments.New(db.DB)

	r.GET("/departments", departmentHandler.ListDepartments)
	r.GET("/department/:id", departmentHandler.GetDepartment)

	// trainer subjects
	trainerSubjectHandler := trainerSubjects.New(db.DB)

	r.GET("/trainer-subjects", trainerSubjectHandler.ListTrainerSubjects)
	r.GET("/trainer-subject/:id", trainerSubjectHandler.GetTrainerSubject)

	// learning terms
	learningTermHandler := learningTerms.New(db.DB)

	r.GET("/learning-terms", learningTermHandler.ListLearningTerms)
	r.GET("/learning-term/:id", learningTermHandler.GetLearningTerm)

	// rooms
	roomHandler := trainingRooms.New(db.DB)

	r.GET("/rooms", roomHandler.ListTrainingRooms)
	r.GET("/room/:id", roomHandler.GetTrainingRoom)

	// courses
	courseHandler := courses.New(db.DB)

	r.GET("/courses", courseHandler.ListCourses)
	r.GET("/course/:id", courseHandler.GetCourse)

	// course sessions
	courseSessionHandler := courseSessions.New(db.DB)

	r.GET("/course-sessions", courseSessionHandler.ListCourseSessions)
	r.GET("/course-session/:id", courseSessionHandler.GetCourseSession)

	// documentation types
	documentationTypeHandler := documentationType.New(db.DB)
	
	r.GET("/documentation-types", documentationTypeHandler.ListDocTypes)
	r.GET("/documentation-type/:id", documentationTypeHandler.GetDocType)

	// documentation
	documentationHandler := documentations.New(db.DB)

	r.GET("/documentations", documentationHandler.ListDocumentations)
	r.GET("/documentation/:id", documentationHandler.GetDocumentation)

	// evaluation types
	evaluationTypeHandler := evaluationType.New(db.DB)

	r.GET("/evaluation-types", evaluationTypeHandler.ListEvaluationTypes)
	r.GET("/evaluation-type/:id", evaluationTypeHandler.GetEvaluationType)

	r.Run(":8000")
}
