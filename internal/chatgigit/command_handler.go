package chatgigit

import (
	"github.com/bwmarrin/discordgo"
)

func GetCommandHandlersMap() (commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hey": c_hey,
	}

	return
}

func c_hey(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "見にくいからコマンドは廃止になったよ！これからは「@ギィギ」を先頭につけて話しかけてね！",
			Flags: 1 << 6,
		},
	})
}
