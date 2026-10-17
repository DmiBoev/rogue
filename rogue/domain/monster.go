// monstress/AiMonstr
package domain

type EnemyType string

const (
	EnemyTypeZomby     EnemyType = "enemy_zomby"
	EnemyTypeVampire   EnemyType = "enemy_vampire"
	EnemyTypeGhost     EnemyType = "enemy_ghost"
	EnemyTypeOgre      EnemyType = "enemy_ogre"
	EnemyTypeSnakeMage EnemyType = "snake_mage"
)

type Enemy struct {
	Type   EnemyType
	Name   string
	Symbol rune

	//Pos    Position
	//typMonstr string

	Health   int
	Strength int
	Agility  int
	Spite    int //?

	Identified bool //?
}

var EnemyCatalog = map[EnemyType]Enemy{
	EnemyTypeGhost: {
		Type:     EnemyTypeGhost,
		Name:     "Ghost",
		Symbol:   'g',
		Health:   1,
		Strength: 1,
		Agility:  1,
		Spite:    1,
	},
	EnemyTypeOgre: {
		Type:     EnemyTypeOgre,
		Name:     "Ogre",
		Symbol:   'O',
		Health:   1,
		Strength: 1,
		Agility:  1,
		Spite:    1,
	},
	EnemyTypeSnakeMage: {
		Type:     EnemyTypeSnakeMage,
		Name:     "Serpent Mage",
		Symbol:   's',
		Health:   1,
		Strength: 1,
		Agility:  1,
		Spite:    1,
	},
	EnemyTypeVampire: {
		Type:     EnemyTypeVampire,
		Name:     "Vampire",
		Symbol:   'v',
		Health:   1,
		Strength: 1,
		Agility:  1,
		Spite:    1,
	},
	EnemyTypeZomby: {
		Type:     EnemyTypeZomby,
		Name:     "Zomby",
		Symbol:   'z',
		Health:   1,
		Strength: 1,
		Agility:  1,
		Spite:    1,
	},
}
