package sorting

import (
	"ems/l10n"
	"ems/publish"
	"ems/test/utils"
	"fmt"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name    string
	Sorting //position
}

var db *gorm.DB
var pb *publish.Publish

func init() {
	db = utils.TestDB()
	RegisterCallbacks(db)
	l10n.RegisterCallbacks(db)

	pb = publish.New(db)
	if err := pb.ProductionDB().DropTableIfExists(&User{}, &Product{}, &Brand{}).Error; err != nil {
		panic(err)
	}
	//创建production表
	db.AutoMigrate(&User{}, &Product{}, &Brand{})
	//创建publish  _draft表
	pb.AutoMigrate(&Product{})
}

//删除测试，删除后还会数据库表进行重新排序, 如果是draft model还会创建一个publishEvent
func TestDeleteAndReorder(t *testing.T) {
	prepareUsers()
	if !(getUser("user1").GetPosition() == 1 && getUser("user2").GetPosition() == 2 && getUser("user3").GetPosition() == 3 && getUser("user4").GetPosition() == 4 && getUser("user5").GetPosition() == 5) {
		t.Errorf("user's order should be correct after create")
	}
	//删除user2
	user := getUser("user2")
	db.Delete(user)
	/*
			1
		deleted_at 2
			2
			3
			4
	*/
	if !checkPosition("user1", "user3", "user4", "user5") {
		t.Errorf("user2 is deleted, order should be correct")
	}
}

func TestMoveUpPosition(t *testing.T) {
	prepareUsers()
	MoveUp(db, getUser("user5"), 2)
	if !checkPosition("user1", "user2", "user5", "user3", "user4") {
		t.Errorf("user5 should be moved up")
	}
	MoveUp(db, getUser("user5"), 1)
	if !checkPosition("user1", "user5", "user2", "user3", "user4") {
		t.Errorf("user5's postion won't be changed because it is already the last")
	}

	MoveUp(db, getUser("user1"), 1)
	if !checkPosition("user1", "user5", "user2", "user3", "user4") {
		t.Errorf("user1's position won't be changed because it is already on the top")
	}

	MoveUp(db, getUser("user5"), 1)
	if !checkPosition("user5", "user1", "user2", "user3", "user4") {
		t.Errorf("user5 should be moved up")
	}

}

func TestMoveDownPosition(t *testing.T) {
	prepareUsers()
	MoveDown(db, getUser("user1"), 1)
	if !checkPosition("user2", "user1", "user3", "user4", "user5") {
		t.Errorf("user1 should be moved down")
	}

	MoveDown(db, getUser("user1"), 2)
	if !checkPosition("user2", "user3", "user4", "user1", "user5") {
		t.Errorf("user1 should be moved down")
	}

	MoveDown(db, getUser("user5"), 2)
	if !checkPosition("user2", "user3", "user4", "user1", "user5") {
		t.Errorf("user5's postion won't be changed because it is already the last")
	}

	MoveDown(db, getUser("user1"), 1)
	if !checkPosition("user2", "user3", "user4", "user5", "user1") {
		t.Errorf("user1 should be moved down")
	}
}

func TestMoveToPosition(t *testing.T) {
	prepareUsers()
	user := getUser("user5")

	MoveTo(db, user, user.GetPosition()-3)
	if !checkPosition("user1", "user5", "user2", "user3", "user4") {
		t.Errorf("user5 should be moved to position 2")
	}

	MoveTo(db, user, user.GetPosition()-1)
	if !checkPosition("user5", "user1", "user2", "user3", "user4") {
		t.Errorf("user5 should be moved to position 1")
	}
}

//准备测试数据
func prepareUsers() {
	db.Delete(&User{}) //安全删除全部user
	for i := 1; i <= 5; i++ {
		user := User{Name: fmt.Sprintf("user%v", i)}
		db.Save(&user)
	}
}

func getUser(name string) *User {
	var user User
	db.First(&user, "name = ?", name)
	return &user
}

func checkPosition(names ...string) bool {
	var users []User
	var positions []string

	db.Find(&users)
	for _, user := range users {
		positions = append(positions, user.Name)
	}

	//比较数据库返回的名字是否跟原来的一样
	if reflect.DeepEqual(positions, names) {
		return true
	} else {
		fmt.Printf("Expect %v, got %v\n", names, positions)
		return false
	}

}
