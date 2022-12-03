package extension_test

import (
	"bytes"
	"testing"

	"github.com/michimani/invocation-history-extension/extension"
	"github.com/stretchr/testify/assert"
)

func Test_LoggerInfo(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		msg    string
		v      []any
		expect string
	}{
		{
			name:   "simple",
			prefix: "prefix",
			msg:    "simple message",
			expect: "[prefix] INFO: simple message\n",
		},
		{
			name:   "with format",
			prefix: "prefix",
			msg:    "%s message",
			v:      []any{"formatted"},
			expect: "[prefix] INFO: formatted message\n",
		},
		{
			name:   "empty prefix",
			msg:    "message",
			expect: "INFO: message\n",
		},
		{
			name:   "empty message",
			expect: "INFO: \n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			out := bytes.Buffer{}
			l := extension.NewLogger(&out, c.prefix)

			l.Info(c.msg, c.v...)
			asst.Equal(c.expect, out.String())
		})
	}
}

func Test_LoggerWarn(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		msg    string
		v      []any
		expect string
	}{
		{
			name:   "simple",
			prefix: "prefix",
			msg:    "simple message",
			expect: "[prefix] WARN: simple message\n",
		},
		{
			name:   "with format",
			prefix: "prefix",
			msg:    "%s message",
			v:      []any{"formatted"},
			expect: "[prefix] WARN: formatted message\n",
		},
		{
			name:   "empty prefix",
			msg:    "message",
			expect: "WARN: message\n",
		},
		{
			name:   "empty message",
			expect: "WARN: \n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			out := bytes.Buffer{}
			l := extension.NewLogger(&out, c.prefix)

			l.Warn(c.msg, c.v...)
			asst.Equal(c.expect, out.String())
		})
	}
}

func Test_LoggerError(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		msg    string
		v      []any
		expect string
	}{
		{
			name:   "simple",
			prefix: "prefix",
			msg:    "simple message",
			expect: "[prefix] ERROR: simple message\n",
		},
		{
			name:   "with format",
			prefix: "prefix",
			msg:    "%s message",
			v:      []any{"formatted"},
			expect: "[prefix] ERROR: formatted message\n",
		},
		{
			name:   "empty prefix",
			msg:    "message",
			expect: "ERROR: message\n",
		},
		{
			name:   "empty message",
			expect: "ERROR: \n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			out := bytes.Buffer{}
			l := extension.NewLogger(&out, c.prefix)

			l.Error(c.msg, c.v...)
			asst.Equal(c.expect, out.String())
		})
	}
}
