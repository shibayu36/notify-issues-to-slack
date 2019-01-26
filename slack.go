package main

import (
	"fmt"
	"time"

	slack "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/google/go-github/github"
)

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
	payload := slack.Payload{
		Channel:     opt.Channel,
		Text:        opt.Text,
		Username:    opt.Username,
		IconEmoji:   opt.IconEmoji,
		Attachments: attachments,
		LinkNames:   "1",
	}
	err := slack.Send(s.webhookURL, "", payload)
	if len(err) > 0 {
		fmt.Println(err)
	}

	return nil
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
