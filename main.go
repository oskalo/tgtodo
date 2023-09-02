package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Task represents a task item.
type Task struct {
	ID   int
	Text string
}

var tasks []Task
var lastTaskID int

func main() {
	// Отримуємо токен вашого бота з змінної оточення.
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN не встановлено у змінних оточення")
	}

	// Створюємо нового бота з використанням отриманого токену.
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Встановлюємо відлагодження для бота.
	bot.Debug = true

	// Отримуємо оновлення від користувачів.
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Обробка отриманих повідомлень від користувачів.
	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		chatID := update.Message.Chat.ID

		switch {
		case strings.HasPrefix(text, "/add "):
			taskText := strings.TrimPrefix(text, "/add ")
			addTask(bot, chatID, taskText)
		case strings.HasPrefix(text, "/delete "):
			taskIDStr := strings.TrimPrefix(text, "/delete ")
			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				sendMessage(bot, chatID, "Невірний формат команди /delete")
			} else {
				deleteTask(bot, chatID, taskID)
			}
		case text == "/list":
			listTasks(bot, chatID)
		default:
			sendMessage(bot, chatID, "Невідома команда. Використовуйте /add, /delete або /list.")
		}
	}
}

func addTask(bot *tgbotapi.BotAPI, chatID int64, taskText string) {
	lastTaskID++
	newTask := Task{ID: lastTaskID, Text: taskText}
	tasks = append(tasks, newTask)

	sendMessage(bot, chatID, "Задача додана успішно!")
}

func deleteTask(bot *tgbotapi.BotAPI, chatID int64, taskID int) {
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			sendMessage(bot, chatID, "Задача видалена успішно!")
			return
		}
	}

	sendMessage(bot, chatID, "Задача не знайдена.")
}

func listTasks(bot *tgbotapi.BotAPI, chatID int64) {
	if len(tasks) == 0 {
		sendMessage(bot, chatID, "Список задач порожній.")
		return
	}

	message := "Список задач:\n"
	for _, task := range tasks {
		message += strconv.Itoa(task.ID) + ". " + task.Text + "\n"
	}

	sendMessage(bot, chatID, message)
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
