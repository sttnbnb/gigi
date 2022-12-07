package eula

import "github.com/bwmarrin/discordgo"

func GetCommandsArray() (commands []*discordgo.ApplicationCommand) {
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "generateeula",
			Description: "EULA生成コマンド",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "original_message_id",
					Description: "EULA生成元メッセージID",
					Required:    true,
				},
			},
		},
	}

	return
}
