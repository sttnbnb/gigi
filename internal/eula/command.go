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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "agree_label",
					Description: "同意ボタンラベル",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "disagree_label",
					Description: "拒否ボタンラベル",
					Required:    true,
				},
			},
		},
	}

	return
}
