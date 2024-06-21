package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xfyer/discord-bot/bot"
	"github.com/0xfyer/discord-bot/state"
	"github.com/bwmarrin/discordgo"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	stop := make(chan os.Signal, 1)

	discord_secret := os.Getenv("DISCORD_SECRET")

	bot, err := bot.New(discord_secret)
	if err != nil {
		return err
	}

	bot.DefaultHandlers()

	go bot.Open(stop)

	err = bot.DefaultCommands()
	if err != nil {
		return err
	}

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	return nil
}
