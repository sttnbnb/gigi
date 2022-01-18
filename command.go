package main

import (
	"fmt"

	"github.com/shmn7iii/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{

		{
			Name:        "bosyu",
			Description: "募集コマンド",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "募集内容",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "atto",
					Description: "募集人数",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "対象ロール",
					Required:    false,
				},
			},
		},

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
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildVoice,
					},
					Required: false,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"bosyu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				i.ApplicationCommandData().Options[0].StringValue(),
			}
			msgformat :=
				` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> string_option: %s
`
			if len(i.ApplicationCommandData().Options) >= 2 {
				margs = append(margs, i.ApplicationCommandData().Options[1].IntValue())
				msgformat += "> integer_option: %d\n"
			}
			if len(i.ApplicationCommandData().Options) >= 3 {
				margs = append(margs, i.ApplicationCommandData().Options[2].RoleValue(nil, "").ID)
				msgformat += "> role-option: <@&%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},

		"osiire": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				i.ApplicationCommandData().Options[0].UserValue(nil).ID,
			}
			msgformat :=
				` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> user-option: <@%s>
`
			if len(i.ApplicationCommandData().Options) >= 2 {
				margs = append(margs, i.ApplicationCommandData().Options[1].ChannelValue(nil).ID)
				msgformat += "> channel-option: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
	}
)
