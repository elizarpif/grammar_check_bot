package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/elizarpif/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (string, error) {
	s, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(s)), err
}

func New() (*Bot, error) {
	token, err := tokenFromFile("bot_token")
	if err != nil {
		return nil, err
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	return &Bot{bot: bot}, nil
}

type Response struct {
	Fixed    string `json:"fixed"`
	Original string `json:"original"`
}

func checkGrammar(text string) (string, error) {
	url := "https://grammar-genius.p.rapidapi.com/dev/grammar/"

	payload := strings.NewReader(fmt.Sprintf("{\n    \"text\": \"%s\"\n}", text))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("API_KEY"))
	req.Header.Add("X-RapidAPI-Host", "grammar-genius.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.Fixed, nil
}

func (b *Bot) Start(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := b.bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		if update.Message == nil {
			continue
		}

		fixed, err := checkGrammar(update.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Can't heck grammar, err: ", err))
			if _, err := b.bot.Send(msg); err != nil {
				logger.GetLogger(ctx).Printf("can't send msg\n")
			}

			continue
		}

		var retMsg string

		if fixed == update.Message.Text {
			retMsg = fmt.Sprintf("Nice! It looks good")
		} else {
			retMsg = fmt.Sprintf("Fixed message: %s", fixed)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, retMsg)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := b.bot.Send(msg); err != nil {
			logger.GetLogger(ctx).Printf("can't send msg\n")
		}
	}
}
