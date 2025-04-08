package handler

import (
	"eduhub/server/internal/services/attendance"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

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
	studentIDstr := c.QueryParam("studentID")
	courseIDstr := c.QueryParam("courseID")
	lectureIDstr := c.QueryParam("lectureID")
	studentID, _ := strconv.Atoi(studentIDstr)
	courseID, _ := strconv.Atoi(courseIDstr)
	lectureID, _ := strconv.Atoi(lectureIDstr)
	ok, err := a.attendanceService.MarkAttendance(ctx, studentID, courseID, lectureID)
	if ok {
		//
		return nil
	}
	return err

}

func (a *AttendanceHandler) GetAttendanceByCourse(c echo.Context) error {
	// ctx := c.Request().Context()
	courseIDstr := c.QueryParam("courseID")
	courseID, _ := strconv.Atoi(courseIDstr)

	attendance, err := a.attendanceService.GetAttendanceByCourse(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, map[string]interface{}{"response": attendance})

}

func (a *AttendanceHandler) GetAttendanceForStudent(c echo.Context) error {
	studentIDstr := c.QueryParam("studentID")
	studentID, _ := strconv.Atoi(studentIDstr)
	attendance, err := a.attendanceService.GetAttendanceByStudent(studentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"response": attendance})
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
	return c.JSON(http.StatusOK, map[string]interface{}{"response": attendance})
}
