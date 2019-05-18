package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/hernanrocha/fin-chat/service/controller"
	_ "github.com/hernanrocha/fin-chat/service/docs"
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

	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080
}
