package liveness

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (l *Liveness) Crash(w http.ResponseWriter, r *http.Request) {
	after := r.PathValue("after")
	d, err := strconv.Atoi(after)
	if err != nil || d < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("after must be an positive integer"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("crash planned after %d seconds", d)))

	l.wg.Add(1)
	go func() {

		select {
		case <-time.After(time.Duration(d) * time.Second):
			log.Printf("crashing after %d seconds", d)
			l.ch <- struct{}{} // simulate dead
		case <-l.finishCh:
			log.Print("graceful shutdown from crash handler")
		}
		l.wg.Done()
	}()
}
