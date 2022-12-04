package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Response struct {
	Message     string       `json:"message"`
	Invocations []invocation `json:"invocations"`
}

func handleRequest(ctx context.Context) (*Response, error) {
	log.Println("start handler")
	defer log.Println("end handler")

	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Failed to get Lambda context from context.")
	}

	history, err := getInvocationHistory()
	if err != nil {
		return nil, err
	}

	return &Response{
		Message:     fmt.Sprintf("Current request ID is %s", lc.AwsRequestID),
		Invocations: history,
	}, nil
}

func init() {
	log.Println("cold start")
}

func main() {
	runtime.Start(handleRequest)
}

// Struct of response from Invocation History Extension - GET /invocations API
type resultFromExtension struct {
	Invocations []invocation `json:"invocations"`
}

type invocation struct {
	AWSRequestID string     `json:"awsRequestId"`
	InvocatedAt  *time.Time `json:"invocatedAt"`
}

const invocationsEndpoint = "http://localhost:1203/invocations"

// Get invocation history from Invocation History Extension IPC.
func getInvocationHistory() ([]invocation, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", invocationsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	// call Extension API
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return nil, err
	}
	bodyString := buf.String()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get invocation history by using extension. statusCode:%d body:%s", res.StatusCode, bodyString)
	}

	exRes := resultFromExtension{}
	if err := json.Unmarshal([]byte(bodyString), &exRes); err != nil {
		return nil, err
	}

	return exRes.Invocations, nil
}
