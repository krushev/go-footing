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
	FindAll() (error, *[]models.User)
	FindActive() (error, *[]models.User)
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

func (us *userService) FindAll() (error, *[]models.User) {
	var users []models.User
	result := us.db.Find(&users)

	return result.Error, &users
}

func (us *userService) FindActive() (error, *[]models.User) {
	var users []models.User
	result := us.db.Find(&users, models.User{Active: true})

	return result.Error, &users
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