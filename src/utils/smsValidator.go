package utils

import (
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

func ValidateFieldForSMSPhone(fl validator.FieldLevel) bool {
	fieldAsString := fl.Field().String()

	return setupAndEvalParams(fieldAsString, validateLength)
}

func setupAndEvalParams(field string, options ...func(string) bool) bool {
	isValid := true
	for _, option := range options {
		isValid = option(field)
	}
	return isValid
}

func validateLength(field string) bool {
	return len(field) == 11
}

func validateNoLetters(field string) bool {
	return regexp.MustCompile("\\D").MatchString(field)
}
