package other

import "net/http"

func GetJsonStatusMessage(code int, message string) interface{}{
	return struct {
		Status string `json:"status"`
		Message string `json:"message"`
	}{http.StatusText(code), message}
}
