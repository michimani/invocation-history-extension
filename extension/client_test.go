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
	"github.com/stretchr/testify/require"
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

func Test_Client_PollingEvent(t *testing.T) {
	cases := []struct {
		name          string
		httpClient    *http.Client
		expect        bool
		expectMessage string
		wantError     bool
		historyNil    bool
	}{
		{
			name: "ok: event type INVOKE",
			httpClient: hcmock.New(&hcmock.MockInput{
				Status:     "OK",
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"INVOKE","deadlineMs":123456,"requestId":"aws-request-id","invokedFunctionArn":"function-arn","tracing":{"type":"X-Amzn-Trace-Id","value":"tracing-value"}}`),
			}),
			expect:        true,
			expectMessage: "Succeeded to save new history. awsRequestId:aws-request-id",
		},
		{
			name: "ok: event type SHUTDOWN",
			httpClient: hcmock.New(&hcmock.MockInput{
				Status:     "OK",
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"SHUTDOWN","shutdownReason":"test-down","deadlineMs":123456}`),
			}),
			expect:        false,
			expectMessage: "Received shutdown event. reason:test-down",
		},
		{
			name: "ng: undefined event type",
			httpClient: hcmock.New(&hcmock.MockInput{
				Status:     "OK",
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"UNDEFINED","shutdownReason":"test-down","deadlineMs":123456}`),
			}),
			expect:        false,
			expectMessage: "Cannot handle event. eventType:UNDEFINED",
			wantError:     true,
		},
		{
			name: "ng: event next API returns not 200 status",
			httpClient: hcmock.New(&hcmock.MockInput{
				Status:     "Not Found",
				StatusCode: 404,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"UNDEFINED","shutdownReason":"test-down","deadlineMs":123456}`),
			}),
			expect:        false,
			expectMessage: "An error occurred at calling /extension/event/next API. statusCode:404",
			wantError:     true,
		},
		{
			name: "ng: extension.EventNext returns error",
			httpClient: hcmock.New(&hcmock.MockInput{
				Status:     "OK",
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`///`),
			}),
			expect:    false,
			wantError: true,
		},
		{
			name: "ng: History is nil",
			httpClient: hcmock.New(&hcmock.MockInput{
				Status:     "OK",
				StatusCode: 200,
				Headers: []hcmock.Header{
					{Key: "Lambda-Extension-Event-Identifier", Value: "lambda-extension-event-identifier"},
				},
				BodyBytes: []byte(`{"eventType":"INVOKE","deadlineMs":123456,"requestId":"aws-request-id","invokedFunctionArn":"function-arn","tracing":{"type":"X-Amzn-Trace-Id","value":"tracing-value"}}`),
			}),
			expect:     false,
			wantError:  true,
			historyNil: true,
		},
	}

	t.Setenv("AWS_LAMBDA_RUNTIME_API", "test-host")

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			rqir := require.New(tt)

			out := bytes.Buffer{}
			l := extension.NewLogger(&out, "test")
			client, err := extension.NewClient(c.httpClient, l)
			rqir.NoError(err)
			rqir.NotNil(client)
			extension.SetExtensionIdentifier(client, "test-ex-id")

			if c.historyNil {
				extension.History = nil
			}

			b, err := client.PollingEvent(context.Background())
			if c.wantError {
				asst.Error(err)
				asst.Equal(c.expect, b)
				asst.Contains(err.Error(), c.expectMessage, err.Error())
				return
			}

			asst.NoError(err)

			asst.Equal(c.expect, b)
			asst.Contains(out.String(), c.expectMessage, out.String())
		})
	}
}
