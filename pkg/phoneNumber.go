package root

type ProcRes struct {
	Field  string `bson:"field"`
	ValErr string `bson:"valErr"`
}

type PhoneNumber struct {
	ID       string   `csv:"id" bson:"id"`
	SmsPhone string   `csv:"sms_phone" bson:"sms_phone" validate:"custom"`
	ProcRes  *ProcRes `bson:"process_results,omitempty"`
}
