package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/google/go-github/v32/github"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{}
	app.Name = "notify-issues-to-slack"
	app.Usage = "The CLI tool to notify Github issues and pull requests to Slack with color."
	app.Version = fmt.Sprintf("%s (rev: %s/%s)", version, revision, runtime.Version())
	app.UsageText = "notify-issues-to-slack -github-token=... -slack-webhook-url=... -query=... [-danger-filter=...] [-warning-filter=...] [-channel=...] [-text=...] [-username=...] [-icon-emoji=...] [-github-api-url=...]"
	cli.AppHelpTemplate = fmt.Sprintf(`%s
SEE ALSO:
   Please see https://github.com/shibayu36/notify-issues-to-slack for detailed usage.
`, cli.AppHelpTemplate)

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
			Name:  "danger-filter",
			Usage: "Colorize the issue's attachment danger. You can use Github search queries",
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
		query := convertRelativeTimeQuery(c.String("query"))

		gc := &githubClient{
			apiURL: c.String("github-api-url"),
			token:  c.String("github-token"),
		}
		issues, err := gc.searchGithubIssues(query)
		if err != nil {
			return err
		}

		warningIssues := []*github.Issue{}
		if wf := c.String("warning-filter"); wf != "" {
			warningIssues, err = gc.searchGithubIssues(query + " " + convertRelativeTimeQuery(wf))
			if err != nil {
				return err
			}
		}

		dangerIssues := []*github.Issue{}
		if df := c.String("danger-filter"); df != "" {
			dangerIssues, err = gc.searchGithubIssues(query + " " + convertRelativeTimeQuery(df))
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
		})
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
