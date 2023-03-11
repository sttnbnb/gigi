package chatgigit

import (
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var (
	chatSystemPrompt string
)

func init() {
	loadChatSystemPrompt()
}

// 外部に用意した設定ファイルを読み込む
func loadChatSystemPrompt() {
	f, err := os.Open("assets/chat_gigit/chat_system_prompt.txt")
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}

	data := make([]byte, 1024)
	count, err := f.Read(data)
	if err != nil {
		log.Fatalf("Cannot read file: %v", err)
	}

	chatSystemPrompt = string(data[:count])
}

// 記録用に会話をログに出力する
func logConversation(chatInputMessages []openai.ChatCompletionMessage, chatOutputMessageContent string) {
	log.Println("")
	log.Println("[ChatGiGiT]")

	for _, inputMessage := range chatInputMessages[1:] {
		log.Println("Role: " + inputMessage.Role)
		log.Println("Content: " + inputMessage.Content)
	}

	log.Println("Role: assistant")
	log.Println(chatOutputMessageContent)
}
