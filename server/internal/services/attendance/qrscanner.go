package attendance

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
)

// studentId in the request body
// course id in the request body and lecture id obtained from qr code
// need to check if student is enrolled in the course
// need to check if the student is enrolled in the lecture
type QRCodeData struct {
	CollegeID int       `json:"college_id"`
	CourseID  int       `json:"course_id"`
	LectureID int       `json:"lecture_id"`
	TimeStamp time.Time `json:"time_stamp"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (a *attendanceService) GenerateQRCode(ctx context.Context, collegeID int, courseID int, lectureID int) (string, error) {
	// generate qr code
	// return qr code
	// need to check if student is enrolled in the course for marking attendance

	now := time.Now()
	expiresAt := now.Add(30 * time.Minute)
	qrCodeData := QRCodeData{
		CollegeID: collegeID,
		CourseID:  courseID,
		LectureID: lectureID,
		TimeStamp: now,
		ExpiresAt: expiresAt,
	}
	jsonData, err := json.Marshal(qrCodeData)
	if err != nil {
		return "", err
	}
	qrBytes, err := qrcode.Encode(string(jsonData), qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)
	qrDataURI := fmt.Sprintf("data:image/png;base64,%s", qrBase64)
	return qrDataURI, nil

}

func (a *attendanceService) ProcessQRCode(ctx context.Context, collegeID int, studentID int, qrCodeContent string) error {
	var qrData QRCodeData
	if err := json.Unmarshal([]byte(qrCodeContent), &qrData); err != nil {
		return errors.New("invalid qr code")
	}

	enrolled, err := a.VerifyStudentEnrollment(ctx, studentID, qrData.CollegeID, qrData.CourseID)

	if err != nil {
		return err
	}
	if time.Now().After(qrData.ExpiresAt) {
		return errors.New("qr code has expired")

	}
	if !enrolled {
		return errors.New("student is not enrolled in the course")
	}
	if qrData.CollegeID != collegeID {
		return errors.New("college id does not match")
	}
	marked, err := a.MarkAttendance(ctx, qrData.CollegeID, studentID, qrData.CourseID, qrData.LectureID)
	if err != nil {
		return err
	}
	if !marked {
		return errors.New("attendance not marked")

	}

	return nil
}
