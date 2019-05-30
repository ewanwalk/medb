package v1

import (
	"encoder-backend/pkg/http/v1/socket"
	"net/http"
)

func webSocket(w http.ResponseWriter, r *http.Request) {
	socket.NewClient(w, r)
}
