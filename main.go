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
	app.UsageText = "notify-issues-to-slack -github-token=... -slack-webhook-url=... -query=... [-danger-over=...] [-warning-over=...] [-slack-channel=...] [-slack-text=...] [-slack-username=...] [-slack-icon-emoji=...]"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "github-token",
			Usage: "Github token",
		},
		&cli.StringFlag{
			Name:  "slack-webhook-url",
			Usage: "Slack webhook URL",
		},
		&cli.StringFlag{
			Name:  "query",
			Usage: "Query to search Github issues",
		},
		&cli.StringFlag{
			Name:  "danger-over",
			Usage: "Colorize the issue's attachment danger",
		},
		&cli.StringFlag{
			Name:  "warning-over",
			Usage: "Colorize the issue's attachment warning",
		},
		&cli.StringFlag{
			Name:  "slack-channel",
			Usage: "Slack channel to be posted",
		},
		&cli.StringFlag{
			Name:  "slack-text",
			Usage: "text to post with issues",
		},
		&cli.StringFlag{
			Name:  "slack-username",
			Usage: "Slack username to post",
		},
		&cli.StringFlag{
			Name:  "slack-icon-emoji",
			Usage: "Slack icon emoji to post",
		},
		&cli.StringFlag{
			Name:  "github-api-base",
			Usage: "Github API base URL",
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
