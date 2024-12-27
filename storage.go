package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type WorkSession struct {
	Task      string        `json:"task"`
	Duration  time.Duration `json:"duration"`
	StartTime time.Time     `json:"start_time"`
}

func getStorageFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "work_sessions.json"
	}
	return filepath.Join(homeDir, ".work_sessions.json")
}

func saveSession(task string, duration time.Duration, startTime time.Time) error {
	filename := getStorageFile()

	sessions, err := loadSessions()
	if err != nil {
		sessions = []WorkSession{}
	}

	session := WorkSession{
		Task:      task,
		Duration:  duration,
		StartTime: startTime,
	}

	sessions = append(sessions, session)

	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling sessions: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

func loadSessions() ([]WorkSession, error) {
	filename := getStorageFile()

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []WorkSession{}, nil
		}
		return nil, fmt.Errorf("error reading sessions file: %w", err)
	}

	var sessions []WorkSession
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("error unmarshaling sessions: %w", err)
	}

	return sessions, nil
}
