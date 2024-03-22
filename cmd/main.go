package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"pvz_controller/internal/model"
	"pvz_controller/internal/service"
	"pvz_controller/internal/storage"
)

func main() {
	commandPtr := flag.String("c", "help", "задает команду")
	clientID := flag.Int("client", service.DefaultInt, "идетификатор получаеля")
	orderID := flag.Int("order", service.DefaultInt, "номер заказа")
	datePtr := flag.String("date", service.DefaultStr, "срок хранения")
	numPtr := flag.Int("num", service.DefaultInt, "кол-во заказов")
	pagPtr := flag.Int("p", service.DefaultInt, "номер страницы возвратов")

	flag.Parse()
	command := *commandPtr
	orderStorage, err := storage.NewOrderStorage()
	if err != nil {
		fmt.Println("Не удалось подключиться к хранилищу заказов")
		return
	}
	orderService := service.NewOrderService(&orderStorage)

	pickupStorage, err := storage.NewPickupStorage()
	if err != nil {
		fmt.Println("Не удалось подключиться к хранилищу ПВЗ")
		return
	}
	pickupService := service.NewPickupService(&pickupStorage)

	switch command {
	case "new":
		date, err := service.CheckDate(*datePtr)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = orderService.OrderAccept(model.OrderInput{
			OrderId:     *orderID,
			RecipientId: *clientID,
			ShelfLife:   date,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказ успешно принят!")

	case "return":
		err := orderService.ReturnOrder(*orderID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказ возвращен")

	case "issue":
		err := orderService.IssueOrder(flag.Args())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказы выданы")

	case "list":
		list, err := orderService.OrdersList(*clientID, *numPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%+v\n", list)

	case "refund":
		err := orderService.RefundOrder(*orderID, *clientID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказ возвращен")

	case "refunds":
		refunds, err := orderService.Refunds(*pagPtr)
		if err != nil {
			return
		}
		fmt.Println(refunds)

	case "pvz":
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		listOut := make(chan string)
		addOut := make(chan string)
		defer func() {
			defer close(addOut)
			defer close(listOut)
		}()
		go func() {
			fmt.Println("Добро пожаловать в ПВЗ Management System!")
			fmt.Println("Доступные команды: ADD, LIST, EXIT")
			fmt.Print("Введите команду: ")
			var cmd string
			if _, err := fmt.Scanln(&cmd); err != nil {
				fmt.Println("Ошибка при чтении команды:", err)
				return
			}
			cmd = strings.ToUpper(strings.TrimSpace(cmd))
			switch cmd {
			case "ADD":
				go pickupService.AddPVZ(addOut)
			case "LIST":
				go pickupService.ListPVZ(listOut)
			case "EXIT":
				exit <- syscall.SIGINT
				return
			default:
				fmt.Println("Неизвестная команда")
			}
		}()
		go func() {
			for {
				select {
				case info := <-addOut:
					log.Print("Writer:", info)
				case info := <-listOut:
					log.Print("Reader:", info)
				case <-exit:
					return
				}
			}
		}()
		<-exit
		fmt.Println("Выход...")

	case "help":
		fmt.Println("Команда задается через флаг -c")
		fmt.Println("Доступные команды:")
		fmt.Println("	new - Принять заказ от курьера. Принимается ID заказа (order), ID получателя (client) и срок хранения в формате YY-MM-DD (date).")
		fmt.Println("	return - Вернуть заказ курьеру. Принимается ID заказа (order)")
		fmt.Println("	issue - Выдать заказ клиенту. Принимается слайс ID заказов (order)")
		fmt.Println("	list - Получить список заказов. Принимается ID пользователя (client) и опциональный параметр (num) - первые num заказов клиента")
		fmt.Println("	refund - Принять возврат от клиента. Принимается  ID пользователя (client) и ID заказа (order).")
		fmt.Println("	refunds - Получить список возвратов. Опциональный параметр (p) - номер страницы")
		fmt.Println("	pvz - Интерактивный режим взаимодействия с ПВЗ")

	default:
		fmt.Println("Неизвестная команда")

	}
}
