package model

// FilmStrategy - стратегия для упаковки в пленку
type FilmStrategy struct{}

func (fs FilmStrategy) CheckWeight(OrderInput) error {
	return nil
}

func (fs FilmStrategy) UpdateCost(order *Order) {
	order.Cost += 1
}
