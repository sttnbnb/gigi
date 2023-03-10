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
	replyMessage := buildReplyMessage(s, m)

	// 返信する
	s.ChannelMessageSendReply(m.ChannelID, replyMessage, &discordgo.MessageReference {
		MessageID: m.Message.ID,
		ChannelID: m.ChannelID,
		GuildID: m.GuildID,
	})
}

func buildReplyMessage(s *discordgo.Session, m *discordgo.MessageCreate) (replyMessage string) {
	// 50文字以上を超える場合はリジェクトする
	if len([]rune(m.Message.Content)) > 50 {
		replyMessage = "⚠️ **ERROR**\n文章が長すぎるよ ><\n50文字以内で話しかけてね"
		return
	}

	// 再起的に会話リストを取得する
	var chatInputMessages []openai.ChatCompletionMessage
	buildChatInputMessages(s, &chatInputMessages, m.Message)

	// 累計の会話が５往復を超える場合はリジェクトする
	if len(chatInputMessages) > 10 {
		replyMessage = "⚠️ **ERROR**\n連続でできる会話は５往復までだよ ><\n新しくメンションして話しかけてね"
		return
	}

	// OpenAIから返答を取得
	replyMessage, err := getChatCompletion(chatInputMessages, m.Author.Username)
	if err != nil {
		log.Fatalf("Error while gpt3.5-turbo request: %v", err)
		replyMessage = "⚠️ **ERROR**\n500 Internal Server Error"
	}

	return
}

// リクエスト用のInputを再帰的に構築する
func buildChatInputMessages(s *discordgo.Session, conversationArray *[]openai.ChatCompletionMessage, message *discordgo.Message) {
	messageContent := strings.Replace(message.Content, botMentionString, "", -1)

	// BOT or UserでRoleを分ける
	var role string
	if message.Author.ID == s.State.User.ID {
		role = openai.ChatMessageRoleAssistant
	} else {
		role = openai.ChatMessageRoleUser
	}

	*conversationArray = append(*conversationArray, openai.ChatCompletionMessage{
		Role: role,
		Content: messageContent,
	})

	// 返信先がない=先頭のメッセージならば抜ける
	if message.ReferencedMessage == nil {
		return
	}

	// message.ReferencedMessage では ReferencedMessage フィールドが取得されないため
	// ChannelMessage で明示的にメッセージを取得する
	// see: https://pkg.go.dev/github.com/bwmarrin/discordgo#Message
	referencedMessage, _ := s.ChannelMessage(message.ChannelID, message.ReferencedMessage.ID)
	buildChatInputMessages(s, conversationArray, referencedMessage)
}

// OpenAI API へ問い合わせる
func getChatCompletion(inputMessages []openai.ChatCompletionMessage, username string) (outputMessage string, err error) {
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: readSystemRoleMessage() + "また、会話相手の名前は" + username + "です。",
	}
	inputMessages = append(inputMessages, systemMessage)

	// inputMessages は降順なので反転する
	for i := 0; i < len(inputMessages) / 2; i++ {
    inputMessages[i], inputMessages[len(inputMessages) - i - 1] = inputMessages[len(inputMessages) - i - 1], inputMessages[i]
	}

	gptResponse, err := openAiGptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: inputMessages,
		},
	)
	if err != nil {
		outputMessage = "⚠️ 500 Internal Server Error"
		return
	}

	outputMessage = gptResponse.Choices[0].Message.Content
	return
}

// 外部に用意した設定ファイルを読み込む
func readSystemRoleMessage() string {
	f, err := os.Open("assets/ChatGiGiT_SystemRoleMessage.txt")
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}

	data := make([]byte, 1024)
	count, err := f.Read(data)
	if err != nil {
		log.Fatalf("Cannot read file: %v", err)
	}

	return string(data[:count])
}
