package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uuid"

	"github.com/dontkillmyvibee/library.git/handlers/errors"
	"github.com/dontkillmyvibee/library.git/schemas"
)

func DecodeJSONHelper(w http.ResponseWriter, r *http.Request, req any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(req); err != nil {
		errDto := schemas.NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())

		http.Error(w, errDto.ToJSONString(), http.StatusBadRequest)
		return false
	}

	return true
}

func EncodeJSONHelper(w http.ResponseWriter, status int, req any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(req); err != nil {
		fmt.Println("err:", err)
	}
}

func ParseUUIDPath(w http.ResponseWriter, r *http.Request, pathValue string) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue(pathValue))
	if err != nil {
		errDTO := schemas.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			errors.ErrInvalidUUID.Error(),
		)
		http.Error(w, errDTO.ToJSONString(), http.StatusBadRequest)
		return uuid.Nil(), false
	}

	return id, true
}
