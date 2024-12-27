package main

import (
	"fmt"
	"time"
)

func printDailySummary() {
	sessions, err := loadSessions()
	if err != nil {
		fmt.Printf("Error loading sessions: %v\n", err)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	var todaySessions []WorkSession
	taskDurations := make(map[string]time.Duration)

	for _, session := range sessions {
		if session.StartTime.Truncate(24 * time.Hour).Equal(today) {
			todaySessions = append(todaySessions, session)
			taskDurations[session.Task] += session.Duration
		}
	}

	if len(todaySessions) == 0 {
		fmt.Println("No work sessions recorded today")
		return
	}

	fmt.Println("\nToday's Work Summary:")
	fmt.Println("--------------------")

	var totalDuration time.Duration
	for task, duration := range taskDurations {
		fmt.Printf("%-20s: %s\n", task, duration.Round(time.Second))
		totalDuration += duration
	}

	fmt.Println("--------------------")
	fmt.Printf("Total Time: %s\n\n", totalDuration.Round(time.Second))
}
