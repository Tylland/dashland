package ecs

type Writer interface {
	StartSection(lbl string) error
	WriteBool(name string, v bool, d bool) error
	WriteEnum(name string, v any, d bool, c string) error
	WriteFloat(name string, v float32, d bool) error
	WriteInt(name string, v int, d bool) error
	WriteString(name string, v string, d bool) error
}

type WriterTo interface {
	WriteTo(w Writer) error
}

type Reader interface {
	Rewind() error
	ConsumeLine() error
	ConsumeLabel(lbl string) error
	ReadBool() (bool, bool, error)
	ReadEnum() (uint32, bool, error)
	ReadFloat() (float32, bool, error)
	ReadInt() (int, bool, error)
	ReadLabel() (string, error)
	ReadString() (string, bool, error)
}

type WriterReader interface {
	Writer
	Reader
}

type ReaderFrom interface {
	ReadFrom(r Reader) error
}

type Transferer interface {
	WriterTo
	ReaderFrom
}

func WithChangedValue[T any](fetch func() (T, bool, error), update func(T)) error {
	val, changed, err := fetch()

	if err != nil {
		return err
	}

	if changed {
		update(val)
	}

	return nil
}
