package extension

import (
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

func (i *InvocationHistory) Init() {
	i.Invocations = []*Invocation{}
}

func (i *InvocationHistory) Add(inv *Invocation) {
	if i == nil {
		return
	}

	if i.Invocations == nil {
		i.Init()
	}

	i.Invocations = append(i.Invocations, inv)
}
