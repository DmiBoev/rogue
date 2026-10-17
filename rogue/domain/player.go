// player_characteristics
package domain

type TPlayer string

const ConstPlayer TPlayer = "player"

type TypePlayer struct {
	Type      TPlayer
	Name      string
	Symbol    rune
	MaxHealth int
	Health    int
	Strength  int
	agility   int
	weapon    string //?
}

var Player = map[TPlayer]TypePlayer{
	ConstPlayer: {
		Type:      ConstPlayer,
		Name:      "Geroy",
		Symbol:    'i',
		MaxHealth: 1,
		Health:    1,
		Strength:  1,
		agility:   1,
		weapon:    "123",
	},
}

type Backpack struct {
	NumberCell int
	Item       ItemSpec
}
