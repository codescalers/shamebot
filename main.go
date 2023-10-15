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
	for update := range updates {
		if update.Message.IsCommand() {
			go singleUpdate(update, bot)
		}
	}
}

func singleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	command := update.Message.Command()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch command {
	case "add":
		msg.Text = handleAdd(update)
	case "remove":
		msg.Text = handleStop(update)
	case "list":
		msg.Text = handleList(update)
	case "help":
	}
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func handleAdd(update tgbotapi.Update) string {
	repos[update.Message.Text] = true
	return "added"
}

func handleList(update tgbotapi.Update) string {
	list := ""
	for repo, ok := range repos {
		if ok {
			list += repo + "\n"
		}
	}
	return list
}

func handleStop(update tgbotapi.Update) string {
	repos[update.Message.Text] = false
	return "stoped"
}
