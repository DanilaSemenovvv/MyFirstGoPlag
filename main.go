package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	ID          int       `json:"ID"`
	UserID      int       `json:"userID"`
	SaveDate    time.Time `json:"saveDate"`
	OrderIssued bool      `json:"orderIssued"`
}

func readOrders(file string) ([]Order, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return []Order{}, nil
		}
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func containsOrders(orders []Order, newOrder Order) bool {
	for _, order := range orders {
		if order.ID == newOrder.ID && order.UserID == newOrder.UserID {
			return true
		}
	}

	return false
}

func writeOrder(filename string, orders []Order) error {
	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func ReceivingOrder(newOrder Order, filename string) error {
	orders, err := readOrders(filename)
	if err != nil {
		return err
	}

	if !containsOrders(orders, newOrder) {
		orders = append(orders, newOrder)
		fmt.Println("Заказ принят на ПВЗ")
	} else {
		fmt.Println("Такой заказ уже был внесен в БД ПВЗ")
	}

	return writeOrder(filename, orders)

}

func statusCheck(orders []Order, id int) (bool, int, error) {
	now := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	for i, order := range orders { //в чем разница между конструкцией for _, order = и for order =
		if order.ID == id {
			if order.OrderIssued {
				return false, -1, fmt.Errorf("заказ с ID %d уже выдан", id)
			}
			if order.SaveDate.Before(now) {
				return false, -1, fmt.Errorf("заказ с ID %d еще не просрочен", id)
			}

			return true, i, nil
		}
	}

	return false, -1, fmt.Errorf("заказ с ID %d не найден или уже выдан", id)
}

func deleteOrder(id int, filename string) error {
	orders, err := readOrders(filename)
	if err != nil {
		return err
	}

	status, index, err := statusCheck(orders, id)
	if err != nil {
		return err
	}

	if status {
		orders = append(orders[:index], orders[index+1:]...)
	} else {
		return fmt.Errorf("невозможно удалить заказ с ID %d: не пройдены проверки", id)
	}

	return writeOrder(filename, orders)
}

func getIntInput(text string, scanner *bufio.Scanner) (int, error) {
	fmt.Println(text)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	return strconv.Atoi(input)
}

func getDateInput(text string, scanner *bufio.Scanner) (time.Time, error) {
	fmt.Println(text)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	return time.Parse("02.01.2006", input)
}

func main() {
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("Дорогой пользователь,  выберете действие из списка представленных")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("1. Принять заказ от курьера")
	fmt.Println("2. Вернуть заказ курьеру")
	fmt.Println("3. Выдать заказ/принять возврат клиента")
	fmt.Println("4. Список заказов")
	fmt.Println("5. Список возвратов")
	fmt.Println("6. История заказов")

	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("-------------------------Ваше действие---------------------------")
	fmt.Println("-----------------------------------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		actionChoise := strings.TrimSpace(scanner.Text())
		if actionChoise == "exit" {
			fmt.Println("Досвидания")
			break
		}

		switch actionChoise {
		case "1":
			id, err := getIntInput("Введите ID-заказа: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода ID-заказа: ", err)
				continue
			}

			userID, err := getIntInput("Введите ID-пользователя: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода ID-пользователя: ", err)
				continue
			}
			saveDate, err := getDateInput("Введите дату хранения: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода даты хранения: ", err)
				continue
			}

			orderIssued := false

			order := Order{
				ID:          id,
				UserID:      userID,
				SaveDate:    saveDate,
				OrderIssued: orderIssued,
			}

			err = ReceivingOrder(order, "data.json")
			if err != nil {
				fmt.Println(err)
			}

		case "2":
			id, err := getIntInput("Введите ID заказа для курьера", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода ID-заказа")
				continue
			}

			err = deleteOrder(id, "data.json")
			if err != nil {
				fmt.Println(err)
			}

		case "3":
			fmt.Println(3)
		case "4":
			fmt.Println(4)
		case "5":
			fmt.Println(5)
		case "6":
			fmt.Println(6)
		default:
			fmt.Println("Неверный выбор, попробуйте снова")
		}
		fmt.Println("Выберете следующее действие или ввведите 'exit' для выхода (или нажмите сочетание клавишь ctrl+c)")
	}
}
