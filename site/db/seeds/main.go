package main

import (
	"ems/auth/auth_identity"
	"ems/site/app/models"

	"log"

	"github.com/fatih/color"
)

var (
	AdminUser *models.User
	//Notification = notification.New(&notification.Config{})
	Tables = []interface{}{
		&auth_identity.AuthIdentity{},
		&models.User{}, /*&models.Address{}, &models.Category{}, &models.Color{}, &models.Size{}, &models.Material{}, &models.Collection{},
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

	TruncateTables(Tables...)
	createRecords()
}

func createRecords() {
	color.Green("Start create sample data...")
	//createSetting()
}

func createSetting() {
	setting := models.Setting{}

	if err := DraftDB.Create(&setting).Error; err != nil {
		log.Fatal("create setting (%v) failure, got err %v", setting, err)
	}
}
