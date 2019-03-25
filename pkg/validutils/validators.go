package validutils

import (
	"log"
	"regexp"

	"../../pkg"
	"gopkg.in/go-playground/validator.v9"
)

type fixer = func(pn *root.PhoneNumber, err validator.FieldError, fm *root.FileMeta)

var (
	rChar    = regexp.MustCompile("\\D")
	rNum     = regexp.MustCompile("^((?:27)|0)(=72|82|73|83|74|84)(\\d{7})$")
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
			log.Fatal(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			if fixer, ok := fixMap[err.Field()]; ok {
				fixer(pn, err, fm)
			}
		}
	}
	if err == nil {
		fm.IncreaseCounter("valid")
	}
}

func CheckAndFixSingleNumber(number string) *root.ProcRes {
	pr := &root.ProcRes{IsValid: true}
	if !isFieldValid(number, validateAreCharacters, validateNumberFormat) {
		fixSmsField(number, pr)
	}
	return pr
}

// The Validator Method for phone fields in structs. This is a wrapper method for generic struct level use
func validateFieldForSMSPhone(fl validator.FieldLevel) bool {
	return isFieldValid(fl.Field().String(), validateAreCharacters, validateNumberFormat)
}

func isFieldValid(field string, options ...func(string) bool) bool {
	isValid := true

	for i := 0; i < len(options) && isValid; i++ {
		isValid = options[i](field)
	}

	return isValid
}

func validateAreCharacters(field string) bool {
	// Need to reverse the value because if the match is true it means it has characters therefore it's not valid
	return !rChar.MatchString(field)
}

func validateNumberFormat(field string) bool {
	return rNum.MatchString(field)
}
