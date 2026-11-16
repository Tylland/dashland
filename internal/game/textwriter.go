package game

import (
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/tylland/dashland/internal/ecs"
)

const (
	TokenIndent    = " "
	TokenLabel     = "#"
	TokenChanged   = "*"
	TokenDefault   = " "
	TokenComment   = "//"
	TokenSeparator = " : "
)

const CommentTab = 32

type TextWriter struct {
	w     io.Writer
	bytes int
}

func NewTextWriter(w io.Writer) *TextWriter {
	return &TextWriter{w: w}
}

func (m *TextWriter) write(format string, a ...any) error {
	if m.bytes > 0 {
		format = "\n" + format
	}

	n, err := fmt.Fprintf(m.w, format, a...)

	if err != nil {
		return err
	}

	m.bytes += n

	return nil
}

func (m *TextWriter) WriteEntityID(id ecs.EntityID) error {
	return m.write("%s", id)
}

func (m *TextWriter) WriteEntityType(et ecs.EntityType) error {
	return m.write("%d", et)
}

func (m *TextWriter) WriteEntityCategory(cat ecs.EntityCategory) error {
	return m.write("%d", cat)
}

func (m *TextWriter) WriteVersion(ver int) error {
	return m.write("%d", ver)
}

func (m *TextWriter) WriteLabel(lbl string) error {
	return m.write("%s %s", TokenLabel, lbl)
}

// func (m *TextWriter) StartSection(name string, version int) error {
// 	m.sec = tabwriter.NewWriter(m.w, 1, 32, 0, ' ', 0)
// 	return m.write("%s %s@%d", TokenLabel, name, version)
// }

func (m *TextWriter) WriteString(str string, d bool) error {
	changed := TokenChanged

	if d {
		changed = TokenDefault
	}

	return m.write("%s%s %s", TokenIndent, changed, str)
}

func (m *TextWriter) WriteBool(v bool, d bool) error {
	return m.WriteString(strconv.FormatBool(v), d)
}

func (m *TextWriter) WriteInt(v int, d bool) error {
	bits := math.Float32bits(float32(v))

	dirty := TokenChanged

	if d {
		dirty = TokenDefault
	}

	return m.write("%s%s %d%s%08x", TokenIndent, dirty, v, TokenSeparator, bits)
}

func (m *TextWriter) WriteFloat(v float32, d bool) error {
	bits := math.Float32bits(v) // Get IEEE 754 binary representation

	dirty := TokenChanged

	if d {
		dirty = TokenDefault
	}

	return m.write("%s%s %6f%s%08x", TokenIndent, dirty, v, TokenSeparator, bits)
}

func (m *TextWriter) WriteEnum(enum any, d bool, c string) error {
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

	line := fmt.Sprintf("%s%s %d", TokenIndent, dirty, vint)

	comment := ""
	if c != "" {
		comment = c
	}

	if stringer, ok := enum.(fmt.Stringer); ok {
		if comment != "" {
			comment += TokenIndent
		}

		comment += stringer.String()
	}

	if comment != "" {
		pos := (int(len(line)/CommentTab) + 1) * CommentTab

		line += fmt.Sprintf("%s%s %s", strings.Repeat(TokenIndent, pos-len(line)), TokenComment, comment)
	}

	return m.write("%s", line)
}

func (m *TextWriter) Close() error {
	if c, ok := m.w.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func NewTextFileWriter(fileName string) (*TextWriter, error) {
	file, err := os.Create(fileName)

	if err != nil {
		return nil, err
	}

	return NewTextWriter(file), nil

}
