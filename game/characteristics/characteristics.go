package characteristics

type Characteristics uint16

const (
	None                 = Characteristics(0)
	Void Characteristics = 1 << iota
	RollOff
	CanFall
	Collectable
	Pushable
	IsEnemy
	PlayerObstacle
	EnemyObstacle
)
