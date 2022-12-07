package eula

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var ReadmeRoleID = os.Getenv("README_ROLE_ID")

func GetComponentHandlersMap() (componentsHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ch_agree":    ch_agree,
		"ch_disagree": ch_disagree,
	}

	return
}

func ch_agree(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, ReadmeRoleID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "> 同意しました。",
			Flags:   1 << 6,
		},
	})
	if err != nil {
		panic(err)
	}
}

func ch_disagree(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, ReadmeRoleID)
	channel, _ := s.UserChannelCreate(i.Member.User.ID)
	invite, _ := s.ChannelInviteCreate(i.ChannelID, discordgo.Invite{})
	s.ChannelMessageSend(channel.ID, "マジで押すことないやん :sweat_smile: \n https://discord.gg/"+invite.Code)
	s.GuildMemberDelete(i.GuildID, i.Member.User.ID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "> 拒否しました。",
			Flags:   1 << 6,
		},
	})
	if err != nil {
		panic(err)
	}
}
