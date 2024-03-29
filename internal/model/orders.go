package model

import "time"

// TODO - соотвествующе изменить методы хранилища под корректную обработку новых полей (в идеале перенести
//  взаимодействие в бд через веб-интерфейс)

// Order - модель заказа
type Order struct {
	OrderId     int
	RecipientId int
	IsIssued    bool
	IsRefund    bool
	IssueDate   time.Time
	ShelfLife   time.Time
	Packing     PackingType
	Weight      int
	Cost        int
}

// OrderInput - входящая модель заказа
type OrderInput struct {
	OrderId     int
	RecipientId int
	ShelfLife   time.Time
	Packing     PackingType
	Weight      int
	Cost        int
}
