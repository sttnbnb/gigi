package bosyu

import "github.com/bwmarrin/discordgo"

func GetCommandsArray() (commands []*discordgo.ApplicationCommand) {
	var permissionAdministrator int64 = discordgo.PermissionAdministrator

	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "bosyu",
			Description:              "募集コマンド",
			DefaultMemberPermissions: &permissionAdministrator,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "募集内容",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "対象ロール",
					Required:    true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionInteger,
					Name:         "atto",
					Description:  "募集人数",
					Required:     false,
					Autocomplete: false,
				},
			},
		},
	}

	return
}
