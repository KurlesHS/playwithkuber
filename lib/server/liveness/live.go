package liveness

import "net/http"

func (l *Liveness) Live(w http.ResponseWriter, r *http.Request) {
	l.mu.RLock()
	isLive := l.isLive
	l.mu.RUnlock()
	if isLive {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("NOT OK"))
	}
}
