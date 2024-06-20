package api

import (
	"context"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

type AIClient interface {
	GenerateText(ctx context.Context, content string, instructions string) (string, error)
}
type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient(apiKey)
	return &OpenAIClient{client: client}
}

func (c *OpenAIClient) GenerateText(ctx context.Context, content string, instructions string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: instructions + "\n" + content,
				},
			},
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	log.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}
