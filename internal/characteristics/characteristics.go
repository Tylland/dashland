package characteristics

import "strings"

type Characteristics uint16

const (
	None                    = Characteristics(0)
	CanFall Characteristics = 1 << iota
	Falling
	Pushable
	IsEnemy
	Obstacle
	CanHoldGravity
	GravityRollOff
	Destructable
)

func (c Characteristics) String() string {
	str := make([]string, 0)

	if c&None == None {
		str = append(str, "None")
	}
	if c&CanFall == CanFall {
		str = append(str, "CanFall")
	}
	if c&Falling == Falling {
		str = append(str, "Falling")
	}
	if c&Pushable == Pushable {
		str = append(str, "Pushable")
	}
	if c&IsEnemy == IsEnemy {
		str = append(str, "IsEnemy")
	}
	if c&Obstacle == Obstacle {
		str = append(str, "Obstacle")
	}
	if c&CanHoldGravity == CanHoldGravity {
		str = append(str, "CanHoldGravity")
	}
	if c&GravityRollOff == GravityRollOff {
		str = append(str, "GravityRollOff")
	}
	if c&Destructable == Destructable {
		str = append(str, "Destructable")
	}

	return strings.Join(str, "|")
}
