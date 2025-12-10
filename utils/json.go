package utils

import (
	"encoding/json"
)

func MarshalThis(input ...any) []byte {
	if data, err := json.Marshal(input); err != nil {
		return nil
	} else {
		return data
	}
}
