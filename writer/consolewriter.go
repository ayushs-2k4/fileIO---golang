package writer

import "os"

type ConsoleWriter struct{}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (c *ConsoleWriter) Write(b []byte) {
	_, _ = os.Stdout.Write(b)
}

func (c *ConsoleWriter) Close() {}
