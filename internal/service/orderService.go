package service

import (
	"errors"
	"strconv"
	"time"

	"pvz_controller/internal/model"
	storageFile "pvz_controller/internal/storage"
)

// DefaultInt - скорее всего написал фигню, не знаю как реализуют подобные вещи в проде
const DefaultInt = -1
const DefaultStr = "-1"

// orderInterface - интерфейс, соединиящий storage.Storage и OrderService
type orderInterface interface {
	OrderAccept(input model.OrderInput) error
	List() ([]model.Order, error)
	ChooseOrder(id int) (model.Order, error)
	DeleteOrder(id int) error
	IssueOrder(id int) error
	Refunds(num int) ([]model.Order, error)
	OrderRefund(id int, clientID int) error
}

// OrderService - структура сервиса
type OrderService struct {
	s orderInterface
}

// NewOrderService - создание сервиса
func NewOrderService(s orderInterface) OrderService {
	return OrderService{s: s}
}

// OrderAccept - прием заказа в ПВЗ
func (s OrderService) OrderAccept(input model.OrderInput) error {
	err := CheckOrderId(input.OrderId)
	if err != nil {
		return err
	}
	err = CheckClientId(input.RecipientId)
	if err != nil {
		return err
	}
	if CheckShelfTime(input.ShelfLife) {
		return errors.New("срок хранения в прошлом")
	}
	_, err = s.s.ChooseOrder(input.OrderId)
	if err != nil {
		if errors.Is(err, storageFile.ErrOrderNotFound) {
			return s.s.OrderAccept(input)
		}
		return err
	}
	return errors.New("этот заказ уже существует")
}

// ReturnOrder - возврат заказа курьеру
func (s OrderService) ReturnOrder(id int) error {
	err := CheckOrderId(id)
	if err != nil {
		return err
	}
	order, err := s.s.ChooseOrder(id)
	if err != nil {
		return err
	}
	if !order.IsIssued && CheckShelfTime(order.ShelfLife) {
		err = s.s.DeleteOrder(id)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("невозможно вернуть заказ")
}

// IssueOrder - получения заказа клиентом
func (s OrderService) IssueOrder(idStr []string) error {
	ids := make([]int, 0, len(idStr))
	for _, s := range idStr {
		num, err := strconv.Atoi(s)
		if err != nil {
			return errors.New("ошбика преобразования")
		}
		ids = append(ids, num)
	}
	orders := make([]model.Order, 0, len(ids))
	for _, id := range ids {
		order, err := s.s.ChooseOrder(id)
		if err != nil {
			return err
		}
		if !CheckShelfTime(order.ShelfLife) {
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
func (s OrderService) Refunds(num int) ([]model.Order, error) {
	err := CheckNum(num)
	if err != nil {
		return nil, err
	}
	if num <= 0 {
		return nil, errors.New("неправильно указан номер страницы")
	}
	refunds, err := s.s.Refunds(num)
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

// RefundOrder - возврат заказа в ПВЗ
func (s OrderService) RefundOrder(orderId int, clientId int) error {
	err := CheckOrderId(orderId)
	if err != nil {
		return err
	}
	err = CheckClientId(clientId)
	if err != nil {
		return err
	}
	err = s.s.OrderRefund(orderId, clientId)
	if err != nil {
		return err
	}
	return nil
}

// OrdersList - список всех актуальных заказов
func (s OrderService) OrdersList(clientId int, num int) ([]model.Order, error) {
	err := CheckClientId(clientId)
	if err != nil {
		return nil, err
	}
	if num <= 0 {
		return nil, errors.New("неправильно указан номер")
	}
	list, err := s.s.List()
	if err != nil {
		return nil, err
	}
	onlyClient := make([]model.Order, 0, len(list))
	for _, order := range list {
		if order.RecipientId == clientId && !order.IsIssued {
			onlyClient = append(onlyClient, order)
		}
	}

	if num != DefaultInt {
		if num < len(onlyClient) {
			return onlyClient[:num], nil
		}
		return nil, nil
	}
	return onlyClient, nil
}

func CheckClientId(id int) error {
	if id == DefaultInt {
		return errors.New("не указан ID клиента")
	}
	return nil
}

func CheckOrderId(id int) error {
	if id == DefaultInt {
		return errors.New("не указан ID заказа")
	}
	return nil
}

func CheckDate(s string) (time.Time, error) {
	if s == DefaultStr {
		return time.Time{}, errors.New("не указана дата")
	}
	parse, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return time.Time{}, err
	}
	return parse, nil
}

func CheckNum(i int) error {
	if i == DefaultInt {
		return errors.New("не указан номер")
	}
	return nil
}

func CheckShelfTime(shelf time.Time) bool {
	return time.Now().After(shelf)
}
