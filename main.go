package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/michimani/invocation-history-extension/extension"
)

func main() {
	_, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	l := extension.NewLogger(os.Stdout, extension.ExtensionName)

	go func() {
		s := <-sigs
		cancel()

		l.Info("Received Signal: %v", s)
		l.Info("Exiting")
	}()
}
