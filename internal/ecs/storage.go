package ecs

type Writer interface {
	WriteEntityID(id EntityID) error
	WriteEntityType(et EntityType) error
	WriteBlockType(v int) error
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
	ReadEntityID() (EntityID, error)
	ReadEntityType() (EntityType, error)
	ReadBlockType() (int, error)
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
