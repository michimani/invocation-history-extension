package extension_test

import (
	"testing"
	"time"

	"github.com/michimani/invocation-history-extension/extension"
	"github.com/stretchr/testify/assert"
)

func Test_InvocationHistory_Init(t *testing.T) {
	cases := []struct {
		name      string
		ih        *extension.InvocationHistory
		wantError bool
	}{
		{
			name: "ok",
			ih:   &extension.InvocationHistory{},
		},
		{
			name:      "ng: nil",
			ih:        nil,
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.ih.Init()
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)

			asst.NotNil(c.ih)
			asst.Empty(c.ih.Invocations)
		})
	}
}

func Test_InvocationHistory_Add(t *testing.T) {
	newTime := func(t time.Time) *time.Time {
		return &t
	}

	type e struct {
		len         int
		invocations []extension.Invocation
	}

	cases := []struct {
		name      string
		ih        *extension.InvocationHistory
		added     *extension.Invocation
		expect    e
		wantError bool
	}{
		{
			name:  "ok",
			ih:    &extension.InvocationHistory{},
			added: &extension.Invocation{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
			expect: e{
				len: 1,
				invocations: []extension.Invocation{
					{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
				},
			},
		},
		{
			name: "ok: 2",
			ih: &extension.InvocationHistory{
				Invocations: []*extension.Invocation{
					{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
				},
			},
			added: &extension.Invocation{AWSRequestID: "id-2", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 2, time.UTC))},
			expect: e{
				len: 2,
				invocations: []extension.Invocation{
					{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
					{AWSRequestID: "id-2", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 2, time.UTC))},
				},
			},
		},
		{
			name:  "ok: i.Invocation is nil",
			ih:    &extension.InvocationHistory{Invocations: nil},
			added: &extension.Invocation{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
			expect: e{
				len: 1,
				invocations: []extension.Invocation{
					{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
				},
			},
		},
		{
			name: "ok: added is nil",
			ih: &extension.InvocationHistory{
				Invocations: []*extension.Invocation{
					{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
				},
			},
			added: nil,
			expect: e{
				len: 1,
				invocations: []extension.Invocation{
					{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
				},
			},
		},
		{
			name:      "ng: i is nil",
			ih:        nil,
			added:     &extension.Invocation{AWSRequestID: "id-1", InvocatedAt: newTime(time.Date(2022, 12, 27, 0, 0, 0, 1, time.UTC))},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.ih.Add(c.added)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)

			asst.Equal(c.expect.len, len(c.ih.Invocations))
			for i, ei := range c.expect.invocations {
				asst.Equal(ei, *c.ih.Invocations[i])
			}
		})
	}
}
