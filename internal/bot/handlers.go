package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // Ignore messages sent by bot
		return
	}

	// Uncomment for debugging all messages
	//log.Printf("Received message: %s", m.Content)

	// Standard (non slash) commands
	switch m.Content {
	case "!ping":
		log.Println("Received !ping command")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// Handling for Anti-Pk Requests
	if step, ok := ongoingInteractions[m.Author.ID]; ok {
		b.handleAntiPkRequest(s, m, step)
	}

	// Watch for wiseoldman setrsn command
	if m.Interaction != nil && m.Interaction.Name == "setrsn" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> sure you also do `/setname` to get full access to the Discord!", m.Interaction.User.ID))
	}
}

func (b *Bot) handleAntiPkRequest(s *discordgo.Session, m *discordgo.MessageCreate, step string) {
	switch step {
	case "question1":
		userResponses[m.Author.ID] = append(userResponses[m.Author.ID], m.Content)
		saveAntiPkResponseData()
		delete(ongoingInteractions, m.Author.ID)

		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Printf("Error creating DM channel: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(channel.ID, "2. What gear setup will you run while anti-pking?")
		if err != nil {
			log.Printf("Error sending DM: %v", err)
			return
		}

		ongoingInteractions[m.Author.ID] = "question2"
	case "question2":
		userResponses[m.Author.ID] = append(userResponses[m.Author.ID], m.Content)
		saveAntiPkResponseData()
		delete(ongoingInteractions, m.Author.ID)

		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Printf("Error creating DM channel: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(channel.ID, "3. How often will you be available to anti-pk?")
		if err != nil {
			log.Printf("Error sending DM: %v", err)
			return
		}

		ongoingInteractions[m.Author.ID] = "question3"
	case "question3":
		userResponses[m.Author.ID] = append(userResponses[m.Author.ID], m.Content)
		saveAntiPkResponseData()
		delete(ongoingInteractions, m.Author.ID)

		content := fmt.Sprintf("<@%s> has sent anti-pk application.\n", m.Author.ID)
		embed := &discordgo.MessageEmbed{
			Title: "Responses",
			Description: fmt.Sprintf("Combat Levels: %s\nGear Setup: %s\nAvailability: %s",
				userResponses[m.Author.ID][0],
				userResponses[m.Author.ID][1],
				userResponses[m.Author.ID][2]),
			Color: 0x00ff00,
		}

		_, err := s.ChannelMessageSendComplex(b.cfg.AntiPkResponseChannelID, &discordgo.MessageSend{
			Content: content,
			Embed:   embed,
		})
		if err != nil {
			log.Printf("Error sending embed: %v", err)
		}

		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Printf("Error creating DM channel: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(channel.ID, "Your application has been received, a staff member will review and accept/deny shortly.")
		if err != nil {
			log.Printf("Error sending DM: %v", err)
			return
		}

		// Clear the user's responses
		delete(userResponses, m.Author.ID)
		saveAntiPkResponseData()
	}
}
