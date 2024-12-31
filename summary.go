package main

import (
	"fmt"
	"sort"
	"time"
)

func parseDate(dateStr string) (time.Time, error) {
	// Try different date formats
	formats := []string{
		"2006-01-02", // YYYY-MM-DD
		"2006/01/02", // YYYY/MM/DD
		"01-02",      // MM-DD (assume current year)
		"01/02",      // MM/DD (assume current year)
		"Jan 2",      // MMM DD
		"January 2",  // MMMM DD
	}

	now := time.Now()

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			// If year is not specified (year is 0000), use current year
			if t.Year() == 0 {
				return time.Date(now.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location()), nil
			}
			// Convert to local timezone
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location()), nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format. Use YYYY-MM-DD, MM-DD, or MMM DD")
}

func printDailySummary(targetDate time.Time) {
	sessions, err := loadSessions()
	if err != nil {
		fmt.Printf("Error loading sessions: %v\n", err)
		return
	}

	var todaySessions []WorkSession
	taskDurations := make(map[string]time.Duration)

	for _, session := range sessions {
		sessionDay := time.Date(
			session.StartTime.Year(),
			session.StartTime.Month(),
			session.StartTime.Day(),
			0, 0, 0, 0,
			session.StartTime.Location(),
		)

		if sessionDay.Equal(targetDate) {
			todaySessions = append(todaySessions, session)
			taskDurations[session.Task] += session.Duration
		}
	}

	if len(todaySessions) == 0 {
		fmt.Printf("No work sessions recorded for %s\n", targetDate.Format("Monday, January 2, 2006"))
		return
	}

	// Get tasks and sort them
	var tasks []string
	for task := range taskDurations {
		tasks = append(tasks, task)
	}
	sort.Strings(tasks)

	fmt.Printf("\nWork Summary for %s:\n", targetDate.Format("Monday, January 2, 2006"))
	fmt.Println("--------------------")

	var totalDuration time.Duration
	for _, task := range tasks {
		duration := taskDurations[task]
		fmt.Printf("%-20s: %s\n", task, duration.Round(time.Second))
		totalDuration += duration
	}

	fmt.Println("--------------------")
	fmt.Printf("Total Time: %s\n\n", totalDuration.Round(time.Second))
}
