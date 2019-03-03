package validutils

import (
	"fmt"
	"reflect"
	"regexp"

	root "../../pkg"
	validator "gopkg.in/go-playground/validator.v9"
)

func fixVal(pn *root.PhoneNumber, err validator.FieldError) {
	fixTrimming(pn, err)
}
func fixTrimming(pn *root.PhoneNumber, err validator.FieldError) {
	val := fmt.Sprintf("%v", err.Value())
	// Remove all letters from phonenumber
	replacement := regexp.MustCompile("\\D").ReplaceAllString(val, "")
	// If the length is still not correct create appropriate error message
	if len(replacement) != 11 {
		procResCreate(pn, err.Field(), "After trimming number length was still more than 11.")
		fmt.Println(pn.ProcRes)
		return
	}
	// If trimming helped replace with new value and add new result msg
	reflect.ValueOf(pn).Elem().FieldByName(err.Field()).SetString(replacement)
	procResCreate(pn, err.Field(), fmt.Sprintf("Original value %v was trimmed down to %v.", err.Value(), replacement))
}
func procResCreate(pn *root.PhoneNumber, field string, valErr string) {
	pn.ProcRes = &root.ProcRes{Field: field, ValErr: valErr}
}
