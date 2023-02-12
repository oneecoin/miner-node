package socket

import "net/http"

func InitServer() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":4000", nil)
}
