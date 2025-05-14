package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
)

func main() {
	ctx := context.Background()

	oaiClient := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	params := responses.ResponseNewParams{
		Model:           openai.ChatModelGPT4o,
		Temperature:     openai.Float(0.7),
		MaxOutputTokens: openai.Int(512),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("What were the original design goals for the Go programming language?"),
		},
	}

	resp, err := oaiClient.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(resp.OutputText())
}
