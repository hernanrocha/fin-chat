package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/hernanrocha/fin-chat/service/hub"
)

func SetupRouter(hub hub.HubInterface) *gin.Engine {
	// Controllers
	c := NewRoomController()
	m := NewMessageController(hub)
	ws := NewWebSocketController(hub)
	auth := NewAuthController()
	authMiddleware, _ := auth.JWTMiddleware()

	// Default Engine
	r := gin.Default()

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Auth JWT
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/register", auth.Register)

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.Use(authMiddleware.MiddlewareFunc())

		// Refresh time can be longer than token timeout
		// auth.GET("/refresh_token", authMiddleware.RefreshHandler)

		v1.POST("/rooms", c.CreateRoom)
		v1.GET("/rooms", c.ListRooms)
		v1.GET("/rooms/:id", c.GetRoom)

		v1.GET("/rooms/:id/messages", m.ListRoomMessages)
		v1.POST("/rooms/:id/messages", m.CreateMessage)
	}

	// WebSocket
	r.GET("/ws", ws.WebSocket)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
