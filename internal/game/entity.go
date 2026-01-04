package game

import "github.com/tylland/dashland/internal/ecs"

const (
	EntityCategoryObject ecs.EntityCategory = 0
	EntityCategoryEnemy  ecs.EntityCategory = 100
)

const (
	EntityPlayer    ecs.EntityType = iota
	EntityDiamond                  = 4
	EntityBoulder                  = 5
	EntityExplosion                = 6
	EntityFlash                    = 7
	EntityDoor                     = 11
	EntityFirefly                  = 101
	EntityButterfly                = 102
)
