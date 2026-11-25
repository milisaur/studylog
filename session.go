package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Session struct {
	Module    string    `json:"module"`
	Active    bool      `json:"active"`
	StartTime time.Time `json:"start_time"`
}

type CompletedSession struct {
	Module    string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Note      string
}

func start(module string) {
	session := Session{Module: module, Active: true, StartTime: time.Now()}
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	file, err := os.OpenFile(activeSessionPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if os.IsExist(err) {
		fmt.Println("Error: A session is already active.")
		return
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Started session for module %s\n", module)
}

func stop(notes ...string) {
	session, err := loadSession()
	if os.IsNotExist(err) {
		fmt.Println("No active session")
		return
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	endTime := time.Now()
	completed := CompletedSession{
		Module:    session.Module,
		StartTime: session.StartTime,
		EndTime:   endTime,
		Duration:  endTime.Sub(session.StartTime),
		Note:      strings.Join(notes, " "),
	}
	if completed.Duration < 0 {
		fmt.Println("Error: Invalid session duration")
		return
	}
	err = saveSession(completed)
	if err != nil {
		fmt.Println("Error saving session:", err)
		return
	}
	err = os.Remove(activeSessionPath)
	if err != nil {
		fmt.Println("Session saved, but cleanup failed:", err)
		fmt.Println("Do not run stop again before resolving this issue.")
		return
	}
	fmt.Printf("Completed session: %s\n", completed.Module)
	fmt.Printf("Started: %s\n", completed.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Ended: %s\n", completed.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Duration: %s\n", completed.Duration.Truncate(time.Second))
	if completed.Note != "" {
		fmt.Printf("Note: %s\n", completed.Note)
	}
	fmt.Println("Session saved to sessions.csv")
}

func status() {
	session, err := loadSession()
	if os.IsNotExist(err) {
		fmt.Println("No active session")
		return
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	elapsed := time.Since(session.StartTime)
	fmt.Printf("Module: %s\n", session.Module)
	fmt.Printf("Active: %t\n", session.Active)
	fmt.Printf("Started: %s\n", session.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Elapsed: %s\n", elapsed.Truncate(time.Second))
}
