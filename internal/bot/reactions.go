package bot

import (
	"encoding/json"
	"log"
	"os"

	"agilityfc-bot/config"

	"github.com/bwmarrin/discordgo"
)

var (
	runnerRequestData              = make(map[string]RequestData) // Map to store message ID and request data pairs
	ongoingInteractions            = make(map[string]string)      // Map to store ongoing interactions
	userResponses                  = make(map[string][]string)    // Map to store user responses
)

type RequestData struct {
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
}

func init() {
	file, err := os.Open(config.RUNNER_REQUEST_DATA_FILE)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&runnerRequestData)
		if err != nil {
			log.Printf("Error decoding JSON data: %v", err)
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Error opening JSON file: %v", err)
	}

	// Load anti-pk response data from JSON file
	file, err = os.Open(config.ANTI_PK_RESPONSE_DATA_FILE)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&userResponses)
		if err != nil {
			log.Printf("Error decoding JSON data: %v", err)
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Error opening JSON file: %v", err)
	}
}

func (b *Bot) MessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	b.handleRunnerApprovalReaction(s, r)
	b.handleRunnerRequestReaction(s, r)
	b.handleAntiPkRequestReaction(s, r)
}

func (b *Bot) handleRunnerApprovalReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if request, ok := runnerRequestData[r.MessageID]; ok {
		if r.Emoji.Name == b.cfg.CheckMarkEmoji {
			err := s.GuildMemberRoleAdd(r.GuildID, request.UserID, b.cfg.RunnerRoleID)
			if err != nil {
				log.Printf("Error adding role: %v", err)
			} else {
				log.Printf("Added 'Runner' role to user %s", request.UserID)
				s.ChannelMessageSend(b.cfg.RunnerRequestChannelID, "<@"+request.UserID+"> has been granted the 'Runner' role by <@"+r.UserID+">")
			}
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
		} else if r.Emoji.Name == b.cfg.RedXEmoji {
			// Remove the request from the data and delete the message
			delete(runnerRequestData, r.MessageID)
			saveRunnerRequestData()
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
		}
	}
}

func (b *Bot) handleRunnerRequestReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageID == b.cfg.RunnerRequestSpecificMessageID {
		log.Printf("Reaction added: %s by %s on %s", r.Emoji.Name, r.UserID, r.MessageID)

		messageContent := "<@" + r.UserID + "> has requested to become a runner. " +
			b.cfg.CheckMarkEmoji + " for approval, " + b.cfg.RedXEmoji + " to reject."

		msg, err := s.ChannelMessageSend(b.cfg.RunnerRequestChannelID, messageContent)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		} else {
			// Add reactions for approval and rejection
			s.MessageReactionAdd(b.cfg.RunnerRequestChannelID, msg.ID, b.cfg.CheckMarkEmoji)
			s.MessageReactionAdd(b.cfg.RunnerRequestChannelID, msg.ID, b.cfg.RedXEmoji)

			// Save request data
			runnerRequestData[msg.ID] = RequestData{
				UserID:    r.UserID,
				RequestID: r.MessageID,
			}
			saveRunnerRequestData()
		}
	}
}

func (b *Bot) handleAntiPkRequestReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageID == b.cfg.AntiPkRequestMessageID {
		if r.UserID == s.State.User.ID {
			return
		}

		log.Printf("Reaction added: %s by %s on %s", r.Emoji.Name, r.UserID, r.MessageID)

		channel, err := s.UserChannelCreate(r.UserID)
		if err != nil {
			log.Printf("Error creating DM channel: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(channel.ID, "Hi! Please answer the following questions:")
		if err != nil {
			log.Printf("Error sending DM: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(channel.ID, "1. List your Combat Levels (Attack, Strength, Defense, Ranged, Magic, Prayer)")
		if err != nil {
			log.Printf("Error sending DM: %v", err)
			return
		}

		ongoingInteractions[r.UserID] = "question1"
	}
}

func saveRunnerRequestData() {
	file, err := os.Create(config.RUNNER_REQUEST_DATA_FILE)
	if err != nil {
		log.Printf("Error creating JSON file: %v", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(runnerRequestData)
	if err != nil {
		log.Printf("Error encoding JSON data: %v", err)
	}
}

func saveAntiPkResponseData() {
	file, err := os.Create(config.ANTI_PK_RESPONSE_DATA_FILE)
	if err != nil {
		log.Printf("Error creating JSON file: %v", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(userResponses)
	if err != nil {
		log.Printf("Error encoding JSON data: %v", err)
	}
}