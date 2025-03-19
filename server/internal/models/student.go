package models 


type Student struct {
	StudentID uint `json:"student_id" bun:"student_id,auto_increment"`
	User user
	Batch     int      `json:"Batch" bun:"batch"`
	Year      int      `json:"Year" bun:"year"`
	Sem       int      `json:"Sem" bun:"sem"`
	Subjects  Subjects `json:"Subjects" bun:"subjects,json"`
}
