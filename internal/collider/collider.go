package collider

import "strings"

type CollisionLayer uint32

const (
	LayerNone CollisionLayer = 0
	LayerWall CollisionLayer = 1 << iota
	LayerPlayer
	LayerEnemy
	LayerBedrock
	LayerGround
	LayerProjectile
	LayerItem
	LayerCollectable
	// LayerWater
	// LayerTrap
	LayerTrigger
	LayerAll CollisionLayer = 0xFFFF
)

func (l CollisionLayer) String() string {
	var str []string

	if l == LayerAll {
		return "All"
	}

	if l&LayerWall == LayerWall {
		str = append(str, "Wall")
	}
	if l&LayerPlayer == LayerPlayer {
		str = append(str, "Player")
	}
	if l&LayerEnemy == LayerEnemy {
		str = append(str, "Enemy")
	}
	if l&LayerBedrock == LayerBedrock {
		str = append(str, "Bedrock")
	}
	if l&LayerGround == LayerGround {
		str = append(str, "Ground")
	}
	if l&LayerProjectile == LayerProjectile {
		str = append(str, "Projectile")
	}
	if l&LayerItem == LayerItem {
		str = append(str, "Item")
	}
	if l&LayerCollectable == LayerCollectable {
		str = append(str, "Collectable")
	}
	if l&LayerTrigger == LayerTrigger {
		str = append(str, "Trigger")
	}

	if len(str) > 0 {
		return strings.Join(str, " | ")
	}

	return "None"
}
