package interfaces

import "net/http"

type ServerController interface {
	GetServerInfo(w http.ResponseWriter, r *http.Request)
	GetServerMetrics(w http.ResponseWriter, r *http.Request)
}
