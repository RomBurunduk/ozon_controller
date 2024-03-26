package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

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

	pvzRepo := postgresql.NewArticles(database)
	//implementation := ServerService{repo: pvzRepo}
	implementation := service.NewServerService(pvzRepo)

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

	http.Handle("/", createRouter(implementation))
	if err = http.ListenAndServe(os.Getenv("unsecurePort"), nil); err != nil {
		log.Fatal(err)
	}
}

func createRouter(implementation service.ServerService) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.BasicAuthMiddleware())
	router.Use(middlewares.Logging())

	router.HandleFunc("/pvz", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implementation.Create(w, req)
		case http.MethodGet:
			implementation.ListAll(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc(fmt.Sprintf("/pvz/{%s:[0-9]+}", service.QueryParamKey), func(w http.ResponseWriter, req *http.Request) {
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
	})
	return router
}
