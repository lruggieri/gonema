package visualResource

import (
	"encoding/json"
	"time"
)

type Movie struct{
	resource

	Duration time.Duration `json:"duration"`
}

func(m *Movie) String() string{
	jsonMovie, err := json.Marshal(m)
	if err != nil{
		return ""
	}
	return string(jsonMovie)
}