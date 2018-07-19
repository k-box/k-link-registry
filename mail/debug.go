package mail

import (
	"fmt"
	"io"
	"os"
)

// DebugMailer writes mail to a defined output, if no output is defined
// it will fall back to stdout.
type DebugMailer struct {
	Destination io.Writer
}

// Email satisfies the Emailer interface
func (m DebugMailer) Email(recepient, subject, html, text string) error {
	if m.Destination == nil {
		m.Destination = os.Stdout
	}

	fmt.Fprintf(m.Destination,
		`TO: %s\nSUBJECT: %s
		%s\n`, recepient, subject, text)

	return nil
}
