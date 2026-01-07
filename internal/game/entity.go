package game

import "github.com/tylland/dashland/internal/ecs"

const (
	EntityCategoryObject ecs.EntityCategory = 0
	EntityCategoryEnemy  ecs.EntityCategory = 100
)

const (
	EntityUnknown   ecs.EntityType = iota
	EntityPlayer                   = 1
	EntityDiamond                  = 4
	EntityBoulder                  = 5
	EntityExplosion                = 6
	EntityFlash                    = 7
	EntityDoor                     = 11
	EntityExitDoor                 = 12
	EntityFirefly                  = 101
	EntityButterfly                = 102
)
