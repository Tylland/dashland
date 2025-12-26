package ecs

import "slices"

type World struct {
	entities     []*Entity
	entityNames  map[string]*Entity
	removalQueue []*Entity
	events       []*Event
	singeltons   map[string]Component
	systems      []System
}

func NewWorld() *World {
	return &World{
		entities:     make([]*Entity, 0),
		entityNames:  make(map[string]*Entity),
		removalQueue: make([]*Entity, 0),
		events:       make([]*Event, 100),
		singeltons:   make(map[string]Component),
		systems:      make([]System, 0),
	}
}

func (w *World) AddEntity(entity *Entity) {
	w.entities = append(w.entities, entity)
}

func (w *World) AddEntityNamed(name string, entity *Entity) {
	w.AddEntity(entity)
	w.entityNames[name] = entity
}

func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) GetEntity(name string) *Entity {
	return w.entityNames[name]
}

func (w *World) EnqueueRemoval(entity *Entity) {
	w.removalQueue = append(w.removalQueue, entity)
}

func (w *World) RemovalQueue() []*Entity {
	return w.removalQueue
}

func (w *World) ResetRemovalQueue() {
	w.removalQueue = make([]*Entity, 0)
}

func (w *World) RemoveEntity(entity *Entity) {
	i := slices.Index(w.entities, entity)

	if i > -1 {
		w.entities = append(w.entities[:i], w.entities[i+1:]...)
	}

	for k, v := range w.entityNames {
		if v == entity {
			delete(w.entityNames, k)
		}
	}
}

func (w *World) AddSingleton(component Component) {
	name := ComponentName(component)
	w.singeltons[name] = component
}

func (w *World) Singleton(name string) (Component, bool) {
	component, ok := w.singeltons[name]
	return component, ok
}

func (w *World) AddEvent(name string, data any) {
	w.events = append(w.events, &Event{Name: name, Data: data})
}

func (w *World) Events() []*Event {
	return w.events
}

func (w *World) ClearEvents() {
	w.events = make([]*Event, 100)
}

func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) AddSystems(systems ...System) {
	w.systems = append(w.systems, systems...)
}

func (w *World) Update(deltaTime float32) {
	for _, system := range w.systems {
		system.Update(w, deltaTime)
	}
}

// func GetComponent[T Component](world World, entity Entity) T {
// 	name := fmt.Sprintf("%T", *new(T))
// 	return world.GetComponent(entity, name).(T)
// }
