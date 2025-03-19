package models 


type Student struct {
	StudentID uint `json:"student_id" bun:"student_id,auto_increment"`
	UserID int `bun:"not_null"`
	User *User `rel:"belongs-to,join:user_id=id"`
	RollNo    string   `json:"RollNo" bun:"roll_no,pk"`
	Batch     int      `json:"Batch" bun:"batch"`
	Year      int      `json:"Year" bun:"year"`
	Sem       int      `json:"Sem" bun:"sem"`
	Subjects  Subjects `json:"Subjects" bun:"subjects,json"`
}
