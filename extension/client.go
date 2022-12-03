package extension

import (
	"context"
	"fmt"
	"net/http"

	"github.com/michimani/aws-lambda-api-go/alago"
	"github.com/michimani/aws-lambda-api-go/extension"
)

type Client struct {
	alagoClient         *alago.Client
	logger              *Logger
	extensionIdentifier string
}

func NewClient(hc *http.Client, l *Logger) (*Client, error) {
	ac, err := alago.NewClient(&alago.NewClientInput{
		HttpClient: hc,
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		alagoClient: ac,
		logger:      l,
	}, nil
}

func (c *Client) Register(ctx context.Context, extensionName string) error {
	out, err := extension.Register(ctx, c.alagoClient, &extension.RegisterInput{
		LambdaExtensionName: extensionName,
		Events: []extension.EventType{
			extension.EventTypeInvoke,
		},
	})

	if err != nil {
		return fmt.Errorf("An expected error occurred at extension registration. err:%v", err)
	}

	if out.StatusCode != http.StatusOK {
		return fmt.Errorf("An error occurred at extension registration. statusCode:%d errType:%s errMessage:%s",
			out.StatusCode, out.Error.ErrorType, out.Error.ErrorMessage)
	}

	c.logger.Info("Succeeded to register extension.")
	c.extensionIdentifier = out.LambdaExtensionIdentifier

	return nil
}
