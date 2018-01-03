package app

import (
	"fmt"
	"log"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // bringing in postgres dialect for Gorm
)

// App main application
type App struct {
	Router    *mux.Router
	APIRouter *mux.Router
	DB        *gorm.DB
}

// Initialize - init the app
func (a *App) Initialize(host, user, password, dbname string) {
	log.Printf("Initializing app...")
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)

	var err error
	a.DB, err = gorm.Open("postgres", connectionString)
	defer a.DB.Close()
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.APIRouter = mux.NewRouter()
	a.initializeAPIRoutes()
}

// Run - run the app
func (a *App) Run(addr string) {
	log.Printf("App running on port %s", addr)
	n := negroni.Classic()

	n.UseHandler(a.Router)
	n.Run(addr)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/status", a.StatusHandler).Methods("GET")
	a.Router.HandleFunc("/api/login", a.NotImplementedHandler).Methods("POST")
	a.Router.HandleFunc("/api/logout", a.NotImplementedHandler).Methods("GET")
}

func (a *App) initializeAPIRoutes() {
	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	a.APIRouter.HandleFunc("/api/products", a.GetProducts).Methods("GET")
	a.APIRouter.HandleFunc("/api/products", a.CreateProduct).Methods("POST")
	a.APIRouter.HandleFunc("/api/product/{id:[0-9]+}", a.GetProduct).Methods("GET")
	a.APIRouter.HandleFunc("/api/product/{id:[0-9]+}", a.UpdateProduct).Methods("PUT")
	a.APIRouter.HandleFunc("/api/product/{id:[0-9]+}", a.DeleteProduct).Methods("DELETE")
	a.APIRouter.HandleFunc("/api/users", a.GetUsers).Methods("GET")
	a.APIRouter.HandleFunc("/api/users", a.CreateUser).Methods("POST")
	a.APIRouter.HandleFunc("/api/user/{id:[0-9]+}", a.GetUser).Methods("GET")
	a.APIRouter.HandleFunc("/api/user/{id:[0-9]+}", a.UpdateUser).Methods("PUT")
	a.APIRouter.HandleFunc("/api/user/{id:[0-9]+}", a.DeleteUser).Methods("DELETE")

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(a.APIRouter))
	a.Router.PathPrefix("/api").Handler(an)
}
