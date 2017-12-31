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
	"github.com/nickduskey/api-demo/auth"
)

// App main application
type App struct {
	Router    *mux.Router
	ApiRouter *mux.Router
	DB        *sql.DB
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

	jwtErr := auth.InitAuth()
	if jwtErr != nil {
		log.Println("Error initializing the JWT's!")
		log.Fatal(jwtErr)
	}

	a.Router = mux.NewRouter()
	a.ApiRouter = mux.NewRouter()
	a.initializeRoutes()
	a.initializeApiRoutes()
}

// Run - run the app
func (a *App) Run(addr string) {
	log.Printf("App running on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/login", a.LoginUser).Methods("POST")
	a.Router.HandleFunc("/logout", a.LogoutUser).Methods("GET")
}

func (a *App) initializeApiRoutes() {
	a.Router.HandleFunc("/api/v0/products", a.GetProducts).Methods("GET")
	a.Router.HandleFunc("/api/v0//products", a.CreateProduct).Methods("POST")
	a.Router.HandleFunc("/api/v0//product/{id:[0-9]+}", a.GetProduct).Methods("GET")
	a.Router.HandleFunc("/api/v0//product/{id:[0-9]+}", a.UpdateProduct).Methods("PUT")
	a.Router.HandleFunc("/api/v0//product/{id:[0-9]+}", a.DeleteProduct).Methods("DELETE")
	a.Router.HandleFunc("/api/v0//users", a.GetUsers).Methods("GET")
	a.Router.HandleFunc("/api/v0//users", a.CreateUser).Methods("POST")
	a.Router.HandleFunc("/api/v0//user/{id:[0-9]+}", a.GetUser).Methods("GET")
	a.Router.HandleFunc("/api/v0//user/{id:[0-9]+}", a.UpdateUser).Methods("PUT")
	a.Router.HandleFunc("/api/v0//user/{id:[0-9]+}", a.DeleteUser).Methods("DELETE")
}
