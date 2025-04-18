package models

import "github.com/uptrace/bun"

type Student struct {
	bun.BaseModel    `bun:"table:students"`
	StudentID        int      `json:"student_id" bun:",pk,autoincrement"`
	// UserID           int      `json:"user_id" bun:",notnull,unique"`
	KratosIdentityID string   `json:"-" bun:",notnull,unique"`
	CollegeID        int      `json:"college_id" bun:"college_id,notnull"`
	RollNo           string   `json:"roll_no" bun:",unique"`
	Batch            int      `json:"batch"`
	Year             int      `json:"year"`
	Sem              int      `json:"sem"`
	Subjects         Subjects `json:"subjects" bun:",json"`
	IsActive         bool     `json:"is_active" bun:",default:true"`
	// Relations
	// User *User `bun:"rel:belongs-to,join:user_id=id"`
}
