package validutils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"../../pkg"
	"gopkg.in/go-playground/validator.v9"
)

func fixSmsPhoneStruct(pn *root.PhoneNumber, err validator.FieldError, fm *root.FileMeta) {
	replacements := removeChars(fmt.Sprintf("%v", err.Value()))
	replacement, ferr := findSingleValidNumber(replacements)
	if ferr != nil {
		pn.ProcRes = &root.ProcRes{IsValid: false, Field: err.Field(), ValErr: ferr.Error()}
		fm.IncreaseCounter("unfixable")
		return
	}
	reflect.ValueOf(pn).Elem().FieldByName(err.Field()).SetString(*replacement)
	fm.IncreaseCounter("fixed")
	pn.ProcRes = &root.ProcRes{IsValid: true, Field: err.Field(), ValErr: fmt.Sprintf(succFix, err.Value(), *replacement)}

}

func fixSmsField(number string, pr *root.ProcRes) {
	replacements := removeChars(number)
	replacement, err := findSingleValidNumber(replacements)
	if err != nil {
		pr.ValErr = err.Error()
		pr.IsValid = false
	} else {
		pr.ValErr = fmt.Sprintf(succFix, number, replacement)
		pr.IsValid = true
	}
}

func findSingleValidNumber(replacements []string) (*string, error) {

	for _, repl := range replacements {
		if isFieldValid(repl, validateNumberFormat) {
			return &repl, nil
		}
	}
	return nil, errors.New("there was no valid number in the corrupted field")
}

func removeChars(val string) []string {
	return strings.FieldsFunc(val,splitHelper)
}

// helper function to remove all non-digit characters from the incoming string
func splitHelper (c rune) bool {
	return !unicode.IsNumber(c)
}
