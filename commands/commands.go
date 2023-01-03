package commands

import (
	"buy-list/storage/postgresql"
	"strconv"
	"strings"
)

func CommandHandler(user_msg string, user_id int64, chat_id int64, user_nickname string, user_name string, db *postgresql.Connection) string {
	msg := "Команда не распознана"
	if user_msg == "/help" || user_msg == "/start" {
		msg = "😳ВАЖНО😳 Перед тем, как впервые воспользоваться ботом, вам нужно написать команду /reg, чтобы бот вас запомнил.\n\nНажмите '/' или на кнопку 'Меню' слева внизу, чтобы увидеть доступные команды"
		return msg
	} else if user_msg == "/reg" {
		msg = db.CreateUser(user_nickname, user_name, user_id, chat_id)
		return msg
	} else if user_msg == "/addinlist" {
		msg = "Введите название продукта(одним словом), вес(в граммах) и дату напоминания через пробел\nПример: Чипсы 80 31-01-2023 09:00\n\nЛибо просто введите название продукта😊"
		db.SetStatus(1, user_id, chat_id)
		return msg
	} else if user_msg == "/addinfridge" {
		msg = "Если вы хотите перенести продукт из списка продуктов в холодильник, то введите 'купил', название продукта(одним словом) и срок хранения через пробел\nПример: купил Чипсы 31-01-2022 31-01-2023\n\nЕсли же вы хотите добавить новый продукт в холодильник, то введите название продукта и срок хранения через пробел\nПример: Чипсы 31-01-2022 31-01-2023\n\nРегистр важен! (хлеб и Хлеб это разные продукты)"
		db.SetStatus(2, user_id, chat_id)
		return msg
	} else if user_msg == "/open" {
		msg = "Введите название продукта(одним словом), который вы открыли и новый срок хранения\nПример: Чипсы 24-01-2022 31-01-2023\n\nРегистр важен! (хлеб и Хлеб это разные продукты)"
		db.SetStatus(3, user_id, chat_id)
		return msg
	} else if user_msg == "/finish" {
		msg = "Введите название продукта(одним словом), который вы приготовили/выбросили\n\nРегистр важен! (хлеб и Хлеб это разные продукты)"
		db.SetStatus(4, user_id, chat_id)
		return msg
	} else if user_msg == "/list" {
		msg = "Введите цифру сортировки списка продуктов: \n1.По алфавиту [Список покупок]\n2.По истечению срока годности [Холодильник]\n3.Список всех продуктов"
		db.SetStatus(5, user_id, chat_id)
		return msg
	} else if user_msg == "/listused" {
		msg = "Список ранее использованных продуктов"
		products, _ := db.GetList(user_id, chat_id, 4)
		for i := 0; i < len(products); i++ {
			msg += "\n" + strconv.Itoa(i+1) + ". " + products[i].Name
			if products[i].Weight != 0 {
				msg += ", " + strconv.FormatFloat(products[i].Weight, 'f', 0, 64) + "гр."
			}
			rtime := products[i].Rest_time.String()
			if products[i].Rest_time > 0 {
				rtime = strings.ReplaceAll(rtime, "h", " часов, ")
				rtime = strings.ReplaceAll(rtime, "m", " минут, ")
				rtime = strings.ReplaceAll(rtime, "s", " секунд ")
				msg += ", испортится через: " + rtime
			} else {
				msg += ", срок годности вышел"
			}
		}
		db.SetStatus(0, user_id, chat_id)
		return msg
	} else if user_msg == "/statustime" {
		msg = "Введите за какой промежуток промежуток времени вы хотите просмотреть статистику\nПример: 30-12-2022 31-12-2022"
		db.SetStatus(6, user_id, chat_id)
		return msg
	}
	return msg
}
