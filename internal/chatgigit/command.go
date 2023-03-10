package chatgigit

import "github.com/bwmarrin/discordgo"

func GetCommandsArray() (commands []*discordgo.ApplicationCommand) {
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "hey",
			Description: "ギギと一緒にお話ししよう！",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "話しかけてね！",
					Required:    true,
				},
			},
		},
	}

	return
}
