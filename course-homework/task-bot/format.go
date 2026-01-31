package main

import (
	"fmt"
	"strings"
)

func formatTasks(tasks []*Task, viewer User, cmd string) string {
	if len(tasks) == 0 {
		return "Нет задач"
	}

	var sb strings.Builder

	for i, t := range tasks {
		// Добавляем перенос строки между задачами (но не перед первой)
		if i > 0 {
			sb.WriteString("\n\n")
		}

		titleLine := fmt.Sprintf("%d. %s by @%s", t.ID, t.Title, t.Owner.Username)
		sb.WriteString(titleLine)

		if t.Assignee != nil {
			if t.Assignee.ID == viewer.ID {
				// Задача на МНЕ
				// В команде /my строку "assignee: я" писать не надо, в остальных - надо
				if cmd != "my" {
					sb.WriteString("\nassignee: я")
				}
				// Добавляем кнопки управления
				sb.WriteString(fmt.Sprintf("\n/unassign_%d /resolve_%d", t.ID, t.ID))
			} else {
				sb.WriteString(fmt.Sprintf("\nassignee: @%s", t.Assignee.Username))
			}
		} else {
			// Задача НИЧЬЯ -> кнопка "Взять"
			sb.WriteString(fmt.Sprintf("\n/assign_%d", t.ID))
		}
	}

	return sb.String()
}
