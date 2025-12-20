package characteristics

import "strings"

type Characteristics uint16

const (
	None                 = Characteristics(0)
	Void Characteristics = 1 << iota
	RollOff
	CanFall
	Collectable
	Pushable
	IsEnemy
	EnemyObstacle
	IsPlayer
	PlayerObstacle
	CanHoldGravity
)

func (c Characteristics) String() string {
	str := make([]string, 0)

	if c&None == None {
		str = append(str, "None")
	}
	if c&None == None {
		str = append(str, "None")
	}
	if c&Void == Void {
		str = append(str, "Void")
	}
	if c&RollOff == RollOff {
		str = append(str, "RollOff")
	}
	if c&CanFall == CanFall {
		str = append(str, "CanFall")
	}
	if c&Collectable == Collectable {
		str = append(str, "Collectable")
	}
	if c&Pushable == Pushable {
		str = append(str, "Pushable")
	}
	if c&IsEnemy == IsEnemy {
		str = append(str, "IsEnemy")
	}
	if c&PlayerObstacle == PlayerObstacle {
		str = append(str, "PlayerObstacle")
	}
	if c&EnemyObstacle == EnemyObstacle {
		str = append(str, "EnemyObstacle")
	}

	return strings.Join(str, "|")
}
