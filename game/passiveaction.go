package game

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"math"
)

func (g *Game) pickUpCollectables() {
	ship := g.ship
	for _, collectable := range g.collectables { // compare it only with all subsequent object, if they match (not with itself and not vice versa)
		if ship.BoundingBox().Intersects(collectable.BoundingBox()) {
			collectable.SetPosition(collectable.FollowPosition(Vec2d{X: ship.Position().X - ViewPortX, Y: ship.Position().Y - ViewPortY},
				Vec2d{X: collectable.Position().X - ViewPortX, Y: collectable.Position().Y - ViewPortY}))
			x := collectable.Position().Sub(ship.Position()).Abs()
			if x.X < 1 && x.Y < 1 && !collectable.IsCollected() {
				g.ship.Inventory().AddToInventory(collectable.Type())
				collectable.SetCollected(true)
				collectable = nil
			}
		}
	}
}


func (g *Game) applyCollisions() {
	// apply our physical hit-test
	for a, objA := range g.objects { // take each object
		for b := a + 1; b < len(g.objects); b++ { // compare it only with all subsequent object, if they match (not with itself and not vice versa)
			objB := g.objects[b]
			if objA.BoundingBox().Intersects(objB.BoundingBox()) { // do a and b collide with each other?
				collisionDir := objA.Position().Sub(objB.Position()).Norm()      // the vector of the collision is in general the difference of the two positions
				totalEnergy := math.Abs(objA.Energy()) + math.Abs(objB.Energy()) // the total energy is absolute value of both ships (not physically correct, because it should be actually a force vector)
				massDistributionA := objA.Mass() / (objA.Mass() + objB.Mass())   // e.g 5 / (5 + 10) = 0.3 or 5 / (5+5)= 0.5
				energyShipA := totalEnergy * (1 - massDistributionA)             // the lighter the ship, the more energy it gets => use the inverse: if a ship only weights 25% it gets 75% of the energy
				energyShipB := totalEnergy * massDistributionA                   // ship b just gets the smaller proportion: ship has 75% of the mass => it gets 25% of the energy
				collisionDirA := collisionDir.Scale(energyShipA, energyShipA)
				collisionDirB := collisionDir.Scale(-energyShipB, -energyShipB) // we need to negate one ship direction, depending of the collision dir
				objA.ApplyDamage(massDistributionA * 100)
				objB.ApplyDamage(massDistributionA * 100)
				objA.Applyforce(collisionDirA)
				objB.Applyforce(collisionDirB)
				objA.React()
				objB.React()
			}
		}
	}
}

// applyTorpedos calculates any gameObjects hits.
func (g *Game) applyTorpedos() {
	for _, t := range g.ship.Torpedos() {
		for i, j := range g.objects {
			if i == 0 {
				continue
			}
			if t.IsActive() && t.BoundingBox().Intersects(j.BoundingBox()) {
				t.Explode()
				collisionDir := (j.Position().Sub(t.Position())).Norm()
				knockback := t.Damage * 10
				collission := collisionDir.Div(j.Mass()/knockback, j.Mass()/knockback)
				j.Applyforce(collission)
				j.ApplyDamage(t.Damage)
				j.React()
				return
			}
		}
	}
}
