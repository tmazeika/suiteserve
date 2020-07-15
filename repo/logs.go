package repo

type LogLevelType string

const (
	LogLevelTypeTrace LogLevelType = "trace"
	LogLevelTypeDebug LogLevelType = "debug"
	LogLevelTypeInfo  LogLevelType = "info"
	LogLevelTypeWarn  LogLevelType = "warn"
	LogLevelTypeError LogLevelType = "error"
)

type UnsavedLogLine struct {
	Case      string       `json:"case"`
	Index     int64        `json:"index"`
	Level     LogLevelType `json:"level"`
	Trace     string       `json:"trace,omitempty" bson:",omitempty"`
	Message   string       `json:"message,omitempty" bson:",omitempty"`
	Timestamp int64        `json:"timestamp"`
}

type LogLine struct {
	SavedEntity    `bson:",inline"`
	UnsavedLogLine `bson:",inline"`
}

type LogPage struct {
	NextId *string   `json:"next_id" bson:"next_id,omitempty"`
	Lines  []LogLine `json:"lines,omitempty" bson:",omitempty"`
}
