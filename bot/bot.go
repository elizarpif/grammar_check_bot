package bot

import (
	"context"
	"fmt"

	"github.com/elizarpif/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/elizarpif/speechrecog/grammar"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	grammar *grammar.Grammar
}

func New(botToken string, g *grammar.Grammar) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	return &Bot{bot: bot, grammar: g}, nil
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

		fixed, err := b.grammar.Check(update.Message.Text)
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
			retMsg = fmt.Sprintf("Fixed message:\n%s", fixed)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, retMsg)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := b.bot.Send(msg); err != nil {
			logger.GetLogger(ctx).Printf("can't send msg\n")
		}
	}
}
