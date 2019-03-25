package validutils

import (
	"testing"
	"../../pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

const (
	badNumCharsFix   = "27735405794_DELETED_148898755"
	badNumCharsUnFix = "21735405794_DELETED_148898755"
)

var (
	fm   = root.NewFileMeta( primitive.NewObjectID(), "test")
	pr   = root.ProcRes{}
	pnBC = root.PhoneNumber{ID: "12345", FileID: "test", SmsPhone: badNumCharsFix, ProcRes: &pr }
	pnB = root.PhoneNumber{ID: "12345", FileID: "test", SmsPhone: badNumCharsUnFix, ProcRes: &pr }
)

func TestFixSmsPhoneStructSucc(t *testing.T) {
	vErr := validate.Struct(pnBC)
	for _, err := range vErr.(validator.ValidationErrors) {
		fixSmsPhoneStruct(&pnBC,err,fm)
		t.Log(pnBC.ProcRes)
		if pnBC.ProcRes.IsValid != true || pnBC.ProcRes.Field != err.Field() {
			t.Errorf("Expected phone to be fixed, instead got %v for phone", pnBC.SmsPhone)
		}
	}
}

func TestFixSmsPhoneStructFail(t *testing.T) {
	vErr := validate.Struct(pnB)
	for _, err := range vErr.(validator.ValidationErrors) {
		fixSmsPhoneStruct(&pnB,err,fm)
		t.Log(err.Field())
		t.Log(pnB.ProcRes)
		if pnB.ProcRes.IsValid != false || pnB.ProcRes.Field != err.Field() {
			t.Errorf("Expected phone to be fail, instead got %v for phone", pnB.SmsPhone)
		}
	}
}

func TestRemoveChars(t *testing.T) {
	replacements := removeChars(badNumCharsFix)
	t.Log(replacements)
	if len(replacements) != 2 {
		t.Errorf("Expected two numbers valid number instead got %v", len(replacements))
	}
}

func TestFindValidNumber(t *testing.T) {
	replacements := removeChars(badNumCharsFix)
	t.Log(replacements)
	valNum, _ :=  findSingleValidNumber(replacements)
	t.Log(*valNum)
	if *valNum != "27735405794" {
		t.Errorf("Expected single valid number instead got %v", valNum)
	}
}