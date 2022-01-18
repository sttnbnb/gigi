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

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ch_sanka": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// TODO: BSIDの取得と参加処理

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "参加を申し込みました。",
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"ch_torikesi": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// TODO: BSIDの取得と取り消し処理

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "参加を取り消しました。",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"ch_kanri": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// TODO: BSIDの取得

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "操作を選択してください。",
					Flags:   1 << 6,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    "ch_select",
									Placeholder: "Choose Action 👇",
									Options: []discordgo.SelectMenuOption{
										{
											Label: "〆",
											Value: "sime",
											Emoji: discordgo.ComponentEmoji{
												Name: "🚦",
											},
											Description: "募集を締め切ります。",
										},
										{
											Label: "集合",
											Value: "syuugou",
											Emoji: discordgo.ComponentEmoji{
												Name: "🛎",
											},
											Description: "専用ロールで集合をかけます。",
										},
										{
											Label: "無効化",
											Value: "mukou",
											Emoji: discordgo.ComponentEmoji{
												Name: "🗑",
											},
											Description: "募集を無効化します。",
										},
									},
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"ch_select": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var response *discordgo.InteractionResponse

			data := i.MessageComponentData()
			switch data.Values[0] {
			case "sime":

				// TODO: BSIDの取得と締め切り処理

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "締め切りました。",
						Flags:   1 << 6,
					},
				}
			case "syuugou":

				// TODO: BSIDの取得と集合処理

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "集合かけました。",
						Flags:   1 << 6,
					},
				}
			case "mukou":

				// TODO: BSIDの取得と無効化処理

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "無効化しました。",
						Flags:   1 << 6,
					},
				}
			}
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				panic(err)
			}
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bosyu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{

					// TODO: embedの生成と送付
					Content: "ぼしゅするよ",

					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "参加",
									Style:    discordgo.DangerButton,
									Disabled: false,
									CustomID: "ch_sanka",
								},
								discordgo.Button{
									Label:    "取り消し",
									Style:    discordgo.PrimaryButton,
									Disabled: false,
									CustomID: "ch_torikesi",
								},
								discordgo.Button{
									Label:    "管理",
									Style:    discordgo.SecondaryButton,
									Disabled: false,
									CustomID: "ch_kanri",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
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
