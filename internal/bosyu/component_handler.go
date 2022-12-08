package bosyu

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GetComponentHandlersMap() (componentsHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ch_sousin":   ch_sousin,
		"ch_sanka":    ch_sanka,
		"ch_torikesi": ch_torikesi,
		"ch_kanri":    ch_kanri,
		"ch_select":   ch_select,
	}

	return
}

func syuugou(s *discordgo.Session, guildID string, channelID string, messageID string) {
	message, _ := s.ChannelMessage(channelID, messageID)
	embedDescription := []rune(message.Embeds[0].Description)
	roleName := ""
	if len(embedDescription) <= 9 {
		roleName = string(embedDescription) + "..."
	} else {
		roleName = string(embedDescription[:9]) + "..."
	}
	roles, _ := s.GuildRoles(guildID)
	for _, role := range roles {
		if role.Name == roleName {
			s.ChannelMessageSend(channelID, role.Mention()+" 集合！！！！！")
		}
	}
}

func sime(s *discordgo.Session, guildID string, channelID string, messageID string) {
	message, err := s.ChannelMessage(channelID, messageID)
	if err != nil {
		log.Println(messageID)
		log.Printf("Critical error occurred: %v", err)
	}

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

func mukou(s *discordgo.Session, guildID string, channelID string, messageID string) {
	// sime
	sime(s, guildID, channelID, messageID)

	// delete role
	message, _ := s.ChannelMessage(channelID, messageID)
	embedDescription := []rune(message.Embeds[0].Description)
	roleName := ""
	if len(embedDescription) <= 9 {
		roleName = string(embedDescription) + "..."
	} else {
		roleName = string(embedDescription[:9]) + "..."
	}
	roles, _ := s.GuildRoles(guildID)
	for _, role := range roles {
		if role.Name == roleName {
			s.GuildRoleDelete(guildID, role.ID)
		}
	}
}

func ch_sousin(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			Content: "> 送信しました。\n> 定員へ到達、または「管理/〆」ボタンから募集を締め切ることができます。\n> また「管理/集合」から参加者専用ロールでメンションを飛ばせます。\n> 企画が終了したら「管理/無効化」からロールの削除と募集の無効化をしてください。",
			Flags:   1 << 6,
		},
	})
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_sanka(s *discordgo.Session, i *discordgo.InteractionCreate) {

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

	embedDescription := []rune(embed.Description)
	roleName := ""
	if len(embedDescription) <= 9 {
		roleName = string(embedDescription) + "..."
	} else {
		roleName = string(embedDescription[:9]) + "..."
	}
	roles, _ := s.GuildRoles(i.GuildID)
	for _, role := range roles {
		if role.Name == roleName {
			s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, role.ID)
		}
	}

	s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed)

	if simed {
		sime(s, i.GuildID, i.Message.ChannelID, i.Message.ID)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "> 参加を申し込みました :v: ",
			Flags:   1 << 6,
		},
	})
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_torikesi(s *discordgo.Session, i *discordgo.InteractionCreate) {

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

	embedDescription := []rune(embed.Description)
	roleName := ""
	if len(embedDescription) <= 9 {
		roleName = string(embedDescription) + "..."
	} else {
		roleName = string(embedDescription[:9]) + "..."
	}
	roles, _ := s.GuildRoles(i.GuildID)
	for _, role := range roles {
		if role.Name == roleName {
			s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, role.ID)
		}
	}

	s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "> 参加を取り消しました :wave: ",
			Flags:   1 << 6,
		},
	})
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_kanri(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var err error

	if i.Member.Permissions == 4398046511103 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
	} else {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "> 権限がありません。",
				Flags:   1 << 6,
			},
		})
	}

	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}

func ch_select(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response *discordgo.InteractionResponse
	data := i.MessageComponentData()
	messageID := i.Message.Content[46:65]

	switch data.Values[0] {
	case "sime":
		sime(s, i.GuildID, i.Message.ChannelID, messageID)

		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "> 締め切りました。",
				Flags:   1 << 6,
			},
		}
	case "syuugou":
		syuugou(s, i.GuildID, i.Message.ChannelID, messageID)

		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "> 集合をかけました。",
				Flags:   1 << 6,
			},
		}
	case "mukou":
		mukou(s, i.GuildID, i.Message.ChannelID, messageID)

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
		log.Printf("Critical error occurred: %v", err)
	}
}
