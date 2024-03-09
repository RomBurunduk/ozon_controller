package main

import (
	"flag"
	"fmt"

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
	db, err := storage.New()
	if err != nil {
		fmt.Println("Не удалось подключиться к хранилищу")
		return
	}
	serv := service.New(&db)

	switch command {
	case "new":
		date, err := service.CheckDate(*datePtr)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = serv.OrderAccept(model.OrderInput{
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
		err := serv.ReturnOrder(*orderID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказ возвращен")

	case "issue":
		err := serv.IssueOrder(flag.Args())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказы выданы")

	case "list":
		list, err := serv.OrdersList(*clientID, *numPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%+v\n", list)

	case "refund":
		err := serv.RefundOrder(*orderID, *clientID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Заказ возвращен")

	case "refunds":
		refunds, err := serv.Refunds(*pagPtr)
		if err != nil {
			return
		}
		fmt.Println(refunds)

	case "help":
		fmt.Println("Команда задается через флаг -c")
		fmt.Println("Доступные команды:")
		fmt.Println("	new - Принять заказ от курьера. Принимается ID заказа (order), ID получателя (client) и срок хранения в формате YY-MM-DD (date).")
		fmt.Println("	return - Вернуть заказ курьеру. Принимается ID заказа (order)")
		fmt.Println("	issue - Выдать заказ клиенту. Принимается слайс ID заказов (order)")
		fmt.Println("	list - Получить список заказов. Принимается ID пользователя (client) и опциональный параметр (num) - первые num заказов клиента")
		fmt.Println("	refund - Принять возврат от клиента. Принимается  ID пользователя (client) и ID заказа (order).")
		fmt.Println("	refunds - Получить список возвратов. Опциональный параметр (p) - номер страницы")

	default:
		fmt.Println("Неизвестная команда")

	}
}
