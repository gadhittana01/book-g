package main

import (
	"github.com/gadhittana01/go-modules/utils"
	"github.com/go-chi/chi"
)

// @title           Book Service API
// @version         1.0x
// @description     API Spec for Book Service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /
// @schemes http https

// @securityDefinitions.apikey  authorization
// @in header
// @name Authorization
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
