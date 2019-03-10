package validutils

import (
	"fmt"

	root "../../pkg"
	validator "gopkg.in/go-playground/validator.v9"
)

type fixer = func(pn *root.PhoneNumber, err validator.FieldError, fm *root.FileMeta)

var (
	validate = validator.New()
	fixMap   = make(map[string]fixer)
)

func init() {
	validate.RegisterValidation("custom", validateFieldForSMSPhone)
	fixMap["SmsPhone"] = fixSmsPhoneStruct
}

func CheckAndFixStruct(pn *root.PhoneNumber, fm *root.FileMeta) {
	// The actual validate methods are the ones defined in the struct itself, those are the ones called here
	err := validate.Struct(pn)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fixMap[err.Field()](pn, err, fm)
		}
	}
	if err == nil {
		fm.IncreaseCounter("valid")
	}
}

func CheckAndFixSingleNumber(number string) *root.ProcRes {
	pr := &root.ProcRes{IsValid: true}
	if !isFieldValid(number, validateAreCharacters, validateLength) {
		fixSmsField(number, pr)
	}
	return pr
}

// The Validator Method for phone fields in structs
func validateFieldForSMSPhone(fl validator.FieldLevel) bool {
	return isFieldValid(fl.Field().String(), validateAreCharacters, validateLength)
}

func isFieldValid(field string, options ...func(string) bool) bool {
	isValid := true

	for i := 0; i < len(options) && isValid; i++ {
		isValid = options[i](field)
	}

	return isValid
}

func validateAreCharacters(field string) bool {
	isValid := true
	for _, ch := range field {
		if !(ch >= '0' && ch <= '9') {
			isValid = false
			break
		}
	}
	return isValid
}

func validateLength(field string) bool {
	return len(field) == 11
}
