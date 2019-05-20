package executor

import (
	"bufio"
	"io"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/termui"
)

func printPipe(input io.Reader, output io.Writer, prefix string, lineSuppressor lineSuppressor) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if lineSuppressor.Suppress(line) {
			continue
		}
		termui.Fprintf(output, "%s%s\n", prefix, line)
	}
}

type lineSuppressor interface {
	Suppress(string) bool
}

type substringLineSuppressor struct {
	substrings []string
}

func (s *substringLineSuppressor) Suppress(line string) bool {
	for _, substring := range s.substrings {
		if strings.Contains(line, substring) {
			return true
		}
	}
	return false
}
