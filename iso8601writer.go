package main

import (
	"fmt"
	"io"
	"time"
)

// Iso8601Writer prefixes writes to a nested writer with an ISO 8601 format timestamp
type Iso8601Writer struct {
	upstreamWriter io.Writer
}

func (writer Iso8601Writer) Write(bytes []byte) (int, error) {
	fmt.Print(time.Now().UTC().Format("1970-01-01T01:01:01.001Z "))
	return writer.upstreamWriter.Write(bytes)
}
