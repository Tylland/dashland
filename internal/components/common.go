package components

import "github.com/tylland/dashland/internal/ecs"

func Section(ct ecs.ComponentType) (string, int) {
	return ct.Name, ct.Version
}

func SectionName(ct ecs.ComponentType) string {
	return ct.Name
}
