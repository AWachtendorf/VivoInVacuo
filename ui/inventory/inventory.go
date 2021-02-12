package inventory

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Inventory struct {
	inventory map[ItemType]int
	visible   bool
}

func NewInventory() *Inventory {
	mappimap := make(map[ItemType]int)
	inv := &Inventory{inventory: mappimap}
	return inv
}

func (i *Inventory) AddToInventory(itemtype ItemType) {
	i.inventory[itemtype] += 1
}

func (i *Inventory) AllItems() map[ItemType]int {
	return i.inventory
}

func (i *Inventory) Visible() bool {
	return i.visible
}

func (i *Inventory) SetVisible(state bool) {
	i.visible = state
}

func (i *Inventory) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		i.visible = !i.visible
	}
}
