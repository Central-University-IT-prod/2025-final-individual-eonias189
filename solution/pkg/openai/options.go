package openai

type ChatCompletionConfig struct {
	Model       string
	Temperature float32
	MaxTokens   int
	MaxAttempts int
}

type ChatCompletionOption func(rcfg *ChatCompletionConfig)

func WithModel(model string) ChatCompletionOption {
	return func(rcfg *ChatCompletionConfig) {
		rcfg.Model = model
	}
}

func WithTemperature(temperature float32) ChatCompletionOption {
	return func(rcfg *ChatCompletionConfig) {
		rcfg.Temperature = temperature
	}
}

func WithMaxTokens(maxTokens int) ChatCompletionOption {
	return func(rcfg *ChatCompletionConfig) {
		rcfg.MaxTokens = maxTokens
	}
}

func WithMaxAttempts(maxAttempts int) ChatCompletionOption {
	return func(rcfg *ChatCompletionConfig) {
		rcfg.MaxAttempts = maxAttempts
	}
}

func applyOptions(cfg *ChatCompletionConfig, options []ChatCompletionOption) {
	for _, option := range options {
		option(cfg)
	}
}
