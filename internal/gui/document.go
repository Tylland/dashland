package gui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type ValueKind int

const (
	KindText ValueKind = iota
	KindBool
	KindInt
	KindFloat
	KindEnum
)

type Value struct {
	Kind   ValueKind
	Edited bool
	Name   string
	Text   string // text, "true"/"false", or the decimal part of a number
	Hex    string // numbers only: 8 hex digits of the float32 bits
	Float  bool
}

// Section is a "# label" followed by its values.
type Section struct {
	Label  string
	Values []*Value
}

type Document struct{ Sections []*Section }

// formatNumber renders f as "30" (integer style) or "30.000000" (float style).
// It also reports whether float style had to be used.
func formatNumber(f float64, floatStyle bool) (string, bool) {
	if !floatStyle && f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10), false
	}
	return fmt.Sprintf("%f", f), true
}

// SetDecimal parses a decimal text and updates both the text and the hex.
// ok is false when text is not a number; changed tells whether the value differs.
func (v *Value) SetDecimal(text string) (ok, changed bool) {
	text = strings.TrimSpace(text)
	f, err := strconv.ParseFloat(text, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return false, false
	}
	style := v.Float || strings.ContainsAny(text, ".eE")
	newText, style := formatNumber(f, style)
	newHex := fmt.Sprintf("%08x", math.Float32bits(float32(f)))
	changed = newText != v.Text || newHex != v.Hex
	v.Text, v.Hex, v.Float = newText, newHex, style
	return true, changed
}

// SetHex parses a hex bit pattern and updates both the hex and the decimal text.
func (v *Value) SetHex(text string) (ok, changed bool) {
	u, err := strconv.ParseUint(strings.TrimSpace(text), 16, 32)
	if err != nil {
		return false, false
	}
	f := float64(math.Float32frombits(uint32(u)))
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return false, false
	}
	newText, style := formatNumber(f, v.Float)
	newHex := fmt.Sprintf("%08x", u)
	changed = newText != v.Text || newHex != v.Hex
	v.Text, v.Hex, v.Float = newText, newHex, style
	return true, changed
}
