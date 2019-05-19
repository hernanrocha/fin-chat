package main

import (
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/hernanrocha/fin-chat/service/controller"
	_ "github.com/hernanrocha/fin-chat/service/docs"
	"github.com/hernanrocha/fin-chat/service/models"
)

func setupRouter() *gin.Engine {

	// Default Engine with Logger and
	r := gin.Default()

	// Controller
	c := controller.NewController()

	// CORS
	r.Use(cors.Default())

	// Auth JWT
	authMiddleware, _ := jwtMiddleware()
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/register", c.Register)

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.Use(authMiddleware.MiddlewareFunc())

		// Refresh time can be longer than token timeout
		// auth.GET("/refresh_token", authMiddleware.RefreshHandler)

		v1.POST("/rooms", c.CreateRoom)
		v1.GET("/rooms", c.ListRooms)
		v1.GET("/rooms/:id", c.GetRoom)

		v1.GET("/rooms/:id/messages", c.ListRoomMessages)
		v1.POST("/rooms/:id/messages", c.CreateMessage)
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func jwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	c := controller.NewAuthController()

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "username",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*controller.UserView); ok {
				return jwt.MapClaims{
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("IdentityHandler")
			claims := jwt.ExtractClaims(c)
			return &controller.UserView{
				Username: claims["username"].(string),
			}
		},
		Authenticator: c.Authenticate,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*controller.UserView); ok {
				fmt.Println("AUTHORIZE " + v.Username)
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
		return nil, err
	}

	return authMiddleware, nil
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	fmt.Println("Hello World!")

	os.Setenv("PORT", "8001")

	// Setup database
	db, err := gorm.Open("postgres", "host=localhost port=15432 user=postgres password=postgres dbname=finchat sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Error conecting to database: %s", err))
	}
	defer db.Close()

	models.Setup(db)

	r := setupRouter()

	r.Run() // listen and serve on 0.0.0.0:8080
}
