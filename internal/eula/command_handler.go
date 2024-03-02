package eula

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func GetCommandHandlersMap() (commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"generateeula": c_generateeula,
	}

	return
}

func c_generateeula(s *discordgo.Session, i *discordgo.InteractionCreate) {
	originalMessage, err := s.ChannelMessage(i.ChannelID, i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "> メッセージが存在しません。",
				Flags:   1 << 6,
			},
		})
	}

	eulaText := originalMessage.Content
	agreeLabel := i.ApplicationCommandData().Options[1].StringValue()
	disagreeLabel := i.ApplicationCommandData().Options[2].StringValue()

	_, err = s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: eulaText,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: discordgo.ComponentEmoji{
							Name: "😁",
						},
						Label:    agreeLabel,
						Style:    discordgo.PrimaryButton,
						Disabled: false,
						CustomID: "ch_agree",
					},
					discordgo.Button{
						Emoji: discordgo.ComponentEmoji{
							Name: "😑",
						},
						Label:    disagreeLabel,
						Style:    discordgo.DangerButton,
						Disabled: false,
						CustomID: "ch_disagree",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "> 送信しました。",
			Flags:   1 << 6,
		},
	})
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}
