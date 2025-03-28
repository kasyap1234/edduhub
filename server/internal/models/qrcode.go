package models

import (
	"github.com/uptrace/bun"
	"time"
)

type QRCode struct {
	bun.BaseModel `bun:"table:qrcodes"`
	ID            int       `bun:",pk,autoincrement"`
	Code          string    `json:"code"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type AttendanceRequest struct {
	StudentID int `json:"student_id"`
	QRCodeID string `json:"qr_code_id"`
}
