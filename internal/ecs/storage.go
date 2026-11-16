package ecs

type Writer interface {
	WriteEntityID(id EntityID) error
	WriteEntityType(et EntityType) error
	WriteEntityCategory(c EntityCategory) error
	WriteVersion(ver int) error
	WriteLabel(lbl string) error
	WriteBool(v bool, d bool) error
	WriteEnum(v any, d bool, c string) error
	WriteFloat(v float32, d bool) error
	WriteInt(v int, d bool) error
	WriteString(v string, d bool) error
}

type WriterTo interface {
	WriteTo(w Writer) error
}

type Reader interface {
	ConsumeLine() error
	ConsumeLabel(lbl string) error
	ReadEntityID() (EntityID, error)
	ReadEntityType() (EntityType, error)
	ReadEntityCategory() (EntityCategory, error)
	ReadVersion() (int, error)
	ReadBool() (bool, bool, error)
	ReadEnum() (uint32, bool, error)
	ReadFloat() (float32, bool, error)
	ReadInt() (int, bool, error)
	ReadLabel() (string, error)
	ReadString() (string, bool, error)
}

type ReaderFrom interface {
	ReadFrom(r Reader) error
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
