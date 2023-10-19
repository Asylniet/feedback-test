package migrate

import (
	"fmt"
	"log"

	"github.com/enzhas/feedback_back/initializers"
	"github.com/enzhas/feedback_back/models"
)

func Migrate() {
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.Rating{},
		&models.Role{},
		&models.Subject{},
		&models.Attendance{},
		&models.Schedule{},
		&models.Organization{},
		&models.TotalRating{},
		&models.Todo{},
		&models.Vote{},
		&models.TodoChanges{},
		&models.Bonus{},
		&models.Achievement{},
	)
	//initializers.DB.Migrator().DropColumn(&models.Achievement{}, "organization_id")
	if err != nil {
		log.Fatal("Automigration failed")
	}
	roles := []*models.Role{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "sender"},
		{ID: 3, Name: "receiver"},
		{ID: 4, Name: "manager"},
	}
	for i, role := range roles {
		if initializers.DB.Model(&role).Where("id = ?", i+1).Updates(&role).RowsAffected == 0 {
			initializers.DB.Create(&role)
		}
	}

	newAdmin := models.User{
		Name:     "admin",
		Surname:  "main",
		Email:    "n_kurmash@kbtu.kz",
		Password: "$2a$10$hb.1azHRjq1OHlz/yrCdKeNx1AF9EYOrx34Kjn2u2WYMG97Wyl30a",
		RoleID:   1,
		Verified: true,
		Photo:    "default.png",
		Provider: "admin",
	}

	guestUser := models.User{
		Name:     "Анонимный",
		Surname:  "пользователь",
		Email:    "---",
		Password: "$2a$10$hb.1azHRjq1OHlz/yrCdKeNx1AF9EYOrx34Kjn2u2WYMG97Wyl30a",
		RoleID:   2,
		Verified: true,
		Photo:    "anonym.png",
		Provider: "guest",
	}

	if initializers.DB.Model(&newAdmin).Where("email = ?", "n_kurmash@kbtu.kz").Updates(&newAdmin).RowsAffected == 0 {
		initializers.DB.Create(&newAdmin)
	}

	if initializers.DB.Model(&guestUser).Where("id = ?", guestUser.ID).Updates(&guestUser).RowsAffected == 0 {
		initializers.DB.Create(&guestUser)
	}

	fmt.Println("👍 Migration complete")
}
