package ipc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/michimani/invocation-history-extension/extension"
	"github.com/michimani/invocation-history-extension/types"
)

var History = &types.InvocationHistory{}

const (
	extensionIPCPortEnvKey = "INVOCATION_HISTORY_EXTENSION_HTTP_PORT"
	defaultPort            = "1203"
)

func Start(l *extension.Logger) {
	port := os.Getenv(extensionIPCPortEnvKey)
	if len(port) == 0 {
		port = defaultPort
	}

	go startServer(port, l)
}

func startServer(port string, l *extension.Logger) {
	History.Init()

	http.HandleFunc("/invocations", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(History)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error(err.Error())
			fmt.Fprintf(w, "Internal server error")
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
