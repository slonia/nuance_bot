package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type Configuration struct {
	Token string
}

func main() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	} else {
		fmt.Println(config.Token)
		bot, err := tgbotapi.NewBotAPI(config.Token)
		if err != nil {
			log.Panic(err)
		}
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, err := bot.GetUpdatesChan(u)
		for update := range updates {
			if update.Message == nil {
				continue
			}
			if strings.Index(strings.ToLower(update.Message.Text), "нюанс") > -1 || strings.Index(strings.ToLower(update.Message.Text), "ньюанс") > -1 {
				msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "./chapaev.jpeg")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}

}
