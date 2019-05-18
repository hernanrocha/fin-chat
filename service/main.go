package main

import (
	"fmt"

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

	// Ping API
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", c.Ping)
		v1.POST("/register", c.Register)

		v1.POST("/rooms", c.CreateRoom)
		v1.GET("/rooms", c.ListRooms)
		v1.GET("/rooms/:id", c.GetRoom)
		v1.GET("/rooms/:id/messages", c.ListRoomMessages)

		v1.POST("/messages", c.CreateMessage)
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
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
