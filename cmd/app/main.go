package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"

	"github.com/Skavengerr/image-optimizer/delivery"
	"github.com/Skavengerr/image-optimizer/queue"
	"github.com/Skavengerr/image-optimizer/repository"
	"github.com/Skavengerr/image-optimizer/service"
)

var rabbitmqConn *amqp.Connection

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	rabbitMQURL := viper.GetString("rabbitmq.url")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	messageQueue, err := queue.NewRabbitMQ(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ message queue: %v", err)
	}

	r := mux.NewRouter()

	imageRepository := repository.NewImageRepository(conn)
	imageService := service.NewImageService(imageRepository, conn)
	imageHandler := delivery.NewImageHandler(imageService, messageQueue)

	r.HandleFunc("/upload", imageHandler.UploadImage)
	r.HandleFunc("/getImage/{id}", imageHandler.DownloadImage)

	http.ListenAndServe(":8080", r)
	log.Println("Server started on :8080")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
