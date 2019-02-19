package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/google/go-github/github"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{}
	app.Name = "notify-issues-to-slack"
	app.Version = fmt.Sprintf("%s (rev: %s/%s)", version, revision, runtime.Version())
	app.UsageText = "notify-issues-to-slack -github-token=... -slack-webhook-url=... -query=... [-danger-over=...] [-warning-over=...] [-channel=...] [-text=...] [-username=...] [-icon-emoji=...] [-github-api-url=...]"
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
			Name:  "danger-filter",
			Usage: "Colorize the issue's attachment danger. You can use Github search queries",
		},
		&cli.StringFlag{
			Name:  "warning-over",
			Usage: "Colorize the issue's attachment warning",
		},
		&cli.StringFlag{
			Name:  "warning-filter",
			Usage: "Colorize the issue's attachment warning. You can use Github search queries",
		},
		&cli.StringFlag{
			Name:  "channel",
			Usage: "Slack channel to be posted",
		},
		&cli.StringFlag{
			Name:  "text",
			Usage: "text to post with issues",
		},
		&cli.StringFlag{
			Name:  "issue-text-format",
			Usage: "Text format of each issues used for message attachment",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "Slack username to post",
		},
		&cli.StringFlag{
			Name:  "icon-emoji",
			Usage: "Slack icon emoji to post",
		},
		&cli.StringFlag{
			Name:  "github-api-url",
			Usage: "Github API base URL",
		},
	}
	app.Action = func(c *cli.Context) error {
		query := c.String("query")

		gc := &githubClient{
			apiURL: c.String("github-api-url"),
			token:  c.String("github-token"),
		}
		issues, err := gc.searchGithubIssues(query)
		if err != nil {
			return err
		}

		warningIssues := []github.Issue{}
		if wf := c.String("warning-filter"); wf != "" {
			warningIssues, err = gc.searchGithubIssues(query + " " + wf)
			if err != nil {
				return err
			}
		}

		dangerIssues := []github.Issue{}
		if df := c.String("danger-filter"); df != "" {
			dangerIssues, err = gc.searchGithubIssues(query + " " + df)
			if err != nil {
				return err
			}
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
		sc.postIssuesToSlack(issues, warningIssues, dangerIssues, &slackPostOptions{
			Text:            c.String("text"),
			IssueTextFormat: c.String("issue-text-format"),
			Channel:         c.String("channel"),
			Username:        c.String("username"),
			IconEmoji:       c.String("icon-emoji"),
			DangerOver:      &dangerOver,
			WarningOver:     &warningOver,
		})
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
