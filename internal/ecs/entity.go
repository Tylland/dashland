package ecs

type EntityCategory uint32
type EntityType uint32
type EntityID string

type Entity struct {
	ID         EntityID
	Type       EntityType
	TypeName   string
	Version    int
	Components *Components
}

func NewEntity(id EntityID, entityType EntityType) *Entity {
	return &Entity{
		ID:         id,
		Type:       entityType,
		Components: NewComponents(),
	}
}

func (e *Entity) AddComponent(c Component) {
	e.Components.AddComponent(c)
}

func (e *Entity) RemoveComponent(c Component) {
	e.Components.RemoveComponent(c)
}

func (e *Entity) WriteTo(ew Writer) error {
	ew.WriteEntityID(e.ID)
	ew.WriteEntityType(e.Type)
	ew.WriteVersion(e.Version)

	for _, comp := range e.Components.index {
		if wt, ok := comp.(WriterTo); ok {
			wt.WriteTo(ew)
		}
	}

	return nil
}

func (e *Entity) ReadFrom(r Reader) error {
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

	r.ReadVersion()

	for _, comp := range e.Components.Ordered() {
		if rf, ok := comp.(ReaderFrom); ok {
			rf.ReadFrom(r)
		}
	}

	return nil
}
