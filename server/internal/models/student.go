package models

import "github.com/uptrace/bun"

type Student struct {
	bun.BaseModel `bun:"table:students"`
	ID            int      `json:"student_id" bun:",pk,autoincrement"`
	UserID        int      `json:"user_id" bun:",notnull,unique"`
	RollNo        string   `json:"roll_no" bun:",unique"`
	Batch         int      `json:"batch"`
	Year          int      `json:"year"`
	Sem           int      `json:"sem"`
	Subjects      Subjects `json:"subjects" bun:",json"`
	IsActive      bool     `json:"is_active" bun:",default:true"`
	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
