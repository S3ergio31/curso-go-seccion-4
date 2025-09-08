package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	GetAll() ([]User, error)
	Get(id string) (*User, error)
	Delete(id string) error
	Update(id string, firstName, lastName, email, phone *string) error
}

type repository struct {
	logger *log.Logger
	db     *gorm.DB
}

func (r repository) Create(user *User) error {
	user.ID = uuid.NewString()

	if err := r.db.Create(user).Error; err != nil {
		r.logger.Println(err)
		return err
	}

	r.logger.Println("user created with id: ", user.ID)
	return nil
}

func (r repository) GetAll() ([]User, error) {
	var users []User
	if err := r.db.Model(&users).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r repository) Get(id string) (*User, error) {
	user := User{ID: id}
	if err := r.db.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) Delete(id string) error {
	user := User{ID: id}

	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r repository) Update(id string, firstName, lastName, email, phone *string) error {
	values := make(map[string]interface{}, 0)

	if firstName != nil {
		values["first_name"] = *firstName
	}

	if lastName != nil {
		values["last_name"] = *lastName
	}

	if email != nil {
		values["email"] = *email
	}

	if phone != nil {
		values["phone"] = *phone
	}

	if err := r.db.Model(&User{}).Where("id = ?", id).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func NewRepository(logger *log.Logger, db *gorm.DB) Repository {
	return &repository{logger: logger, db: db}
}
