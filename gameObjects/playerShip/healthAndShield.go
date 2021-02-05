package playerShip

func (s *Ship) ApplyDamage(damage float64) {
	if s.shieldBar.Percentage() <= 10 {
		s.healthBar.ApplyDamage(damage)
	} else {
		s.shieldBar.ApplyDamage(damage)
	}
	if s.healthBar.Percentage() <= 10{
		println("THIS SHIP WOULD BE FUCKING DEAD!")
	}
}

