package main

import (
	"fmt"
	"sort"
	"time"
)

func printDailySummary() {
	sessions, err := loadSessions()
	if err != nil {
		fmt.Printf("Error loading sessions: %v\n", err)
		return
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
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

		if sessionDay.Equal(today) {
			todaySessions = append(todaySessions, session)
			taskDurations[session.Task] += session.Duration
		}
	}

	if len(todaySessions) == 0 {
		fmt.Println("No work sessions recorded today")
		return
	}

	// Get tasks and sort them
	var tasks []string
	for task := range taskDurations {
		tasks = append(tasks, task)
	}
	sort.Strings(tasks)

	fmt.Println("\nToday's Work Summary:")
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
