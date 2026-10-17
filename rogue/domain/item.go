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
	// Wand
	SubtypeWandFire ItemSubtype = "wand_fire"
	SubtypeWandIce  ItemSubtype = "wand_ice"
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

	// Gold
	SubtypeGold ItemSubtype = "gold"
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

var ItemCatalog = map[ItemSubtype]ItemSpec{
	//Food
	SubtypeFoodRation: {
		Type:   ItemTypeFood,
		Name:   "Strong Ration",
		Symbol: '!',
		Effect: 3,
		Amount: 5,
	},
	SubtypeFoodApple: {
		Type:   ItemTypeFood,
		Name:   "Apple",
		Symbol: '!',
		Effect: 1,
		Amount: 2,
	},
	//Poison
	SubtypePotionHeal: {
		Type:   ItemTypePotion,
		Name:   "Heal posion",
		Symbol: '?',
		Effect: 10,
		Amount: 6,
	},
	SubtypePotionConfusion: {
		Type:   ItemTypePotion,
		Name:   "Poison hmm..?",
		Symbol: '?',
		Effect: -2,
		Amount: 3,
	},
	SubtypePotionStrength: {
		Type:   ItemTypePotion,
		Name:   "Strong poison",
		Symbol: '?',
		Effect: 2,
		Amount: 6,
	},
	SubtypePotionBlindness: {
		Type:   ItemTypePotion,
		Name:   "Sun poison",
		Symbol: '?',
		Effect: -2,
		Amount: 0,
	},
	//Armor
	SubtypeArmorChain: {
		Type:   ItemTypeArmor,
		Name:   "",
		Symbol: 'a',
		Effect: 1,
		Amount: 3,
	},
	SubtypeArmorLeather: {
		Type:   ItemTypeArmor,
		Name:   "",
		Symbol: 'a',
		Effect: 2,
		Amount: 5,
	},
	//weapon
	SubtypeWeaponDagger: {
		Type:   ItemTypeWeapon,
		Name:   "Dagger",
		Symbol: ')',
		Effect: 4,
		Amount: 2,
	},
	SubtypeWeaponSword: {
		Type:   ItemTypeWeapon,
		Name:   "Sword",
		Symbol: '!',
		Effect: 3,
		Amount: 3,
	},
	//scroll
	SubtypeScrollCursed: {
		Type:   ItemTypeScroll,
		Name:   "Cursed Scroll",
		Symbol: '&',
		Effect: -2,
	},
	SubtypeScrollIdentify: {
		Type:   ItemTypeScroll,
		Name:   "Nice Scroll",
		Symbol: '&',
		Effect: 1,
		Amount: 10,
	},
	SubtypeScrollTeleport: {
		Type:   ItemTypeScroll,
		Name:   "Teleportaion",
		Symbol: '&',
	},
	//gold
	SubtypeGold: {
		Type:   ItemTypeGold,
		Name:   "Gold",
		Symbol: '-',
		Amount: 1,
	},
	//wand
	SubtypeWandIce: {
		Type:    ItemTypeWand,
		Name:    "Ice wand",
		Symbol:  '|',
		Charges: 3,
		Effect:  3,
		Amount:  5,
	},
	SubtypeWandFire: {
		Type:    ItemTypeWand,
		Symbol:  '|',
		Charges: 6,
		Effect:  1,
		Amount:  5,
	},
}
