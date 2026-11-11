package game

import (
	"fmt"
	"io"
	"math"

	"github.com/tylland/dashland/internal/ecs"
)

const (
	ValueIndent     = " "
	PrefixLabel     = "#"
	PrefixChanged   = "*"
	PrefixDefault   = " "
	NumberDelimiter = " : "
	Comment         = "//"
)

type TextWriter struct {
	w io.Writer
}

func NewTextWriter(w io.Writer) *TextWriter {
	return &TextWriter{w: w}
}

func (m *TextWriter) write(format string, a ...any) error {
	_, err := fmt.Fprintf(m.w, format, a...)

	if err != nil {
		return err
	}

	return nil
}

func (m *TextWriter) WriteEntityID(id ecs.EntityID) error {
	return m.write("%s\n", id)
}

func (m *TextWriter) WriteEntityType(et ecs.EntityType) error {
	return m.write("%d\n", et)
}

func (m *TextWriter) WriteVersion(ver int) error {
	return m.write("%d\n", ver)
}

func (m *TextWriter) WriteLabel(lbl string) error {
	return m.write("%s %s\n", PrefixLabel, lbl)
}

func (m *TextWriter) WriteIntValue(v int, d bool) error {
	bits := math.Float32bits(float32(v))

	dirty := PrefixChanged

	if d {
		dirty = PrefixDefault
	}

	return m.write("%s%s %d%s%08x\n", ValueIndent, dirty, v, NumberDelimiter, bits)
}

func (m *TextWriter) WriteFloatValue(v float32, d bool) error {
	bits := math.Float32bits(v) // Get IEEE 754 binary representation

	dirty := PrefixChanged

	if d {
		dirty = PrefixDefault
	}

	return m.write("%s%s %6f%s%08x\n", ValueIndent, dirty, v, NumberDelimiter, bits)
}
