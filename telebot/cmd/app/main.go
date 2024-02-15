package main

import (
	"context"
	"fmt"
	livenes_client "kurles/kuber/lib/client/liveness"
	"kurles/kuber/lib/logservermux"
	"kurles/kuber/lib/server/liveness"
	"kurles/kuber/telebot/internal/client/sayhello"
	"kurles/kuber/telebot/internal/client/tg"
	"kurles/kuber/telebot/internal/config"
	"kurles/kuber/telebot/internal/usecase/telebot"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.Println("start application version 1.0, build number", config.BuildNumber())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalln(err)
	}
	tb := tg.New(cfg.TgToken, cfg.NotificationChatId)
	go func() {
		if err := tb.Start(); err != nil {
			log.Fatalln(err)
		}
	}()

	sayHelloClient := sayhello.New(cfg.HelloServiceAddr)
	mux := logservermux.New()
	livenes := liveness.New(mux)
	statusCode := 0
	gracefulCh := make(chan struct{})
	crashCh := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)

	lc := livenes_client.New(cfg.HelloServiceAddr)

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
	useCase := telebot.New(ctx, tb, lc, sayHelloClient)
	useCase.Start(ctx)
	wg.Wait()
	log.Println("exit")
	os.Exit(statusCode)
}
