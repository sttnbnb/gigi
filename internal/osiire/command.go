package osiire

import "github.com/bwmarrin/discordgo"

func GetCommandsArray() (commands []*discordgo.ApplicationCommand) {
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "osiire",
			Description: "押し入れコマンド",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "対象ユーザー",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "行き先",
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildVoice,
					},
					Required: false,
				},
			},
		},
	}

	return
}
