package game

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"

	"github.com/tylland/dashland/internal/ecs"
)

const (
	TokenIndent          = " "
	TokenLabel           = "#"
	TokenChanged         = "*"
	TokenDefault         = " "
	TokenComment         = "//"
	TokenNumberDelimiter = " : "
)

const CommentPos = 32

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

func (m *TextWriter) WriteBlockType(v int) error {
	return m.write("%d\n", v)
}

func (m *TextWriter) WriteVersion(ver int) error {
	return m.write("%d\n", ver)
}

func (m *TextWriter) WriteLabel(lbl string) error {
	return m.write("%s %s\n", TokenLabel, lbl)
}

func (m *TextWriter) WriteIntValue(v int, d bool) error {
	bits := math.Float32bits(float32(v))

	dirty := TokenChanged

	if d {
		dirty = TokenDefault
	}

	return m.write("%s%s %d%s%08x\n", TokenIndent, dirty, v, TokenNumberDelimiter, bits)
}

func (m *TextWriter) WriteFloatValue(v float32, d bool) error {
	bits := math.Float32bits(v) // Get IEEE 754 binary representation

	dirty := TokenChanged

	if d {
		dirty = TokenDefault
	}

	return m.write("%s%s %6f%s%08x\n", TokenIndent, dirty, v, TokenNumberDelimiter, bits)
}

func (m *TextWriter) WriteEnumValue(enum any, d bool) error {
	var vint uint32

	v := reflect.ValueOf(enum)

	switch v.Kind() {
	case reflect.Uint16:
		vint = uint32(v.Uint())
	case reflect.Uint32:
		vint = uint32(v.Uint())
	default:
		return fmt.Errorf("unsupported enum type: %T", enum)
	}

	dirty := TokenChanged

	if d {
		dirty = TokenDefault
	}

	valStr := fmt.Sprintf("%s%s %d", TokenIndent, dirty, vint)

	comment := ""
	if stringer, ok := enum.(fmt.Stringer); ok {
		comment = fmt.Sprintf("%s%s %s", strings.Repeat(TokenIndent, CommentPos-len(valStr)), TokenComment, stringer.String())
	}

	return m.write("%s%s\n", valStr, comment)
}
