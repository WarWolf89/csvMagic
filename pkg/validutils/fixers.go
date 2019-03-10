package validutils

import (
	"fmt"
	"reflect"
	"regexp"

	root "../../pkg"
	validator "gopkg.in/go-playground/validator.v9"
)

func fixSmsPhoneStruct(pn *root.PhoneNumber, err validator.FieldError, fm *root.FileMeta) {
	replacement := fixTrimming(fmt.Sprintf("%v", err.Value()))
	// If the length is still not correct create appropriate error message
	if !isFieldValid(replacement, validateLength) {
		pn.ProcRes = &root.ProcRes{IsValid: false, Field: err.Field(), ValErr: lengthError}
		fm.IncreaseCounter("unfixable")
		return
	}
	// If trimming helped replace with new value and add new result msg
	reflect.ValueOf(pn).Elem().FieldByName(err.Field()).SetString(replacement)
	fm.IncreaseCounter("fixed")
	pn.ProcRes = &root.ProcRes{IsValid: true, Field: err.Field(), ValErr: fmt.Sprintf(succFix, err.Value(), replacement)}
}

func fixSmsField(number string, pr *root.ProcRes) {
	replacement := fixTrimming(number)
	fmt.Println(replacement)
	if !isFieldValid(replacement, validateLength) {
		pr.ValErr = lengthError
		pr.IsValid = false
	} else {
		pr.ValErr = fmt.Sprintf(succFix, number, replacement)
		pr.IsValid = true
	}
}

func fixTrimming(val string) string {
	// Remove all letters from phonenumber
	return regexp.MustCompile("\\D").ReplaceAllString(val, "")
}
