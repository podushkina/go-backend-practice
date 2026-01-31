package main

// сюда писать код

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	// Подключаем библиотеку API
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	// @BotFather в телеграме даст вам токен.
	BotToken = "YOUR_TOKEN_HERE"

	// Урл, в который будет стучаться телега
	WebhookURL = "https://your-domain.com"
)

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}
	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	tm := NewTaskManager()
	updates := bot.ListenForWebhook("/")
	go func() {
		log.Println("start listen :8081")
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			log.Println("http err:", err)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		user := User{ID: update.Message.From.ID, Username: update.Message.From.UserName}
		parsed, ok := parseCommand(update.Message.Text)
		if !ok {
			continue
		}
		switch parsed.Cmd {
		case "new":
			task := tm.Create(user, parsed.Title)
			text := fmt.Sprintf(`Задача "%s" создана, id=%d`, task.Title, task.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
		case "assign":
			task, prev, err := tm.Assign(parsed.ID, user)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задача не найдена"))
				continue
			}
			text := fmt.Sprintf(`Задача "%s" назначена на вас`, task.Title)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

			var notifyUserID int64
			if prev != nil {
				if prev.ID != user.ID {
					notifyUserID = prev.ID
				}
			} else {
				if task.Owner.ID != user.ID {
					notifyUserID = task.Owner.ID
				}
			}
			if notifyUserID != 0 {
				notifyText := fmt.Sprintf(`Задача "%s" назначена на @%s`, task.Title, user.Username)
				bot.Send(tgbotapi.NewMessage(notifyUserID, notifyText))
			}
		case "unassign":
			task, _, err := tm.Unassign(parsed.ID, user)
			if err != nil {
				if errors.Is(err, NotYourTask) {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задача не на вас"))
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задача не найдена"))
				}
				continue
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Принято"))
			if task.Owner.ID != user.ID {
				notifyText := fmt.Sprintf(`Задача "%s" осталась без исполнителя`, task.Title)
				bot.Send(tgbotapi.NewMessage(task.Owner.ID, notifyText))
			}
		case "resolve":
			task, _, err := tm.Resolve(parsed.ID, user)
			if err != nil {
				if errors.Is(err, NotYourTask) {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задача не на вас"))
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задача не найдена"))
				}
				continue
			}
			text := fmt.Sprintf(`Задача "%s" выполнена`, task.Title)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
			if task.Owner.ID != user.ID {
				notifyText := fmt.Sprintf(`Задача "%s" выполнена @%s`, task.Title, user.Username)
				bot.Send(tgbotapi.NewMessage(task.Owner.ID, notifyText))
			}
		case "tasks":
			tasks := tm.ListAll()
			text := formatTasks(tasks, user, "tasks")
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

		case "my":
			tasks := tm.ListByAssignee(user.ID)
			text := formatTasks(tasks, user, "my")
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

		case "owner":
			tasks := tm.ListByOwner(user.ID)
			text := formatTasks(tasks, user, "owner")
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
		}
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
