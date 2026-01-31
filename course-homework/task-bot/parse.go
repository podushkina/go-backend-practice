package main

import (
	"strconv"
	"strings"
)

type Parsed struct {
	Cmd   string
	ID    int
	Title string
}

func parseCommand(text string) (Parsed, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Parsed{}, false
	}
	switch text {
	case "/tasks":
		return Parsed{Cmd: "tasks"}, true
	case "/my":
		return Parsed{Cmd: "my"}, true
	case "/owner":
		return Parsed{Cmd: "owner"}, true
	}

	if strings.HasPrefix(text, "/new ") {
		title := strings.TrimPrefix(text, "/new ")
		title = strings.TrimSpace(title)
		if title == "" {
			return Parsed{}, false
		}
		return Parsed{Cmd: "new", Title: title}, true
	}
	if strings.HasPrefix(text, "/assign_") {
		rest := strings.TrimPrefix(text, "/assign_")
		n, err := strconv.Atoi(rest)
		if err != nil {
			return Parsed{}, false
		}
		return Parsed{Cmd: "assign", ID: n}, true
	}

	if strings.HasPrefix(text, "/unassign_") {
		rest := strings.TrimPrefix(text, "/unassign_")
		n, err := strconv.Atoi(rest)
		if err != nil {
			return Parsed{}, false
		}
		return Parsed{Cmd: "unassign", ID: n}, true
	}

	if strings.HasPrefix(text, "/resolve_") {
		rest := strings.TrimPrefix(text, "/resolve_")
		n, err := strconv.Atoi(rest)
		if err != nil {
			return Parsed{}, false
		}
		return Parsed{Cmd: "resolve", ID: n}, true
	}

	return Parsed{}, false
}
