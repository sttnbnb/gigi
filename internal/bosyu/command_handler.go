package bosyu

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func GetCommandHandlersMap() (commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bosyu": c_bosyu,
	}

	return
}

func c_bosyu(s *discordgo.Session, i *discordgo.InteractionCreate) {

	embed := discordgo.MessageEmbed{
		Title:       ":mega: 募集",
		Description: i.ApplicationCommandData().Options[0].StringValue(),
		Color:       0x00f900,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.User.Username,
			IconURL: i.Member.User.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "BSID: " + i.ID,
		},
	}

	if len(i.ApplicationCommandData().Options) >= 3 {
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("参加者 | @%d", i.ApplicationCommandData().Options[2].IntValue()),
				Value:  "ｲﾅｲﾖ",
				Inline: false,
			},
		}
	} else {
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "参加者 | 0人",
				Value:  "ｲﾅｲﾖ",
				Inline: false,
			},
		}
	}

	var (
		name        string
		color       int   = 0
		hoist       bool  = false
		permission  int64 = 0
		mentionable bool  = true
	)
	embedDescription := []rune(embed.Description)
	if len(embedDescription) <= 9 {
		name = string(embedDescription) + "..."
	} else {
		name = string(embedDescription)[:9] + "..."
	}
	_, err := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
		Name:        name,
		Color:       &color,
		Hoist:       &hoist,
		Permissions: &permission,
		Mentionable: &mentionable,
	})
	if err != nil {
		log.Printf("Can't create role: %v", name)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				&embed,
			},
			Content: fmt.Sprintf(
				"<@&%s>",
				i.ApplicationCommandData().Options[1].RoleValue(nil, "").ID,
			),
			Flags: 1 << 6,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "送信",
							Style:    discordgo.SuccessButton,
							Disabled: false,
							CustomID: "ch_sousin",
						},
						discordgo.Button{
							Label:    "「送信」でチャンネルへ送信",
							Style:    discordgo.SecondaryButton,
							Disabled: true,
							CustomID: "ch_dummy",
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}
