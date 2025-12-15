package ecs

type World struct {
	entities    []*Entity
	entityNames map[string]*Entity
	components  map[*Entity]*Components
	systems     []System
}

func NewWorld() *World {
	return &World{
		entities:    make([]*Entity, 0),
		entityNames: make(map[string]*Entity),
		components:  make(map[*Entity]*Components),
		systems:     make([]System, 0),
	}
}

func (w *World) AddEntity(entity *Entity, components *Components) {
	w.entities = append(w.entities, entity)
	w.components[entity] = components
}

func (w *World) AddEntityNamed(name string, entity *Entity, components *Components) {
	w.AddComponent(entity, components)
	w.entityNames[name] = entity
}

func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) GetEntity(name string) (*Entity, bool) {
	entity, ok := w.entityNames[name]
	return entity, ok
}

func (w *World) AddComponent(entity *Entity, component Component) {
	w.components[entity].AddComponent(component)
}

func (w *World) GetComponent(entity *Entity, componentName string) (Component, bool) {
	return w.components[entity].GetComponent(componentName)
}

func (w *World) GetComponents(entity *Entity) *Components {
	return w.components[entity]
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
