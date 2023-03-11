package chatgigit

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// メッセージ送信をすべて検知
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { return }

	// BOTへの返信、もしくはメンション付きメッセージで会話とみなす
	if (m.Message.ReferencedMessage != nil && m.Message.ReferencedMessage.Author.ID == s.State.User.ID) || strings.HasPrefix(m.Content, botMentionString) {
		reply(s, m)
		return
	}
}

func reply(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 「入力中...」の表示
	s.ChannelTyping(m.ChannelID)

	// 返信内容の生成
	replyMessageSend := buildReplyMessageSend(s, m)

	// 返信する
	s.ChannelMessageSendComplex(m.ChannelID, replyMessageSend)
}
