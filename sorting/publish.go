package sorting

import (
	"ems/publish"
	"encoding/json"

	"github.com/jinzhu/gorm"
	"fmt"
	"strings"
	"errors"
)

//参考publish/event_test.go
//具体的使用参考callbacks, createPublishEvent
type changedSortingPublishEvent struct {
	Table       string
	PrimaryKeys []string
}
//更新_draft表的position到production表中
func (e changedSortingPublishEvent) Publish(db *gorm.DB, event publish.PublishEventInterface) (err error) {
	if event, ok := event.(*publish.PublishEvent); ok {
		scope := db.NewScope("")
		if err = json.Unmarshal([]byte(event.Argument), &e); err == nil {
			var conditions []string
			originalTable := scope.Quote(publish.OriginalTableName(e.Table))
			draftTable := scope.Quote(publish.DraftTableName(e.Table))
			for _, primaryKey := range e.PrimaryKeys {
				conditions = append(conditions, fmt.Sprintf("%v.%v = %v.%v", originalTable, primaryKey, draftTable, primaryKey))
				sql := fmt.Sprintf("UPDATE %v SET position = (select position FROM %v WHERE %v);", originalTable, draftTable, strings.Join(conditions, " AND "))
				return db.Exec(sql).Error
			}

		}
		return err
	}
	return errors.New("invalid publish event")
}

func (e changedSortingPublishEvent) Discard(db *gorm.DB, event publish.PublishEventInterface) (err error) {
	if event, ok := event.(*publish.PublishEvent); ok {
		scope := db.NewScope("")
		if err = json.Unmarshal([]byte(event.Argument), &e); err == nil {
			var conditions []string
			originalTable := scope.Quote(publish.OriginalTableName(e.Table))
			draftTable := scope.Quote(publish.DraftTableName(e.Table))
			for _, primaryKey := range e.PrimaryKeys {
				conditions = append(conditions, fmt.Sprintf("%v.%v = %v.%v", originalTable, primaryKey, draftTable, primaryKey))
			}
			sql := fmt.Sprintf("UPDATE %v SET position = (select position FROM %v WHERE %v);", draftTable, originalTable, strings.Join(conditions, " AND "))
			return db.Exec(sql).Error
		}
		return err
	}
	return errors.New("invalid publish event")
}

func init() {
	publish.RegisterEvent("changed_sorting", changedSortingPublishEvent{})
}