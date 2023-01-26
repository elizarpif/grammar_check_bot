package main

import (
	"context"
)

func main() {
	bot, err := New()
	if err != nil {
		panic(err)
	}

	// start
	bot.Start(context.Background())
}
