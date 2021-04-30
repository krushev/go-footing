package services

import (
	"errors"
	"github.com/krushev/go-footing/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BootstrapService interface {
	InitUsers()
}

type bootstrapService struct {
	us UserService
}

func NewBootstrapService(us UserService) *bootstrapService {
	return &bootstrapService{us}
}

func (bs *bootstrapService) InitUsers() {
	_, err := bs.us.IsEmpty()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Fatal(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		admin := models.User{ Name: "Admin", Email: "admin@host.xyz", Password: "admin", Roles: "admin" }
		_, err := bs.us.Create(admin)
		if err != nil {
			zap.S().Fatalf("Error during user (%s) initialisation", admin.Name)
		}
		zap.S().Infof("'%s' has been created successfully", admin.Name)

		user := models.User{ Name: "User", Email: "user@host.xyz", Password: "user", Roles: "user" }
		_, err = bs.us.Create(user)
		if err != nil {
			zap.S().Fatalf("Error during user (%s) initialisation", user.Name)
		}
		zap.S().Infof("'%s' has been created successfully", user.Name)
	}
}
