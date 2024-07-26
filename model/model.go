package model

import "github.com/shawn-10x/100pfps/db"

func SetupMigrations() {
	db.GetDB().AutoMigrate(
		&Profile{},
		&Tag{},
		&Ban{},
		&Admin{},
	)
}
