package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
			rand.Seed(time.Now().UTC().UnixNano())
			if update.Message.Chat.Type == "private" && *update.Message.Photo != nil && len(*update.Message.Photo) > 0 {
				for _, image := range *update.Message.Photo {
					url, _ := bot.GetFileDirectURL(image.FileID)
					response, e := http.Get(url)
					if e != nil {
						log.Fatal(e)
					}

					defer response.Body.Close()
					nameParts := strings.Split(url, "/")
					name := strconv.Itoa(rand.Int()) + nameParts[len(nameParts)-1]
					//open a file for writing
					file, err := os.Create("./images/" + name)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println()
					// Use io.Copy to just dump the response body to the file. This supports huge files
					_, err = io.Copy(file, response.Body)
					if err != nil {
						log.Fatal(err)
					}
					file.Close()
				}
			} else if strings.Index(strings.ToLower(update.Message.Text), "нюанс") > -1 || strings.Index(strings.ToLower(update.Message.Text), "ньюанс") > -1 {
				msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, randomImage())
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}

func randomImage() string {
	var files []string
	fileInfo, err := ioutil.ReadDir("./images/")
	if err != nil {
		log.Panic(err)
	}

	for _, file := range fileInfo {
		name := file.Name()
		if name[0:1] != "." {
			files = append(files, name)
		}
	}
	fmt.Println(files)
	return "./images/" + files[rand.Intn(len(files))]
}
