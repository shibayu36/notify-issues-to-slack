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
	Text      string
	Channel   string
	Username  string
	IconEmoji string
}

func (s *slackClient) postIssuesToSlack(issues []github.Issue, opt *slackPostOptions) error {
	attachments := []slack.Attachment{}
	for _, i := range issues {
		user := i.GetAssignee()
		if user == nil {
			user = i.GetUser()
		}

		title := fmt.Sprintf("@%s %s", user.GetLogin(), i.GetTitle())

		// TODO: Set color by -danger-over and -warning-over flag
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
