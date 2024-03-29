package model

// PackageStrategy - стратегия для упаковки в пакет
type PackageStrategy struct{}

func (ps PackageStrategy) CheckWeight(order OrderInput) error {
	if order.Weight > 10 {
		return ErrHeavyOrder
	}
	return nil
}

func (ps PackageStrategy) UpdateCost(order *Order) {
	order.Cost += 5
}
