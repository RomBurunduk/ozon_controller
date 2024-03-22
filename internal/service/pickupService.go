package service

import (
	"fmt"

	"pvz_controller/internal/model"
)

type pickupInterface interface {
	WritePVZ(pvz model.Pickups) error
	ReadPVZ() ([]model.Pickups, error)
}

type PickupService struct {
	s pickupInterface
}

func NewPickupService(s pickupInterface) PickupService {
	return PickupService{s: s}
}

const startJob = "Начал работу"
const finishJob = "Закончил работу"

func (p PickupService) AddPVZ(c chan<- string) {
	c <- startJob
	var pvz model.Pickups
	fmt.Println("Введите данные о ПВЗ:")
	fmt.Print("Имя: ")
	_, err := fmt.Scanln(&pvz.Name)
	if err != nil {
		fmt.Println("Ошбика при чтении")
		return
	}
	fmt.Print("Адрес: ")
	_, err = fmt.Scanln(&pvz.Address)
	if err != nil {
		fmt.Println("Ошбика при чтении")
		return
	}
	fmt.Print("Связаться с: ")
	_, err = fmt.Scanln(&pvz.Contact)
	if err != nil {
		fmt.Println("Ошбика при чтении")
		return
	}

	if err := p.s.WritePVZ(pvz); err != nil {
		fmt.Println("Ошибка добавления ПВЗ:", err)
	} else {
		fmt.Println("ПВЗ успешно добавлен.")
	}
	c <- finishJob
}

func (p PickupService) ListPVZ(c chan<- string) {

	c <- startJob
	pvzList, err := p.s.ReadPVZ()
	if err != nil {
		fmt.Println("Ошибка чтения ПВЗ:", err)
		return
	}

	if len(pvzList) == 0 {
		fmt.Println("ПВЗ не обнаружены.")
		return
	}

	fmt.Println("Список ПВЗ:")
	for _, pvz := range pvzList {
		fmt.Printf("Имя: %s\nАдрес: %s\nКонтакты: %s\n\n", pvz.Name, pvz.Address, pvz.Contact)
	}
	c <- finishJob
}
