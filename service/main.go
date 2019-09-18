package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/hernanrocha/fin-chat/service/controller"
	_ "github.com/hernanrocha/fin-chat/service/docs"
	"github.com/hernanrocha/fin-chat/service/hub"
	"github.com/hernanrocha/fin-chat/service/models"
)

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
	log.Println("Starting web server...")
	os.Setenv("PORT", "8001")

	// Setup Postgres database
	dbconn := getEnv("DB_CONNECTION", "host=localhost port=15432 user=postgres password=postgres dbname=finchat sslmode=disable")
	db, err := gorm.Open("postgres", dbconn)
	failOnError(err, "Error conecting to database")
	defer db.Close()

	// Run migration
	models.Setup(db)

	/*
		// Setup RabbitMQ
		rabbitconn := getEnv("RABBIT_CONNECTION", "amqp://rabbitmq:rabbitmq@localhost:5672/")
		conn, err := amqp.Dial(rabbitconn)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		msg := messenger.NewRabbitMessenger(ch)
		failOnError(err, "Failed to declare a queue")
	*/

	// Run Messages Hub
	h := hub.NewHub()
	h.Run()

	/*
		// Add CmdMessageHandler
		handler, err := handler.NewCmdMessageHandler("cmd-rabbit", msg, h, models.GetDB())
		failOnError(err, "Error starting command message handler")
		h.AddClient(handler)

		// Run CmdResponse Consumer
		go msg.StartConsumer(handler.CmdResponseHandler)
	*/

	// Setup router
	r := controller.SetupRouter(h)
	err = r.Run() // listen and serve on 0.0.0.0:8001
	failOnError(err, "Failed starting server")
}
