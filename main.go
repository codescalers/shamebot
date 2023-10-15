package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const token = "6358239344:AAH52_RwnCjrSFfGr7CPvil4pDpgcmwZlbc"

//	var commands = map[string]string{
//		"/add":    "adds new github repo",
//		"/remove": "removes a github repo",
//		"/list":   "list all monitored repos",
//	}
var repos = map[string]bool{}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	command := ""
	for update := range updates {
		if update.Message.IsCommand() {
			command = update.Message.Command()
			singleUpdate(update, bot)
		} else {
			if command == "add" {
				addRepo(update, bot)
			} else if command == "stop" {
				stopRepo(update, bot)
			}
			command = ""
		}
	}
}

func singleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	command := update.Message.Command()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch command {
	case "add":
		msg.Text = "Enter the a new repo"
	case "stop":
		msg.Text = "Enter the repo name"
	case "list":
		msg.Text = getList()
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func addRepo(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	repos[update.Message.Text] = true
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" is being monitored now")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func getList() string {
	list := ""
	for repo, ok := range repos {
		if ok {
			list += repo + "\n"
		}
	}
	return list
}

func stopRepo(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	repos[update.Message.Text] = false
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" is removed from the list")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
