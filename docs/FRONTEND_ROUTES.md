# Frontend API Documentation

Every route exposed by the backend (`cmd/server/main.go`, port `8000`), the model it
returns, and — for POST routes — which GET route to call first so the user can pick
values from real data instead of typing raw IDs.

## Authentication

All routes except `POST /login` require a JWT. Send it as:

```
Authorization: Bearer <token>
```

The token is issued by `/login` and contains `user_id` and `client_id`. Every
authenticated response is already scoped to the caller's `client_id` — the frontend
never needs to send `client_id` explicitly.

Flask parallel: this works like `@jwt_required` from Flask-JWT-Extended, but instead
of decorators on each view, one middleware (`auth.JWTMiddleware`) is attached to the
whole route group; handlers read the identity from `gin.Context` the way Flask views
read `g.user`.

---

## Route Catalog

### Persons / Authentication

| Method | Route | Model / Payload | Notes |
|---|---|---|---|
| POST | `/login` (public) | `LoginRequest` `{username, password}` | Returns JWT. |
| GET | `/persons` | `models.Person` | List persons. |
| GET | `/person/:id` | `models.Person` | One person with all associations. |
| POST | `/new-person` | `CreatePersonRequest` | See picker sources below. |
| POST | `/logout` | — | Client just discards the token. |
| POST | `/register-user` | `RegisterUserRequest` `{person_id, username, password, group_id}` | Turn an existing person into a user. |
| POST | `/change-password` | `ChangePasswordRequest` `{old_password, new_password}` | Current user. |
| POST | `/person/:id/reset-password` | `ResetPasswordRequest` `{new_password}` | Admin only (group 1). |
| GET | `/users` | `UserSummary` | Registered users — picker for reset-password target. |
| GET | `/allowed-groups` | `models.Group` | Groups current admin may assign — picker for `group_id` in `/register-user`. |

**POST `/new-person`** — fetch these GET routes to render dropdowns:

| Payload field | Render from | Route |
|---|---|---|
| `governorate_id` | `Governorate.name` | `GET /governorates` |
| `qualification_id` | `Qualification.name` | `GET /qualifications` |
| `job_id` | `Job.name` | `GET /jobs` |
| `organization_id` | `Organization.name` | `GET /organizations` |
| `department_id` | `Department.name` | `GET /departments` |
| `work_site_id` | `WorkSite.name` | `GET /work-sites` |
| `job_type_group_id` | `JobTypeGroup.name` | `GET /job-type-groups` |
| `current_learning_path_id` | `LearningPath.name` | `GET /learning-paths` |
| `current_grade_id` | `Grade.name` | `GET /grades` |
| `trainer_certifying_organization_id` | `Organization.name` | `GET /organizations` |

**POST `/register-user`** — fetch these to render dropdowns:

| Payload field | Render from | Route |
|---|---|---|
| `person_id` | `Person.name` (persons without credentials) | `GET /persons` |
| `group_id` | `Group.name` | `GET /allowed-groups` |

---

### Reference / Lookup Tables (read-only pickers)

These exist mainly so forms can render dropdowns; each list route doubles as the
picker source for any POST payload field referencing it.

| Method | Route | Model |
|---|---|---|
| GET | `/organization-types`, `/organization-type/:id` | `models.OrganizationType` |
| GET | `/organizations`, `/organization/:id` | `models.Organization` |
| GET | `/countries`, `/country/:id` | `models.Country` |
| GET | `/clients`, `/client/:id` | `models.Client` |
| GET | `/governorates`, `/governorate/:id` | `models.Governorate` |
| GET | `/job-type-groups`, `/job-type-group/:id` | `models.JobTypeGroup` |
| GET | `/work-centers`, `/work-center/:id` | `models.WorkCenter` |
| GET | `/work-sites`, `/work-site/:id` | `models.WorkSite` |
| GET | `/work-groups`, `/work-group/:id` | `models.WorkGroup` |
| GET | `/qualification-types`, `/qualification-type/:id` | `models.QualificationType` |
| GET | `/qualifications`, `/qualification/:id` | `models.Qualification` |
| GET | `/jobs`, `/job/:id` | `models.Job` |
| GET | `/grades`, `/grade/:id` | `models.Grade` |
| GET | `/learning-paths`, `/learning-path/:id` | `models.LearningPath` |
| GET | `/learning-subjects`, `/learning-subject/:id` | `models.LearningSubject` |
| GET | `/learning-terms`, `/learning-term/:id` | `models.LearningTerm` |
| GET | `/path-grades`, `/path-grade/:id` | `models.PathGrade` |
| GET | `/path-grade-subjects`, `/path-grade-subject/:id` | `models.PathGradeSubject` |
| GET | `/path-grade-subject-terms`, `/path-grade-subject-term/:id` | `models.PathGradeSubjectTerm` |
| GET | `/departments`, `/department/:id` | `models.Department` |
| GET | `/trainer-subjects`, `/trainer-subject/:id` | `models.TrainerSubject` |
| GET | `/rooms`, `/room/:id` | `models.TrainingRoom` |
| GET | `/documentation-types`, `/documentation-type/:id` | `models.DocumentationType` |
| GET | `/evaluation-types`, `/evaluation-type/:id` | `models.EvaluationType` |
| GET | `/participant-types`, `/participant-type/:id` | `models.ParticipantType` |
| GET | `/groups`, `/group/:id` | `models.Group` |

---

### Path Grade Candidates (full CRUD)

| Method | Route | Model / Payload |
|---|---|---|
| GET | `/path-grade-candidates` | `models.PathGradeCandidate` |
| GET | `/path-grade-candidate/:id` | `models.PathGradeCandidate` |
| POST | `/new-path-grade-candidate` | `models.PathGradeCandidate` |
| PUT | `/path-grade-candidate/:id` | `models.PathGradeCandidate` |
| DELETE | `/path-grade-candidate/:id` | — |

**POST/PUT `/new-path-grade-candidate`** — dropdown sources:

| Field | Render from | Route |
|---|---|---|
| `person_id` | `Person.name` | `GET /persons` |
| `path_grade_id` | `PathGrade` (via its `learning_path` + `grade` names) | `GET /path-grades` |
| `term_id` | `LearningTerm.name` | `GET /learning-terms` |
| `final_evaluation_id` | `EvaluationType.name` | `GET /evaluation-types` |

---

### Courses

| Method | Route | Model / Payload |
|---|---|---|
| GET | `/courses` | `models.Course` |
| GET | `/course/:id` | `models.Course` |
| GET | `/get-planned-courses` | `models.Course` — optional query params `path_grade_subject_id`, `path_grade_subject_term_id`, `learning_subject_id` |
| GET | `/get-executed-courses` | `models.Course` — same optional query params |
| POST | `/new-course` | `models.CreateCourseRequest` |

**POST `/new-course`** — dropdown sources:

| Field | Render from | Route |
|---|---|---|
| `path_grade_subject_id` | `PathGradeSubject.name` | `GET /path-grade-subjects` |
| `path_grade_subject_term_id` | `PathGradeSubjectTerm` | `GET /path-grade-subject-terms` |
| `learning_subject_id` | `LearningSubject.name` | `GET /learning-subjects` |
| `room_id` | `TrainingRoom.name` | `GET /rooms` |
| `trainer_id` / `backup_trainer_id` | `Person.name` (filter `is_trainer`) | `GET /persons` |
| `coordinator_id` | `Person.name` | `GET /persons` |
| `funding_organization_id` | `Organization.name` | `GET /organizations` |
| `planned_id` | planned courses | `GET /get-planned-courses` |

Business rules enforced server-side (mirror them in the form UI):

- Exactly **one** of the three subject IDs must be set.
- Lifecycle states: planned (`is_planned` only), executed (`is_executed` only), or
  executed-from-planned (both flags + `planned_id` required).

---

### Course Sessions (full CRUD)

| Method | Route | Model / Payload |
|---|---|---|
| GET | `/course-sessions` | `models.CourseSession` |
| GET | `/course-session/:id` | `models.CourseSession` |
| POST | `/new-course-session` | `models.CourseSessionRequest` |
| PUT | `/course-session/:id` | `models.CourseSessionRequest` (full update) |
| DELETE | `/course-session/:id` | — |

**POST/PUT** — dropdown sources:

| Field | Render from | Route |
|---|---|---|
| `course_id` | `Course.name` (executed courses only) | `GET /get-executed-courses` |

---

### Read-only Domain Sections

| Method | List route | Single route | Model |
|---|---|---|---|
| GET | `/documentations` | `/documentation/:id` | `models.Documentation` |
| GET | `/documentation-types` | `/documentation-type/:id` | `models.DocumentationType` |
| GET | `/course-participants` | `/course-participant/:id` | `models.CourseParticipant` |
| GET | `/course-session-participants` | `/course-session-participant/:id` | `models.CourseSessionParticipant` |
| GET | `/course-participant-final-exams` | `/course-participant-final-exam/:id` | `models.CourseParticipantFinalExam` |
| GET | `/performance-evaluation-categories` | `/performance-evaluation-category/:id` | `models.PerformanceEvaluationCategory` |
| GET | `/performance-evaluation-items` | `/performance-evaluation-item/:id` | `models.PerformanceEvaluationItem` |
| GET | `/performance-evaluation-participant-categories` | `/performance-evaluation-participant-category/:id` | `models.PerformanceEvaluationParticipantCategory` |
| GET | `/performance-evaluation-participant-category-items` | `/performance-evaluation-participant-category-item/:id` | `models.PerformanceEvaluationParticipantCategoryItem` |
| GET | `/performance-evaluation-course-details` | `/performance-evaluation-course-detail/:id` | `models.PerformanceEvaluationCourseDetail` |
| GET | `/performance-evaluation-course-participant-details` | `/performance-evaluation-course-participant-detail/:id` | `models.PerformanceEvaluationCourseParticipantDetail` |

---

## Model Field Reference (JSON keys)

JSON keys come straight from the `json` struct tags — the snake_case names below are
exactly what the frontend receives.

- **Client** — `id`, `name`
- **Country** — `id`, `name`, `client_id`, `governorates[]`
- **Governorate** — `id`, `name`, `country_id`, `client_id`
- **OrganizationType** — `id`, `name`, `client_id`
- **Organization** — `id`, `name`, `code`, `address`, `organization_type_id`, `governorate_id`, `client_id`, `work_centers[]`, `departments[]`
- **WorkGroup** — `id`, `code`, `name`, `client_id`
- **JobTypeGroup** — `id`, `code`, `name`, `work_group_id`, `client_id`
- **WorkCenter** — `id`, `name`, `governorate_id`, `organization_id`, `client_id`
- **WorkSite** — `id`, `code`, `name`, `work_center_id`, `client_id`
- **Department** — `id`, `code`, `name`, `organization_id`, `client_id`
- **QualificationType** — `id`, `name`, `client_id`
- **Qualification** — `id`, `code`, `name`, `qualification_type_id`, `client_id`
- **Job** — `id`, `code`, `name`, `client_id`
- **Grade** — `id`, `code`, `name`, `grade_interval`, `client_id`
- **LearningPath** — `id`, `code`, `name`, `qualification_type_id`, `client_id`
- **LearningTerm** — `id`, `name`, `client_id`
- **LearningSubject** — `id`, `name`, `client_id`
- **PathGrade** — `id`, `learning_path_id`, `grade_id`, `client_id`, `path_grade_subjects[]`
- **PathGradeSubject** — `id`, `name`, `code`, `path_grade_id`, `learning_subject_id`, `client_id`
- **PathGradeSubjectTerm** — `id`, `path_grade_subjects_id`, `term_id`, `code`, `notes`, `client_id`
- **EvaluationType** — `id`, `name`, `client_id`
- **ParticipantType** — `id`, `name`, `client_id`
- **Group** — `id`, `name`
- **TrainingRoom** — `id`, `name`, `work_site_id`, `client_id`
- **Person** — `id`, `code`, `name`, `governorate_id`, `address`, `telephone_whatsapp`, `qualification_id`, `job_id`, `organization_id`, `department_id`, `work_site_id`, `job_type_group_id`, `current_learning_path_id`, `current_grade_id`, `date_of_current_grade`, `contract_date`, `is_active`, `is_trainer`, `trainer_certifying_organization_id`, `user_name`, `group_id`, `client_id`
- **UserSummary** — `id`, `code`, `name`, `username`, `group_id`, `is_active`
- **PathGradeCandidate** — `id`, `person_id`, `path_grade_id`, `term_id`, `final_evaluation_id`, `client_id`
- **Course** — `id`, `name`, `path_grade_subject_id` / `path_grade_subject_term_id` / `learning_subject_id` (exactly one), `duration_in_days`, `starting_date`, `end_date`, `number_of_internal_participants`, `number_of_external_participants`, `room_id`, `funding_organization_id`, `trainer_id`, `backup_trainer_id`, `coordinator_id`, `cost`, `is_planned`, `is_executed`, `planned_id`, `client_id`, `sessions[]`, `documentations[]`, `participants[]`
- **CourseSession** — `id`, `session_date`, `start_time`, `end_time`, `course_id`, `client_id`
- **CourseParticipant** — `id`, `course_id`, `person_id`, `client_id`, `final_exams[]`, `performance_evaluation_details[]`
- **CourseSessionParticipant** — `id`, `course_session_id`, `person_id`, `client_id`
- **CourseParticipantFinalExam** — `id`, `exam_date`, `organized_by`, `course_participant_id`, `documentation_id`, `exam_score`, `client_id`
- **Documentation** — `id`, `name`, `file_path`, `documentation_type_id`, `learning_subject_id`, `path_grade_subject_id`, `path_grade_subject_term_id`, `course_id`, `course_session_id`, `client_id`
- **DocumentationType** — `id`, `name`, `client_id`
- **PerformanceEvaluationCategory** — `id`, `name`, `client_id`
- **PerformanceEvaluationItem** — `id`, `name`, `minimum_score`, `maximum_score`, `client_id`
- **PerformanceEvaluationParticipantCategory** — `id`, `participant_type_id`, `performance_evaluation_category_id`, `client_id`, `items[]`
- **PerformanceEvaluationParticipantCategoryItem** — `id`, `performance_evaluation_participant_category_id`, `performance_evaluation_item_id`, `client_id`
- **PerformanceEvaluationCourseDetail** — `id`, `course_id`, `score`, `notes`, `entered_by`, `entered_at`, `performance_evaluation_participant_category_item_id`, `client_id`
- **PerformanceEvaluationCourseParticipantDetail** — `id`, `course_participant_id`, `score`, `notes`, `entered_by`, `entered_at`, `performance_evaluation_participant_category_item_id`, `client_id`

