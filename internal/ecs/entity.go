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

func (e *Entity) ComponentTuple[T1 any, T2 any]() (*T1, *T2, bool) {
	if comp1, ok := e.Component[T1](); ok {
		if comp2, ok := e.Component[T2](); ok {
			return comp1, comp2, true
		}
	}

	return nil, nil, false
}

func (e *Entity) ComponentTriple[T1 Component, T2 Component, T3 Component]() (*T1, *T2, *T3, bool) {
	if comp1, ok := e.Component[T1](); ok {
		if comp2, ok := e.Component[T2](); ok {
			if comp3, ok := e.Component[T3](); ok {
				return comp1, comp2, comp3, true
			}
		}
	}

	return nil, nil, nil, false
}

func (e *Entity) ComponentQuadruple[T1 Component, T2 Component, T3 Component, T4 Component]() (*T1, *T2, *T3, *T4, bool) {
	if comp1, ok := e.Component[T1](); ok {
		if comp2, ok := e.Component[T2](); ok {
			if comp3, ok := e.Component[T3](); ok {
				if comp4, ok := e.Component[T4](); ok {
					return comp1, comp2, comp3, comp4, true
				}
			}
		}
	}

	return nil, nil, nil, nil, false
}

func (e *Entity) ComponentQuintuple[T1 Component, T2 Component, T3 Component, T4 Component, T5 Component]() (*T1, *T2, *T3, *T4, *T5, bool) {
	if comp1, ok := e.Component[T1](); ok {
		if comp2, ok := e.Component[T2](); ok {
			if comp3, ok := e.Component[T3](); ok {
				if comp4, ok := e.Component[T4](); ok {
					if comp5, ok := e.Component[T5](); ok {
						return comp1, comp2, comp3, comp4, comp5, true
					}
				}
			}
		}
	}

	return nil, nil, nil, nil, nil, false
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
