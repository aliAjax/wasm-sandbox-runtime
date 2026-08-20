package httptransport

import "net/http"

func probeHandler(ok func() bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if !ok() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		jsonWrite(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}
func versionHandler(version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		jsonWrite(w, http.StatusOK, map[string]string{"version": version})
	})
}
