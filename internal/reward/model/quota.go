package model

import "time"

const (
	QuotaPeriodDay   = "day"
	QuotaPeriodWeek  = "week"
	QuotaPeriodMonth = "month"
)

func ValidQuotaPeriod(period string) bool {
	return period == QuotaPeriodDay || period == QuotaPeriodWeek || period == QuotaPeriodMonth
}

func (r *Reward) HasQuota() bool {
	return r.QuotaCount > 0 && ValidQuotaPeriod(r.QuotaPeriod)
}

// QuotaWindow returns the calendar period containing now: today, the week
// starting on Monday, or the month, in the reward's timezone. Both bounds are
// returned in UTC. Without a timezone it uses the server's (TZ).
func (r *Reward) QuotaWindow(now time.Time) (start, end time.Time) {
	loc := time.Local
	if r.QuotaTimezone != "" {
		if l, err := time.LoadLocation(r.QuotaTimezone); err == nil {
			loc = l
		}
	}
	local := now.In(loc)
	day := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	switch r.QuotaPeriod {
	case QuotaPeriodWeek:
		start = day.AddDate(0, 0, -((int(day.Weekday()) + 6) % 7))
		end = start.AddDate(0, 0, 7)
	case QuotaPeriodMonth:
		start = time.Date(local.Year(), local.Month(), 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 1, 0)
	default:
		start = day
		end = day.AddDate(0, 0, 1)
	}
	return start.UTC(), end.UTC()
}
