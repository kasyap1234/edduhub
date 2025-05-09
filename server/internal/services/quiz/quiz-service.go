package quiz

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type QuizService interface {
	CreateQuiz(ctx context.Context, quiz *models.Quiz) error
	GetQuizByID(ctx context.Context, collegeID int, quizID int) (*models.Quiz, error)
	UpdateQuiz(ctx context.Context, quiz *models.Quiz) error
	DeleteQuiz(ctx context.Context, collegeID int, quizID int) error
	FindQuizzesByCourse(ctx context.Context, collegeID int, courseID int, limit, offset uint64) ([]*models.Quiz, error)

	CreateQuestion(ctx context.Context, question *models.Question) error
	GetQuestionByID(ctx context.Context, questionID int) (*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, questionID int) error
	FindQuestionByQuiz(ctx context.Context, quizID int, limit, offset uint64, withOptions bool) ([]*models.Question, error)
	CountQuestionByQuiz(ctx context.Context, quizID int) (int, error)

	CreateAnswerOption(ctx context.Context, option *models.AnswerOption) error
	GetAnswerOptionByID(ctx context.Context, optionID int) (*models.AnswerOption, error)
	UpdateAnswerOption(ctx context.Context, option *models.AnswerOption) error
	DeleteAnswerOption(ctx context.Context, optionID int) error	
	FindAnswerOptionsByQuestion(ctx context.Context, questionID int) ([]*models.AnswerOption, error)

	StartQuizAttempt(ctx context.Context,attempt *models.QuizAttempt)error 
	GetQuizAttemptByID(ctx context.Context,collegeID int,attemptID int)(*models.QuizAttempt,error)
	SubmitQuizAttempt(ctx context.Context,collegeID int,attemptID int)error 
	GradeQuizAttempt(ctx context.Context,collegeID )
}


type quizService struct {
	quizRepo repository.QuizRepository
	validate  validator.Validate
}

func NewQuizService(quizRepo repository.QuizRepository)QuizService{
	&quizService{
		quizRepo: quizRepo, 
		validate: validator.New()
	}
}

func(q *quizService)CreateQuiz(ctx context.Context,quiz*models.Quiz)error {
	if err :=q.validate.Struct(quiz); err !=nil{
		return fmt.Errorf("validation failed for quiz %w",err)		
	}
	return q.quizRepo.CreateQuiz(ctx,quiz)
}

func(q*quizService)GetQuizByID(ctx context.Context,collegeID int,quizID int)(*models.Quiz,error){
	return q.quizRepo.GetQuizByID(ctx,collegeID,quizID)
}

func(q *quizService)UpdateQuiz(ctx context.Context,collegeID int,quizID int)(*models.Quiz,error){
	if err :=q.validate.Struct(quiz); err !=nil{
		return fmt.Errorf("validation failed for quiz %w",err)
	}
	if quiz.ID ==0 {
		return fmt.Errorf("quiz id required")
	}
	return q.quizRepo.UpdateQuiz(ctx,quiz)
}

func(q*quizService)DeleteQuiz(ctx context.Context,collegeID int,quiz int)error {
	if _,err:=q.GetQuizByID(ctx,collegeID,quizID); err !=nil{
		return fmt.Errorf("quiz does not exist")
	}
	if err=q.validate.Struct(quiz); err !=nil{
		return fmt.Errorf("validation failed %w",err)
	}
	return s.quizRepo.DeleteQuiz(ctx,collegeID,quizID)
}

func(q*quizService)F