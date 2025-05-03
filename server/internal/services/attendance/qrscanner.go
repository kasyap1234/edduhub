package attendance

import (
	"context"
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
	qrBase64 := fmt.Sprintf("data ,%s", qrBytes)
	return qrBase64, nil
}

// process qr and take values from it to mark attendance(process qr and chaning state)
func (a *attendanceService) ProcessQRCode(ctx context.Context, collegeID int, studentID int, qrCodeContent string) error {
	var qrData QRCodeData

	if err := json.Unmarshal([]byte(qrCodeContent), &qrData); err != nil {
		return errors.New("invalid qr code")
	}
	timestamp: = qrData.TimeStamp

	if timeStamp <time.Now(){
		return errros.New("qr expried")
	}

	marked, err := a.MarkAttendance(ctx,collegeID,studentID,qrData.CourseID,qrData.LectureID)
	if !marked {
		return errors.New("unable to mark attendance using processqrcode")
	}
	if err != nil {
		return err
	}
	return nil
}

/// process qr takes qr input and marks attendance
// mark attendance verifies student and marks attendance
