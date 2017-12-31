package main

import (
	"os"

	"github.com/nickduskey/api-demo/app"
	"github.com/subosito/gotenv"
)

func main() {
	a := app.App{}
	gotenv.Load()
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(os.Getenv("APP_PORT"))
}
