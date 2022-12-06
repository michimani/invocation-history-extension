package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/michimani/invocation-history-extension/extension"
	"github.com/michimani/invocation-history-extension/ipc"
)

var extensionName = filepath.Base(os.Args[0])

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	l := extension.NewLogger(os.Stdout, extension.NameForLog)
	go func() {
		s := <-sigs
		cancel()

		l.Info("Received Signal: %v", s)
		l.Info("Exiting")
	}()

	hc := &http.Client{Timeout: 0}
	extensionClient, err := extension.NewClient(hc, l)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	// register extension
	if err := extensionClient.Register(ctx, extensionName); err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	ipc.Start(l)

	if err := processEvents(ctx, extensionClient); err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
}

func processEvents(ctx context.Context, c *extension.Client) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if waitNextEvent, err := c.PollingEvent(ctx); err != nil {
				return err
			} else if !waitNextEvent {
				// received shutdown event
				return nil
			}
		}
	}
}
