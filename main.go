package main

import (
	"github.com/shawn-10x/100pfps/db"
	"github.com/shawn-10x/100pfps/model"
	"github.com/shawn-10x/100pfps/router"
	"github.com/shawn-10x/100pfps/validators"
	"github.com/shawn-10x/100pfps/views"
)

func main() {
	db.Connect()
	model.SetupMigrations()
	e := router.New()
	validators.SetupValidators(e)
	views.SetupViews(e)
	e.Logger.Fatal(e.Start(":8080"))
}
