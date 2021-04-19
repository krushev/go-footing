package security

import (
	"encoding/json"
	jwt "github.com/LdDl/fiber-jwt/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/krushev/go-footing/models"
	"github.com/krushev/go-footing/services"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	for k, v := range buf {
		buf[k] = letterBytes[v%byte(len(letterBytes))]
	}
	return string(buf)
}

func InitAuth(us services.UserService) *jwt.FiberJWTMiddleware {
	identityKey := "login"
	authMiddleware, err := jwt.New(&jwt.FiberJWTMiddleware{
		Realm:            "fiber",
		Key:              []byte(randStringBytes(270)),
		Timeout:          time.Minute * time.Duration(viper.GetInt("jwt.timeout")),
		MaxRefresh:       time.Minute * time.Duration(viper.GetInt("jwt.maxRefresh")),
		IdentityKey:      identityKey,
		SigningAlgorithm: "HS256",
		PayloadFunc: func(userId interface{}) jwt.MapClaims {
			user, _ := us.GetByEmail(userId.(string))
			return jwt.MapClaims{
				"login": userId.(string),
				"ID":    user.ID,
				"name":  user.Name,
				"roles": user.Roles,
			}
		},
		IdentityHandler: func(ctx *fiber.Ctx) interface{} {
			claims := jwt.ExtractClaims(ctx)
			return &models.User{
				Email: claims["login"].(string),
				Name:  claims["name"].(string),
			}
		},
		Authenticator: func(ctx *fiber.Ctx) (interface{}, error) {
			loginData := login{}
			bodyBytes := ctx.Context().PostBody()
			if err := json.Unmarshal(bodyBytes, &loginData); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginData.Username
			password := loginData.Password
			user, err := us.GetByEmail(userID)
			if err != nil {
				return userID, jwt.ErrFailedAuthentication
			}
			err = us.VerifyPassword(user.Password, password)
			if err == nil && user.Active {
				return userID, nil
			}
			return userID, jwt.ErrFailedAuthentication
		},
		Authorizator: func(userId interface{}, ctx *fiber.Ctx) bool {
			user, err := us.GetByEmail(userId.(*models.User).Email)
			if err != nil {
				return false
			}
			if user.Active {
				return true
			}
			return false
		},
		Unauthorized: func(ctx *fiber.Ctx, code int, message string) error {
			if message == jwt.ErrFailedAuthentication.Error() {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"Error": string(ctx.Context().URI().Path()) + ";Unauthorized"})
			}
			return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"Error": string(ctx.Context().URI().Path()) + message})
		},
		RefreshResponse: func(ctx *fiber.Ctx, code int, token string, expire time.Time) error {
			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		zap.S().Error("Can not init auth")
		return nil
	}
	return authMiddleware
}
