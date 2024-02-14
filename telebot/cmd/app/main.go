package main

import (
	"context"
	"fmt"
	livenes_client "kurles/kuber/lib/client/liveness"
	"kurles/kuber/lib/logservermux"
	"kurles/kuber/lib/server/liveness"
	"kurles/kuber/telebot/internal/client/sayhello"
	"kurles/kuber/telebot/internal/client/tg"
	"kurles/kuber/telebot/internal/usecase/telebot"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	TOKEN_ENV_VAR                = "TOKEN"
	NOTIFICATION_CHAT_ID_ENV_VAR = "NOTIFICATION_CHAT_ID"
	HELLO_SERVICE_ADDR_ENV_VAR   = "HELLO_SERVICE_ADDR"
	PORT_ENV_VAR                 = "PORT"
)

const (
	defaultTgToken = "6989054031:AAEHMm2RalGlOdJThNHapXt3oB-3VE0I8hA"
)

func toInt(s string) int64 {
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return res
}

func getStrEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func buildNumber() string {
	version := "__build_number__"
	if strings.HasPrefix(version, "__build") {
		version = "undetermined"
	}
	return version
}

// chat with bot
// https://gist.github.com/2minchul/6d344a0f1f85ead1530803df2e4f9894
func main() {
	log.Println("start application version 1.0, build number", buildNumber())
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	port := os.Getenv(PORT_ENV_VAR)
	if len(port) == 0 {
		log.Fatalln("PORT environment variable is not set")
	}
	token := getStrEnv(TOKEN_ENV_VAR, defaultTgToken)
	notificationChatId := toInt(os.Getenv(NOTIFICATION_CHAT_ID_ENV_VAR))
	if notificationChatId == 0 {
		notificationChatId = -1001694212750
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err)
	}
	tb := tg.New(token, notificationChatId)
	go func() {
		if err := tb.Start(); err != nil {
			log.Fatalln(err)
		}
	}()
	helloServiceAddr := os.Getenv(HELLO_SERVICE_ADDR_ENV_VAR)
	sayHelloClient := sayhello.New(helloServiceAddr)
	mux := logservermux.New()
	livenes := liveness.New(mux)
	statusCode := 0
	gracefulCh := make(chan struct{})
	crashCh := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)

	lc := livenes_client.New(helloServiceAddr)

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

	/*


		pref := telebot.Settings{
			Token:  token,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		}
		b, err := telebot.NewBot(pref)
		if err != nil {
			log.Fatalln(err)
		}

		b.Handle("/start", func(m *telebot.Message) {
			b.Send(m.Chat, "Hello, "+m.Sender.FirstName+"!")
			<-time.After(4 * time.Second)

		})

		b.Handle(telebot.OnText, func(m *telebot.Message) {
			if m.Text == "hi" {
				t, _ := b.Send(m.Chat, "Hello, "+m.Sender.FirstName+"!")
				_ = t
				chat := telebot.Chat{
					ID:   notificationChatId,
					Type: telebot.ChatPrivate,
				}
				b.Send(&chat, "test chat")
			}
		})

		b.Handle(telebot.OnChannelPost, func(m *telebot.Message) {
			who := "Admin"
			if m.Sender != nil {
				who = m.Sender.Username
			}
			message := who + ", you entered " + m.Text

			log.Println(message)
			b.Send(m.Chat, message)
		})
		b.Start()
	*/
}
