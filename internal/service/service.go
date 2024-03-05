package service

import (
	"errors"
	"pvz_controller/internal/model"
	"time"
)

// storage - интерфейс, соединиящий storage.Storage и Service
type storage interface {
	OrderAccept(input model.OrderInput) error
	List() ([]model.Order, error)
	ChooseOrder(int) (model.Order, error)
	DeleteOrder(int) error
	IssueOrder(int) error
	Refunds() ([]model.Order, error)
	OrderRefund(id int, clientID int) error
}

// Service - структура сервиса
type Service struct {
	s storage
}

// New - создание сервиса
func New(s storage) Service {
	return Service{s: s}
}

// OrderAccept - прием заказа в ПВЗ
func (s Service) OrderAccept(input model.OrderInput) error {
	if time.Now().After(input.ShelfLife) {
		return errors.New("срок хранения в прошлом")
	}
	return s.s.OrderAccept(input)
}

// ReturnOrder - возврат заказа курьеру
func (s Service) ReturnOrder(id int) error {
	order, err := s.s.ChooseOrder(id)
	if err != nil {
		return err
	}
	if !order.IsIssued && time.Now().After(order.ShelfLife) {
		err = s.s.DeleteOrder(id)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("невозможно вернуть заказ")
	}
}

// IssueOrder - получения заказа клиентом
func (s Service) IssueOrder(ids ...int) error {
	orders := make([]model.Order, 0)
	for _, id := range ids {
		order, err := s.s.ChooseOrder(id)
		if err != nil {
			return err
		}
		if time.Now().Before(order.ShelfLife) {
			orders = append(orders, order)
		}
	}
	if len(orders) == 0 {
		return errors.New("срок хранения заказов вышел")
	}
	clientId := orders[0].RecipientId
	for _, order := range orders {
		if order.RecipientId != clientId {
			return errors.New("не все заказы принадлежат одному получателю")
		}
	}
	for _, order := range orders {
		err := s.s.IssueOrder(order.OrderId)
		if err != nil {
			return err
		}
	}
	return nil
}

// Refunds - список возвратов
func (s Service) Refunds() ([]model.Order, error) {
	refunds, err := s.s.Refunds()
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

// RefundOrder - возврат заказа в ПВЗ
func (s Service) RefundOrder(orderId int, clientId int) error {
	err := s.s.OrderRefund(orderId, clientId)
	if err != nil {
		return err
	}
	return nil
}

// OrdersList - список всех актуальных заказов
func (s Service) OrdersList(clientId int) ([]model.Order, error) {
	list, err := s.s.List()
	if err != nil {
		return nil, err
	}
	onlyClient := make([]model.Order, 0, len(list))
	for _, order := range list {
		if order.RecipientId == clientId {
			onlyClient = append(onlyClient, order)
		}
	}
	return onlyClient, nil
}
