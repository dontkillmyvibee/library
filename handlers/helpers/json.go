package helpers

import (
	"encoding/json"
	"net/http"

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
