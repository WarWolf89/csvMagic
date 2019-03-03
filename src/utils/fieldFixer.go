package utils

import (
	"fmt"
	"reflect"
	"regexp"

	root "../../pkg"
	validator "gopkg.in/go-playground/validator.v9"
)

func fixVal(pn *root.PhoneNumber, err validator.FieldError) {
	val := fmt.Sprintf("%v", err.Value())
	replacement := regexp.MustCompile("\\D").ReplaceAllString(val, "")
	field := err.Field()
	if len(replacement) != 11 {
		pn.ProcRes = &root.ProcRes{Field: err.Field(), ValErr: "After trimming number length was still more than 11"}
		fmt.Println(pn.ProcRes)
		return
	}
	reflect.ValueOf(pn).Elem().FieldByName(field).SetString(replacement)
}
