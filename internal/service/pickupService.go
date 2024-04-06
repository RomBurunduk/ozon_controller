package service

import (
	"bufio"
	"fmt"
	"os"

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
	defer func() {
		c <- finishJob
	}()
	var pvz model.Pickups
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите данные о ПВЗ:")
	fmt.Print("Имя: ")
	//_, err := fmt.Scanln(&pvz.Name)
	//if err != nil {
	//	fmt.Println("Ошбика при чтении")
	//	return
	//}
	scanner.Scan()
	pvz.Name = scanner.Text()
	fmt.Print("Адрес: ")
	scanner.Scan()
	pvz.Address = scanner.Text()
	fmt.Print("Связаться с: ")
	scanner.Scan()
	pvz.Contact = scanner.Text()
	if err := p.s.WritePVZ(pvz); err != nil {
		fmt.Println("Ошибка добавления ПВЗ:", err)
	} else {
		fmt.Println("ПВЗ успешно добавлен.")
	}

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
