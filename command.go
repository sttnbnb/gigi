package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shmn7iii/discordgo"
)

func syuugou(s *discordgo.Session, channelID string, messageID string) {
	message, _ := s.ChannelMessage(channelID, messageID)
	desc := message.Embeds[0].Description
	roleName := ""
	if len(desc) <= 9 {
		roleName = desc + "..."
	} else {
		roleName = desc[:9] + "..."
	}
	roles, _ := s.GuildRoles(*GuildID)
	for _, role := range roles {
		if role.Name == roleName {
			s.ChannelMessageSend(channelID, role.Mention()+" 集合！！！！！")
		}
	}
}

func sime(s *discordgo.Session, channelID string, messageID string) {
	message, _ := s.ChannelMessage(channelID, messageID)
	embed := message.Embeds[0]
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "参加者 | 〆!!",
			Value:  embed.Fields[0].Value,
			Inline: false,
		},
	}
	embed.Color = 0xff2600

	str := ""
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      messageID,
		Channel: channelID,
		Content: &str,
		Embeds:  []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "管理",
						Style:    discordgo.SecondaryButton,
						Disabled: false,
						CustomID: "ch_kanri",
					},
				},
			},
		},
	})
}

func mukou(s *discordgo.Session, channelID string, messageID string) {
	// sime
	sime(s, channelID, messageID)

	// delete role
	message, _ := s.ChannelMessage(channelID, messageID)
	desc := message.Embeds[0].Description
	roleName := ""
	if len(desc) <= 9 {
		roleName = desc + "..."
	} else {
		roleName = desc[:9] + "..."
	}
	roles, _ := s.GuildRoles(*GuildID)
	for _, role := range roles {
		if role.Name == roleName {
			s.GuildRoleDelete(*GuildID, role.ID)
		}
	}
}

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
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "対象ロール",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "atto",
					Description: "募集人数",
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
		"ch_sousin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// send message
			s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
				Content: i.Message.Content,
				Embed:   i.Message.Embeds[0],
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
			})
			// response
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "> 送信しました。\n> 定員に達する、または「管理/〆」ボタンか募集を締め切ることができます。\n> また「管理/集合」から参加者専用ロールでメンションを飛ばせます。\n> 企画が終了したら「管理/無効化」からロールの削除と募集の無効化をしてください。",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"ch_sanka": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			simed := false
			embed := i.Message.Embeds[0]

			if strings.HasPrefix(strings.Split(embed.Fields[0].Name, " ")[2], "@") {
				// @ ari
				members := embed.Fields[0].Value
				if members == "ｲﾅｲﾖ" {
					members = ""
				}

				add := "- " + i.Member.User.Username + " #" + i.Member.User.Discriminator
				atto, _ := strconv.Atoi(strings.Replace(strings.Split(embed.Fields[0].Name, " ")[2], "@", "", -1))

				if !strings.Contains(members, add) {
					members += "\n" + add
					atto -= 1
				}

				embed.Fields = []*discordgo.MessageEmbedField{
					{
						Name:   fmt.Sprintf("参加者 | @%d", atto),
						Value:  members,
						Inline: false,
					},
				}

				// sime
				if atto == 0 {
					simed = true
				}

			} else {
				// @ nasi
				members := embed.Fields[0].Value
				if members == "ｲﾅｲﾖ" {
					members = ""
				}

				add := "- " + i.Member.User.Username + " #" + i.Member.User.Discriminator
				atto, _ := strconv.Atoi(strings.Replace(strings.Split(embed.Fields[0].Name, " ")[2], "人", "", -1))

				if !strings.Contains(members, add) {
					members += "\n" + add
					atto += 1
				}

				embed.Fields = []*discordgo.MessageEmbedField{
					{
						Name:   fmt.Sprintf("参加者 | %d人", atto),
						Value:  members,
						Inline: false,
					},
				}
			}

			roleName := ""
			if len(embed.Description) <= 9 {
				roleName = embed.Description + "..."
			} else {
				roleName = embed.Description[:9] + "..."
			}
			roles, _ := s.GuildRoles(*GuildID)
			for _, role := range roles {
				if role.Name == roleName {
					s.GuildMemberRoleAdd(*GuildID, i.Member.User.ID, role.ID)
				}
			}

			s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed)

			if simed {
				sime(s, i.Message.ChannelID, i.Message.ID)
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "> 参加を申し込みました。",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"ch_torikesi": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			embed := i.Message.Embeds[0]

			if strings.HasPrefix(strings.Split(i.Message.Embeds[0].Fields[0].Name, " ")[2], "@") {
				// @ ari
				members := i.Message.Embeds[0].Fields[0].Value

				add := "- " + i.Member.User.Username + " #" + i.Member.User.Discriminator
				atto, _ := strconv.Atoi(strings.Replace(strings.Split(i.Message.Embeds[0].Fields[0].Name, " ")[2], "@", "", -1))

				if strings.Contains(members, add) {
					members = strings.Replace(members, add, "", -1)
					atto += 1
				}

				if members == "" {
					members = "ｲﾅｲﾖ"
				}

				embed.Fields = []*discordgo.MessageEmbedField{
					{
						Name:   fmt.Sprintf("参加者 | @%d", atto),
						Value:  members,
						Inline: false,
					},
				}
			} else {
				// @ nasi
				members := i.Message.Embeds[0].Fields[0].Value

				add := "- " + i.Member.User.Username + " #" + i.Member.User.Discriminator
				atto, _ := strconv.Atoi(strings.Replace(strings.Split(i.Message.Embeds[0].Fields[0].Name, " ")[2], "人", "", -1))

				if strings.Contains(members, add) {
					members = strings.Replace(members, add, "", -1)
					atto -= 1
				}

				if members == "" {
					members = "ｲﾅｲﾖ"
				}

				embed.Fields = []*discordgo.MessageEmbedField{
					{
						Name:   fmt.Sprintf("参加者 | %d人", atto),
						Value:  members,
						Inline: false,
					},
				}
			}

			roleName := ""
			if len(embed.Description) <= 9 {
				roleName = embed.Description + "..."
			} else {
				roleName = embed.Description[:9] + "..."
			}
			roles, _ := s.GuildRoles(*GuildID)
			for _, role := range roles {
				if role.Name == roleName {
					s.GuildMemberRoleRemove(*GuildID, i.Member.User.ID, role.ID)
				}
			}

			s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "> 参加を取り消しました。",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				panic(err)
			}
		},

		"ch_kanri": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "> 操作を選択してください。\n> `ID: " + i.Message.ID + "`",
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
			messageID := i.Message.Content[46:64]

			switch data.Values[0] {
			case "sime":
				sime(s, i.Message.ChannelID, messageID)

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "> 締め切りました。",
						Flags:   1 << 6,
					},
				}
			case "syuugou":
				syuugou(s, i.Message.ChannelID, messageID)

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "> 集合かけました。",
						Flags:   1 << 6,
					},
				}
			case "mukou":
				mukou(s, i.Message.ChannelID, messageID)

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "> 無効化しました。",
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

			Role, _ := s.GuildRoleCreate(*GuildID)
			if len(i.ApplicationCommandData().Options[0].StringValue()) <= 9 {
				s.GuildRoleEdit(*GuildID, Role.ID, i.ApplicationCommandData().Options[0].StringValue()+"...", 0, false, 0, true)
			} else {
				s.GuildRoleEdit(*GuildID, Role.ID, i.ApplicationCommandData().Options[0].StringValue()[:9]+"...", 0, false, 0, true)
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
				panic(err)
			}
		},

		"osiire": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
			if len(i.ApplicationCommandData().Options) >= 2 {
				*OsiireChannelID = i.ApplicationCommandData().Options[1].ChannelValue(nil).ID
			}

			// move
			if *OsiireChannelID != "" {
				s.GuildMemberMove(*GuildID, userID, OsiireChannelID)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("> <@%s> を押し入れに入れました :thumbsup:", userID),
					Flags:   1 << 6,
				},
			})
		},
	}
)
