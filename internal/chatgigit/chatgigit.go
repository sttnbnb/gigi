package chatgigit

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-openai"
)

var (
	openAiGptClient      *openai.Client = openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	botUserMentionString string         = "<@" + os.Getenv("BOT_USER_ID") + ">"
	botRoleMentionString string         = "<@&" + os.Getenv("BOT_ROLE_ID") + ">"
)

// 返信用のMessageSendを構築
func buildReplyMessageSend(s *discordgo.Session, m *discordgo.MessageCreate) (replyMessageSend *discordgo.MessageSend) {
	var replyMessageContent string
	replyMessageReference := &discordgo.MessageReference{
		MessageID: m.Message.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}

	// 50文字を超える場合はリジェクトする
	if len([]rune(m.Message.Content)) > 50 {
		replyMessageContent = "⚠️ **ERROR**\n文章が長すぎるよ ><\n50文字以内で話しかけてね"
		replyMessageSend = &discordgo.MessageSend{
			Content:   replyMessageContent,
			Reference: replyMessageReference,
		}
		return
	}

	// 再起的に会話リストを取得する
	var chatInputMessages []openai.ChatCompletionMessage
	buildChatInputMessages(s, &chatInputMessages, m.Message)

	// 累計の会話が５往復を超える場合はリジェクトする
	if len(chatInputMessages) > 10 {
		replyMessageContent = "⚠️ **ERROR**\n連続でできる会話は５往復までだよ ><\n新しくメンションして話しかけてね"
		replyMessageSend = &discordgo.MessageSend{
			Content:   replyMessageContent,
			Reference: replyMessageReference,
		}
		return
	}

	// OpenAIから返答を取得
	replyMessageContent = getChatCompletion(chatInputMessages, m.Author.Username)
	replyMessageSend = &discordgo.MessageSend{
		Content:   replyMessageContent,
		Reference: replyMessageReference,
	}

	return
}

// 返信を再帰的に遡り[]openai.ChatCompletionMessageを構築する
func buildChatInputMessages(s *discordgo.Session, chatInputMessages *[]openai.ChatCompletionMessage, originMessage *discordgo.Message) {
	originMessageContent := strings.Replace(originMessage.Content, botUserMentionString, "", -1)

	var role string
	if originMessage.Author.ID == s.State.User.ID {
		role = openai.ChatMessageRoleAssistant
	} else {
		role = openai.ChatMessageRoleUser
	}

	*chatInputMessages = append(*chatInputMessages, openai.ChatCompletionMessage{
		Role:    role,
		Content: originMessageContent,
	})

	if originMessage.ReferencedMessage == nil {
		return
	}

	// message.ReferencedMessage では取得した ReferencedMessage の ReferencedMessage フィールドが取得されないため
	// ChannelMessage で明示的にメッセージを取得する
	// see: https://pkg.go.dev/github.com/bwmarrin/discordgo#Message
	referencedMessage, _ := s.ChannelMessage(originMessage.ChannelID, originMessage.ReferencedMessage.ID)
	buildChatInputMessages(s, chatInputMessages, referencedMessage)
}

// OpenAI API へ問い合わせる
func getChatCompletion(chatInputMessages []openai.ChatCompletionMessage, username string) (chatOutputMessageContent string) {
	// プロンプトの構築
	chatSystemPromptMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: chatSystemPrompt + "また、会話相手の名前は" + username + "です。",
	}
	chatInputMessages = append(chatInputMessages, chatSystemPromptMessage)

	// inputMessages は降順になっているので反転する
	for i := 0; i < len(chatInputMessages)/2; i++ {
		chatInputMessages[i], chatInputMessages[len(chatInputMessages)-i-1] = chatInputMessages[len(chatInputMessages)-i-1], chatInputMessages[i]
	}

	res, err := openAiGptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    chatInputMessages,
			Temperature: 0.7,
		},
	)
	if err != nil {
		log.Printf("Error while OpenAI API request: %v", err)
		chatOutputMessageContent = "⚠️ **ERROR** [500 Internal Server Error]\n" + err.Error()
		return
	}

	chatOutputMessageContent = res.Choices[0].Message.Content
	logConversation(chatInputMessages, chatOutputMessageContent)

	return
}
