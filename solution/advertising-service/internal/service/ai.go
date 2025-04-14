package service

import (
	"advertising/pkg/openai"
	"context"
	"fmt"
	"strings"
)

type AIService struct {
	chat openai.Chat
}

func NewAIService(chat openai.Chat) *AIService {
	return &AIService{
		chat: chat,
	}
}

func (ais *AIService) GenerateAdText(ctx context.Context, adTitle string) (string, error) {
	op := "AIService.GenerateAdText"

	prompt := fmt.Sprintf(`generate text for advertisment post with title "%s", use language of title`, adTitle)

	adText, err := ais.chat.Completion(ctx, []openai.Message{
		{
			Role:    openai.RoleUser,
			Content: prompt,
		},
	}, openai.WithTemperature(0.8))
	if err != nil {
		return "", fmt.Errorf("%s: chat.Completion: %w", op, err)
	}

	return adText, nil
}

func (ais *AIService) ModerateAdText(ctx context.Context, adText string) (bool, []string, error) {
	op := "AIService.ModerateAdText"

	systemPrompt := "Analyze user`s message strictly according to the following rules: " +
		`1. If there are any obscene or illegal phrases, list them separated by commas ` +
		`2. If there are no violations, respond with "true"`

	res, err := ais.chat.Completion(ctx, []openai.Message{
		{
			Role:    openai.RoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.RoleUser,
			Content: adText,
		},
	}, openai.WithMaxTokens(10000))
	if err != nil {
		return false, nil, fmt.Errorf("%s: chat.Completion: %w", op, err)
	}

	res = strings.TrimSpace(res)
	if strings.EqualFold(res, "true") {
		return true, []string{}, nil
	} else {
		illegalPhrases := strings.Split(res, ",")

		for i, phrase := range illegalPhrases {
			illegalPhrases[i] = strings.TrimSpace(phrase)
		}

		return false, illegalPhrases, nil
	}
}
