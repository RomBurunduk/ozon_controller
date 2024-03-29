package model

// BoxStrategy - стратегия для упаковки в коробку
type BoxStrategy struct{}

func (bs BoxStrategy) CheckWeight(order OrderInput) error {
	if order.Weight > 30 {
		return ErrHeavyOrder
	}
	return nil
}

func (bs BoxStrategy) UpdateCost(order *Order) {
	order.Cost += 20
}
