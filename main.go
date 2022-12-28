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

		user_msg := update.Message.Text
		user_name := update.Message.From.FirstName
		user_id := update.Message.From.ID
		chat_id := update.Message.Chat.ID
		status := db.GetStatus(user_id, chat_id)

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
				if (reflect.TypeOf(text[0]) == reflect.TypeOf(user_msg)) && (errFloat == nil) && (errDate == nil) {
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
				p.CreateProduct(user_id, chat_id, text[0], 0, true, false, time.Now(), time.Now(), false)
				err := db.AddInF(&p)

				msg = "Продукт добавлен в список!\nБез таймера."
				if err != nil {
					msg = "Ошибка при добавлении."
				}
			}
			if len(text) == 3 && text[1] == "с" {
				today := time.Now().Format("02-01-2006")
				today += " " + text[2]
				log.Println(today)
				date, dateerr := time.Parse("02-01-2006 15:04", today)
				log.Println(date)
				if dateerr == nil {
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
				if (fromerr == nil) && (toerr == nil) {
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
				if (fromerr == nil) && (toerr == nil) {
					err := db.SetFridge(user_id, chat_id, from, to, text[1])

					msg = text[1] + " перенесён из списка в холодильник!"
					if err != nil {
						msg = "Ошибка при переносе."
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
				if (fromerr == nil) && (toerr == nil) {
					err := db.OpenProduct(user_id, chat_id, from, to, text[0])

					msg = text[0] + " открыт, срок хранения обновлен!"
					if err != nil {
						msg = "Ошибка при открытии продукта."
					}
				}

			}
			db.SetStatus(0, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}
		if user_msg == "/m" {
			msg := "лю крысулечьку катю"
			SendMessageToUser(bot, update, msg)
		}

		if user_msg == "/help" {
			msg := "😳ВАЖНО😳 Перед тем, как впервые воспользоваться ботом, вам нужно написать команду /reg, чтобы бот вас запомнил.\n\nНажмите '/' или на кнопку 'Меню' слева внизу, чтобы увидеть доступные команды\n\nЭтот бот умеет:\n1. Добавлять продукты в список покупок\n2.Переносить из списка/добавлять продукты в холодильник\n3. Показывать список продуктов"
			SendMessageToUser(bot, update, msg)
		}

		if user_msg == "/reg" {
			msg := db.CreateUser(user_name, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}

		if user_msg == "/addinlist" {
			msg := "Введите название продукта(одним словом), вес(в граммах) и дату напоминания через пробел\nПример: Чипсы 80 31-01-2022 09:00\n\nЛибо просто введите название продукта😊"
			db.SetStatus(1, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}

		if user_msg == "/addinfridge" {
			msg := "Если вы хотите перенести продукт из списка продуктов в холодильник, то введите 'купил', название продукта(одним словом) и срок хранения через пробел\nПример: купил Чипсы 31-01-2022 31-01-2023\n\nЕсли же вы хотите добавить новый продукт в холодильник, то введите название продукта и срок хранения через пробел\nПример: Чипсы 31-01-2022 31-01-2023"
			db.SetStatus(2, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}

		if user_msg == "/open" {
			msg := "Введите название продукта(одним словом), который вы открыли и новый срок хранения\nПример: Чипсы 24-01-2022 31-01-2022"
			db.SetStatus(3, user_id, chat_id)
			SendMessageToUser(bot, update, msg)
		}

		if user_msg == "/list" {
			msg := "Список продуктов:"
			products, _ := db.GetList(user_id, chat_id)
			for i := 0; i < len(products); i++ {
				msg += "\n" + strconv.Itoa(i) + ". " + products[i].Name
				if products[i].Weight != 0 {
					msg += ", " + strconv.FormatFloat(products[i].Weight, 'f', 0, 64) + "гр."
				}
				if products[i].TimerEnable {
					msg += ", Таймер включен"
				}
			}
			SendMessageToUser(bot, update, msg)
		}
	}

}
