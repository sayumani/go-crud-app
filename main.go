package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sayooj/trivago/item"
	"github.com/sayooj/trivago/router"
	"github.com/sayooj/trivago/utils"
)

//Server struct
type Server struct {
	db *sql.DB
}

var server = Server{}

func init() {
	//loading environment variables
	e := godotenv.Load()
	if e != nil {
		log.Fatal("Error in loading env file ")
	}

	//setting up db connection
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSL_MODE")
	var err error
	server.db, err = utils.CreateDbConnection(host, port, user, password, dbname, sslmode)
	if errors.Is(err, utils.ErrDBConnectionError) {
		log.Fatal(err)
	}
}

func main() {
	server.runServer()
}

//runServer starts the server
func (s Server) runServer() {
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, initializeRoutes()))
}

func initializeRoutes() http.Handler {
	//logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	//repositories
	ir := item.NewItemsRepository(server.db)

	//usecases
	iu := item.NewItemsUseCase(ir)

	//handlers
	ih := item.NewItemsHandler(iu, log)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/", func(r chi.Router) {
		r.Mount("/item", router.ItemsRoutes(ih))
	})
	return r
}
