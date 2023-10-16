package app

import (
	"log"

	shamebot "github.com/codescalers/shamebot/internal/bot"
	manager "github.com/codescalers/shamebot/internal/bot-manager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const token = ""

func StartApp() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	projects := make(chan int64)
	manager.ManageProjects(projects)

	command := ""
	for update := range updates {
		if update.Message.IsCommand() {
			command = update.Message.Command()
			shamebot.ManageCommand(update, bot)
		} else {
			if command == "add" {
				shamebot.AddRepo(update, bot)
			} else if command == "stop" {
				shamebot.StopRepo(update, bot)
			}
			command = ""
		}
	}
}
