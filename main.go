package main

import (
	"os"

	"github.com/shawn-10x/100pfps/db"
	"github.com/shawn-10x/100pfps/middleware"
	"github.com/shawn-10x/100pfps/model"
	"github.com/shawn-10x/100pfps/router"
	"github.com/shawn-10x/100pfps/validators"
	"github.com/shawn-10x/100pfps/views"
)

func main() {
	db.Connect()
	model.SetupMigrations()

	admin := model.Admin{
		User:     os.Getenv("ADMIN_USER"),
		Password: os.Getenv("ADMIN_PASSWORD"),
		Role:     model.Owner,
	}
	if err := admin.CreateIfNotExists(); err != nil {
		panic(err)
	}

	e := router.New()
	middleware.SetupMiddlewares(e)
	validators.SetupValidators(e)
	views.SetupViews(e)
	e.Logger.Fatal(e.Start(":8080"))
}
