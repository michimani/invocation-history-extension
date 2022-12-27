package extension

import (
	"fmt"
	"time"
)

var History = &InvocationHistory{}

type Invocation struct {
	AWSRequestID string     `json:"awsRequestId"`
	InvocatedAt  *time.Time `json:"invocatedAt"`
}

type InvocationHistory struct {
	Invocations []*Invocation `json:"invocations"`
}

func (i *InvocationHistory) Init() error {
	if i == nil {
		return fmt.Errorf("i is nil")
	}

	i.Invocations = []*Invocation{}
	return nil
}

func (i *InvocationHistory) Add(inv *Invocation) error {
	if i == nil {
		return fmt.Errorf("i is nil")
	}

	if inv == nil {
		return nil
	}

	if i.Invocations == nil {
		_ = i.Init()
	}

	i.Invocations = append(i.Invocations, inv)
	return nil
}
