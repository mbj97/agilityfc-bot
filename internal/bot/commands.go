package bot

import (
	"agilityfc-bot/config"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
    commands = []*discordgo.ApplicationCommand{
        {
            Name:        "setname",
            Description: "Set your OSRS name and get added to Runner role",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionString,
                    Name:        "accountname",
                    Description: "Account name",
                    Required:    true,
                },
            },
        },
    }

    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, config *config.Config){
        "setname": func(s *discordgo.Session, i *discordgo.InteractionCreate, config *config.Config) {
            options := i.ApplicationCommandData().Options
            accountName := options[0].StringValue()

            // Change the user's nickname
            err := s.GuildMemberNickname(i.GuildID, i.Member.User.ID, accountName)
            if err != nil {
                log.Printf("Error changing nickname: %v", err)
                s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                    Type: discordgo.InteractionResponseChannelMessageWithSource,
                    Data: &discordgo.InteractionResponseData{
                        Content: "Failed to change nickname. Please check my permissions.",
                        Flags:   64, // Private response
                    },
                })
                return
            }

            // Add the user to the Member role
            err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, config.MemberRoleID)
            if err != nil {
                log.Printf("Error adding role: %v", err)
                s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                    Type: discordgo.InteractionResponseChannelMessageWithSource,
                    Data: &discordgo.InteractionResponseData{
                        Content: "Nickname changed, but failed to add role. Please check my permissions.",
                        Flags:   64, // Private response
                    },
                })
                return
            }

            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: "Thanks for registering your name! Now go to <#1259636963149873244> to gain full access to the Discord",
                    Flags:   64, // Private response
                },
            })

            // Send a direct message to the user
            channel, err := s.UserChannelCreate(i.Member.User.ID)
            if err != nil {
                log.Printf("Error creating DM channel: %v", err)
                return
            }

            _, err = s.ChannelMessageSend(channel.ID, "Your OSRS name has been set to "+accountName+" and you have been added to the Runner role.")
            if err != nil {
                log.Printf("Error sending DM: %v", err)
            }
        },
    }
)
