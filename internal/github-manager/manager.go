package github

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/go-github/github"
)

type MorntoringRequest struct {
	ProjectId int64
	ChatId    int64
}

func ManageProjects(requestsChan chan MorntoringRequest, bot *tgbotapi.BotAPI) {
	client := github.NewClient(nil)

	for request := range requestsChan {
		shamedIssues := ""

		project, _, err := client.Projects.GetProject(context.Background(), request.ProjectId)
		if err != nil {
			msg := tgbotapi.NewMessage(request.ChatId, fmt.Sprintf("failed to get project info %d", request.ProjectId))
			if _, err := bot.Send(msg); err != nil {
				log.Println("faild to respond")
			}
			continue
		}

		val, err := shameProject(client, *project)
		if err != nil {
			log.Println(err)
		}

		shamedIssues += val
		msg := tgbotapi.NewMessage(request.ChatId, fmt.Sprintf("Shamed Issues for the project %s are:\n%s", project.GetName(), shamedIssues))
		if _, err := bot.Send(msg); err != nil {
			log.Println("faild to respond")
		}

	}
}

func shameProject(client *github.Client, project github.Project) (shamedIssues string, err error) {
	columns, _, err := client.Projects.ListProjectColumns(context.Background(), *project.ID, nil)
	if err != nil {
		return
	}

	for _, col := range columns {
		switch col.GetName() {
		case "In Progress", "Blocked":
			val, err := getShamedIssues(client, col.GetID())
			if err != nil {
				return "", err
			}

			shamedIssues += val
		}
	}
	return
}

func getShamedIssues(client *github.Client, colId int64) (shamedIssues string, err error) {
	issues, _, err := client.Projects.ListProjectCards(context.Background(), colId, nil)
	if err != nil {
		return
	}

	for _, issue := range issues {
		if issue.GetUpdatedAt().Time.Compare(time.Now().AddDate(0, 0, -1)) < 1 {
			shamedIssues += "issue: " + issue.GetNote() + "\n"
		}
	}
	return
}
