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
	SubtypeScrollTeleport ItemSubtype = "scroll_teleport"
	SubtypeScrollIdentify ItemSubtype = "scroll_identify"
	SubtypeScrollCursed   ItemSubtype = "scroll_cursed"

	// Оружие
	SubtypeWeaponDagger ItemSubtype = "weapon_dagger"
	SubtypeWeaponSword  ItemSubtype = "weapon_sword"

	// Броня
	SubtypeArmorLeather ItemSubtype = "armor_leather"
	SubtypeArmorChain   ItemSubtype = "armor_chain"

	// Еда
	SubtypeFoodRation ItemSubtype = "food_ration"
	SubtypeFoodApple  ItemSubtype = "food_apple"
)

type ItemSpec struct {
	Type    ItemType
	Subtype ItemSubtype

	Name   string
	Symbol rune
	//Pos    Position

	Charges int // для жезлов
	Amount  int // для золота, еды
	Effect  int // сила эффекта (например, сколько HP лечит)

	Identified bool // опознан ли предмет (в Rogue зелья и свитки сначала неопознаны)
}

var ItemCatalog map[ItemSubtype]ItemSpec{

}