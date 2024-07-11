package main

import (
    "log"
    "agilityfc-bot/config"
    "agilityfc-bot/internal/bot"
)

func main() {
    cfg, err := config.LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    discordBot, err := bot.NewBot(cfg)
    if err != nil {
        log.Fatalf("Error creating bot: %v", err)
    }

    err = discordBot.Start()
    if err != nil {
        log.Fatalf("Error starting bot: %v", err)
    }

    discordBot.Wait()
}
