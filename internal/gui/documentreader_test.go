package gui

import (
	"math"
	"testing"
)

func TestDocumentReaderReadsSectionsAndValues(t *testing.T) {
	doc := &Document{Sections: []*Section{
		{
			Label: "header",
			Values: []*Value{
				{Name: "ID", Text: "entity-1"},
				{Name: "Type", Kind: KindEnum, Text: "3", Edited: true},
				{Name: "Position", Kind: KindInt, Text: "42"},
			},
		},
		{
			Label: "properties",
			Values: []*Value{
				{Name: "Enabled", Kind: KindBool, Text: "true"},
				{Name: "Layer", Kind: KindEnum, Text: "5", Display: "Ground", Edited: true},
				{Name: "Speed", Kind: KindFloat, Text: "2.5", Hex: "40200000"},
			},
		},
	}}
	reader := NewDocumentReader(doc)

	if err := reader.ConsumeLabel("header"); err != nil {
		t.Fatal(err)
	}
	if got, changed, err := reader.ReadString(); err != nil || changed || got != "entity-1" {
		t.Fatalf("ReadString() = (%q, %t, %v), want (%q, false, nil)", got, changed, err, "entity-1")
	}
	if got, changed, err := reader.ReadEnum(); err != nil || !changed || got != 3 {
		t.Fatalf("ReadEnum() = (%d, %t, %v), want (3, true, nil)", got, changed, err)
	}
	if got, changed, err := reader.ReadInt(); err != nil || changed || got != 42 {
		t.Fatalf("ReadInt() = (%d, %t, %v), want (42, false, nil)", got, changed, err)
	}

	if err := reader.ConsumeLabel("properties"); err != nil {
		t.Fatal(err)
	}
	if got, changed, err := reader.ReadBool(); err != nil || changed || !got {
		t.Fatalf("ReadBool() = (%t, %t, %v), want (true, false, nil)", got, changed, err)
	}
	if got, changed, err := reader.ReadEnum(); err != nil || !changed || got != 5 {
		t.Fatalf("ReadEnum() = (%d, %t, %v), want (5, true, nil)", got, changed, err)
	}
	if got, changed, err := reader.ReadFloat(); err != nil || changed || got != math.Float32frombits(0x40200000) {
		t.Fatalf("ReadFloat() = (%g, %t, %v), want (2.5, false, nil)", got, changed, err)
	}
}

func TestDocumentReaderReportsInvalidValues(t *testing.T) {
	reader := NewDocumentReader(&Document{Sections: []*Section{
		{Label: "properties", Values: []*Value{{Name: "Layer", Kind: KindEnum, Text: "Ground"}}},
	}})
	if err := reader.ConsumeLabel("properties"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := reader.ReadEnum(); err == nil {
		t.Fatal("ReadEnum() expected an error for a non-numeric enum value")
	}
	if err := reader.ConsumeLabel("missing"); err == nil {
		t.Fatal("ConsumeLabel() expected an error for a missing section")
	}
}
