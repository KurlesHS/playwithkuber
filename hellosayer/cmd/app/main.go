package main

import (
	"context"
	"fmt"
	http_c "kurles/kuber/hellosayer/internal/controller/http"
	"kurles/kuber/hellosayer/usecase/hellosayer"
	"kurles/kuber/lib/server/liveness"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const PORT_ENV_VAR = "PORT"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := os.Getenv(PORT_ENV_VAR)
	if len(port) == 0 {
		log.Fatalln("PORT environment variable is not set")
	}

	fmt.Println("hello")
	mux := http.NewServeMux()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err)
	}
	livenes := liveness.New(mux)
	useCase := hellosayer.New()
	http_c.New(useCase, mux)

	wg := sync.WaitGroup{}
	wg.Add(1)
	statusCode := 0
	gracefulCh := make(chan struct{})
	crashCh := make(chan struct{})

	go func() {
		err := http.Serve(listener, mux)
		if err == http.ErrServerClosed {
			gracefulCh <- struct{}{}
		} else {
			crashCh <- struct{}{}
		}
	}()

	go func() {
		select {
		case <-livenes.CrashChannel():
			log.Println("crash simulation")
			statusCode = 1
		case <-ctx.Done():
			log.Println("graceful shutdown (by sigint or sigterm)")
		case <-gracefulCh: // graceful shutdown
			log.Println("graceful shutdown")
		case <-crashCh: // crash shutdown
			log.Println("crash shutdown (something bad with the http server)")
			statusCode = 1
		}
		livenes.GracefulShutdown()
		log.Println("graceful shutdown complete")
		wg.Done() // done
	}()

	wg.Wait()
	log.Println("exit")
	os.Exit(statusCode)
}
