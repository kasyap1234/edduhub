package grades

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"

	"github.com/go-playground/validator/v10"
)

type GradeServices interface {
	CreateGrade(ctx context.Context, grade *models.Grade) error
	GetGradeByID(ctx context.Context, gradeId int, collegeID int) (*models.Grade, error)
	UpdateGrade(ctx context.Context, grade *models.Grade) error
	DeleteGrade(ctx context.Context, gradeID int, collegeID int) error
	GetGrades(ctx context.Context, filter models.GradeFilter) ([]*models.Grade, error)
}

type gradeServices struct {
	gradeRepo repository.GradeRepository
	validate  validator.Validate
}

func NewGradeServices(gradeRepo repository.GradeRepository) GradeServices {
	&gradeServices{
		gradeRepo: gradeRepo,
		validate:  *validator.New(),
	}
}

func (g *gradeServices) CreateGrade(ctx context.Context, grade *models.Grade) error {

}

func (g *gradeServices) GetGradeByID(ctx context.Context, gradeID int, collegeID int) (*models.Grade, error) {

}

func (g *gradeServices) UpdateGrade(ctx context.Context, grade *models.Grade) error {

}

func (g *gradeServices) DeleteGrade(ctx context.Context, gradeID int, collegeID int) error {

}

func (g *gradeServices) GetGrades(ctx context.Context, filters models.GradeFilter) ([]*models.Grade, error) {

}
