package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOTAPITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		file, err := getPictureFile(bot, &update)
		if err != nil {
			log.Printf("Failed to get file: %v", err)
			continue
		}

		log.Printf("FILEPATH IS: %s", file.FilePath)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "YOLO")
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func getPictureFile(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (*tgbotapi.File, error) {
	if update.Message == nil {
		return nil, errors.New("No new message")
	}
	receivedMessage := update.Message

	fileID := ""
	if receivedMessage.Photo != nil {
		// The last element is the original picture
		fileID = (*receivedMessage.Photo)[len((*receivedMessage.Photo))-1].FileID
	} else if receivedMessage.Document != nil {
		log.Printf("Got document file url: %v", receivedMessage.Document.FileID)
		fileID = receivedMessage.Document.FileID
	} else {
		return nil, errors.New("No picture in the message")
	}

	file, err := bot.GetFile(tgbotapi.FileConfig{fileID})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get file %v got %v", fileID, err))
	}

	return &file, nil
}
