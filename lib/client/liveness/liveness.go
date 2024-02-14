package liveness

import (
	"context"
	"fmt"
	httphelpers "kurles/kuber/lib/client/http_helpers"
	"strings"
	"time"
)

type LivenessClient struct {
	addr string
}

func New(addr string) *LivenessClient {
	addr = strings.TrimRight(addr, "/")
	return &LivenessClient{addr: addr}
}

func (c *LivenessClient) Live(ctx context.Context) error {
	_, err := httphelpers.GetRequest(ctx, c.addr+"/live")
	return err
}

func (c *LivenessClient) Crash(ctx context.Context, after time.Duration) error {
	sec := after.Seconds()
	_, err := httphelpers.GetRequest(ctx, fmt.Sprintf("%s/crash/%d", c.addr, int(sec)))
	return err
}

func (c *LivenessClient) Dead(ctx context.Context, after time.Duration) error {
	sec := after.Seconds()
	_, err := httphelpers.GetRequest(ctx, fmt.Sprintf("%s/simulate/%d", c.addr, int(sec)))
	return err
}
