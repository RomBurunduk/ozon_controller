package storage

import "time"

// OrderDTO - транспортная модель заказа
type OrderDTO struct {
	OrderId     int
	RecipientId int
	IsIssued    bool
	IsDeleted   bool
	IsRefund    bool
	IssueDate   time.Time
	ShelfLife   time.Time
}
