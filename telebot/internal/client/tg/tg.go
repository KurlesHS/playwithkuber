package tg

import (
	"context"
	"kurles/kuber/telebot/internal/model"
	"sync"
	"time"

	"github.com/tucnak/telebot"
)

type Tg struct {
	token  string
	chatId int64
	b      *telebot.Bot
	mu     sync.RWMutex
	chs    []chan model.Message
}

func New(token string, chatId int64) *Tg {
	return &Tg{token: token, chatId: chatId}
}

func (tg *Tg) handler(m *telebot.Message) {
	who := "Admin"
	if m.Sender != nil {
		who = m.Sender.Username
	}
	tg.informAboutMessage(m.Text, who)
}

func (tg *Tg) SendMessage(ctx context.Context, text string) error {
	chat := telebot.Chat{
		ID:   tg.chatId,
		Type: telebot.ChatPrivate,
	}
	_, err := tg.b.Send(&chat, text)
	return err
}

func (tg *Tg) Stop() {
	if tg.b != nil {
		tg.b.Stop()
		tg.b = nil
	}
	tg.mu.Lock()
	defer tg.mu.Unlock()
	for _, ch := range tg.chs {
		close(ch)
	}
	tg.chs = nil
}

func (tg *Tg) GetMessageChannel(ctx context.Context) <-chan model.Message {
	ch := make(chan model.Message, 1)
	tg.mu.Lock()
	defer tg.mu.Unlock()
	tg.chs = append(tg.chs, ch)
	return ch
}

func (tg *Tg) informAboutMessage(message, username string) {
	go func() {
		tg.mu.RLock()
		defer tg.mu.RUnlock()
		m := model.Message{
			Message: message,
			Sender:  username,
		}
		for _, ch := range tg.chs {
			ch <- m
		}
	}()
}

func (tg *Tg) Start() error {
	pref := telebot.Settings{
		Token:  tg.token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	var err error
	tg.b, err = telebot.NewBot(pref)
	if err != nil {
		return err
	}
	tg.b.Handle(telebot.OnChannelPost, tg.handler)
	tg.b.Start()
	return nil
}
