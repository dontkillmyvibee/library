package schemas

import (
	"encoding/json"
	"fmt"
	"time"
)

type Error struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewError(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
		Time:    time.Now(),
	}
}

func (e Error) ToJSONString() string {
	bytes, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		// тут надо бы как то логировать, но пока не знаю как
		fmt.Println(err)
	}

	return string(bytes)
}
