package model

import "time"

// Order - модель заказа
type Order struct {
	OrderId     int
	RecipientId int
	IsIssued    bool
	IsRefund    bool
	IssueDate   time.Time
	ShelfLife   time.Time
}

// OrderInput - входящая модель заказа
type OrderInput struct {
	OrderId     int
	RecipientId int
	ShelfLife   time.Time
}
