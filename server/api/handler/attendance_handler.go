package handler

import (
	"eduhub/server/internal/helpers"
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

func (a *AttendanceHandler) GenerateQRCode(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}
	courseIDStr := c.Param("courseID")
	lectureIDstr := c.Param("lectureID")
	courseID, err := strconv.Atoi(courseIDStr)
	lectureID, err := strconv.Atoi(lectureIDstr)
	if err != nil {
		return helpers.Error(c, err, 400)
	}
	qrCode, err := a.attendanceService.GenerateQRCode(ctx, collegeID, courseID, lectureID)
	if err != nil {
		return helpers.Error(c, err, 400)
	}
	return helpers.Success(c, qrCode, 200)
}

// func (a *AttendanceHandler) ProcessQRCode(c echo.Context) error {
// 	ctx := c.Request().Context()
// 	collegeID, err := helpers.ExtractCollegeID(c)
// 	if err != nil {
// 		return helpers.Error(c, err, 400)
// 	}

// 	// TODO extract student ID from context , need to link studentid from kratos with db ;
// 	studentID, err := helpers.ExtractStudentID(c)
// 	if err != nil {
// 		return helpers.Error(c, err, 400)

// 	}
	


func (a *AttendanceHandler) MarkAttendance(c echo.Context) error {
	ctx := c.Request().Context()

	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return helpers.Error(c, "Invalid college ID", http.StatusBadRequest)
	}

	studentID, err := helpers.ExtractStudentID(c)
	if err != nil {
		return helpers.Error(c, "Invalid student ID", http.StatusBadRequest)
	}

	courseIDStr := c.QueryParam("courseID")
	lectureIDStr := c.QueryParam("lectureID")

	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		return helpers.Error(c, "Invalid course ID", http.StatusBadRequest)
	}

	lectureID, err := strconv.Atoi(lectureIDStr)
	if err != nil {
		return helpers.Error(c, "Invalid lecture ID", http.StatusBadRequest)
	}
	ok,err :=a.attendanceService.VerifyStudentEnrollment(ctx,collegeID,studentID,courseID)
	
	ok, err = a.attendanceService.MarkAttendance(ctx, collegeID, studentID, courseID, lectureID)
	if err != nil {
		return helpers.Error(c, err.Error(), http.StatusForbidden)
	}

	if !ok {
		return helpers.Error(c, "Failed to mark attendance", http.StatusInternalServerError)
	}

	return helpers.Success(c, map[string]interface{}{
		"message": "Attendance marked successfully",
		"status":  true,
	}, http.StatusOK)
}

func (a *AttendanceHandler) GetAttendanceByCourse(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}
	courseIDstr := c.QueryParam("courseID")
	courseID, _ := strconv.Atoi(courseIDstr)

	attendance, err := a.attendanceService.GetAttendanceByCourse(ctx, collegeID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)

	}
	return c.JSON(http.StatusOK, attendance)

}

// use structs instead of maps while returning c.JSON
func (a *AttendanceHandler) GetAttendanceForStudent(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}
	studentID, err := helpers.ExtractStudentID(c)
	if err != nil {
		return err
	}

	attendance, err := a.attendanceService.GetAttendanceByStudent(ctx, collegeID, studentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
	}
	return c.JSON(http.StatusOK, attendance)
}

func (a *AttendanceHandler) GetAttendanceByStudentAndCourse(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}
	studentID, err := helpers.ExtractStudentID(c)
	if err != nil {
		return err
	}
	courseIDstr := c.QueryParam("courseID")
	courseID, err := strconv.Atoi(courseIDstr)
	if err != nil {
		return err
	}
	attendance, err := a.attendanceService.GetAttendanceByStudentAndCourse(ctx, collegeID, studentID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, attendance)
}
