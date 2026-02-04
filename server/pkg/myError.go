package pkg

import "errors"

// 自定义错误
var (
	ErrorPasswordWrong  = errors.New("the password is wrong")
	ErrorIdExists       = errors.New("the ID has already been used")
	ErrorIDNotExist     = errors.New("the ID does not exist")
	ErrorPasswordSimple = errors.New("the password is too simple")
	ErrorIDOnline       = errors.New("the ID is already online")
	ErrorData           = errors.New("the data is wrong")
	ErrorIllegalData    = errors.New("this data is illegal")
)
