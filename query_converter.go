package main

import (
	"regexp"
	"strings"
	"time"

	tparse "github.com/karrick/tparse/v2"
)

var (
	// e.g.) created:<=now-7d, then $1 is `created:<=` and $2 is `now-7d`.
	// See Also: https://help.github.com/en/articles/understanding-the-search-syntax#query-for-dates
	relativeTimeQueryReg = regexp.MustCompile(`^((?:created|updated|closed|merged):[<>=]*)(.+)$`)
)

func convertRelativeTimeQuery(query string) string {
	conditions := strings.Split(query, " ")
	convertedConditions := []string{}
	for _, c := range conditions {
		matches := relativeTimeQueryReg.FindStringSubmatch(c)
		if len(matches) != 3 {
			convertedConditions = append(convertedConditions, c)
		} else {
			absolute, err := tparse.ParseNow(time.RFC3339, matches[2])
			if err != nil {
				convertedConditions = append(convertedConditions, c)
			} else {
				convertedConditions = append(convertedConditions, matches[1]+absolute.UTC().Format(time.RFC3339))
			}
		}
	}
	return strings.Join(convertedConditions, " ")
}
