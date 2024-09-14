package main

import (
	"github.com/gadhittana-01/book-go/utils"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	config := utils.CheckAndSetConfig("./config", "app")
	DBpool := utils.ConnectDBPool(config.DBConnString)
	DB := utils.ConnectDB(config.DBConnString)

	if err := utils.RunMigrationPool(DB, config); err != nil {
		panic(err)
	}

	app, err := InitializeApp(r, DBpool, config)
	if err != nil {
		panic(err)
	}

	app.Start()
}
