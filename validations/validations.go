package validations

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

//

// NewError generate a new error for a model's field
func NewError(resource interface{}, column, err string) error {
	return &Error{Resource: resource, Column: column, Message: err}
}

// Error用来保存校验的错误，它包含Resource（行), 列，错误信息
type Error struct {
	Resource interface{}
	Column   string
	Message  string
}

// Label 返回一个错误model, 主键, 列名
// Label is a label including model type, primary key and column name
func (err Error) Label() string {
	scope := gorm.Scope{Value: err.Resource}
	return fmt.Sprintf("%v_%v_%v", scope.GetModelStruct().ModelType.Name(), scope.PrimaryKeyValue(), err.Column)
}

// Error show error message
func (err Error) Error() string {
	return fmt.Sprintf("%v", err.Message)
}