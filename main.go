package main

import (
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	a := App{}
	gotenv.Load()
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8080")
}
