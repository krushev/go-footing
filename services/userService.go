package services

import (
	"github.com/asaskevich/govalidator"
	"github.com/krushev/go-footing/db"
	"github.com/krushev/go-footing/models"
	"github.com/krushev/go-footing/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type UserService interface {
	FindAll() (*[]models.User, error)
	FindActive() (*[]models.User, error)
	Search(q string, pagination *models.Pagination) (*[]models.User, error)
	GetById(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user models.User) (*models.User, error)
	EncryptPassword(password string) (string, error)
	VerifyPassword(hashed, password string) error
	IsEmpty() (bool, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(conn db.Connection) *userService {
	return &userService{db: conn.GetDB()}
}

func (us *userService) FindAll() (*[]models.User, error) {
	var users []models.User
	result := us.db.Find(&users)

	return &users, result.Error
}

func (us *userService) Search(q string, pagination *models.Pagination) (*[]models.User, error) {
	var users []models.User
	offset := (pagination.Page - 1) * pagination.Size
	queryBuilder := us.db.Limit(pagination.Size).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Model(&models.User{}).
		Where("lower(name) LIKE lower(?)", "%" + q + "%").
		Or("lower(email) LIKE lower(?)", "%" + q + "%").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (us *userService) FindActive() (*[]models.User, error) {
	var users []models.User
	result := us.db.Find(&users, models.User{Active: true})

	return &users, result.Error
}

func (us *userService) GetById(id string) (*models.User, error) {
	var user models.User
	result := us.db.First(&user, id)

	return &user, result.Error
}

func (us *userService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := us.db.Where("email = ?", email).First(&user)

	return &user, result.Error
}

func (us *userService) IsEmpty() (bool, error) {
	var user models.User
	result := us.db.First(&user)
	return result.RowsAffected <= 0, result.Error
}

func (us *userService) Create(user models.User) (*models.User, error) {
	if strings.TrimSpace(user.Password) == "" {
		return &user, util.ErrEmptyPassword
	}

	user.Email = util.NormalizeEmail(user.Email)
	if !govalidator.IsEmail(user.Email) {
		return &user, util.ErrInvalidEmail
	}

	if !IsValid(user.Roles) {
		return &user, util.ErrInvalidRole
	}

	exists, err  := us.GetByEmail(user.Email)
	if err == gorm.ErrRecordNotFound {
		user.Password, err = us.EncryptPassword(user.Password)
		if err != nil {
			return &user, err
		}
		result := us.db.Create(&user)
		return &user, result.Error
	}

	if exists != nil {
		return &user, util.ErrEmailAlreadyExists
	}
	return &user, err
}

func (us *userService) EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (us *userService) VerifyPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}