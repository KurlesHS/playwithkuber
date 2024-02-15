package liveness

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (l *Liveness) SimulateDead(w http.ResponseWriter, r *http.Request) {
	after := r.PathValue("after")
	d, err := strconv.Atoi(after)
	if err != nil || d < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("after must be an positive integer"))
		return
	}

	hostname, _ := os.Hostname()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("dead simulate after %d seconds (pod %s)", d, hostname)))

	go func() {
		select {
		case <-time.After(time.Duration(d) * time.Second):
			l.mu.Lock()
			l.isLive = false
			l.mu.Unlock()
		case <-l.finishCh:
			log.Print("graceful shutdown from crash handler")
			return
		}
	}()

}
