package osiire

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetCommandHandlersMap() (commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"osiire": c_osiire,
	}

	return
}

func c_osiire(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := ""

	userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
	osiireID := ""
	if len(i.ApplicationCommandData().Options) >= 2 {
		// sitei sareteru
		osiireID = i.ApplicationCommandData().Options[1].ChannelValue(nil).ID
	} else {
		// sarete nai
		channels, _ := s.GuildChannels(i.GuildID)
		for _, channel := range channels {
			if channel.Name == "押し入れ" || channel.Name == "おしいれ" {
				osiireID = channel.ID
			}
		}
	}

	// move
	if osiireID != "" {
		err := s.GuildMemberMove(i.GuildID, userID, &osiireID)
		if err != nil {
			response = "> 押し入れ閉まってた :pensive: "
		} else {
			response = fmt.Sprintf("> <@%s> を押し入れに入れました :thumbsup:", userID)
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
