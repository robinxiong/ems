package main

import (
	"ems/auth/auth_identity"

	"log"

	"ems/site/models/products"
	"ems/site/models/users"

	"fmt"

	"github.com/fatih/color"
)

var (
	AdminUser *users.User
	//Notification = notification.New(&notification.Config{})
	Tables = []interface{}{
		&auth_identity.AuthIdentity{},
		&users.User{},
		&products.Color{}, /*&models.Address{}, &models.Category{}, &models.Color{}, &models.Size{}, &models.Material{}, &models.Collection{},
		&models.Product{}, &models.ProductImage{}, &models.ColorVariation{}, &models.SizeVariation{},
		&models.Store{}, &models.Order{}, &models.OrderItem{}, &models.Setting{},
		&adminseo.MySEOSetting{},
		&models.Article{}, &models.MediaLibrary{},
		&banner_editor.QorBannerEditorSetting{},
		&asset_manager.AssetManager{},
		&i18n_database.Translation{},

		&notification.NotificationMessage{},
		&help.QorHelpEntry{},*/
	}
)

func main() {
	//Notification.RegisterChannel(database.New(&database.Config{db.DB}))  //创建notification_message表

	//TruncateTables(&products.Color{})
	createRecords()
}

func createRecords() {
	color.Green("Start create sample data...")

	createColors()
	fmt.Println("--> Created colors.")
}

//product color
func createColors() {
	for _, c := range Seeds.Colors {
		color := products.Color{}
		color.Name = c.Name
		color.Code = c.Code
		if err := DraftDB.Create(&color).Error; err != nil {
			log.Fatalf("create color (%v) failure, got err %v", color, err)
		}
	}
}
