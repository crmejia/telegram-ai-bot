package main

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	ollapi "github.com/ollama/ollama/api"
)

func main() {
	output := os.Stdout
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		fmt.Fprintln(output, err)
		return
	}
	bot.Debug = true

	client, err := ollapi.ClientFromEnvironment()
	if err != nil {
		fmt.Fprintln(output, err)
		return
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		req := &ollapi.GenerateRequest{
			Model:  "llama3",
			Stream: new(bool),
			Prompt: update.Message.Text,
		}
		ctx := context.Background()
		msg := tgbotapi.MessageConfig{}
		respFunc := func(resp ollapi.GenerateResponse) error {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, resp.Response)
			msg.ReplyToMessageID = update.Message.MessageID
			return nil
		}
		err = client.Generate(ctx, req, respFunc)

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
