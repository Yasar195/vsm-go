package db

import (
	"fmt"
	"os"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")

	if DB == nil {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// err = db.Migrator().DropTable(&schema.Users{})
		// if err != nil {
		// 	println("Warning: Could not drop Users table:", err.Error())
		// }

		// err = db.AutoMigrate(&schema.Users{}, &schema.Visitor{}, &schema.Visits{}, &schema.Notifications{})

		// if err != nil {
		// 	panic("failed to migrate database schema")
		// }

		go func() {
			var count int64
			result := db.Model(&schema.Users{}).Count(&count)

			if result.Error != nil {
				panic("failed to count Users table rows")
			}

			if count == 0 {
				println("No users found, inserting default user...")

				emailConfig := utility.EmailConfig{
					SMTPHost:     os.Getenv("SMTP_HOST"),
					SMTPPort:     587,
					SMTPUsername: os.Getenv("ADMIN_EMAIL"),
					SMTPPassword: os.Getenv("ADMIN_PASSWORD"),
					FromEmail:    os.Getenv("ADMIN_EMAIL"),
				}

				emailService := utility.NewEmailService(emailConfig)

				err := emailService.SendEmail(os.Getenv("ADMIN_EMAIL"), "admin created", "Hi\nNew admin create\n\nemail: imyasar07@gmail.com\npassword: admin@123")
				if err != nil {
					panic("error sending email")
				}

				hash, err := utility.HashPassword("admin@123")

				if err != nil {
					panic("hashing password failed")
				}

				defaultUser := schema.Users{
					Username:  "admin",
					UserEmail: "imyasar07@gmail.com",
					Password:  hash,
					UserType:  "staff",
				}

				err = db.Create(&defaultUser).Error

				if err != nil {
					panic("error creating default user")
				}
			}

			DB = db
			fmt.Println("✅ Database connected")
		}()
	}
}
