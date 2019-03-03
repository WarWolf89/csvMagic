package root

import (
	"sync"
)

type ProcRes struct {
	Field  string `bson:"field"`
	ValErr string `bson:"valErr"`
}

type PhoneNumber struct {
	ID       string   `csv:"id" bson:"_id"`
	FileID   string   `bson:"file_id"`
	SmsPhone string   `csv:"sms_phone" bson:"sms_phone" validate:"custom"`
	ProcRes  *ProcRes `bson:"process_results,omitempty"`
}

type FileMeta struct {
	sync.Mutex `bson:"-"`
	UUID       string           `bson:"_id"`
	Name       string           `bson:"name"`
	Counters   map[string]int64 `bson:"stats"`
	ExecTime   float64          `bson:"execution_time"`
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
