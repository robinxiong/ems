package sorting

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
)

type sortingInterface interface {
	GetPosition() int
	SetPosition(int)
}
type sortingDescInterface interface {
	GetPosition() int
	SetPosition(int)
	SortingDesc()
}

// Sorting ascending mode
type Sorting struct {
	Position int `sql:"DEFAULT:NULL"`
}

// GetPosition get current position
func (position Sorting) GetPosition() int {
	return position.Position
}

// SetPosition set position, only set field value, won't save
func (position *Sorting) SetPosition(pos int) {
	position.Position = pos
}

// SortingDESC descending mode
type SortingDESC struct {
	Sorting
}

// SortingDesc make your model sorting desc by default
func (SortingDESC) SortingDesc() {}

func init() {
	//todo: add admin.RegisterViewPath
}

func newModel(value interface{}) interface{} {
	return reflect.New(reflect.Indirect(reflect.ValueOf(value)).Type()).Interface()
}

//传入当前操作的DB, 数据库行, 移动的们罩
func MoveUp(db *gorm.DB, value sortingInterface, pos int) error {
	return move(db, value, -pos)
}

func MoveDown(db *gorm.DB, value sortingInterface, pos int) error {
	return move(db, value, pos)
}

func MoveTo(db *gorm.DB, value sortingInterface, pos int) error {
	return move(db, value, pos-value.GetPosition())
}

//移动当前记录到指定位置，同时重新排序其它相同的类型的行(多个主键，但除去id主键， language_code)
func move(db *gorm.DB, value sortingInterface, pos int) (err error) {
	var startedTransaction bool
	var tx = db.Set("publish:publish_event", true)

	if t := tx.Begin(); t.Error == nil {
		startedTransaction = true
		tx = t
	}

	scope := db.NewScope(value)

	//比如language_code也是主键之一的时候，它并不是将id和language_code作为主键去查找，
	//而是排除掉id主键，查找language_code相面的列，然后重排所有带en-US行
	for _, field := range scope.PrimaryFields() {
		if field.DBName != "id" {
			tx = tx.Where(fmt.Sprintf("%s = ?", field.DBName), field.Field.Interface())
		}
	}

	currentPos := value.GetPosition()

	var results *gorm.DB

	if pos > 0 {
		//向下移动
		//当前位置之后，到currentPos+pos的记录，向上移动
		results = tx.Model(newModel(value)).
			Where("position > ? AND position <= ?", currentPos, currentPos+pos).
			UpdateColumn("position", gorm.Expr("position-?", 1))
	} else {
		results = tx.Model(newModel(value)).
			Where("position < ? AND position >= ?", currentPos, currentPos+pos).
			UpdateColumn("position", gorm.Expr("position + ?", 1))
	}


	if err = results.Error; err == nil {
		var rowsAffected = int(results.RowsAffected)
		if pos < 0 {
			rowsAffected = -rowsAffected
		}

		//更新当前行的位置, 如果之前的行没有变化，则不更新，如果其它的行发生变化，才更新当前行的位置
		value.SetPosition(currentPos + rowsAffected)
		err = tx.Model(value).UpdateColumn("position", gorm.Expr("position + ?", rowsAffected)).Error
	}

	// Create Publish Event
	createPublishEvent(tx, value)

	if startedTransaction {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
	return err
}
