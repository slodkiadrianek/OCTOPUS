package middleware

import (
	"fmt"
	"net/http"
	"time"
)

const (
	GREEN = "\x1b[32m"
	RESET = "\x1b[0m"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualDate := time.Now()
		logTime := actualDate.Format("2006-01-02 15:04:05")
		fmt.Println(GREEN + "[INFO: " + logTime + "] " + r.Method + "-" + r.URL.Path + "-" + r.RemoteAddr + RESET)
		next.ServeHTTP(w, r)
	})
}
