package handler

import (
	"net/http"

	"eduhub/server/internal/helpers"
	"eduhub/server/internal/services/attendance"

	"github.com/labstack/echo/v4"
)

type AttendanceHandler struct {
	attendanceService attendance.AttendanceService
}

type QRCodeRequest struct {
	QRCodeData string `json:"qrcode_data"`
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

func (a *AttendanceHandler) ProcessAttendance(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return err
	}
	studentId, err := helpers.ExtractStudentID(c)
	if err != nil {
		return err
	}
	var qrcodeData QRCodeRequest
	if c.Bind(&qrcodeData); err != nil {
		return helpers.Error(c, "invalid request body", 400)
	}
	err = a.attendanceService.ProcessQRCode(ctx, collegeID, studentId, qrcodeData.QRCodeData)
	return err
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

	courseID, err := helpers.GetIDFromParam(c, "courseID")
	if err != nil {
		return helpers.Error(c, "Invalid course ID", http.StatusBadRequest)
	}

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
	if err != nil {
		return helpers.Error(c, "invalid collegeID", 400)
	}
	studentID, err := helpers.GetIDFromParam(c, "studentID")
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
	studentID, err := helpers.GetIDFromParam(c, "studentID")
	if err != nil {
		return err
	}

	courseID, err := helpers.GetIDFromParam(c, "courseID")
	if err != nil {
		return err
	}

	attendance, err := a.attendanceService.GetAttendanceByStudentAndCourse(ctx, collegeID, studentID, courseID)
	if err != nil {
		return helpers.Error(c, "unable to get attendance", http.StatusInternalServerError)
	}
	return helpers.Success(c, attendance, 200)
}

func (a *AttendanceHandler) UpdateAttendance(c echo.Context) error {
	ctx := c.Request().Context()
	collegeID, err := helpers.ExtractCollegeID(c)
	if err != nil {
		return helpers.Error(c, "Invalid collegeID", 400)

	}
	studentID, err := helpers.ExtractStudentID(c)
	if err != nil {
		return helpers.Error(c, "Invalid studentID", 400)
	}
	courseID, err := helpers.GetIDFromParam(c, "courseID")
	if err != nil {
		return helpers.Error(c, "invalid courseID ", 400)
	}

	ok, err := a.attendanceService.UpdateAttendance(ctx, collegeID, studentID, courseID, studentID)
	if !ok {
		return helpers.Error(c, "Unable update attendance", 500)
	}
	if err != nil {
		return helpers.Error(c, "error in update attendance", 502)
	}

	return helpers.Success(c, "Success", 200)
}
