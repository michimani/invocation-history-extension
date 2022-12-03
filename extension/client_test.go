package extension_test

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/michimani/http-client-mock/hcmock"
	"github.com/michimani/invocation-history-extension/extension"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	cases := []struct {
		name       string
		httpClient *http.Client
		envKey     string
		wantErr    bool
	}{
		{
			name:    "ok",
			envKey:  "AWS_LAMBDA_RUNTIME_API",
			wantErr: false,
		},
		{
			name: "ok: use custom http client",
			httpClient: &http.Client{
				Timeout: 0,
			},
			envKey:  "AWS_LAMBDA_RUNTIME_API",
			wantErr: false,
		},
		{
			name:    "error: env key is not set",
			envKey:  "",
			wantErr: true,
		},
	}

	l := extension.NewLogger(os.Stdout, "test")

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			if c.envKey != "" {
				tt.Setenv(c.envKey, "test-env-value")
			}

			client, err := extension.NewClient(c.httpClient, l)

			if c.wantErr {
				asst.Error(err)
				asst.Nil(client)
				return
			}

			asst.NoError(err)
			asst.NotNil(client)
		})
	}
}

func Test_Register(t *testing.T) {
	cases := []struct {
		name       string
		httpClient *http.Client
		host       string
		expectExID string
		wantErr    bool
	}{
		{
			name: "ok",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Identifier", Value: "lambda-extension-identifier"},
				},
				BodyBytes: []byte(`{"functionName":"my-function","functionVersion":"$LATEST","handler":"lambda_handler"}`),
			}),
			host:       "test-host",
			expectExID: "lambda-extension-identifier",
			wantErr:    false,
		},
		{
			name: "error: not OK status code",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 400,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Identifier", Value: "lambda-extension-identifier"},
				},
				BodyBytes: []byte(`{"errorMessage":"test-error-message", "errorType":"test-error-type"}`),
			}),
			host:       "test-host",
			expectExID: "",
			wantErr:    true,
		},
		{
			name: "error: expected error",
			httpClient: hcmock.New(&hcmock.MockInput{
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Identifier", Value: "lambda-extension-identifier"},
				},
				BodyBytes: []byte(`{"functionName":"my-function","functionVersion":"$LATEST","handler":"lambda_handler"}`),
			}),
			host:       "\U00000001",
			expectExID: "",
			wantErr:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			tt.Setenv("AWS_LAMBDA_RUNTIME_API", c.host)

			out := bytes.Buffer{}
			l := extension.NewLogger(&out, "test")
			client, err := extension.NewClient(c.httpClient, l)

			asst.NoError(err)

			rerr := client.Register(context.Background(), "test-extension-name")
			if c.wantErr {
				asst.Error(rerr, rerr)
			} else {
				asst.NoError(rerr, rerr)
			}

			asst.Equal(c.expectExID, extension.ExportedExtensionIdentifier(client))
		})
	}
}
