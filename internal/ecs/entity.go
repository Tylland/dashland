package ecs

type EntityCategory uint32
type EntityType uint32
type EntityID string

type EntityHeader struct {
	ID       EntityID
	Type     EntityType
	Category EntityCategory
	Version  int
}

func NewEntityHeader(id EntityID, entType EntityType, category EntityCategory) *EntityHeader {
	return &EntityHeader{
		ID:       id,
		Type:     entType,
		Category: category,
		Version:  1,
	}
}

func (e *EntityHeader) WriteTo(w Writer) error {
	const version = 1

	w.StartSection("header", version)
	w.WriteString("ID", string(e.ID), true)
	w.WriteEnum("Type", e.Type, true, "")
	w.WriteEnum("Category", e.Category, true, "")
	w.WriteInt("Version", e.Version, true)

	return nil
}

func (e *EntityHeader) ReadFrom(r Reader) error {
	r.ConsumeSection("header")

	if value, _, err := r.ReadString(); err == nil {
		e.ID = EntityID(value)
	} else {
		return err
	}

	if value, _, err := r.ReadEnum(); err == nil {
		e.Type = EntityType(value)
	} else {
		return err
	}

	if value, _, err := r.ReadEnum(); err == nil {
		e.Category = EntityCategory(value)
	} else {
		return err
	}

	if value, _, err := r.ReadInt(); err == nil {
		e.Version = value
	} else {
		return err
	}

	return nil
}

type Entity struct {
	*EntityHeader
	Components *Components
}

func NewEntity(header *EntityHeader) *Entity {
	return &Entity{
		EntityHeader: header,
		Components:   NewComponents(),
	}
}

func (e *Entity) AddComponent(c Component) {
	e.Components.AddComponent(c)
}

func (e *Entity) RemoveComponent(c Component) {
	e.Components.RemoveComponent(c)
}

func (e *Entity) WriteTo(w Writer) error {
	e.EntityHeader.WriteTo(w)

	for _, comp := range e.Components.Ordered() {
		if wt, ok := comp.(WriterTo); ok {
			wt.WriteTo(w)
		}
	}

	return nil
}

func (e *Entity) ReadFrom(r Reader) error {
	e.EntityHeader.ReadFrom(r)

	for _, comp := range e.Components.Ordered() {
		if rf, ok := comp.(ReaderFrom); ok {
			rf.ReadFrom(r)
		}
	}

	return nil
}
