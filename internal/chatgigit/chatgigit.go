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
	openAiGptClient *openai.Client = openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	botMentionString string = "<@" + os.Getenv("BOT_USER_ID") + ">"
	botRoleMentionString string = "<@&" + os.Getenv("BOT_ROLE_ID") + ">"
)

// 返信用のMessageSendを構築
func buildReplyMessageSend(s *discordgo.Session, m *discordgo.MessageCreate) (replyMessageSend *discordgo.MessageSend) {
	var replyMessageContent string

	// 50文字を超える場合はリジェクトする
	if len([]rune(m.Message.Content)) > 50 {
		replyMessageContent = "⚠️ **ERROR**\n文章が長すぎるよ ><\n50文字以内で話しかけてね"
		return
	}

	// 再起的に会話リストを取得する
	var chatInputMessages []openai.ChatCompletionMessage
	buildChatInputMessages(s, &chatInputMessages, m.Message)

	// 累計の会話が５往復を超える場合はリジェクトする
	if len(chatInputMessages) > 10 {
		replyMessageContent = "⚠️ **ERROR**\n連続でできる会話は５往復までだよ ><\n新しくメンションして話しかけてね"
		return
	}

	// OpenAIから返答を取得
	replyMessageContent, err := getChatCompletion(chatInputMessages, m.Author.Username)
	if err != nil {
		log.Fatalf("Error while gpt3.5-turbo request: %v", err)
		replyMessageContent = "⚠️ **ERROR**\n500 Internal Server Error"
	}

	replyMessageSend = &discordgo.MessageSend{
		Content: replyMessageContent,
		Reference: &discordgo.MessageReference {
			MessageID: m.Message.ID,
			ChannelID: m.ChannelID,
			GuildID: m.GuildID,
		},
	}

	return
}

// リクエスト用の[]openai.ChatCompletionMessageを再帰的に構築する
func buildChatInputMessages(s *discordgo.Session, chatInputMessages *[]openai.ChatCompletionMessage, originMessage *discordgo.Message) {
	originMessageContent := strings.Replace(originMessage.Content, botMentionString, "", -1)

	// BOT or UserでRoleを分ける
	var role string
	if originMessage.Author.ID == s.State.User.ID {
		role = openai.ChatMessageRoleAssistant
	} else {
		role = openai.ChatMessageRoleUser
	}

	*chatInputMessages = append(*chatInputMessages, openai.ChatCompletionMessage{
		Role: role,
		Content: originMessageContent,
	})

	// 返信先がない=先頭のメッセージならば抜ける
	if originMessage.ReferencedMessage == nil {
		return
	}

	// message.ReferencedMessage では ReferencedMessage フィールドが取得されないため
	// ChannelMessage で明示的にメッセージを取得する
	// see: https://pkg.go.dev/github.com/bwmarrin/discordgo#Message
	referencedMessage, _ := s.ChannelMessage(originMessage.ChannelID, originMessage.ReferencedMessage.ID)
	buildChatInputMessages(s, chatInputMessages, referencedMessage)
}

// OpenAI API へ問い合わせる
func getChatCompletion(chatInputMessages []openai.ChatCompletionMessage, username string) (chatOutputMessageContent string, err error) {
	chatSystemPromptMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: chatSystemPrompt + "また、会話相手の名前は" + username + "です。",
	}
	chatInputMessages = append(chatInputMessages, chatSystemPromptMessage)

	// inputMessages は降順なので反転する
	for i := 0; i < len(chatInputMessages) / 2; i++ {
    chatInputMessages[i], chatInputMessages[len(chatInputMessages) - i - 1] = chatInputMessages[len(chatInputMessages) - i - 1], chatInputMessages[i]
	}

	openAiApiResponse, err := openAiGptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: chatInputMessages,
		},
	)
	if err != nil {
		chatOutputMessageContent = "⚠️ 500 Internal Server Error"
		return
	}

	chatOutputMessageContent = openAiApiResponse.Choices[0].Message.Content
	logConversation(chatInputMessages, chatOutputMessageContent)

	return
}
