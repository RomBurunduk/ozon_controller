package model

import "errors"

type PackingType string

const (
	Package PackingType = "package"
	Box     PackingType = "box"
	Film    PackingType = "film"
)

var ErrHeavyOrder = errors.New("слишком тяжелый заказ")
