package audit

import (
	"encoding/json"
	"os"
	"time"
)

// Entry is intentionally value-free. Adding a value field is forbidden (spec §7).
type Entry struct {
	Time      time.Time `json:"time"`
	Tool      string    `json:"tool"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Target    string    `json:"target"`
	Decision  string    `json:"decision"`
}

// Logger appends value-free audit entries to a jsonl file.
type Logger struct{ path string }

// New returns a Logger writing to path.
func New(path string) *Logger { return &Logger{path: path} }

// Record appends one entry, stamping the time if unset.
func (l *Logger) Record(e Entry) error {
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	return err
}
