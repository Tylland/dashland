package ecs

type World struct {
	entities    []*Entity
	entityNames map[string]*Entity
	components  map[EntityID]*Components
	singeltons  map[string]Component
	systems     []System
}

func NewWorld() *World {
	return &World{
		entities:    make([]*Entity, 0),
		entityNames: make(map[string]*Entity),
		components:  make(map[EntityID]*Components),
		singeltons:  make(map[string]Component),
		systems:     make([]System, 0),
	}
}

func (w *World) AddEntity(entity *Entity, components *Components) {
	w.entities = append(w.entities, entity)
	w.components[entity.ID] = components
}

func (w *World) AddEntityNamed(name string, entity *Entity, components *Components) {
	w.AddEntity(entity, components)
	w.entityNames[name] = entity
}

func (w *World) Entities() []*Entity {
	return w.entities
}

func (w *World) GetEntity(name string) *Entity {
	return w.entityNames[name]
}

func (w *World) AddComponent(entity *Entity, component Component) {
	w.components[entity.ID].AddComponent(component)
}

func (w *World) AddSingleton(component Component) {
	name := ComponentName(component)
	w.singeltons[name] = component
}

func (w *World) Singleton(name string) (Component, bool) {
	component, ok := w.singeltons[name]
	return component, ok
}

func (w *World) GetComponent(entity *Entity, componentName string) (Component, bool) {
	return w.components[entity.ID].GetComponent(componentName)
}

func (w *World) GetComponents(entity *Entity) *Components {
	return w.components[entity.ID]
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
