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

func NewAttendanceHandler(attendance attendance.AttendanceService) *AttendanceHandler {
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
	// courseIDStr := c.Param("courseID")
	// lectureIDstr := c.Param("lectureID")
	// courseID, err := strconv.Atoi(courseIDStr)
	// if err != nil {
	// 	return helpers.Error(c, err.Error(), 400)
	// }
	// lectureID, err := strconv.Atoi(lectureIDstr)
	// if err != nil {
	// 	return helpers.Error(c, err, 400)
	// }
	courseID, err := helpers.GetIDFromParam(c, "courseID")
	if err != nil {
		return err
	}
	lectureID, err := helpers.GetIDFromParam(c, "lectureID")
	if err != nil {
		return err
	}
	qrCode, err := a.attendanceService.GenerateQRCode(ctx, collegeID, courseID, lectureID)
	if err != nil {
		return helpers.Error(c, err, 400)
	}
	return helpers.Success(c, qrCode, 200)
}

func (a *AttendanceHandler) MarkAttendance(c echo.Context) error {
	ctx := c.Request().Context()

	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}

	studentID, err := helpers.ExtractStudentID(c)
	if err != nil {
		return err
	}

	courseID, err := helpers.GetIDFromParam(c, "courseID")
	if err != nil {
		return err
	}
	lectureID, err := helpers.GetIDFromParam(c, "lectureID")
	if err != nil {
		return err
	}

	ok, err := a.attendanceService.MarkAttendance(ctx, collegeID, studentID, courseID, lectureID)
	if err != nil {
		return helpers.Error(c, err.Error(), http.StatusInternalServerError)
	}

	if !ok {
		return helpers.Error(c, "Failed to mark attendance", http.StatusInternalServerError)
	}

	return helpers.Success(c, "attendance marked", http.StatusOK)
}

func (a *AttendanceHandler) GetAttendanceByCourse(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}
	courseIDstr := c.QueryParam("courseID")
	courseID, err := strconv.Atoi(courseIDstr)
	if err != nil {
		return helpers.Error(c, "Invalid course ID", http.StatusBadRequest)
	}
	// courseID, err := helpers.GetIDFromParam(c, "courseID")
	attendance, err := a.attendanceService.GetAttendanceByCourse(ctx, collegeID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)

	}
	return helpers.Success(c, attendance, http.StatusOK)

}

// use structs instead of maps while returning c.JSON
func (a *AttendanceHandler) GetAttendanceForStudent(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	studentID, err := helpers.ExtractStudentID(c)
	if err != nil {
		return helpers.Error(c, "invalid studentID", 400)
	}

	attendance, err := a.attendanceService.GetAttendanceByStudent(ctx, collegeID, studentID)
	if err != nil {
		return helpers.Error(c, "unable to get attendance by student", http.StatusInternalServerError)

	}
	return helpers.Success(c, attendance, http.StatusOK)
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
	// courseID, err := helpers.GetIDFromParam(c, "courseID")
	// if err != nil {
	// 	return err
	// }
	attendance, err := a.attendanceService.GetAttendanceByStudentAndCourse(ctx, collegeID, studentID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, attendance)
}
