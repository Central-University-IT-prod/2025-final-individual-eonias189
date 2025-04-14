package openai

import (
	"advertising/pkg/logger"
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type Chat interface {
	Completion(ctx context.Context, messages []Message, options ...ChatCompletionOption) (string, error)
}

type chatImpl struct {
	cli   *openai.Client
	model string
}

func (ci *chatImpl) Completion(ctx context.Context, messages []Message, options ...ChatCompletionOption) (string, error) {
	cfg := &ChatCompletionConfig{
		Model:       ci.model,
		Temperature: 0.5,
		MaxTokens:   5000,
		MaxAttempts: 5,
	}

	applyOptions(cfg, options)

	openaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, message := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	var (
		resp   openai.ChatCompletionResponse
		err    error
		attemt int
	)

	for {
		if attemt == cfg.MaxAttempts {
			return "", ErrNoResponse
		}

		logger.FromCtx(ctx).Debug("request to openai")
		resp, err = ci.cli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:       cfg.Model,
			Messages:    openaiMessages,
			MaxTokens:   cfg.MaxTokens,
			Temperature: cfg.Temperature,
		})
		if err != nil {
			return "", err
		}

		if len(resp.Choices) != 0 {
			break
		}
		attemt++

	}

	return resp.Choices[0].Message.Content, nil

}

func NewChat(cfg Config) Chat {
	openaiCfg := openai.DefaultConfig(cfg.OpenAIApiKey)
	openaiCfg.BaseURL = cfg.OpenAIBaseUrl

	cli := openai.NewClientWithConfig(openaiCfg)

	return &chatImpl{
		cli:   cli,
		model: cfg.OpenAIModel,
	}
}
