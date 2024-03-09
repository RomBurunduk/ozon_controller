package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"time"

	"pvz_controller/internal/model"
)

type Storage struct {
	storage *os.File
}

var ErrOrderNotFound = errors.New("нет такого заказа")

const storageName = "storageName"

// New - соеднинение с файлом данных
func New() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, fs.ModePerm)
	if err != nil {
		return Storage{}, err
	}
	return Storage{storage: file}, err
}

// OrderAccept принимает заказ получателя по model.OrderInput
func (s *Storage) OrderAccept(input model.OrderInput) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	newOrder := OrderDTO{
		OrderId:     input.OrderId,
		RecipientId: input.RecipientId,
		IsIssued:    false,
		IsDeleted:   false,
		IsRefund:    false,
		IssueDate:   time.Time{},
		ShelfLife:   input.ShelfLife,
	}
	all = append(all, newOrder)
	err = s.writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// List возвращает слайс актуальных заказов
func (s *Storage) List() ([]model.Order, error) {
	all, err := s.listAll()
	if err != nil {
		return nil, err
	}
	onlyActive := make([]model.Order, 0, len(all))
	for _, dto := range all {
		if !dto.IsDeleted {
			onlyActive = append(onlyActive, model.Order{
				OrderId:     dto.OrderId,
				RecipientId: dto.RecipientId,
				IsIssued:    dto.IsIssued,
				IssueDate:   dto.IssueDate,
				IsRefund:    dto.IsRefund,
				ShelfLife:   dto.ShelfLife,
			})
		}
	}
	return onlyActive, nil
}

// ChooseOrder возвращает заказ по его ID
func (s *Storage) ChooseOrder(id int) (model.Order, error) {
	all, err := s.listAll()
	if err != nil {
		return model.Order{}, err
	}
	for _, dto := range all {
		if dto.OrderId == id && !dto.IsDeleted {
			return model.Order{
				OrderId:     dto.OrderId,
				RecipientId: dto.RecipientId,
				IsIssued:    dto.IsIssued,
				IsRefund:    dto.IsRefund,
				ShelfLife:   dto.ShelfLife,
			}, nil
		}
	}
	return model.Order{}, ErrOrderNotFound
}

// DeleteOrder удаляет (soft-delete) заказ по его ID - по сути заказ больше не находится на ПВЗ
func (s *Storage) DeleteOrder(id int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	for i, dto := range all {
		if dto.OrderId == id {
			all[i].IsDeleted = true
			err = s.writeBytes(all)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("нет такого заказа")
}

// IssueOrder - приемка заказа получателем
func (s *Storage) IssueOrder(id int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	for i, dto := range all {
		if dto.OrderId == id {
			all[i].IsIssued = true
			all[i].IssueDate = time.Now()
			err = s.writeBytes(all)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("заказ не найден")
}

// listAll - возвращает слайс ВСЕХ заказов в формате DTO
func (s *Storage) listAll() ([]OrderDTO, error) {
	reader := bufio.NewReader(s.storage)
	RawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	_, err = s.storage.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var orders []OrderDTO
	if len(RawBytes) == 0 {
		return orders, nil
	}
	err = json.Unmarshal(RawBytes, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Refunds - возвращает слайс возвратов
func (s *Storage) Refunds(num int) ([]model.Order, error) {
	all, err := s.listAll()
	if err != nil {
		return nil, err
	}
	onlyRefunds := make([]model.Order, 0, len(all))
	for _, dto := range all {
		if !dto.IsDeleted && dto.IsRefund {
			onlyRefunds = append(onlyRefunds, model.Order{
				OrderId:     dto.OrderId,
				RecipientId: dto.RecipientId,
				IsIssued:    dto.IsIssued,
				IssueDate:   dto.IssueDate,
				IsRefund:    dto.IsRefund,
				ShelfLife:   dto.ShelfLife,
			})
		}
	}
	const pageSize = 3
	startIndex := (num - 1) * pageSize
	endIndex := num * pageSize

	if startIndex >= len(onlyRefunds) {
		return nil, nil
	}

	if endIndex > len(onlyRefunds) {
		endIndex = len(onlyRefunds)
	}
	return onlyRefunds[startIndex:endIndex], nil
}

// OrderRefund - возврат заказа в ПВЗ
func (s *Storage) OrderRefund(id int, clientID int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	for i, dto := range all {
		if dto.OrderId == id {
			if clientID != dto.RecipientId {
				return errors.New("несоответсвие ID клиента и ID заказа")
			}
			if !dto.IsIssued {
				return errors.New("заказ еще не выдан")
			}
			if time.Now().Sub(dto.IssueDate).Hours() >= 48 {
				return errors.New("срок возврата истек")
			}
			all[i].IsRefund = true
			err := s.writeBytes(all)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("заказ не найден")
}

// writeBytes - запись в файл
// TODO - записывать данные не виде списка, а в виде словаря orderID -> order, чтобы в дальнейшем поиск заказа был за O(1), а не O(n)
func (s *Storage) writeBytes(all []OrderDTO) error {
	bytes, err := json.Marshal(all)
	if err != nil {
		return err
	}
	err = os.WriteFile(storageName, bytes, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
