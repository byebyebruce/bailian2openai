package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/byebyebruce/bailian2openai"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

// CLI run a CLI(Command Line Interface)
func CLI(p *bailian2openai.Proxy, systemPrompt string) {
	messages := []openai.ChatCompletionMessage{}
	if len(systemPrompt) > 0 {
		messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: systemPrompt})
		fmt.Println(color.CyanString("System:"))
		fmt.Println(color.CyanString(systemPrompt))
	}

	for {
		fmt.Println(color.GreenString("You:"))

		q := ""
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			q = strings.TrimSpace(scanner.Text())
			if len(q) > 0 {
				break
			}
		}

		m := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: q}

		resp, err := p.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Messages: append(messages, m),
		})
		if err != nil {
			fmt.Println(color.RedString("ERROR:%s", err.Error()))
			continue
		}

		chatGPT := resp.Choices[0].Message.Content
		fmt.Println(color.BlueString("AI:"))
		fmt.Println(color.YellowString(chatGPT))

		aiMessage := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: chatGPT}
		messages = append(messages, m, aiMessage)
		fmt.Println()
	}

}
