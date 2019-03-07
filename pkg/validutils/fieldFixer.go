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
		pn.ProcRes = &root.ProcRes{IsValid: false, Field: err.Field(), ValErr: "After trimming number length was still more than 11."}
		fm.IncreaseCounter("unfixable")
		return
	}
	// If trimming helped replace with new value and add new result msg
	reflect.ValueOf(pn).Elem().FieldByName(err.Field()).SetString(replacement)
	fm.IncreaseCounter("fixed")
	pn.ProcRes = &root.ProcRes{IsValid: true, Field: err.Field(), ValErr: fmt.Sprintf("Original value %v was trimmed down to %v.", err.Value(), replacement)}
}

func fixSmsField(number string, pr *root.ProcRes) {
	replacement := fixTrimming(number)
	if !isFieldValid(replacement, validateLength) {
		pr.Field = "sms_phone"
		pr.ValErr = "After trimming number length was still more than 11."
		pr.IsValid = false
	}
	{
		pr.Field = "sms_phone"
		pr.ValErr = fmt.Sprintf("Original value %v was trimmed down to %v.", number, replacement)
		pr.IsValid = true
	}
}

func fixTrimming(val string) string {
	// Remove all letters from phonenumber
	return regexp.MustCompile("\\D").ReplaceAllString(val, "")
}
