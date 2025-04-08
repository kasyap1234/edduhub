package student

type StudentService interface {
}

type studentService struct {
}

func NewstudentService() StudentService {
	return &studentService{}
}
