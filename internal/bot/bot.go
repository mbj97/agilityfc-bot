package bot

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/bwmarrin/discordgo"
    "agilityfc-bot/config"
)

type Bot struct {
    session *discordgo.Session
    cfg  *config.Config
}

func NewBot(cfg *config.Config) (*Bot, error) {
    dg, err := discordgo.New("Bot " + cfg.Token)
    if err != nil {
        return nil, err
    }

    // Define intents to receive message events and modify user nickname and role
    dg.Identify.Intents = discordgo.IntentsAll

    bot := &Bot{
        session: dg,
        cfg:  cfg,
    }

    dg.AddHandler(bot.messageCreate)
    dg.AddHandler(bot.interactionCreate)
	dg.AddHandler(bot.MessageReactionAdd)

    return bot, nil
}

func (b *Bot) Start() error {
    log.Println("Starting bot...")
    err := b.session.Open()
    if err != nil {
        return err
    }

    // Register the slash commands after opening the session
    for _, v := range commands {
        _, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", v)
        if err != nil {
            log.Fatalf("Cannot create '%v' command: %v", v.Name, err)
        }
    }

    return nil
}

func (b *Bot) Wait() {
    fmt.Println("Bot is now running. Press CTRL+C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    b.session.Close()
}

func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
        handler(s, i, b.cfg)
    }
}

func (b *Bot) SendMessage(channelID, message string) error {
    _, err := b.session.ChannelMessageSend(channelID, message)
    return err
}