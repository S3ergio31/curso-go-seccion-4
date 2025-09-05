package user

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
}

type repository struct {
	logger *log.Logger
	db     *gorm.DB
}

func (r repository) Create(user *User) error {
	r.logger.Printf("Create")
	return nil
}

func NewRepository(logger *log.Logger, db *gorm.DB) Repository {
	return &repository{logger: logger, db: db}
}
