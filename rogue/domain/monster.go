package domain

type EnemyType string

const (
	EnemyTypeZomby     EnemyType = "zombie"
	EnemyTypeVampire   EnemyType = "vampire"
	EnemyTypeGhost     EnemyType = "ghost"
	EnemyTypeOgre      EnemyType = "ogre"
	EnemyTypeSnakeMage EnemyType = "snake_mage"
	EnemyTypeMimic     EnemyType = "mimic"
)

type Enemy struct {
	Type                                        EnemyType
	Name                                        string
	Symbol                                      rune
	Position                                    Position
	Health, MaxHealth, Strength, Agility, Spite int
	Engaged, Hidden, FirstDodged, Counter       bool
	Rest, Pattern                               int
}
