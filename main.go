package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/michimani/invocation-history-extension/extension"
	"github.com/michimani/invocation-history-extension/ipc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	l := extension.NewLogger(os.Stdout, extension.ExtensionName)

	go func() {
		s := <-sigs
		cancel()

		l.Info("Received Signal: %v", s)
		l.Info("Exiting")
	}()

	ipc.Start(l)

	processEvents(ctx, l)
}

func processEvents(ctx context.Context, l *extension.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			l.Info("Waiting for event...")
		}
	}
}
