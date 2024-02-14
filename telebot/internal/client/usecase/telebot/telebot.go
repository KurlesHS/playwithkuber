package telebot

import (
	"context"
	"fmt"
	"kurles/kuber/telebot/internal/model"
	"strconv"
	"strings"
	"time"
)

type Telebot interface {
	SendMessage(ctx context.Context, text string) error
	GetMessageChannel(ctx context.Context) <-chan model.Message
}

type Liveness interface {
	Crash(ctx context.Context, after time.Duration) error
	Dead(ctx context.Context, after time.Duration) error
}

type SayHello interface {
	SayHello(ctx context.Context, name string) (string, error)
}

type TelebotUseCase struct {
	t      Telebot
	l      Liveness
	s      SayHello
	doneCh chan struct{}
}

func New(ctx context.Context, t Telebot, l Liveness, s SayHello) *TelebotUseCase {
	return &TelebotUseCase{
		t:      t,
		l:      l,
		s:      s,
		doneCh: make(chan struct{}),
	}
}

func (t *TelebotUseCase) ProcessMessage(ctx context.Context, msg model.Message) {
	idx := strings.Index(msg.Message, " ")
	var cmd, data string
	if idx == -1 {
		cmd = msg.Message
	} else {
		data = msg.Message[idx+1:]
		cmd = msg.Message[:idx]
	}
	switch cmd {
	case "hello":
		res, err := t.s.SayHello(ctx, data)
		if err == nil {
			t.t.SendMessage(ctx, res)
		}
	case "crash":
		sec, err := strconv.Atoi(data)
		if err != nil {
			t.t.SendMessage(ctx, "invalid number")
			return
		}
		t.t.SendMessage(ctx, "say hello service will crash after "+time.Duration(sec).String())
		t.l.Crash(ctx, time.Duration(sec)*time.Second)
	case "dead":
		sec, err := strconv.Atoi(data)
		if err != nil {
			t.t.SendMessage(ctx, "invalid number")
			return
		}
		t.t.SendMessage(ctx, "say hello service will simulate dead after "+(time.Duration(sec)*time.Second).String())
		t.l.Dead(ctx, time.Duration(sec)*time.Second)
	default:
		t.t.SendMessage(ctx, fmt.Sprintf("unknown command '%s'", cmd))
	}
}

func (t *TelebotUseCase) MessageWaiter(ch <-chan model.Message, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.doneCh:
			return
		case msg := <-ch:
			t.ProcessMessage(ctx, msg)
		}
	}
}

func (t *TelebotUseCase) Start(ctx context.Context) {
	ch := t.t.GetMessageChannel(ctx)
	go t.MessageWaiter(ch, ctx)
}

func (t *TelebotUseCase) Stop(ctx context.Context) {
	close(t.doneCh)
}
