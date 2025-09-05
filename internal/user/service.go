package user

import "log"

type Service interface {
	Create(firstName, lastName, email, phone string) error
}

type service struct {
	logger     *log.Logger
	repository Repository
}

func (s service) Create(firstName, lastName, email, phone string) error {
	s.repository.Create(&User{})

	return nil
}

func NewService(repository Repository, logger *log.Logger) Service {
	return &service{logger: logger, repository: repository}
}
