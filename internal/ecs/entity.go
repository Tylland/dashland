package ecs

import "slices"

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
	index      map[ComponentKey]Component
	components []Component
}

func NewEntity(header *EntityHeader) *Entity {
	return &Entity{
		EntityHeader: header,
		index:        make(map[ComponentKey]Component),
	}
}

func (e *Entity) AddComponent(comp Component) {
	e.index[keyOf(comp)] = comp
	e.components = append(e.components, comp)
}

func (e *Entity) RemoveComponent[T any]() {
	k := key[T]()

	if comp, ok := e.index[k]; ok {
		delete(e.index, k)

		if i := slices.Index(e.components, comp); i != -1 {
			e.components = slices.Delete(e.components, i, i+1)
		}
	}
}

func (e *Entity) Component[T any]() (*T, bool) {
	if comp, ok := e.index[key[T]()]; ok {
		return comp.(*T), true
	}

	return nil, false
}

func (e *Entity) WriteTo(w Writer) error {
	e.EntityHeader.WriteTo(w)

	for _, comp := range e.components {
		if wt, ok := comp.(WriterTo); ok {
			wt.WriteTo(w)
		}
	}

	return nil
}

func (e *Entity) ReadFrom(r Reader) error {
	e.EntityHeader.ReadFrom(r)

	for _, comp := range e.components {
		if rf, ok := comp.(ReaderFrom); ok {
			rf.ReadFrom(r)
		}
	}

	return nil
}

// Deprecated: use entity.Component direct.
func GetComponent[T Component](entity *Entity) *T {
	if comp, ok := entity.Component[T](); ok {
		return comp
	}

	return nil
}

// Deprecated: use entity.RemoveComponent direct.
func RemoveComponent[T Component](entity *Entity) {
	entity.RemoveComponent[T]()
}
