package sayhello

import (
	"context"
	httphelpers "kurles/kuber/lib/client/http_helpers"
	"strings"
)

type SayHello struct {
	addr string
}

func (s SayHello) makeAddress(path string) string {
	return s.addr + path
}

func New(addr string) *SayHello {
	addr = strings.TrimRight(addr, "/")
	res := &SayHello{
		addr: addr,
	}
	return res
}

func (s *SayHello) SayHello(ctx context.Context, name string) (string, error) {
	res, err := httphelpers.GetRequest(ctx, s.makeAddress("/hello/"+name))
	if err != nil {
		return "", err
	}
	return string(res), nil
}
