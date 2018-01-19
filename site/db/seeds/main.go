package main

import (
	"ems/auth/auth_identity"
	"ems/banner_editor"
	"ems/help"
	i18n_database "ems/i18n/backends/database"
	"ems/media/asset_manager"
	"ems/notification"
	"ems/site/app/models"
	"ems/site/config/admin"
	adminseo "ems/site/config/seo"

	"github.com/fatih/color"
	"log"
)

var (
	AdminUser    *models.User
	Notification = notification.New(&notification.Config{})
	Tables       = []interface{}{
		&auth_identity.AuthIdentity{},
		&models.User{}, &models.Address{}, &models.Category{}, &models.Color{}, &models.Size{}, &models.Material{}, &models.Collection{},
		&models.Product{}, &models.ProductImage{}, &models.ColorVariation{}, &models.SizeVariation{},
		&models.Store{}, &models.Order{}, &models.OrderItem{}, &models.Setting{},
		&adminseo.MySEOSetting{},
		&models.Article{}, &models.MediaLibrary{},
		&banner_editor.QorBannerEditorSetting{},
		&asset_manager.AssetManager{},
		&i18n_database.Translation{},

		&notification.NotificationMessage{},
		&admin.QorWidgetSetting{},
		&help.QorHelpEntry{},
	}
)

func main() {
	//Notification.RegisterChannel(database.New(&database.Config{db.DB}))  //创建notification_message表

	TruncateTables(Tables...)
	createRecords()
}

func createRecords() {
	color.Green("Start create sample data...")
	createSetting()
}

func createSetting(){
	setting := models.Setting{}
	setting.GiftWrappingFee = Seeds.Setting.GiftWrappingFee
	setting.CODFee = Seeds.Setting.CODFee
	setting.TaxRate = Seeds.Setting.TaxRate
	setting.Address = Seeds.Setting.Address
	setting.Region = Seeds.Setting.Region
	setting.City = Seeds.Setting.City
	setting.Country = Seeds.Setting.Country
	setting.Zip = Seeds.Setting.Zip
	setting.Latitude = Seeds.Setting.Latitude
	setting.Longitude = Seeds.Setting.Longitude

	if err := DraftDB.Create(&setting).Error; err != nil {
		log.Fatal("create setting (%v) failure, got err %v", setting, err)
	}
}
