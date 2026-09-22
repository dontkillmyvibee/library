package helpers

import (
	"net/http"

	"github.com/dontkillmyvibee/library.git/handlers/errors"
	"github.com/dontkillmyvibee/library.git/schemas"
)

func InitError(w http.ResponseWriter, status int, err error) {
	errDTO := schemas.NewError(
		status,
		http.StatusText(status),
		err.Error(),
	)

	http.Error(w, errDTO.ToJSONString(), status)
}

func InternalServerError(w http.ResponseWriter) {
	InitError(w, http.StatusInternalServerError, errors.ErrInternalServer)
}
