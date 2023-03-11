package chatgigit

import (
	"log"
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
	err := s.ChannelTyping(m.ChannelID)
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}

	// 返信内容の生成
	replyMessageSend := buildReplyMessageSend(s, m)

	// 返信する
	_, err = s.ChannelMessageSendComplex(m.ChannelID, replyMessageSend)
	if err != nil {
		log.Printf("Critical error occurred: %v", err)
	}
}
