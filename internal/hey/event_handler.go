package hey

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

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { return }
	if !strings.HasPrefix(m.Content, botMentionString) { return }

	s.ChannelTyping(m.ChannelID)
	replyMessage := composeReplyMessage(m.Content, m.Author.Username)
	s.ChannelMessageSendReply(m.ChannelID, replyMessage, &discordgo.MessageReference {
		MessageID: m.Message.ID,
		ChannelID: m.ChannelID,
		GuildID: m.GuildID,
	})
}

func composeReplyMessage(messageContent string, username string) (replyMessage string) {
	mentionedMessageContent := strings.Replace(messageContent, botMentionString, "", -1)

	if len([]rune(mentionedMessageContent)) > 50 {
		replyMessage = "文章が長すぎるよ ><\n50文字以内で話しかけてね"
		return
	}

	replyMessage, err := getChatCompletion(mentionedMessageContent, username)
	if err != nil {
		log.Fatalf("Error while gpt3.5-turbo request: %v", err)
	}

	return
}

func getChatCompletion(inputMessage string, username string) (outputMessage string, err error) {
	chatRequestSystemMessage := readSystemRoleMessage() + "また、会話相手の名前は" + username + "です。"

	gptResponse, err := openAiGptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: chatRequestSystemMessage,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: inputMessage,
				},
			},
		},
	)
	if err != nil {
		outputMessage = "⚠️ 500 Internal Server Error"
		return
	}

	outputMessage = gptResponse.Choices[0].Message.Content

	log.Println("[Hey gigi]")
	log.Println("InputContent: " + inputMessage)
	log.Println("CommandResponseContent: " + outputMessage)

	return
}

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
