package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // need to load postgres
)

// App main application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize - init the app
func (a *App) Initialize(user, password, dbname string) {
	log.Printf("Initializing app...")
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run - run the app
func (a *App) Run(addr string) {
	log.Printf("App running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/products", a.GetProducts).Methods("GET")
	a.Router.HandleFunc("/products", a.CreateProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.GetProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.UpdateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.DeleteProduct).Methods("DELETE")
	a.Router.HandleFunc("/users", a.GetUsers).Methods("GET")
	a.Router.HandleFunc("/users", a.CreateUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.GetUser).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.UpdateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.DeleteUser).Methods("DELETE")
}
