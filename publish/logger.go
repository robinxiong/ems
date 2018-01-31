package publish

import (
	"log"
	"os"
	"github.com/jinzhu/gorm"
	"fmt"
	"reflect"
	"strings"
)

type LoggerInterface interface{
	Print(...interface{})
}


var Logger LoggerInterface

func init() {
	Logger = log.New(os.Stdout, "\r\n", 0)
}


func stringify(object interface{})string {

	if obj, ok := object.(interface{
		Stringify()string
	}); ok {
		return obj.Stringify()
	}

	//object为model struct 查找object中的Description
	scope := gorm.Scope{Value:object}
	for _, column := range []string{"Description", "Name", "Title", "Code"} {
		if field, ok := scope.FieldByName(column); ok {
			return fmt.Sprintf("%v", field.Field.Interface())
		}
	}

	if scope.PrimaryField() != nil {
		if scope.PrimaryKeyZero() {
			return ""
		}
		return fmt.Sprintf("%v#%v", scope.GetModelStruct().ModelType.Name(), scope.PrimaryKeyValue())
	}
	return fmt.Sprint(reflect.Indirect(reflect.ValueOf(object)).Interface())
}

//以 [1, "zh-cn"];[2, "en-us"]
func stringifyPrimaryValues(primaryValues [][][]interface{}, columns ...string) string {
	var values []string
	for _, primaryValue := range primaryValues {
		var primaryKeys []string
		for _, value := range primaryValue {
			if len(columns) == 0 {
				primaryKeys = append(primaryKeys, fmt.Sprint(value[1]))
			} else {
				for _, column := range columns {
					if column == fmt.Sprint(value[0]) {
						primaryKeys = append(primaryKeys, fmt.Sprint(value[1]))
					}
				}
			}
		}
		if len(primaryKeys) > 1 {
			values = append(values, fmt.Sprintf("[%v]", strings.Join(primaryKeys, ", ")))
		} else {
			values = append(values, primaryKeys...)
		}
	}
	return strings.Join(values, "; ")
}