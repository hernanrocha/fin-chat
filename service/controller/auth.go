package controller

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hernanrocha/fin-chat/service/models"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

// AuthController ...
type AuthController struct {
	db *gorm.DB
}

// NewAuthController ...
func NewAuthController() *AuthController {
	return &AuthController{
		db: models.GetDB(),
	}
}

// Authenticate godoc
// @Summary Login
// @Description Login with Username and Password
// @Tags Authentication
// @Produce json
// @Param login body viewmodels.LoginRequest true "Login Credentials"
// @Success 200 {object} viewmodels.LoginResponse
// @Router /login [post]
func (c *AuthController) Authenticate(ctx *gin.Context) (interface{}, error) {
	var json viewmodels.LoginRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	var user models.User
	if err := c.db.Where("username = ?", json.Username).Find(&user).Error; err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	if json.Password != user.Password {
		return nil, jwt.ErrFailedAuthentication
	}

	return &viewmodels.UserView{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

// Register godoc
// @Summary Register User
// @Description Register User in database
// @Tags Authentication
// @Param user body viewmodels.RegisterRequest true "User Data"
// @Produce  json
// @Success 200 {object} viewmodels.RegisterResponse
// @Router /register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var json viewmodels.RegisterRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Username:  json.Username,
		Password:  json.Password,
		Email:     json.Email,
		FirstName: json.FirstName,
		LastName:  json.LastName,
	}

	if err := c.db.Create(user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &viewmodels.RegisterResponse{
		viewmodels.UserView{
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) JWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "username",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*viewmodels.UserView); ok {
				return jwt.MapClaims{
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			log.Println("IdentityHandler")
			claims := jwt.ExtractClaims(c)
			return &viewmodels.UserView{
				Username: claims["username"].(string),
			}
		},
		Authenticator: c.Authenticate,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			_, ok := data.(*viewmodels.UserView)
			return ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
		return nil, err
	}

	return authMiddleware, nil
}
