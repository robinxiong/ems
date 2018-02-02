package validations

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

var skipValidations = "validations:skip_validations"

func validate(scope *gorm.Scope) {
	//是否调用的UpdateColumns方法, 更新指定的列，如果不是，则调用校验
	//如果db中设置了validations:skip_validations, 则跳过所以的在此db上的操作
	if _, ok := scope.Get("gorm:update_column"); !ok {
		if result, ok := scope.DB().Get(skipValidations); !(ok && result.(bool)) {
			if !scope.HasError() {
				scope.CallMethod("Validate")  //CallMethod调用的方法，返回了error，则保存到scope.Err中, 详细的可以参考 gorm/scope.go callMethod
				if scope.Value != nil {
					resource := scope.IndirectValue().Interface()
					//调用ValidateStruct进行校验, 反回一个error, 它可能是单个error, 也可能是govalidator.Errors (实现了Error() string方法)
					_, validatorErrors := govalidator.ValidateStruct(resource)
					if validatorErrors != nil {
						if errors, ok := validatorErrors.(govalidator.Errors); ok {
							for _, err := range flatValidatorErrors(errors) {
								scope.DB().AddError(formattedError(err, resource))
							}
						} else {
							scope.DB().AddError(validatorErrors)
						}
					}
				}
			}
		}
	}
}

func formattedError(err govalidator.Error, resource interface{}) error {
	message := err.Error()
	attrName := err.Name
	if strings.Index(message, "non zero value required") >= 0 {
		message = fmt.Sprintf("%v can't be blank", attrName)
	} else if strings.Index(message, "as length") >= 0 {
		reg, _ := regexp.Compile(`\(([0-9]+)\|([0-9]+)\)`)
		submatch := reg.FindSubmatch([]byte(err.Error()))
		message = fmt.Sprintf("%v is the wrong length (should be %v~%v characters)", attrName, string(submatch[1]), string(submatch[2]))
	} else if strings.Index(message, "as numeric") >= 0 {
		message = fmt.Sprintf("%v is not a number", attrName)
	} else if strings.Index(message, "as email") >= 0 {
		message = fmt.Sprintf("%v is not a valid email address", attrName)
	}
	return NewError(resource, attrName, message)

}

func flatValidatorErrors(validatorErrors govalidator.Errors) []govalidator.Error {
	resultErrors := []govalidator.Error{}
	for _, validatorError := range validatorErrors.Errors() {
		if errors, ok := validatorError.(govalidator.Errors); ok {
			for _, e := range errors {
				resultErrors = append(resultErrors, e.(govalidator.Error))
			}
		}
		if e, ok := validatorError.(govalidator.Error); ok {
			resultErrors = append(resultErrors, e)
		}
	}
	return resultErrors
}

func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()

	callback.Create().Before("gorm:before_create").Register("validations:validate", validate)
}
