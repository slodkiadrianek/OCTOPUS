package interfaces

import "net/http"

type AppController interface {
	GetInfoAboutApps(w http.ResponseWriter, r *http.Request)
	GetInfoAboutApp(w http.ResponseWriter, r *http.Request)
	CreateApp(w http.ResponseWriter, r *http.Request)
	UpdateApp(w http.ResponseWriter, r *http.Request)
	DeleteApp(w http.ResponseWriter, r *http.Request)
	GetAppStatus(w http.ResponseWriter, r *http.Request)
}

type DockerController interface {
	PauseContainer(w http.ResponseWriter, r *http.Request)
	RestartContainer(w http.ResponseWriter, r *http.Request)
	StartContainer(w http.ResponseWriter, r *http.Request)
	UnpauseContainer(w http.ResponseWriter, r *http.Request)
	StopContainer(w http.ResponseWriter, r *http.Request)
	ImportDockerContainers(w http.ResponseWriter, r *http.Request)
}

type RouteController interface {
	AddWorkingRoutes(w http.ResponseWriter, r *http.Request)
}

type WsController interface {
	Logs(w http.ResponseWriter, r *http.Request)
}
