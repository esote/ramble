package http

import (
	"net/http"
)

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
