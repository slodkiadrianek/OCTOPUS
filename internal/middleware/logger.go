package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	green = "\x1b[32m"
	reset = "\x1b[0m"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		actualDate := time.Now()
		logTime := actualDate.Format("2006-01-02 15:04:05")
		next.ServeHTTP(w, r)
		durationOfTheRoute := time.Since(start) / time.Millisecond
		formattedDurationOfTheRoute := strconv.FormatInt(int64(durationOfTheRoute), 10) + "ms"
		fmt.Println(green + "[INFO: " + logTime + "] " + r.Method + "-" + r.URL.Path + "-" + r.
			RemoteAddr + "-" + formattedDurationOfTheRoute + reset)
	})
}
