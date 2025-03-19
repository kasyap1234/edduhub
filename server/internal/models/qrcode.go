package models


type QRCode struct {
	bun.BaseModel `bun:"table:qrcodes"`
	ID            int    `bun:",pk,autoincrement"`
	Code          string `json:"code"`
	ExpiresAt     int64  `json:"expires_at"` 

}
