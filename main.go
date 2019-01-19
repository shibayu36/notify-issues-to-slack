package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	slack "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{}
	app.Name = "notify-issues-to-slack"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "github-token",
			Usage: "Github token",
		},
		&cli.StringFlag{
			Name:  "slack-webhook-url",
			Usage: "Slack webhook URL",
		},
	}
	app.Action = func(c *cli.Context) error {
		issues, err := searchGithubIssues(c.String("github-token"))
		if err != nil {
			return err
		}
		postIssuesToSlack(c.String("slack-webhook-url"), issues)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func searchGithubIssues(token string) ([]github.Issue, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	i, _, err := client.Search.Issues(ctx, "repo:playframework/playframework label:\"help wanted\" label:\"topic:documentation\" state:open", &github.SearchOptions{Sort: "created", Order: "asc"})
	// fmt.Println(res)
	if err != nil {
		return nil, err
	}
	// fmt.Println(i)
	return i.Issues, nil
}

func postIssuesToSlack(webhookURL string, issues []github.Issue) {
	attachments := []slack.Attachment{}
	for _, i := range issues {
		user := i.GetAssignee()
		if user == nil {
			user = i.GetUser()
		}

		title := fmt.Sprintf("@%s %s", user.GetLogin(), i.GetTitle())

		duration := time.Now().Sub(i.GetCreatedAt())
		var color string
		if duration.Hours() > 24*365 {
			color = "danger"
		} else if duration.Hours() > 24*100 {
			color = "warning"
		} else {
			color = "good"
		}
		a := slack.Attachment{
			Title:     &title,
			TitleLink: i.HTMLURL,
			Color:     &color,
		}
		attachments = append(attachments, a)
	}
	payload := slack.Payload{
		Text:        "レビュー依頼君",
		Username:    "reviewiraikun",
		Channel:     "shibayu36-private",
		Attachments: attachments,
	}
	err := slack.Send(webhookURL, "", payload)
	if len(err) > 0 {
		fmt.Println(err)
	}
}
