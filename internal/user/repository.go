package user

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	GetAll(filters Filters, offset, limit int) ([]User, error)
	Get(id string) (*User, error)
	Delete(id string) error
	Update(id string, firstName, lastName, email, phone *string) error
	Count(filters Filters) (int, error)
}

type repository struct {
	logger *log.Logger
	db     *gorm.DB
}

func (r repository) Create(user *User) error {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Println(err)
		return err
	}

	r.logger.Println("user created with id: ", user.ID)
	return nil
}

func (r repository) GetAll(filters Filters, offset, limit int) ([]User, error) {
	var users []User

	tx := r.db.Model(&users)

	tx = applyFilters(tx, filters)

	tx = tx.Limit(limit).Offset(offset)

	if err := tx.Order("created_at desc").Find(&users).Error; err != nil {
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

func (r repository) Count(filters Filters) (int, error) {
	var count int64

	tx := r.db.Model(User{})

	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func NewRepository(logger *log.Logger, db *gorm.DB) Repository {
	return &repository{logger: logger, db: db}
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}

	return tx
}
