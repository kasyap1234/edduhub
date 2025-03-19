package models 


type User struct {
	ID        int      `json:"ID" bun:"id,autoincrement"`
	CollegeID int      `json:"CollegeID" bun:"college_id"`
	
	Name 	string `json:"name" bun:"name"`
	Role string `json:"role" bun:"role"`
	Student *Student `bun:"rel:has-one,join:id=user_id,nullzero"`
	
}

