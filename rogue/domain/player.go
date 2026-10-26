package domain

type Hero struct {
	Position                             Position
	MaxHealth, Health, Strength, Agility int
	Weapon                               *Item
	Inventory                            []Item
	Effects                              []Effect
	Keys                                 uint8
	Sleep                                int
}

type Effect struct {
	Stat         Stat
	Bonus, Turns int
}

// DamageRange is shared by combat and the equipment display.
func (h *Hero) DamageRange(weapon *Item) (int, int) {
	base := h.Strength
	if weapon != nil {
		base += weapon.Bonus + h.Strength/3
	}
	return base, base + 2
}

func (h *Hero) Items(kind ItemType) []Item {
	var items []Item
	for _, item := range h.Inventory {
		if item.Type == kind {
			items = append(items, item)
		}
	}
	return items
}

func (h *Hero) modify(stat Stat, amount int) {
	switch stat {
	case StatHealth:
		// Vampire drains must not let expiration reduce maximum HP below one.
		h.MaxHealth = max(1, h.MaxHealth+amount)
		if h.Health > 0 {
			h.Health = max(1, min(h.MaxHealth, h.Health+amount))
		}
	case StatStrength:
		h.Strength += amount
	case StatAgility:
		h.Agility += amount
	}
}
