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
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	authorUsername := m.Author.Username
	mentionedGuildID := m.GuildID
	mentionedChannelID := m.ChannelID
	mentionedMessage := m.Message

	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "<@842360360006582283>") {
		return
	}

	s.ChannelTyping(mentionedChannelID)

	mentionedMessageContent := strings.Replace(m.Content, "<@842360360006582283>", "", -1)

	if len([]rune(mentionedMessageContent)) > 50 {
		s.ChannelMessageSendReply(
			m.ChannelID,
			"文章が長すぎるよ ><\n50文字以内で話しかけてね",
			&discordgo.MessageReference {
				MessageID: mentionedMessage.ID,
				ChannelID: mentionedChannelID,
				GuildID: mentionedGuildID,
		})
		return
	}

	commandResponseContent := "⚠️ 500 Internal Server Error"

	chatRequestSystemMessage := `ゲーム「モンスターハンター」のモンスターである「ギィギ」との会話をシミュレートします。
															彼は子供です。彼の発言サンプルを以下に列挙します。

															やあ！ボクの名前は「ギィギ」だよ！
															今日はいい天気だね！

															上記を参考に口調のみを模倣し、「ギィギ」になりきり回答を構築してください。
															自己紹介の必要はありません。`
	chatRequestSystemMessage	+= "また、会話相手の名前は" + authorUsername + "です。"

	chatGPTResp, err := openAiGptClient.CreateChatCompletion(
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
					Content: mentionedMessageContent,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Error while gpt3.5-turbo request: %v", err)
	} else {
		commandResponseContent = chatGPTResp.Choices[0].Message.Content
		log.Println("[Hey gigi]")
		log.Println("InputContent: " + mentionedMessageContent)
		log.Println("CommandResponseContent: " + commandResponseContent)
	}

	s.ChannelMessageSendReply(
		m.ChannelID,
		commandResponseContent,
		&discordgo.MessageReference {
			MessageID: mentionedMessage.ID,
			ChannelID: mentionedChannelID,
			GuildID: mentionedGuildID,
	})
}
