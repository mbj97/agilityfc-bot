package main

import (
	"agilityfc-bot/config"
	"agilityfc-bot/internal/bot"
	"agilityfc-bot/internal/dynamo"
	"agilityfc-bot/internal/server"
	"fmt"
	"os"
)

func main() {
    cfg, err := config.LoadConfig("config.json")
    if err != nil {
        fmt.Println("Error loading config:", err)
        os.Exit(1)
    }

    dynamoService, err := dynamo.NewDynamoDBService()
    if err != nil {
        fmt.Println("Error initializing DynamoDB service:", err)
        os.Exit(1)
    }

    discordBot, err := bot.NewBot(cfg)
    if err != nil {
        fmt.Println("Error creating bot:", err)
        os.Exit(1)
    }

    go func() {
        if err := discordBot.Start(); err != nil {
            fmt.Println("Error starting bot:", err)
            os.Exit(1)
        }
    }()

    srv := server.NewServer(discordBot, dynamoService)
    srv.StartHTTPServer()
}
