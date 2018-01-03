package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/nickduskey/api-demo/app"
	"github.com/subosito/gotenv"
)

// db is the root db instance
var db *gorm.DB
var err error

func main() {
	a := app.App{}
	gotenv.Load()
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("APP_DB_HOST"), os.Getenv("APP_DB_USERNAME"), os.Getenv("APP_DB_NAME"), os.Getenv("APP_DB_PASSWORD"))
	fmt.Println(connectionString)
	db, err = gorm.Open("postgres", connectionString)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	a.Initialize(db)

	a.Run(os.Getenv("APP_PORT"))
}
