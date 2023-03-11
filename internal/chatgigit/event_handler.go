package chatgigit

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// メッセージ送信をすべて検知
func MessageCreateEventHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { return }

	// 以下で会話とみなし発火
	// ・BOTへのメンション
	// ・BOTロールへのメンション
	// ・BOTへの返信（エラーメッセージを除く）
	if strings.HasPrefix(m.Content, botUserMentionString) ||
		strings.HasPrefix(m.Content, botRoleMentionString) ||
		(m.Message.ReferencedMessage != nil &&
			m.Message.ReferencedMessage.Author.ID == s.State.User.ID &&
			!strings.HasPrefix(m.Message.ReferencedMessage.Content, "⚠️ **ERROR**")) {
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
