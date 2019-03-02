package utils

import (
	"fmt"
	"reflect"
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

type PhoneNumber struct {
	ID       string `csv:"id"`
	SmsPhone string `csv:"sms_phone" validate:"custom"`
}

func FixVal(pn *PhoneNumber, err validator.FieldError) {

	val := fmt.Sprintf("%v", err.Value())
	replacement := regexp.MustCompile("\\D").ReplaceAllString(val, "")
	field := err.Field()
	// if len(replacement) != 11 {
	// 	fmt.Printf("couldn't fix error %s \n ", val)
	// 	return
	// }
	reflect.ValueOf(pn).Elem().FieldByName(field).SetString(replacement)

}
