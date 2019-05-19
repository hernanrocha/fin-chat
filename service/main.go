package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/hernanrocha/fin-chat/rabbit"
	"github.com/hernanrocha/fin-chat/service/controller"
	_ "github.com/hernanrocha/fin-chat/service/docs"
	"github.com/hernanrocha/fin-chat/service/models"
)

func setupRouter(rb rabbit.RabbitChannel) *gin.Engine {

	// Default Engine with Logger and
	r := gin.Default()

	// Run Chat Hub
	hub := controller.NewHub()
	go hub.Run()

	// Run Bot Consumer
	go handleRabbitResponse(rb, hub)

	// Controller
	c := controller.NewController(hub, rb)

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

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

	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, hub)
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func handleRabbitResponse(rb rabbit.RabbitChannel, hub *controller.Hub) {
	db := models.GetDB()

	user := models.User{
		Username: "Bot",
		Email:    "bot@mail.com",
	}
	if err := db.Where("username = ?", "Bot").First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println("Creating bot user")
			db.Create(&user)
		}
	}

	msgs, err := rb.Consume()
	if err != nil {
		return
	}

	log.Printf("Waiting for messages from rabbit...")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		room, _ := strconv.Atoi(d.CorrelationId)
		message := &models.Message{
			Text:   string(d.Body),
			RoomID: uint(room),
			UserID: user.ID,
		}

		if err := db.Create(message).Error; err != nil {
			log.Println("Error creating new message from bot: ", err)
			continue
		}

		mv := controller.MessageView{
			ID:        message.ID,
			Text:      message.Text,
			RoomID:    message.RoomID,
			Username:  user.Username,
			CreatedAt: message.CreatedAt,
		}

		fmt.Println("Broadcasting message...")
		hub.BroadcastChan <- mv
	}

	/*
		if err := db.Where("username = ?", userView.(*UserView).Username).Find(&user).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	*/
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wshandler(w http.ResponseWriter, r *http.Request, h *controller.Hub) {
	fmt.Println("NEW WEBSOCKET")
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	h.AddClientChan <- conn

	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			h.RemoveClientChan <- conn
			return
		}
	}
}

// @title Swagger FinChat API
// @version 1.0
// @description This is a simple bot-based chat.

// @contact.name Hernan Rocha
// @contact.email hernanrocha93(at)gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8001
// @BasePath /

func main() {
	fmt.Println("Starting web server...")
	os.Setenv("PORT", "8001")

	// Setup Postgres database
	dbconn := getEnv("DB_CONNECTION", "host=localhost port=15432 user=postgres password=postgres dbname=finchat sslmode=disable")
	db, err := gorm.Open("postgres", dbconn)
	failOnError(err, "Error conecting to database")
	defer db.Close()

	// Run migration
	models.Setup(db)

	// Setup RabbitMQ
	rabbitconn := getEnv("RABBIT_CONNECTION", "amqp://rabbitmq:rabbitmq@localhost:5672/")
	conn, err := amqp.Dial(rabbitconn)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	rb := rabbit.NewRabbitChannel(ch)
	err = rb.QueueDeclare()
	failOnError(err, "Failed to declare a queue")

	// Setup router
	r := setupRouter(rb)
	r.Run() // listen and serve on 0.0.0.0:8001
}
