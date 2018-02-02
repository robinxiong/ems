1. 在model struct中可以定义validator方法
2. validations调用govalidator
3. 在db中设置validations:skip_validations后，可以跳过validations
4. 将错误输出到scope.DB中，使之后的回调可以检查，并跳过