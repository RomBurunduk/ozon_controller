package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"

	"pvz_controller/internal/app/receiver"
	"pvz_controller/internal/app/sender"
	"pvz_controller/internal/infrastructure/kafka"
	"pvz_controller/internal/pkg/db"
	"pvz_controller/internal/pkg/middlewares"
	"pvz_controller/internal/pkg/repository/postgresql"
	"pvz_controller/internal/service"
)

func ServiceWithDb() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pvzRepo := postgresql.NewPVZRepo(database)
	implementation := service.NewServerService(pvzRepo, database.GetPool(ctx))

	server := http.Server{
		Addr:    os.Getenv("securePort"),
		Handler: createRouter(implementation),
	}
	go func() {
		err = server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil {
			log.Fatal(err)
		}
	}()

	http.Handle("/", server.Handler)
	if err = http.ListenAndServe(os.Getenv("unsecurePort"), nil); err != nil {
		log.Fatal(err)
	}
}

func createRouter(implementation service.ServerInterface) *mux.Router {
	brokers := []string{
		"127.0.0.1:9091",
		"127.0.0.1:9092",
		"127.0.0.1:9093",
	}

	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		return nil
	}
	loggingService := sender.NewService(sender.NewKafkaSender(kafkaProducer, "logging"))
	router := mux.NewRouter()
	router.Use(middlewares.BasicAuthMiddleware())
	router.Use(middlewares.Logging(loggingService))
	consume(brokers)
	router.HandleFunc("/pvz", PVZHandler(implementation))

	router.HandleFunc(fmt.Sprintf("/pvz/{%s:[0-9]+}", service.QueryParamKey), PVZKeyHandler(implementation))
	return router
}

func PVZKeyHandler(implementation service.ServerInterface) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implementation.GetByID(w, req)
		case http.MethodDelete:
			implementation.DeleteById(w, req)
		case http.MethodPut:
			implementation.Update(w, req)
		default:
			fmt.Println("error")
		}
	}
}

func PVZHandler(implementation service.ServerInterface) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implementation.Create(w, req)
		case http.MethodGet:
			implementation.ListAll(w, req)
		case http.MethodDelete: // предполагается curl на адрес localhost/pvz с указанием массива id в флаге -d
			implementation.DeleteList(w, req)
		default:
			fmt.Println("error")
		}
	}
}

func consume(brokers []string) {
	consumer, err := kafka.NewConsumer(brokers)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers := map[string]receiver.HandleFunc{
		"logging": func(message *sarama.ConsumerMessage) {
			lm := sender.LoggingMessage{}
			err = json.Unmarshal(message.Value, &lm)
			if err != nil {
				fmt.Println("Consumer error", err)
				return
			}
			fmt.Println("Received Key: ", string(message.Key), " Value: ", lm)
		},
	}

	kafkaService := receiver.NewService(receiver.NewKafkaReceiver(consumer, handlers))

	err = kafkaService.StartConsume("logging")
	if err != nil {
		fmt.Println(err)
		return
	}
}
