package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"

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
		text := anylisePicture("https://api.telegram.org/file/bot" + os.Getenv("TGBOTAPITOKEN") + "/" + file.FilePath)

		if len(text) > 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text[0])
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Successfully sent response")
			}
		}
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

func anylisePicture(address string) []string {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	resp, err := http.Get(address)
	if err != nil {
		log.Fatalf("Failed to get picture: %v", err)
	}

	image, err := vision.NewImageFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	texts, err := client.DetectTexts(ctx, image, nil, 1)
	if err != nil {
		log.Fatalf("Failed to detect texts: %v", err)
	}

	result := []string{}
	for _, text := range texts {
		result = append(result, text.Description)
	}

	return result
}
