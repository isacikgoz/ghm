package printer

import (
	"time"

	"github.com/fatih/color"
	"github.com/justincampbell/timeago"
)

type LastUpdate struct {
	Date time.Time
}

func (u *LastUpdate) String() string {
	since := time.Since(u.Date)
	if since < 7*24*time.Hour {
		return color.GreenString(timeago.FromDuration(since))
	}

	if since < 30*24*time.Hour {
		return color.YellowString(timeago.FromDuration(since))
	}

	return color.RedString(timeago.FromDuration(since))
}
