package mqr

import (
	"encoding/json"
	"time"
)

type Message struct {
	Id   string
	Time time.Time
	Body string //json body in string
}

type JsonBody struct {
	City  string `json:"city"`
	State string `json:"state"`
}

// MarshalBinary -
func (e *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalBinary -
func (e *Message) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	return nil
}
