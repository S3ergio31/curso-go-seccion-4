package enrollment

import (
	"errors"
	"log"

	"github.com/S3ergio31/curso-go-seccion-4/internal/course"
	"github.com/S3ergio31/curso-go-seccion-4/internal/domain"
	"github.com/S3ergio31/curso-go-seccion-4/internal/user"
)

type Service interface {
	Create(userID, courseID string) (*domain.Enrollment, error)
}

type service struct {
	userService   user.Service
	courseService course.Service
	logger        *log.Logger
	repository    Repository
}

func (s service) Create(userID, courseID string) (*domain.Enrollment, error) {
	enrollment := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "P",
	}

	if _, err := s.userService.Get(userID); err != nil {
		return nil, errors.New("user id does not exists")
	}

	if _, err := s.courseService.Get(userID); err != nil {
		return nil, errors.New("course id does not exists")
	}

	if err := s.repository.Create(enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func NewService(
	repository Repository,
	logger *log.Logger,
	userService user.Service,
	courseService course.Service,
) Service {
	return &service{
		logger:        logger,
		repository:    repository,
		userService:   userService,
		courseService: courseService,
	}
}
