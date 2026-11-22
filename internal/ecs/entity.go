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

func (e *EntityHeader) WriteTo(ew Writer) error {
	ew.StartSection("header")
	ew.WriteEntityID(e.ID)
	ew.WriteEntityType(e.Type)
	ew.WriteEntityCategory(e.Category)
	ew.WriteVersion(e.Version)

	return nil
}

func (e *EntityHeader) ReadFrom(r Reader) error {
	if id, err := r.ReadEntityID(); err == nil {
		e.ID = id
	} else {
		return err
	}

	if et, err := r.ReadEntityType(); err == nil {
		e.Type = et
	} else {
		return err
	}

	if category, err := r.ReadEntityCategory(); err == nil {
		e.Category = category
	} else {
		return err
	}

	if ver, err := r.ReadVersion(); err == nil {
		e.Version = ver
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
	//e.EntityHeader.ReadFrom(r)

	for _, comp := range e.Components.Ordered() {
		if rf, ok := comp.(ReaderFrom); ok {
			rf.ReadFrom(r)
		}
	}

	return nil
}
