package controller

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hernanrocha/fin-chat/service/models"
)

// Controller example
type AuthController struct {
	db *gorm.DB
}

// NewController example
func NewAuthController() *AuthController {
	return &AuthController{
		db: models.GetDB(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserView struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginResponse struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

// Authenticate godoc
// @Summary Login
// @Description Login with Username and Password
// @Tags Authentication
// @Produce  json
// @Param login body controller.LoginRequest true "Login Credentials"
// @Success 200 {object} controller.LoginResponse
// @Router /login [post]
func (c *AuthController) Authenticate(ctx *gin.Context) (interface{}, error) {
	var json LoginRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	var user models.User
	if err := c.db.Where("username = ?", json.Username).Find(&user).Error; err != nil {
		return "", jwt.ErrFailedAuthentication
	}

	if json.Password == user.Password {
		return &UserView{
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}
