package main

import (
	"log"
	"os"
	"time"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{}
	app.Name = "notify-issues-to-slack"
	app.UsageText = "notify-issues-to-slack -github-token=... -slack-webhook-url=... -query=... [-danger-over=...] [-warning-over=...] [-slack-channel=...] [-slack-text=...] [-slack-username=...] [-slack-icon-emoji=...] [-github-api-url=...]"
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
			Name:  "github-api-url",
			Usage: "Github API base URL",
		},
	}
	app.Action = func(c *cli.Context) error {
		gc := &githubClient{
			apiURL: c.String("github-api-url"),
			token:  c.String("github-token"),
		}
		issues, err := gc.searchGithubIssues(c.String("query"))
		if err != nil {
			return err
		}

		var dangerOver, warningOver time.Duration
		if c.String("danger-over") != "" {
			dangerOver, err = time.ParseDuration(c.String("danger-over"))
			if err != nil {
				return err
			}
		}
		if c.String("warning-over") != "" {
			warningOver, err = time.ParseDuration(c.String("warning-over"))
			if err != nil {
				return err
			}
		}

		sc := &slackClient{webhookURL: c.String("slack-webhook-url")}
		sc.postIssuesToSlack(issues, &slackPostOptions{
			Text:        c.String("slack-text"),
			Channel:     c.String("slack-channel"),
			Username:    c.String("slack-username"),
			IconEmoji:   c.String("slack-icon-emoji"),
			DangerOver:  &dangerOver,
			WarningOver: &warningOver,
		})
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
