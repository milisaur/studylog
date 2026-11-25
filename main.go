package main

import (
	"fmt"
	"os"
)

func main() {
	name := "Studylog"
	version := "0.1.0"

	fmt.Printf("%s v%s\n", name, version)

	if len(os.Args) < 2 {
		fmt.Println("Usage: studylog <command>")
		return
	}

	if err := initPaths(); err != nil {
		fmt.Fprintln(os.Stderr, "Configuration error:", err)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		if len(os.Args) < 3 {
			fmt.Println("Usage: studylog start <module>")
			return
		}
		start(os.Args[2])
	case "stop":
		stop(os.Args[2:]...)
	case "status":
		status()
	case "today":
		today()
	case "week":
		week()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
