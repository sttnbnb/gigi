package hey

import (
	"context"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-openai"
)

var (
	openAiGptClient *openai.Client = openai.NewClient(os.Getenv("OPENAI_TOKEN"))
)

func GetCommandHandlersMap() (commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hey": c_hey,
	}

	return
}

func c_hey(s *discordgo.Session, i *discordgo.InteractionCreate) {
	InputContent := i.ApplicationCommandData().Options[0].StringValue()
	CommandResponseContent := "⚠️ 500 Internal Server Error"

	if len([]rune(InputContent)) > 50 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "文章が長すぎるよ ><\n50文字以内で話しかけてね",
				Flags: 1 << 6,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	username := ""
	if i.Member != nil {
		username = i.Member.User.Username
	} else {
		username = i.User.Username
	}
	chatRequestSystemMessage := `ゲーム「モンスターハンター」のキャラクター「ギィギ」との会話をシミュレートします。
															彼は子供です。彼の発言サンプルを以下に列挙します。

															やあ！ボクの名前は「ギィギ」だよ！
															今日はいい天気だね！

															上記を参考に口調のみを模倣し、「モンスターハンター」のキャラクターになりきり回答を構築してください。
															自己紹介の必要はありません。`
	chatRequestSystemMessage	+= "また、会話相手の名前は" + username + "です。"

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
					Content: InputContent,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Error while gpt3.5-turbo request: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &CommandResponseContent,
		})
	}

	CommandResponseContent = chatGPTResp.Choices[0].Message.Content
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &CommandResponseContent,
	})

	log.Println("[Hey gigi]")
	log.Println("InputContent: " + InputContent)
	log.Println("CommandResponseContent: " + CommandResponseContent)
}
