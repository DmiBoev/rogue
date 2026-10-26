package domain

type ItemType string

const (
	ItemTypeWeapon ItemType = "weapon"
	ItemTypePotion ItemType = "potion"
	ItemTypeScroll ItemType = "scroll"
	ItemTypeFood   ItemType = "food"
	ItemTypeKey    ItemType = "key"
)

type Stat string

const (
	StatHealth   Stat = "health"
	StatStrength Stat = "strength"
	StatAgility  Stat = "agility"
)

type Item struct {
	Type     ItemType
	Name     string
	Stat     Stat
	Bonus    int
	Color    uint8
	Position Position
}

func (i Item) Symbol() rune {
	switch i.Type {
	case ItemTypeWeapon:
		return ')'
	case ItemTypePotion:
		return '!'
	case ItemTypeScroll:
		return '?'
	case ItemTypeKey:
		return 'k'
	default:
		return '%'
	}
}
