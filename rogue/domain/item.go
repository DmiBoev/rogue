// scroll/potion/weapon
package domain

type ItemType string

const (
	ItemTypeWeapon ItemType = "weapon"
	ItemTypeArmor  ItemType = "armor"
	ItemTypePotion ItemType = "potion"
	ItemTypeScroll ItemType = "scroll"
	ItemTypeWand   ItemType = "wand"
	ItemTypeFood   ItemType = "food"
	ItemTypeGold   ItemType = "gold"
	ItemTypeAmulet ItemType = "amulet"
)

type ItemSubtype string

const (

	// Зелья
	SubtypePotionHeal      ItemSubtype = "potion_heal"
	SubtypePotionStrength  ItemSubtype = "posion_strenght"
	SubtypePotionConfusion ItemSubtype = "posion_confusion"
	SubtypePotionBlindness ItemSubtype = "potion_blidness"

	// Свитки
	SubtypeScrollTeleport ItemSubtype
	SubtypeScrollIdentify ItemSubtype
	SubtypeScrollCursed   ItemSubtype

	// Оружие
	SubtypeWeaponDagger ItemSubtype
	SubtypeWeaponSword  ItemSubtype

	// Броня
	SubtypeArmorLeather ItemSubtype
	SubtypeArmorChain   ItemSubtype

	// Еда
	SubtypeFoodRation ItemSubtype
	SubtypeFoodApple  ItemSubtype
)

type Item struct {
	Type    ItemType
	Subtype ItemSubtype

	Name   string
	Symbol rune
	Pos    Position

	Charges int // для жезлов
	Amount  int // для золота, еды
	Effect  int // сила эффекта (например, сколько HP лечит)

	Identified bool // опознан ли предмет (в Rogue зелья и свитки сначала неопознаны)
}

type ItemSpec struct {
	Type   ItemType
	Name   string // "зелье лечения"
	Symbol rune   // '!'
	Effect int    // сила эффекта
	//что нужно
}
