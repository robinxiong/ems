package main

import (
	"ems/test/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
)

// `User` belongs to `Profile`, `ProfileID` is the foreign key
type User struct {
	gorm.Model
	Profile   Profile `gorm:"save_associations:false"`
	ProfileID int
}

type Profile struct {
	gorm.Model
	Name string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := utils.TestDB()
	db.DropTableIfExists(&User{})
	db.DropTableIfExists(&Profile{})

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Profile{})

	user := User{
		Profile: Profile{
			Name: "test",
		},
	}

	db.Create(&user)
	var p Profile
	db.Model(&user).Association("Profile").Find(&p)
	scope := db.NewScope(&user)
	fmt.Printf("%v\n", scope)
	fmt.Println(scope.InstanceID())
}
