package main

import (
	"training-app/db"
	"training-app/internal/auth"
	"training-app/internal/clients"
	"training-app/internal/countries"
	"training-app/internal/courseParticipantFinalExam"
	"training-app/internal/courseParticipants"
	"training-app/internal/courseSessionParticipants"
	"training-app/internal/courseSessions"
	"training-app/internal/courses"
	"training-app/internal/departments"
	"training-app/internal/documentationType"
	"training-app/internal/documentations"
	"training-app/internal/evaluationType"
	"training-app/internal/governorates"
	"training-app/internal/grades"
	"training-app/internal/groups"
	"training-app/internal/jobTypeGroups"
	"training-app/internal/jobs"
	"training-app/internal/learningPaths"
	"training-app/internal/learningSubjects"
	"training-app/internal/learningTerms"
	"training-app/internal/organizations"
	"training-app/internal/organizationtypes"
	"training-app/internal/participantType"
	pathGradeSubjectTerm "training-app/internal/pathGradeSubjectTerms"
	"training-app/internal/pathGradeSubjects"
	"training-app/internal/pathGrades"
	"training-app/internal/performanceEvaluationCategories"
	"training-app/internal/performanceEvaluationCourseDetails"
	"training-app/internal/performanceEvaluationCourseParticipantDetails"
	"training-app/internal/performanceEvaluationItems"
	"training-app/internal/performanceEvaluationParticipantCategory"
	"training-app/internal/performanceEvaluationParticipantCategoryItems"
	"training-app/internal/persons"
	"training-app/internal/qualificationTypes"
	"training-app/internal/qualifications"
	"training-app/internal/trainerSubjects"
	"training-app/internal/trainingRooms"
	"training-app/internal/workCenters"
	"training-app/internal/workGroups"
	"training-app/internal/workSites"

	"github.com/gin-gonic/gin"

	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}
	db.Connect()

	// Secret used to sign and verify JWTs.
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	r := gin.Default()

	// r.GET("/", handlers.Home)
	// r.GET("/db-check", handlers.DbCheck)

	// routes that require a valid JWT.
	// The middleware verifies the token and stores user_id and
	// client_id in the request context, so every handler below
	// can read them with c.GetInt("client_id").
	// /login is the only public route: it issues the token.
	authed := r.Group("/")
	authed.Use(auth.JWTMiddleware(jwtSecret))

	// organization types
	organizationTypeHandler := organizationtypes.New(db.DB)

	authed.GET("/organization-types", organizationTypeHandler.ListOrganizationTypes)
	authed.GET("/organization-type/:id", organizationTypeHandler.GetOrganizationType)

	// organizations
	organizationHandler := organizations.New(db.DB)

	authed.GET("/organizations", organizationHandler.ListOrganizations)
	authed.GET("/organization/:id", organizationHandler.GetOrganization)

	// countries
	countryHandler := countries.New(db.DB)

	authed.GET("/countries", countryHandler.ListCountries)
	authed.GET("/country/:id", countryHandler.GetCountry)

	// clients
	clientHandler := clients.New(db.DB)

	authed.GET("/clients", clientHandler.ListClients)
	authed.GET("/client/:id", clientHandler.GetClient)

	// governorates
	govHandler := governorates.New(db.DB)

	authed.GET("/governorates", govHandler.ListGovernorates)
	authed.GET("/governorate/:id", govHandler.GetGovernorate)

	// gob type groups
	jobTypeGroupHandler := jobTypeGroups.New(db.DB)

	authed.GET("/job-type-groups", jobTypeGroupHandler.ListJobTypeGroups)
	authed.GET("/job-type-group/:id", jobTypeGroupHandler.GetJobTypeGroup)

	// work centers
	workCenterHandler := workCenters.New(db.DB)

	authed.GET("/work-centers", workCenterHandler.ListWorkCenters)
	authed.GET("/work-center/:id", workCenterHandler.GetWorkCenter)

	// work sites
	workSiteHandler := workSites.New(db.DB)

	authed.GET("/work-sites", workSiteHandler.ListWorkSites)
	authed.GET("/work-site/:id", workSiteHandler.GetWorkSite)

	// work groups
	workGroupHandler := workGroups.New(db.DB)

	authed.GET("/work-groups", workGroupHandler.ListWorkGroups)
	authed.GET("/work-group/:id", workGroupHandler.GetWorkGroup)

	// qualification types
	qualificationTypeHandler := qualificationTypes.New(db.DB)

	authed.GET("/qualification-types", qualificationTypeHandler.ListQualificationTypes)
	authed.GET("/qualification-type/:id", qualificationTypeHandler.GetQualificationType)

	//qualifications
	qualificationHandler := qualifications.New(db.DB)

	authed.GET("/qualifications", qualificationHandler.ListQualifications)
	authed.GET("/qualification/:id", qualificationHandler.GetQualification)

	// jobs
	jobHandler := jobs.New(db.DB)

	authed.GET("/jobs", jobHandler.ListJobs)
	authed.GET("/job/:id", jobHandler.GetJob)

	// grades
	gradeHandler := grades.New(db.DB)

	authed.GET("/grades", gradeHandler.ListGrades)
	authed.GET("/grade/:id", gradeHandler.GetGrade)

	// learning paths
	learningPathHandler := learningPaths.New(db.DB)

	authed.GET("/learning-paths", learningPathHandler.ListLearningPaths)
	authed.GET("/learning-path/:id", learningPathHandler.GetLearningPath)

	// learning subjects
	learningSubjectHandler := learningSubjects.New(db.DB)

	authed.GET("/learning-subjects", learningSubjectHandler.ListLearningSubjects)
	authed.GET("/learning-subject/:id", learningSubjectHandler.GetLearningSubject)

	// path grade subjects
	pathGradeSubjectHandler := pathGradeSubjects.New(db.DB)

	authed.GET("/path-grade-subjects", pathGradeSubjectHandler.ListPathGradeSubjects)
	authed.GET("/path-grade-subject/:id", pathGradeSubjectHandler.GetPathGradeSubject)

	// path grge subject terms
	pathGradeSubjectTermHandler := pathGradeSubjectTerm.New(db.DB)

	authed.GET("/path-grade-subject-terms", pathGradeSubjectTermHandler.ListPathGradeSubjectTerms)
	authed.GET("/path-grade-subject-term/:id", pathGradeSubjectTermHandler.GetPathGradeSubjectTerm)

	// persons
	personHandler := persons.New(db.DB, jwtSecret)

	authed.GET("/persons", personHandler.ListPersons)
	authed.GET("/person/:id", personHandler.GetPerson)
	authed.POST("/new-person", personHandler.Create)

	// groups
	groupHandler := groups.New(db.DB)

	authed.GET("/groups", groupHandler.ListGroups)
	authed.GET("/group/:id", groupHandler.GetGroup)

	// authentication
	// login is public: it issues the token required by
	// every authenticated route below.
	r.POST("/login", personHandler.Login)

	{
		// logout: the token itself is stateless, so the
		// client discards it after calling this route.
		authed.POST("/logout", personHandler.Logout)

		// register user: turn an existing person into a user
		authed.POST("/register-user", personHandler.RegisterUser)

		// groups the authenticated user is allowed to assign users to
		authed.GET("/allowed-groups", groupHandler.ListAllowedGroups)

		// change password: any authenticated user, the old
		// password must match first.
		authed.POST("/change-password", personHandler.ChangePassword)

		// reset password: admins (group 1) set a new password
		// for another user without the old one.
		authed.POST("/person/:id/reset-password", personHandler.ResetPassword)

		// users: the client's registered persons, so the
		// admin picks a real user as a reset target.
		authed.GET("/users", personHandler.ListUsers)
	}

	// departments
	departmentHandler := departments.New(db.DB)

	authed.GET("/departments", departmentHandler.ListDepartments)
	authed.GET("/department/:id", departmentHandler.GetDepartment)

	// trainer subjects
	trainerSubjectHandler := trainerSubjects.New(db.DB)

	authed.GET("/trainer-subjects", trainerSubjectHandler.ListTrainerSubjects)
	authed.GET("/trainer-subject/:id", trainerSubjectHandler.GetTrainerSubject)

	// learning terms
	learningTermHandler := learningTerms.New(db.DB)

	authed.GET("/learning-terms", learningTermHandler.ListLearningTerms)
	authed.GET("/learning-term/:id", learningTermHandler.GetLearningTerm)

	// rooms
	roomHandler := trainingRooms.New(db.DB)

	authed.GET("/rooms", roomHandler.ListTrainingRooms)
	authed.GET("/room/:id", roomHandler.GetTrainingRoom)

	// courses
	courseHandler := courses.New(db.DB)

	authed.GET("/courses", courseHandler.ListCourses)
	authed.GET("/course/:id", courseHandler.GetCourse)

	// course sessions
	courseSessionHandler := courseSessions.New(db.DB)

	authed.GET("/course-sessions", courseSessionHandler.ListCourseSessions)
	authed.GET("/course-session/:id", courseSessionHandler.GetCourseSession)

	// documentation types
	documentationTypeHandler := documentationType.New(db.DB)

	authed.GET("/documentation-types", documentationTypeHandler.ListDocTypes)
	authed.GET("/documentation-type/:id", documentationTypeHandler.GetDocType)

	// documentation
	documentationHandler := documentations.New(db.DB)

	authed.GET("/documentations", documentationHandler.ListDocumentations)
	authed.GET("/documentation/:id", documentationHandler.GetDocumentation)

	// evaluation types
	evaluationTypeHandler := evaluationType.New(db.DB)

	authed.GET("/evaluation-types", evaluationTypeHandler.ListEvaluationTypes)
	authed.GET("/evaluation-type/:id", evaluationTypeHandler.GetEvaluationType)

	// course participants
	courseParticipantHandler := courseParticipants.New(db.DB)

	authed.GET("/course-participants", courseParticipantHandler.ListCourseParticipants)
	authed.GET("/course-participant/:id", courseParticipantHandler.GetCourseParticipant)

	// course session participants
	courseSessionParticipantHandler := courseSessionParticipants.New(db.DB)

	authed.GET("/course-session-participants", courseSessionParticipantHandler.ListCourseSessionParticipants)
	authed.GET("/course-session-participant/:id", courseSessionParticipantHandler.GetCourseSessionParticipant)

	// course participant final exams
	courseParticipantFinalExamHandler := courseParticipantFinalExam.New(db.DB)

	authed.GET("/course-participant-final-exams", courseParticipantFinalExamHandler.ListCourseParticipantsFinalExams)
	authed.GET("/course-participant-final-exam/:id", courseParticipantFinalExamHandler.GetCourseParticipantFinalExam)

	// performance evaluation categories
	performanceEvaluationCategoryHandler := performanceEvaluationCategories.New(db.DB)

	authed.GET("/performance-evaluation-categories", performanceEvaluationCategoryHandler.ListPerformanceEvaluationCategories)
	authed.GET("/performance-evaluation-category/:id", performanceEvaluationCategoryHandler.GetPerformanceEvaluationCategory)

	// performance evaluation items
	performanceEvaluationItemHandler := performanceEvaluationItems.New(db.DB)

	authed.GET("/performance-evaluation-items", performanceEvaluationItemHandler.ListPerformanceEvaluationItems)
	authed.GET("/performance-evaluation-item/:id", performanceEvaluationItemHandler.GetPerformanceEvaluationItem)

	// participant types
	participantTypeHandler := participantType.New(db.DB)

	authed.GET("/participant-types", participantTypeHandler.ListParticipantTypes)
	authed.GET("/participant-type/:id", participantTypeHandler.GetParticipantType)

	// path grades
	pathGradeHandler := pathGrades.New(db.DB)

	authed.GET("/path-grades", pathGradeHandler.ListPathGrades)
	authed.GET("/path-grade/:id", pathGradeHandler.GetPathGrade)

	// performance evaluation participant categories
	performanceEvaluationParticipantCategoryHandler := performanceEvaluationParticipantCategory.New(db.DB)

	authed.GET("/performance-evaluation-participant-categories", performanceEvaluationParticipantCategoryHandler.ListPerformanceEvaluationParticipantCategories)
	authed.GET("/performance-evaluation-participant-category/:id", performanceEvaluationParticipantCategoryHandler.GetPerformanceEvaluationParticipantCategory)

	// performance evaluation participant category items
	performanceEvaluationParticipantCategoryItemHandler := performanceEvaluationParticipantCategoryItems.New(db.DB)

	authed.GET("/performance-evaluation-participant-category-items", performanceEvaluationParticipantCategoryItemHandler.ListPerformanceEvaluationParticipantCategoryItems)
	authed.GET("/performance-evaluation-participant-category-item/:id", performanceEvaluationParticipantCategoryItemHandler.GetPerformanceEvaluationParticipantCategoryItem)

	// performance evaluation course details
	performanceEvaluationCourseDetailHandler := performanceEvaluationCourseDetails.New(db.DB)

	authed.GET("/performance-evaluation-course-details", performanceEvaluationCourseDetailHandler.ListPerformanceEvaluationCourseDetails)
	authed.GET("/performance-evaluation-course-detail/:id", performanceEvaluationCourseDetailHandler.GetPerformanceEvaluationCourseDetail)

	// performance evaluation course participant details
	performanceEvaluationCourseParticipantDetailHandler := performanceEvaluationCourseParticipantDetails.New(db.DB)

	authed.GET("/performance-evaluation-course-participant-details", performanceEvaluationCourseParticipantDetailHandler.ListPerformanceEvaluationCourseParticipantDetails)
	authed.GET("/performance-evaluation-course-participant-detail/:id", performanceEvaluationCourseParticipantDetailHandler.GetPerformanceEvaluationCourseParticipantDetail)

	r.Run(":8000")
}
