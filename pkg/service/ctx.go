package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SignContext() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	return ctx
}
