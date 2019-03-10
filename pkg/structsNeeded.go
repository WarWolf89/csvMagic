package root

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcRes struct {
	IsValid bool
	Field   string `json:",omitempty"`
	ValErr  string `json:",omitempty"`
}

type PhoneNumber struct {
	ID       string   `csv:"id"`
	FileID   string   `bson:"file_id"`
	SmsPhone string   `csv:"sms_phone" bson:"sms_phone" validate:"custom"`
	ProcRes  *ProcRes `bson:"process_results,omitempty"`
}

type FileMeta struct {
	sync.Mutex `bson:"-"`
	ID         primitive.ObjectID `bson:"_id"  json:"file_id"`
	Name       string
	Counters   map[string]int64 `bson:"stats" json:"stats"`
	ExecTime   float64          `bson:"execution_time" json:"execution_time"`
	Errors     []string         `bson:"runtime_errors,omitempty" json:"-"`
}

func NewFileMeta(uuid primitive.ObjectID, name string) *FileMeta {
	return &FileMeta{ID: uuid, Name: name, Counters: make(map[string]int64), Errors: []string{}}
}

func (fm *FileMeta) IncreaseCounter(key string) {
	fm.Lock()
	defer fm.Unlock()
	if c, exists := fm.Counters[key]; !exists {
		fm.Counters[key] = 1
	} else {
		fm.Counters[key] = c + 1
	}
}
