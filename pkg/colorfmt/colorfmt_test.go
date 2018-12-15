package colorfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithTag(t *testing.T) {
	type testRun struct {
		text     string
		expected string
	}

	tests := []testRun{
		testRun{"", ""},
		testRun{"XXX {yellow} XXX", "XXX \x1b[0;33m XXX"},
		testRun{"-{{yellow}-", "-{{yellow}-"},
	}

	for _, run := range tests {
		buf := new(bytes.Buffer)
		New(buf).Print(run.text)
		require.Equal(t, run.expected, buf.String())
	}
}
