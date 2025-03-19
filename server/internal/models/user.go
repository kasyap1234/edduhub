package models 


type User struct {
	ID        int      `json:"ID" bun:"id,autoincrement"`
	CollegeID int      `json:"CollegeID" bun:"college_id"`
	RollNo    string   `json:"RollNo" bun:"roll_no,pk"`
	Name 	string `json:"name" bun:"name"`
	Role string `json:"role" bun:"role"`

}

