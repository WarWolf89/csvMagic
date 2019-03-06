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
	if len(replacement) != 11 {
		pn.procResCreate(err.Field(), "After trimming number length was still more than 11.")
		fm.IncreaseCounter("unfixable")
		return
	}
	// If trimming helped replace with new value and add new result msg
	reflect.ValueOf(pn).Elem().FieldByName(err.Field()).SetString(replacement)
	fm.IncreaseCounter("fixed")
	pn.procResCreate(err.Field(), fmt.Sprintf("Original value %v was trimmed down to %v.", err.Value(), replacement))
}
func (pn *root.PhoneNumber) procResCreate(field string, valErr string) {
	pn.ProcRes = &root.ProcRes{Field: field, ValErr: valErr}
}

func fixTrimming(val string) string {
	// Remove all letters from phonenumber
	return regexp.MustCompile("\\D").ReplaceAllString(val, "")
}
