package models

import (
	"time"
)

type User struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Role  string `json:"role" gorm:"type:enum('admin','super_admin','student','faculty')"`
}

type Student struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	UserID     string  `json:"user_id"`
	User       User    `gorm:"foreignKey:UserID"`
	RollNumber string  `json:"roll_number" gorm:"unique"`
	Year       int     `json:"year"`
	Semester   int     `json:"semester"`
	GPA        float64 `json:"gpa"`
}

type Attendance struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	StudentID  string    `json:"student_id"`
	RollNumber string    `json:"roll_number"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status" gorm:"type:enum('present','absent','late')"`
	MarkedBy   string    `json:"marked_by"`
	Student    Student   `gorm:"foreignKey:StudentID"`
}

type Fee struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	StudentID   string    `json:"student_id"`
	RollNumber  string    `json:"roll_number"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type" gorm:"type:enum('tuition','hostel','transport','other')"`
	DueDate     time.Time `json:"due_date"`
	PaidDate    time.Time `json:"paid_date,omitempty"`
	Status      string    `json:"status" gorm:"type:enum('pending','paid','overdue')"`
	PaymentMode string    `json:"payment_mode,omitempty"`
	Student     Student   `gorm:"foreignKey:StudentID"`
}
