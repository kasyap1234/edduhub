package repository

import "eduhub/server/internal/models"

type StudentRepository interface {
	GetStudentByRollNo(rollNo string)(*models.Student,error)
	GetStudentByID(studentID int)(*models.Student,error)
	
}
