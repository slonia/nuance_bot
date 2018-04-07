package main

import (
	"encoding/json"
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

var err error
var bot *tgbotapi.BotAPI

func main() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	fatal(err)
	bot, err = tgbotapi.NewBotAPI(config.Token)
	fatal(err)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	fatal(err)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		rand.Seed(time.Now().UTC().UnixNano())
		text := strings.ToLower(update.Message.Text)
		log.Println("Processing message", text)
		if update.Message.Chat.Type == "private" && *update.Message.Photo != nil && len(*update.Message.Photo) > 0 {
			log.Println("Downloading photo")
			for i, image := range *update.Message.Photo {
				if (i+1)%2 == 0 {
					downloadImage(image.FileID)
				}
			}
		} else if strings.Index(text, "нюанс") > -1 || strings.Index(text, "ньюанс") > -1 {
			log.Println("Nuance found. Responding")
			msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, randomImage())
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func randomImage() string {
	var files []string
	fileInfo, err := ioutil.ReadDir("./images/")
	fatal(err)
	for _, file := range fileInfo {
		name := file.Name()
		if name[0:1] != "." {
			files = append(files, name)
		}
	}
	log.Printf("Selecting random image from %v variants\n", len(files))
	return "./images/" + files[rand.Intn(len(files))]
}

func fatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func downloadImage(fileId string) {
	url, _ := bot.GetFileDirectURL(fileId)
	response, err := http.Get(url)
	fatal(err)
	defer response.Body.Close()
	nameParts := strings.Split(url, "/")
	name := strconv.Itoa(rand.Int()) + nameParts[len(nameParts)-1]
	log.Println("Downloading image to", "./images/"+name)
	file, err := os.Create("./images/" + name)
	fatal(err)
	_, err = io.Copy(file, response.Body)
	fatal(err)
	file.Close()
}
