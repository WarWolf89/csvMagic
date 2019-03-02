package utils

import (
	"fmt"
	"reflect"
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

type procRes struct {
	Field  string `bson:"field"`
	ValErr string `bson:"valErr"`
}

type PhoneNumber struct {
	ID       string   `csv:"id" bson:"id"`
	SmsPhone string   `csv:"sms_phone" bson:"sms_phone" validate:"custom"`
	ProcRes  *procRes `bson:"process_results,omitempty"`
}

func fixVal(pn *PhoneNumber, err validator.FieldError) {

	val := fmt.Sprintf("%v", err.Value())
	replacement := regexp.MustCompile("\\D").ReplaceAllString(val, "")
	field := err.Field()
	if len(replacement) != 11 {
		pn.ProcRes = &procRes{Field: err.Field(), ValErr: "After trimming number length was still more than 11"}
		fmt.Println(pn.ProcRes)
		return
	}
	reflect.ValueOf(pn).Elem().FieldByName(field).SetString(replacement)
}
