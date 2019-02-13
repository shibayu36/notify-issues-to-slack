package main

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	slack "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/google/go-github/github"
)

// slackClient has the role of formatting the issues and posting them to slack.
type slackClient struct {
	webhookURL string
}

type slackPostOptions struct {
	Text        string
	Channel     string
	Username    string
	IconEmoji   string
	DangerOver  *time.Duration
	WarningOver *time.Duration
}

func (s *slackClient) postIssuesToSlack(issues []github.Issue, opt *slackPostOptions) error {
	attachments := []slack.Attachment{}
	for _, i := range issues {
		user := i.GetAssignee()
		if user == nil {
			user = i.GetUser()
		}

		title := fmt.Sprintf(
			"%s @%s",
			i.GetTitle(),
			user.GetLogin(),
		)
		color := s.getColorByIssue(i, opt.DangerOver, opt.WarningOver)
		a := slack.Attachment{
			Title:     &title,
			TitleLink: i.HTMLURL,
			Color:     &color,
		}
		attachments = append(attachments, a)
	}

	text, err := s.formatText(opt.Text, issues)
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := slack.Payload{
		Channel:     opt.Channel,
		Text:        text,
		Username:    opt.Username,
		IconEmoji:   opt.IconEmoji,
		Attachments: attachments,
		LinkNames:   "1",
	}
	errs := slack.Send(s.webhookURL, "", payload)
	if len(errs) > 0 {
		fmt.Println(errs)
	}

	return nil
}

func (s *slackClient) formatText(format string, issues []github.Issue) (string, error) {
	t, err := template.New("slack-text").Parse(format)
	if err != nil {
		return "", err
	}

	var text strings.Builder
	err = t.Execute(&text, issues)
	if err != nil {
		return "", err
	}

	return text.String(), nil
}

func (s *slackClient) getColorByIssue(issue github.Issue, dangerOver *time.Duration, warningOver *time.Duration) string {
	durationFromIssueCreated := time.Now().Sub(issue.GetCreatedAt())

	color := "good"
	if dangerOver != nil && durationFromIssueCreated.Hours() > dangerOver.Hours() {
		color = "danger"
	} else if warningOver != nil && durationFromIssueCreated.Hours() > warningOver.Hours() {
		color = "warning"
	}

	return color
}
