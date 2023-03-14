package bosyu

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func GetComponentHandlersMap() (componentsHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ch_sousin":   ch_sousin,
		"ch_sanka":    ch_sanka,
		"ch_torikesi": ch_torikesi,
		"ch_kanri":    ch_kanri,
		"ch_select":   ch_select,
	}

	return
}

func ch_sousin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := sousin(s, i)

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_sanka(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := sanka(s, i)

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_torikesi(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := torikesi(s, i)

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_kanri(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := kanri(s, i)

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_select(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response *discordgo.InteractionResponse
	data := i.MessageComponentData()
	messageID := i.Message.Content[46:65]

	switch data.Values[0] {
	case "syuugou":
		response = syuugou(s, i.GuildID, i.Message.ChannelID, messageID)
	case "sime":
		response = sime(s, i.GuildID, i.Message.ChannelID, messageID)
	case "mukou":
		response = mukou(s, i.GuildID, i.Message.ChannelID, messageID)
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}
