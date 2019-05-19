package controller

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// Controller example
type AuthController struct {
}

// NewController example
func NewAuthController() *AuthController {
	return &AuthController{}
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

func (c *AuthController) Authenticate(ctx *gin.Context) (interface{}, error) {
	var json LoginRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	return &UserView{
		Username:  json.Username,
		Email:     "hernan@gmail.com",
		FirstName: "Hernan",
		LastName:  "Rocha",
	}, nil

	/*
		db := models.GetDB()


		var user models.User
		if err := db.Where("username = ?", json.Username).Find(&user).Error; err != nil {
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
	*/
}
