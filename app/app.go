package app

import (
	"log"
	"strconv"
	"time"

	shamebot "github.com/codescalers/shamebot/internal/bot"
	github "github.com/codescalers/shamebot/internal/github-manager"
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

	addChan := make(chan tgbotapi.Update)
	deleteChan := make(chan tgbotapi.Update)
	go monitorProjects(addChan, deleteChan, bot)

	command := ""
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message.IsCommand() {

			command = update.Message.Command()
			shamebot.ManageCommand(update, bot)
		} else {
			if command == "add" {
				addChan <- update
			} else if command == "stop" {
				deleteChan <- update
			}
			command = ""

		}
	}
}

func monitorProjects(addChan chan tgbotapi.Update, deleteChan chan tgbotapi.Update, bot *tgbotapi.BotAPI) {
	requestsChan := make(chan github.MorntoringRequest)
	projects := map[int64]int64{}
	ticker := time.NewTicker(time.Second * 2)

	go github.ManageProjects(requestsChan, bot)

	for {
		select {
		case update := <-addChan:
			if projectId, err := addProject(update.Message.Text, update, bot); err == nil {
				projects[int64(projectId)] = update.Message.Chat.ID
			}

		case update := <-deleteChan:
			if projectId, err := deleteProject(update.Message.Text, update, bot); err == nil {
				projects[int64(projectId)] = 0
			}

		case <-ticker.C:
			for projectId, chatId := range projects {
				if chatId != 0 {
					requestsChan <- github.MorntoringRequest{ProjectId: projectId, ChatId: chatId}
				}
			}
		}
	}
}

func addProject(project string, update tgbotapi.Update, bot *tgbotapi.BotAPI) (int, error) {
	// parse project
	projectId, err := strconv.Atoi(project)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "invalid project")
		if _, err = bot.Send(msg); err != nil {
			log.Println(err)
		}
		return 0, err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" is being monitored now")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
		return 0, err
	}

	return projectId, nil
}

func deleteProject(project string, update tgbotapi.Update, bot *tgbotapi.BotAPI) (int, error) {
	// parse project
	projectId, err := strconv.Atoi(project)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "invalid project")
		if _, err := bot.Send(msg); err != nil {
			return 0, err
		}

		return 0, err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" is removed from the list")
	if _, err := bot.Send(msg); err != nil {
		return 0, err
	}

	return projectId, nil
}
