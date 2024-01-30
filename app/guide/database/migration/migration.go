package migrations

import (
	"app/database"
)

func RunMigrations() {
	db := database.GetDB()
	db.AutoMigrate(
		&database.Team{},
		&database.Task{},
		&database.CompletedTask{},
		&database.TeamTask{},
		&database.CombineTask{},
	)
}
