package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var repos = map[string]bool{}

func ManageCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

func AddRepo(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// validate the repo name
	repos[update.Message.Text] = true
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" is being monitored now")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func StopRepo(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	repos[update.Message.Text] = false
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" is removed from the list")
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
