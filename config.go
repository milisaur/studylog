package main

import (
	"os"
	"path/filepath"
)

var activeSessionPath string
var completedSessionsPath string

func dataDir() (string, error) {
	dir := os.Getenv("STUDYLOG_DATA_DIR")

	if dir == "" {
		home, err := os.UserHomeDir()

		if err != nil {
			return "", err
		}

		dir = filepath.Join(home, ".local", "share", "studylog")
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return dir, nil
}

func initPaths() error {
	dir, err := dataDir()

	if err != nil {
		return err
	}

	activeSessionPath = filepath.Join(dir, "active.json")
	completedSessionsPath = filepath.Join(dir, "sessions.csv")

	return nil
}
