package validutils

import "testing"
const (validNum = "27735405794" ; faultyWithText = "_DELETED_27735405799")

func TestNumberFormatFaulty(t *testing.T) {
	v := validateNumberFormat("21735405794")
	if v != false {
		t.Error("Expected number to fail, instead passed")
	}
}

func TestNumberFormatGood(t *testing.T) {

	v := validateNumberFormat(validNum)
	if v != true {
		t.Error("Expected number to pass, instead failed")
	}
}

func TestIsValidOnOnlyNumbers(t *testing.T) {
	v := validateNumberFormat(validNum)
	if v != true {
		t.Error("Expected number to pass, instead it failed")
	}
}

func TestIsValidOnCharacters(t *testing.T) {
	v := validateNumberFormat(faultyWithText)
	if v != false {
		t.Error("Expected number to fail, instead it passed")
	}
}

func TestWithMultipleArgsValid(t *testing.T) {
	v := isFieldValid(validNum,validateAreCharacters,validateNumberFormat)
	if v != true {
		t.Error("Expected number to pass for all args, instead it failed")
	}
}

func TestWithMultipleArgsFaulty(t *testing.T) {
	v := isFieldValid(faultyWithText,validateAreCharacters,validateNumberFormat)
	if v != false {
		t.Error("Expected number to fail for all args, instead it passed")
	}
}
