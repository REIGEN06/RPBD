package main

import (
	"buy-list/product"
	"buy-list/storage/postgresql"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func SendMessageToUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, msg string) {
	bot_msg := tgbotapi.NewMessage(update.Message.From.ID, msg)
	bot_msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(bot_msg); err != nil {
		panic(err)
	}
}

// как лучше сделать таймер? как вообще делать его
// запустить программу в goroutine как отдельный процесс

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	db := postgresql.Connect(os.Getenv("POSTGRESQL_TOKEN"))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 90
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		user_nickname := update.SentFrom().UserName
		user_msg := update.Message.Text
		user_name := update.Message.From.FirstName
		user_id := update.Message.From.ID
		chat_id := update.Message.Chat.ID
		status := db.GetStatus(user_id, chat_id)

		log.Printf("nickname: %s \n", user_nickname)
		log.Printf("user_name: %s \n", user_name)
		//log.Printf("user_id: %d \n", user_id)
		//log.Printf("chat_id: %d \n", chat_id)
		//log.Printf("status: %d \n", status)
		log.Printf("msg: %s \n", user_msg)
		log.Println()

		if status == 1 {
			msg := "Продукт указан не верно"
			text := strings.Fields(user_msg)
			if len(text) == 4 {
				weight, errFloat := strconv.ParseFloat(text[1], 32)
				date, errDate := time.Parse("02-01-2006 15:04", text[2]+" "+text[3])
				if date.Before(time.Now()) {
					msg = "Некорректное время"
				}
				if (reflect.TypeOf(text[0]) == reflect.TypeOf(user_msg)) && (errFloat == nil) && (errDate == nil) && date.After(time.Now()) {
					p := product.Product{}
					p.CreateProduct(user_id, chat_id, text[0], weight, true, false, time.Now(), date, true)
					err := db.AddIn(&p)

					msg = "Продукт добавлен в список!\nТаймер заведён."
					if err != nil {
						msg = "Ошибка при добавлении."
					}
				}
			}
			if len(text) == 1 {
				p := product.Product{}
				p.CreateProduct(user_id, chat_id, text[0], 0, true, false, time.Now(), time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), false)
				err := db.AddIn(&p)

				msg = "Продукт добавлен в список!\nБез таймера."
				if err != nil {
					msg = "Ошибка при добавлении."
				}
			}
			if len(text) == 3 && text[1] == "с" {
				today := time.Now().Format("02-01-2006")
				today += " " + text[2]
				date, dateerr := time.Parse("02-01-2006 15:04", today)
				if date.Before(time.Now()) {
					msg = "Выбранное время уже прошло!"
				}
				if dateerr == nil && date.After(time.Now()) {
					p := product.Product{}
					p.CreateProduct(user_id, chat_id, text[0], 0, true, false, time.Now(), date, true)
					err := db.AddIn(&p)
					msg = "Продукт добавлен в список!\nТаймер заведён."
					if err != nil {
						msg = "Ошибка при добавлении в базу данных."
					}
				} else {
					msg = "С датой что-то не так.."
				}
			}

			db.SetStatus(0, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}
		if status == 2 {
			msg := "Продукт указан не верно"
			text := strings.Fields(user_msg)
			if len(text) == 3 {
				from, fromerr := time.Parse("02-01-2006", text[1])
				to, toerr := time.Parse("02-01-2006", text[2])
				if to.Before(from) && to.Before(time.Now()) {
					msg = "Некорректное время"
				}
				if (fromerr == nil) && (toerr == nil) && to.After(from) && to.After(time.Now()) {
					p := product.Product{}
					p.CreateProduct(user_id, chat_id, text[0], 0, false, true, from, to, true)
					err := db.AddIn(&p)

					msg = text[0] + " добавлен в холодильник!"
					if err != nil {
						msg = "Ошибка при добавлении."
					}
				}
			}
			//купил - перенос из листа
			if len(text) == 4 {
				from, fromerr := time.Parse("02-01-2006", text[2])
				to, toerr := time.Parse("02-01-2006", text[3])
				if to.Before(from) && to.Before(time.Now()) {
					msg = "Некорректное время"
				}
				if (fromerr == nil) && (toerr == nil) && to.After(from) && to.After(time.Now()) {
					err := db.SetFridge(user_id, chat_id, from, to, text[1])

					msg = text[1] + " перенесён из списка в холодильник!"
					if err != nil {
						msg = "Ошибка при переносе."
						if err.Error() == "NOT_EXISTS" {
							msg = text[1] + " не существует"
						}
					}
				}
			}
			db.SetStatus(0, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}
		if status == 3 {
			msg := "Продукт указан не верно"
			text := strings.Fields(user_msg)
			if len(text) == 3 {
				from, fromerr := time.Parse("02-01-2006", text[1])
				to, toerr := time.Parse("02-01-2006", text[2])
				if to.Before(from) && to.Before(time.Now()) {
					msg = "Некорректное время"
				}
				if (fromerr == nil) && (toerr == nil) && to.After(from) && to.After(time.Now()) {
					err := db.SetUsed(user_id, chat_id, from, to, text[0])

					msg = text[0] + " открыт, срок хранения обновлен!"
					if err != nil {
						msg = "Ошибка при открытии продукта."
						if err.Error() == "NOT_EXISTS" {
							msg = text[0] + " не существует"
						}
					}
				}

			}
			db.SetStatus(0, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}
		if status == 4 {
			msg := "Продукт указан не верно"
			text := strings.Fields(user_msg)
			if len(text) == 1 {
				msg = text[0] + " успешно приготовлен/выброшен!"
				err := db.SetTrash(user_id, chat_id, text[0])
				if err != nil {
					msg = "Ошибка при приготовлении/выбрасывании."
					if err.Error() == "NOT_EXISTS" {
						msg = text[0] + " не существует"
					}
				}
			}
			db.SetStatus(0, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}
		if status == 5 {
			msg := "Неправильно выбран список"
			text := strings.Fields(user_msg)
			if len(text) == 1 {
				if text[0] == "1" || text[0] == "2" {
					param, _ := strconv.ParseInt(text[0], 10, 64)
					msg = "Список продуктов"
					if text[0] == "1" {
						msg += " по алфавиту"
					} else if text[0] == "2" {
						msg += " в холодильнике по сроку годности"
					}
					products, _ := db.GetList(user_id, chat_id, param)
					for i := 0; i < len(products); i++ {
						msg += "\n" + strconv.Itoa(i+1) + ". " + products[i].Name
						if products[i].Weight != 0 {
							msg += ", " + strconv.FormatFloat(products[i].Weight, 'f', 0, 64) + "гр."
						}
						if products[i].TimerEnable && products[i].Rest_time > 0 {
							msg += ", таймер включен"
						}
						if text[0] == "2" && products[i].TimerEnable {
							if products[i].Rest_time > 0 {
								rtime := products[i].Rest_time.String()
								rtime = strings.ReplaceAll(rtime, "h", " часов, ")
								rtime = strings.ReplaceAll(rtime, "m", " минут, ")
								rtime = strings.ReplaceAll(rtime, "s", " секунд ")
								msg += ", испортится через: " + rtime
							} else {
								msg += ", срок годности вышел"
							}

						}
					}

				}
			}
			db.SetStatus(0, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}

		if user_msg == "/help" || user_msg == "/start" {
			msg := "😳ВАЖНО😳 Перед тем, как впервые воспользоваться ботом, вам нужно написать команду /reg, чтобы бот вас запомнил.\n\nНажмите '/' или на кнопку 'Меню' слева внизу, чтобы увидеть доступные команды"
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/reg" {
			msg := db.CreateUser(user_nickname, user_name, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/addinlist" {
			msg := "Введите название продукта(одним словом), вес(в граммах) и дату напоминания через пробел\nПример: Чипсы 80 31-01-2023 09:00\n\nЛибо просто введите название продукта😊"
			db.SetStatus(1, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/addinfridge" {
			msg := "Если вы хотите перенести продукт из списка продуктов в холодильник, то введите 'купил', название продукта(одним словом) и срок хранения через пробел\nПример: купил Чипсы 31-01-2022 31-01-2023\n\nЕсли же вы хотите добавить новый продукт в холодильник, то введите название продукта и срок хранения через пробел\nПример: Чипсы 31-01-2022 31-01-2023"
			db.SetStatus(2, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/open" {
			msg := "Введите название продукта(одним словом), который вы открыли и новый срок хранения\nПример: Чипсы 24-01-2022 31-01-2022"
			db.SetStatus(3, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/finish" {
			msg := "Введите название продукта(одним словом), который вы приготовили/выбросили"
			db.SetStatus(4, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/list" {
			msg := "Введите цифру сортировки списка продуктов: \n1.По алфавиту [Список покупок]\n2.По истечению срока годности [Холодильник]"
			db.SetStatus(5, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/listused" {
			msg := "Список ранее использованных продуктов"
			products, _ := db.GetList(user_id, chat_id, 3)
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
			SendMessageToUser(bot, update, msg)
		} else if status <= 0 {
			msg := "Команда не распознана"
			SendMessageToUser(bot, update, msg)
		}
	}

}
