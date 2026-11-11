package ecs

type Writer interface {
	WriteEntityID(id EntityID) error
	WriteEntityType(et EntityType) error
	WriteVersion(ver int) error
	WriteLabel(lbl string) error
	WriteIntValue(v int, d bool) error
	WriteFloatValue(v float32, d bool) error
}

type WriterTo interface {
	WriteTo(w Writer) error
}

type Reader interface {
	ReadEntityID() (EntityID, error)
	ReadEntityType() (EntityType, error)
	ReadVersion() (int, error)
	ReadLabel() (string, error)
	ReadIntValue() (int, bool, error)
	ReadFloatValue() (float32, bool, error)
}

type ReaderFrom interface {
	ReadFrom(r Reader) error
}
