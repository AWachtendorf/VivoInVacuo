package inventory

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Inventory stores collected items in a map.
type Inventory struct {
	inventory map[ItemType]int
	visible   bool
}

// NewInventory creates a New Inventory. The map is not initialized as a nil map,
func NewInventory() *Inventory {
	invMap := make(map[ItemType]int)
	inv := &Inventory{inventory: invMap}
	return inv
}

// AddToInventory increases the amount of an ItemType.
func (i *Inventory) AddToInventory(itemtype ItemType) {
	i.inventory[itemtype] += 1
}

// AllItems returns all items.
func (i *Inventory) AllItems() map[ItemType]int {
	return i.inventory
}

// Visible toggles Visibility.
func (i *Inventory) Visible() bool {
	return i.visible
}

// SetVisible sets visibility.
func (i *Inventory) SetVisible(state bool) {
	i.visible = state
}

// Update listens to a key and changes the visibility.
func (i *Inventory) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		i.visible = !i.visible
	}
}
