package log

import (
	"fmt"
	"github.com/nathanielc/morgoth"
	"io"
	"os"
)

type LogNotifier struct {
	out io.Writer
}

func New(file string) (*LogNotifier, error) {

	out, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	ln := &LogNotifier{
		out,
	}
	return ln, nil
}

func (self *LogNotifier) Notify(msg string, w *morgoth.Window) {
	fmt.Fprintf(self.out, "Anomaly: %s Name: %s Tags: %v Start: %s Stop %s\n", msg, w.Name, w.Tags, w.Start, w.Stop)
}
