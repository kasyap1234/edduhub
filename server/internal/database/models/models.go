package models

import "time"

type User struct {
	ID string  `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Email string `json:"email" unique`
	Role string `json:"role" gorm:"type:enum('admin','super_admin','student','faculty')"`
}


type Student struct {
	ID string `json:"id" gorm:"primaryKey"`
	UserID string `json:"user_id"`
	User User `gorm: "foreignKey:UserID"`
	RollNumber string `json:"roll_number"`
	Year int `json:"year"`
	Semester int `json:"semester"`
	GPA float64 `json:"gpa"`

}

type Attendance struct {
ID 	string  `json:"id" gorm:"primaryKey"`
StudentID string `json:student_id"`
RollNumber string `json:"roll_number"`
Date time.Time `json:"date"`
Status string `json:"status"`
MarkedBy string `json:"status"`
Student Student `gorm:"foreignKey"`


}
