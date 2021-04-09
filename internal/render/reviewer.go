package render

import (
	"strings"

	"github.com/fatih/color"
	"github.com/isacikgoz/ghm/internal/fetcher"
)

type Reviewers struct {
	Reviewers []*fetcher.Reviewer
}

func (r *Reviewers) String() string {
	var s []string
	for _, reviewer := range r.Reviewers {
		var formatted string
		switch reviewer.Status {
		case "COMMENTED":
			formatted = color.YellowString(reviewer.Handle)
		case "CHANGES_REQUESTED":
			formatted = color.RedString(reviewer.Handle)
		case "APPROVED":
			formatted = color.GreenString(reviewer.Handle)
		default:
			formatted = reviewer.Handle
		}
		s = append(s, formatted)
	}

	return strings.Join(s, ", ")
}
