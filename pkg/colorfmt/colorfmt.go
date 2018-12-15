// colorfmt
/*

{:green}This is a header
{:green}Click on {link}http://poipoi.com
{bi:white}Flash NEWS
{b:}

modifiers:
- b for bold
- h for high
-   for blink
- u for underline


*/

package colorfmt

import (
	"fmt"
	"io"
)

// Formatter interprets the color tags in a string and prints it to the file provided
type Formatter struct {
	// colorEnabled bool
	file io.Writer
}

// New returns a *Formatter
func New(file io.Writer) *Formatter {
	return &Formatter{file: file}
}

// Print interprets tags and print the colorized text
func (f *Formatter) Print(text string) {
	// formatted := scan(text, '{', '}', colorizeTag)
	// f.file.Write([]byte(formatted))
	fmt.Fprintf(f.file, scan(text, '{', '}', colorizeTag))
}
