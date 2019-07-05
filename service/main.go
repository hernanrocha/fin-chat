package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"

	"github.com/hernanrocha/fin-chat/rabbit"
	"github.com/hernanrocha/fin-chat/service/controller"
	_ "github.com/hernanrocha/fin-chat/service/docs"
	"github.com/hernanrocha/fin-chat/service/models"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

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

		mv := viewmodels.MessageView{
			ID:        message.ID,
			Text:      message.Text,
			RoomID:    message.RoomID,
			Username:  user.Username,
			CreatedAt: message.CreatedAt,
		}

		fmt.Println("Broadcasting message...")
		hub.BroadcastChan <- mv
	}
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

	// Run Chat Hub
	hub := controller.NewHub()
	hub.Run()

	// Run Bot Consumer
	go handleRabbitResponse(rb, hub)

	// Setup router
	r := controller.SetupRouter(rb, hub)
	r.Run() // listen and serve on 0.0.0.0:8001
}
