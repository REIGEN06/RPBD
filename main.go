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

var status = 0

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	db := postgresql.Connect(os.Getenv("POSTGRESQL_TOKEN"))
	//bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 90
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		user_msg := update.Message.Text
		user_id := update.Message.From.ID

		log.Printf("user_name: %s \n", update.SentFrom().UserName)
		log.Printf("msg: %s \n", user_msg)
		log.Println()

		if status == 1 {
			msg := "Продукт указан не верно"
			temp := strings.Fields(user_msg)
			if len(temp) >= 4 {
				weight, errFloat := strconv.ParseFloat(temp[1], 32)
				date, errDate := time.Parse("2006-01-02 15:04", temp[2]+" "+temp[3])
				if (reflect.TypeOf(temp[0]) == reflect.TypeOf(user_msg)) && (errFloat == nil) && (errDate == nil) {
					p := product.Product{}
					p.CreateProduct(user_id, temp[0], weight, date.Format(time.RFC3339))
					err := db.AddInList(&p)
					status = 0
					log.Print(p)
					msg = "Продукт добавлен в список!"
					if err != nil {
						msg = "Ошибка при добавлении в базу данных"
					}
				}
			}
			SendMessageToUser(bot, update, msg)
		}

		if status == 2 {
			//msg := "новый - добавить новый продукт в холодильник\nстарый - перенести продукт в холодильник из списка"

		}

		if user_msg == "/help" {
			msg := "Этот бот умеет добавлять продукты в список или в ваш холодильник, устанавливать таймер на напоминание о покупке\nНажмите '/', чтобы увидеть доступные команды"
			SendMessageToUser(bot, update, msg)
		} else if user_msg == "/addinlist" {
			if status == 0 {
				msg := "Введите название продукта, вес и дату напоминания через пробел\nПример: Чипсы 0.1 2022-01-31 09:00"
				SendMessageToUser(bot, update, msg)
				status = 1
			}
		} else if user_msg == "/addinfridge" {
			if status == 0 {
				msg := "Введите название продукта, вес и дату напоминания через пробел\nПример: Чипсы 0.1 2022-01-31 09:00"
				SendMessageToUser(bot, update, msg)
				status = 2
			}
		} else if user_msg == "/open" {
			if status == 0 {
				msg := "Введите название продукта, который вы открыли"
				SendMessageToUser(bot, update, msg)
				status = 3
			}
		}
	}
}
