package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
)

type generator struct {
	buf bytes.Buffer // Accumulated output.
}

func (g *generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *generator) Println(a ...interface{}) {
	fmt.Fprintln(&g.buf, a...)
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *generator) Format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
