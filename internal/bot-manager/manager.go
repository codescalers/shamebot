package botmanager

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v56/github"
)

func ManageProjects(projectsChan chan int64) {
	client := github.NewClient(nil)

	for projectID := range projectsChan {
		shamedIssues := ""

		project, _, err := client.Projects.GetProject(context.Background(), projectID)
		if err != nil {
			log.Println(err)
		}

		val, err := manageProject(client, *project)
		if err != nil {
			log.Println(err)
		}

		shamedIssues += val
		log.Println(shamedIssues)

	}
}

func manageProject(client *github.Client, project github.Project) (shamedIssues string, err error) {
	columns, _, err := client.Projects.ListProjectColumns(context.Background(), *project.ID, nil)
	if err != nil {
		return
	}

	for _, col := range columns {
		switch col.GetName() {
		case "In Progress", "Blocked":
			val, err := manageCol(client, col.GetID())
			if err != nil {
				return "", err
			}

			shamedIssues += val
		}
	}
	return
}

func manageCol(client *github.Client, colId int64) (shamedIssues string, err error) {
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
