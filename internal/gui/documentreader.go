package gui

import (
	"fmt"
	"math"
	"strconv"

	"github.com/tylland/dashland/internal/ecs"
)

type DocumentReader struct {
	doc     *Document
	section int
	value   int
}

var _ ecs.Reader = (*DocumentReader)(nil)

func NewDocumentReader(doc *Document) *DocumentReader {
	return &DocumentReader{doc: doc, section: -1}
}

func (r *DocumentReader) ReadLabel() (string, error) {
	if r.doc == nil || r.section+1 >= len(r.doc.Sections) {
		return "", fmt.Errorf("read label: no more sections")
	}
	r.section++
	r.value = 0
	section := r.doc.Sections[r.section]
	if section == nil {
		return "", fmt.Errorf("read label: section %d is nil", r.section)
	}
	return section.Label, nil
}

func (r *DocumentReader) ConsumeLabel(expected string) error {
	label, err := r.ReadLabel()
	if err != nil {
		return err
	}
	if label != expected {
		return fmt.Errorf("invalid section label %q: expected %q", label, expected)
	}
	return nil
}

func (r *DocumentReader) nextValue() (*Value, error) {
	if r.doc == nil || r.section < 0 || r.section >= len(r.doc.Sections) {
		return nil, fmt.Errorf("read value: no current section")
	}
	section := r.doc.Sections[r.section]
	if section == nil || r.value >= len(section.Values) {
		return nil, fmt.Errorf("read value: no more values in section %q", section.Label)
	}
	value := section.Values[r.value]
	r.value++
	if value == nil {
		return nil, fmt.Errorf("read value: value %d in section %q is nil", r.value-1, section.Label)
	}
	return value, nil
}

func (r *DocumentReader) ConsumeLine() error {
	_, err := r.nextValue()
	return err
}

func (r *DocumentReader) ReadBool() (bool, bool, error) {
	value, err := r.nextValue()
	if err != nil {
		return false, false, err
	}
	result, err := strconv.ParseBool(value.Text)
	if err != nil {
		return false, value.Edited, fmt.Errorf("parse bool %q for %q: %w", value.Text, value.Name, err)
	}
	return result, value.Edited, nil
}

func (r *DocumentReader) ReadEnum() (uint32, bool, error) {
	value, err := r.nextValue()
	if err != nil {
		return 0, false, err
	}
	result, err := strconv.ParseUint(value.Text, 10, 32)
	if err != nil {
		return 0, value.Edited, fmt.Errorf("parse enum %q for %q: %w", value.Text, value.Name, err)
	}
	return uint32(result), value.Edited, nil
}

func (r *DocumentReader) ReadFloat() (float32, bool, error) {
	value, err := r.nextValue()
	if err != nil {
		return 0, false, err
	}
	if value.Hex != "" {
		bits, err := strconv.ParseUint(value.Hex, 16, 32)
		if err != nil {
			return 0, value.Edited, fmt.Errorf("parse float bits %q for %q: %w", value.Hex, value.Name, err)
		}
		result := math.Float32frombits(uint32(bits))
		if math.IsNaN(float64(result)) || math.IsInf(float64(result), 0) {
			return 0, value.Edited, fmt.Errorf("invalid non-finite float for %q", value.Name)
		}
		return result, value.Edited, nil
	}
	result, err := strconv.ParseFloat(value.Text, 32)
	if err != nil {
		return 0, value.Edited, fmt.Errorf("parse float %q for %q: %w", value.Text, value.Name, err)
	}
	return float32(result), value.Edited, nil
}

func (r *DocumentReader) ReadInt() (int, bool, error) {
	value, err := r.nextValue()
	if err != nil {
		return 0, false, err
	}
	result, err := strconv.ParseInt(value.Text, 10, 32)
	if err != nil {
		return 0, value.Edited, fmt.Errorf("parse int %q for %q: %w", value.Text, value.Name, err)
	}
	return int(result), value.Edited, nil
}

func (r *DocumentReader) ReadString() (string, bool, error) {
	value, err := r.nextValue()
	if err != nil {
		return "", false, err
	}
	return value.Text, value.Edited, nil
}
