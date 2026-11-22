package gui

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/tylland/dashland/internal/ecs"
)

type DocumentWriter struct {
	doc            *Document
	currentSection *Section
}

func NewDocumentWriter(doc *Document) *DocumentWriter {
	return &DocumentWriter{doc: doc}
}

func (w *DocumentWriter) addValue(val *Value) {
	w.currentSection.Values = append(w.currentSection.Values, val)
}

func (w *DocumentWriter) StartSection(lbl string) error {
	w.currentSection = &Section{Label: lbl}
	w.doc.Sections = append(w.doc.Sections, w.currentSection)

	return nil
}

func (w *DocumentWriter) WriteEntityID(id ecs.EntityID) error {
	w.addValue(&Value{
		Name: "ID",
		Kind: KindText,
		Text: string(id)})

	return nil
}
func (w *DocumentWriter) WriteEntityType(et ecs.EntityType) error {
	w.addValue(&Value{
		Name: "Type",
		Kind: KindInt,
		Text: strconv.Itoa(int(et)),
	})
	return nil
}

func (w *DocumentWriter) WriteEntityCategory(c ecs.EntityCategory) error {
	w.addValue(&Value{
		Name: "Category",
		Kind: KindInt,
		Text: strconv.Itoa(int(c))})

	return nil
}

func (w *DocumentWriter) WriteVersion(ver int) error {
	w.addValue(&Value{
		Name: "Version",
		Kind: KindInt,
		Text: strconv.Itoa(ver)})

	return nil
}

func (w *DocumentWriter) WriteBool(name string, v bool, d bool) error {
	w.addValue(&Value{
		Name: name,
		Kind: KindBool,
		Text: strconv.FormatBool(v)})

	return nil
}

func (w *DocumentWriter) WriteEnum(name string, enum any, d bool, c string) error {
	var vint uint32

	v := reflect.ValueOf(enum)

	switch v.Kind() {
	case reflect.Uint16:
		vint = uint32(v.Int())
	case reflect.Uint32:
		vint = uint32(v.Uint())
	default:
		return fmt.Errorf("unsupported enum type: %T", enum)
	}
	w.addValue(&Value{
		Name: name,
		Kind: KindEnum,
		Text: strconv.Itoa(int(vint))})

	return nil
}

func (w *DocumentWriter) WriteFloat(name string, v float32, d bool) error {
	bits := math.Float32bits(float32(v))

	w.addValue(&Value{
		Name:  name,
		Kind:  KindFloat,
		Text:  fmt.Sprintf("%6f", v),
		Hex:   fmt.Sprintf("%08x", bits),
		Float: true})

	return nil
}

func (w *DocumentWriter) WriteInt(name string, v int, d bool) error {
	bits := math.Float32bits(float32(v))
	w.addValue(&Value{
		Name:  name,
		Kind:  KindInt,
		Text:  strconv.Itoa(v),
		Hex:   fmt.Sprintf("%08x", bits),
		Float: false})

	return nil
}

func (w *DocumentWriter) WriteString(name string, v string, d bool) error {
	w.addValue(&Value{
		Name: name,
		Kind: KindText,
		Text: v})

	return nil
}
