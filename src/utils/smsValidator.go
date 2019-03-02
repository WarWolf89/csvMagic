package utils

import (
	"fmt"
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate = validator.New()
)

func init() {
	validate.RegisterValidation("custom", validateFieldForSMSPhone)
}

func CheckAndFixStruct(pn *PhoneNumber) {
	// The actual validate methods are the ones defined in the struct itself, those are teh ones called here
	err := validate.Struct(pn)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fixVal(pn, err)
		}
	}
}

// The Validator Method for phone strings
func validateFieldForSMSPhone(fl validator.FieldLevel) bool {
	return setupAndEvalParams(fl.Field().String(), validateLettersAndLength)
}

func setupAndEvalParams(field string, options ...func(string) bool) bool {
	isValid := true

	for i := 0; i < len(options) && isValid; i++ {
		isValid = options[i](field)
	}

	return isValid
}

func validateLettersAndLength(field string) bool {
	re := regexp.MustCompile("\\D")

	if re.MatchString(field) {
		return false
	}
	if len(field) != 11 {
		return false
	}
	return true
}
