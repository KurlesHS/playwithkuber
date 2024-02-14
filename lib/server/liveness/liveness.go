package liveness

import (
	"log"
	"net/http"
	"sync"
)

type HandlerRegistrar interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type Liveness struct {
	mux      HandlerRegistrar
	ch       chan struct{}
	isLive   bool
	mu       sync.RWMutex
	finishCh chan struct{}
	wg       sync.WaitGroup
}

func New(mux HandlerRegistrar) *Liveness {
	res := &Liveness{
		mux:      mux,
		ch:       make(chan struct{}),
		isLive:   true,
		mu:       sync.RWMutex{},
		finishCh: make(chan struct{}),
	}
	mux.HandleFunc("GET /crash/{after}", res.Crash)
	mux.HandleFunc("GET /simulate/{after}", res.SimulateDead)
	mux.HandleFunc("GET /live", res.Live)
	return res
}

func (l *Liveness) GracefulShutdown() {
	if l.finishCh != nil {
		close(l.finishCh)
		l.finishCh = nil
		log.Println("graceful shutdown already in progress")
	}
	l.wg.Wait()

}

func (l *Liveness) CrashChannel() chan struct{} {
	return l.ch
}
