package main

import (
	"fmt"
	"strings"
	"text/template"

	slack "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/google/go-github/github"
)

// slackClient has the role of formatting the issues and posting them to slack.
type slackClient struct {
	webhookURL string
}

type slackPostOptions struct {
	Text            string
	IssueTextFormat string
	Channel         string
	Username        string
	IconEmoji       string
}

const (
	defaultIssueTextFormat = "{{.GetTitle}} @{{if .GetAssignee }}{{.GetAssignee.GetLogin}}{{else}}{{.GetUser.GetLogin}}{{end}}"
)

func (s *slackClient) postIssuesToSlack(issues []github.Issue, warningIssues []github.Issue, dangerIssues []github.Issue, opt *slackPostOptions) error {
	issueTextFormat := opt.IssueTextFormat
	if issueTextFormat == "" {
		issueTextFormat = defaultIssueTextFormat
	}

	warningIssueMap := s.makeIssueExistsMap(warningIssues)
	dangerIssueMap := s.makeIssueExistsMap(dangerIssues)

	attachments := []slack.Attachment{}
	for _, i := range issues {
		issueText, err := s.formatIssueText(issueTextFormat, i)
		if err != nil {
			fmt.Println(err)
			return err
		}
		color := s.getColorByIssue(i, warningIssueMap, dangerIssueMap)
		a := slack.Attachment{
			Title:     &issueText,
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

func (s *slackClient) formatIssueText(format string, issue github.Issue) (string, error) {
	t, err := template.New("issue-text").Parse(format)
	if err != nil {
		return "", err
	}

	var text strings.Builder
	err = t.Execute(&text, &issue)
	if err != nil {
		return "", err
	}

	return text.String(), nil
}

func (s *slackClient) getColorByIssue(issue github.Issue, warningIssueMap map[int64]bool, dangerIssueMap map[int64]bool) string {
	if dangerIssueMap[issue.GetID()] {
		return "danger"
	} else if warningIssueMap[issue.GetID()] {
		return "warning"
	} else {
		return "good"
	}
}

func (s *slackClient) makeIssueExistsMap(issues []github.Issue) map[int64]bool {
	issueExistsMap := map[int64]bool{}
	for _, i := range issues {
		issueExistsMap[i.GetID()] = true
	}
	return issueExistsMap
}
