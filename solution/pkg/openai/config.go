package openai

type Config struct {
	OpenAIBaseUrl string `env:"OPENAI_BASE_URL" env-default:"https://openrouter.ai/api/v1"`
	OpenAIApiKey  string `env:"OPENAI_API_KEY"`
	OpenAIModel   string `env:"OPENAI_MODEL" env-default:"deepseek/deepseek-chat:free"`
}
