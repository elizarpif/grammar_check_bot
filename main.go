package main

import (
	"context"
)

func main() {
	bot, err := New()
	if err != nil {
		panic(err)
	}

	bot.Start(context.Background())
}
