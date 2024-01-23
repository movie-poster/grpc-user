package utils

import "time"

func HasPassedOneMonth(planStartDate time.Time) bool {
	currentDate := time.Now()
	duration := currentDate.Sub(planStartDate)

	// Defines the duration of one month. 30 days
	oneMonthDuration := 30 * 24 * time.Hour
	return duration >= oneMonthDuration
}
