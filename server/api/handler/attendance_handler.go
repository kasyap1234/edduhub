package handler

import (
	"eduhub/server/internal/helpers"
	"eduhub/server/internal/models"
	"eduhub/server/internal/services/attendance"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error error
}

type AttendanceResponse struct {
	Success bool
	Message string
}
type Response struct {
	Message any
}
type attendanceModel models.Attendance

type AttendanceHandler struct {
	attendanceService attendance.AttendanceService
}

func NewAttedanceHandler(attendance attendance.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: attendance,
	}
}

func (a *AttendanceHandler) MarkAttendance(c echo.Context) error {
	// return a.attendanceService.MarkAttendance(c, studentID , courseID int, lectureID int)
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}

	studentIDstr := c.QueryParam("studentID")
	courseIDstr := c.QueryParam("courseID")
	lectureIDstr := c.QueryParam("lectureID")
	studentID, _ := strconv.Atoi(studentIDstr)
	courseID, _ := strconv.Atoi(courseIDstr)
	lectureID, _ := strconv.Atoi(lectureIDstr)
	ok, _ := a.attendanceService.MarkAttendance(ctx, collegeID, studentID, courseID, lectureID)
	if ok {
		//
		return helpers.Success(c, "attendance marked")
	}
	return helpers.Error(c, "attendance not marked")
}

func (a *AttendanceHandler) GetAttendanceByCourse(c echo.Context) error {
	// ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}
	courseIDstr := c.QueryParam("courseID")
	courseID, _ := strconv.Atoi(courseIDstr)

	attendance, err := a.attendanceService.GetAttendanceByCourse(collegeID,courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)

	}
	return c.JSON(http.StatusOK, attendance)

}

// use structs instead of maps while returning c.JSON
func (a *AttendanceHandler) GetAttendanceForStudent(c echo.Context) error {
	studentIDstr := c.QueryParam("studentID")
	studentID, _ := strconv.Atoi(studentIDstr)
	attendance, err := a.attendanceService.GetAttendanceByStudent(studentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
	}
	return c.JSON(http.StatusOK, attendance)
}

func (a *AttendanceHandler) GetAttendanceByStudentAndCourse(c echo.Context) error {

	studentIDstr := c.QueryParam("studentID")
	courseIDstr := c.QueryParam("courseID")
	courseID, _ := strconv.Atoi(courseIDstr)
	studentID, _ := strconv.Atoi(studentIDstr)

	attendance, err := a.attendanceService.GetAttendanceByStudentAndCourse(studentID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, attendance)
}
