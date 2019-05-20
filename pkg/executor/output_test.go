package executor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_printPipe(t *testing.T) {

	input := &bytes.Buffer{}
	input.WriteString("line1 thisisfine\n")
	input.WriteString("line2 poipoi\n")
	input.WriteString("line2 nope nope\n")

	prefix := "---->"

	lineSuppressor := &substringLineSuppressor{[]string{"poipoi", "nope"}}

	output := &bytes.Buffer{}

	printPipe(input, output, prefix, lineSuppressor)

	require.Equal(t, "---->line1 thisisfine\n", output.String())
}
