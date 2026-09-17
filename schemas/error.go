package schemas

import (
	"encoding/json"
	"fmt"
	"time"
)

type ErrorSchema struct {
	Code        int       `json:"code"`
	Description string    `json:"description"`
	Message     string    `json:"message"`
	Time        time.Time `json:"time"`
}

func NewError(code int, description, message string) ErrorSchema {
	return ErrorSchema{
		Code:        code,
		Description: description,
		Message:     message,
		Time:        time.Now(),
	}
}

func (e ErrorSchema) ToJSONString() string {
	bytes, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		// тут надо бы как то логировать, но пока не знаю как
		fmt.Println(err)
	}

	return string(bytes)
}
