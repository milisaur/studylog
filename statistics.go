package main

import (
	"fmt"
	"sort"
	"time"
)

func today() {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)
	printStatistics("Today's study time", start, end)
}

func week() {
	now := time.Now()
	daysSinceMonday := (int(now.Weekday()) + 6) % 7
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := midnight.AddDate(0, 0, -daysSinceMonday)
	end := start.AddDate(0, 0, 7)
	printStatistics("This week's study time", start, end)
}

func printStatistics(title string, start, end time.Time) {
	sessions, err := loadCompletedSessions()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	totals := make(map[string]time.Duration)
	var total time.Duration
	for _, session := range sessions {
		if !session.StartTime.Before(start) && session.StartTime.Before(end) {
			totals[session.Module] += session.Duration
			total += session.Duration
		}
	}
	modules := make([]string, 0, len(totals))
	for module := range totals {
		modules = append(modules, module)
	}
	sort.Strings(modules)
	fmt.Printf("\n%s\n\n", title)
	for _, module := range modules {
		fmt.Printf("%-8s %s\n", module+":", totals[module])
	}
	fmt.Printf("\nTotal:   %s\n", total)
}
