I. Core Authentication & Authorization (Completion)

Logout: Implement HandleLogout handler and corresponding AuthService logic (invalidate Kratos session).
Token Refresh: Implement RefreshToken handler and AuthService logic (if using refresh tokens alongside Kratos sessions, otherwise potentially remove).
Password Reset: Implement RequestPasswordReset, CompletePasswordReset handlers and AuthService logic (trigger Kratos password reset flow).
Session Validation Middleware: Fully implement and test middleware.ValidateSession to correctly verify Kratos sessions and populate user context.
Authorization Middleware: Implement and test middleware.RequireCollege, middleware.RequireRole, middleware.LoadStudentProfile, middleware.VerifyStudentOwnership.
II. System & Foundational

System Handler: Implement SystemHandler with HealthCheck and SwaggerDocs.
Database Schema & Migrations: Define and implement the full database schema for all entities and set up migrations.
Core Service Interfaces & Implementations: Define interfaces and implement the core logic for all services (User, College, Student, Course, Lecture, Grade, Calendar, etc.).
Configuration & Logging: Ensure robust configuration loading and consistent logging are fully implemented.
III. User & Profile Management

User CRUD: Implement UserHandler and UserService for listing, creating, retrieving, updating, and deleting users.
User Role/Status Management: Implement endpoints for updating user roles and status (active/inactive).
User Profile Management: Implement UserHandler and UserService for getting/updating user profiles (/profile).
Change Password (Profile): Implement the endpoint for logged-in users to change their own password (likely initiating a Kratos flow).
IV. College Management

College Details: Implement CollegeHandler and CollegeService for retrieving and updating college information.
College Statistics: Implement logic to calculate and retrieve college-level statistics.
V. Student Management

Student CRUD: Implement StudentHandler and StudentService for listing, creating, retrieving, updating, and deleting student profiles (linking them to users).
Freeze Student Account: Implement logic to temporarily disable/freeze a student's record.
VI. Course & Enrollment Management

Course CRUD: Implement CourseHandler and CourseService for listing, creating, retrieving, updating, and deleting courses.
Student Enrollment: Implement endpoints and service logic for enrolling students in courses and removing them.
List Enrolled Students: Implement endpoint to list students enrolled in a specific course.
VII. Lecture Management

Lecture CRUD: Implement LectureHandler and LectureService for listing, creating, retrieving, updating, and deleting lectures associated with a course.
VIII. Attendance Management (Completion)

Fix Existing Handlers: Correct parameter extraction and error handling in existing attendance getters.
Update Attendance: Implement UpdateAttendance handler and AttendanceService logic (allow manual correction of attendance records).
Bulk Mark Attendance: Implement MarkBulkAttendance handler and AttendanceService logic (for faculty to mark multiple students at once).
Attendance Reports: Implement GetCourseAttendanceReport, GetStudentAttendanceReport handlers and AttendanceService logic (aggregate data for reporting).
IX. Grade & Assessment Management

Assessment CRUD: Implement GradeHandler and GradeService for creating, retrieving, updating, and deleting assessments (e.g., tests, assignments) for a course.
Submit/Update Scores: Implement endpoint and service logic for submitting/updating scores for students on specific assessments.
Get Grades (Course): Implement endpoint to retrieve all grades/scores for a specific course.
Get Grades (Student): Implement endpoint for a student (or authorized staff) to retrieve their own grades across courses or for a specific course.
X. Calendar & Scheduling

Calendar Event CRUD: Implement CalendarHandler and CalendarService for listing, creating, retrieving, updating, and deleting calendar events (college-wide, course-specific, etc.).
XI. Testing & Documentation

Unit Tests: Write unit tests for handlers and services.
Integration Tests: Write tests for API endpoints.
API Documentation: Ensure Swagger/OpenAPI documentation is complete and accurate.