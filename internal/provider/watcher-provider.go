package provider

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type IWatcherProvider interface {
	Catch()
}

type WatcherProvider struct {
}

func NewWatcherProvider() IWatcherProvider {
	return &WatcherProvider{}
}

func (rcv *WatcherProvider) Catch() {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGABRT)
	signal.Notify(sig, syscall.SIGQUIT)

	s := <-sig

	slog.Warn("shutdown requested", "signal", s)
}
