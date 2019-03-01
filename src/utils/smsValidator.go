package utils

import (
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

// The Validator Method for phone strings
func ValidateFieldForSMSPhone(fl validator.FieldLevel) bool {
	return setupAndEvalParams(fl.Field().String(), validateLength, validateNoLetters)
}

func setupAndEvalParams(field string, options ...func(string) bool) bool {
	isValid := true

	for i := 0; i < len(options) && isValid; i++ {
		isValid = options[i](field)
	}

	return isValid
}

func validateLength(field string) bool {
	return len(field) == 11
}

// turning around boolean value is never a good idea
func validateNoLetters(field string) bool {
	return !regexp.MustCompile("\\D").MatchString(field)
}
