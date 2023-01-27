package main

import (
	"context"

	"github.com/elizarpif/speechrecog/bot"
	"github.com/elizarpif/speechrecog/grammar"
)

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	grammarChecker := grammar.NewGrammar(cfg.Grammar.GrammarUrl, cfg.Grammar.GrammarApiKey, cfg.Grammar.Host)

	tgBot, err := bot.New(cfg.BotToken, grammarChecker)
	if err != nil {
		panic(err)
	}

	// start
	tgBot.Start(context.Background())
}
