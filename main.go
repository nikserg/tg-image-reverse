package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	jsonconfig "github.com/nikserg/go-json-config"
)

type Config struct {
	BotToken             string
}

func main() {
	fmt.Println("Start")
	var config Config
	jsonconfig.ReadConfig("config.json", &config)

	fmt.Printf("TG token %v \n", config.BotToken)
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			fmt.Printf("[%s] %v", update.Message.From.UserName, update.Message)
			for _, photo := range update.Message.Photo {

				picUrl, _ := bot.GetFileDirectURL(photo.FileID)
				searchUrl := "https://yandex.ru/images/search?rpt=imageview&url=" + url.QueryEscape(picUrl)

				fmt.Println(searchUrl)

				response, err := http.Get(searchUrl)
				if err != nil {
					panic(err)
				}
				defer response.Body.Close()

				bodyBytes, err := ioutil.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}
				bodyString := string(bodyBytes)
				var resultMessage string
				if !strings.Contains(bodyString, "Таких же изображений не найдено") {
					resultMessage = "❌ fake"
				} else {

					resultMessage = "✅ not fake"
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, resultMessage)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				break
			}
		}
	}
}
