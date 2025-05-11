package main

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	oaiClient := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	params := responses.ResponseNewParams{
		Model:           openai.ChatModelGPT4o,
		Temperature:     openai.Float(0.7),
		MaxOutputTokens: openai.Int(10240),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("What were the original design goals for the Go programming language?"),
		},
		Store: openai.Bool(true), // This is already TRUE by default
	}

	resp, err := oaiClient.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(resp.OutputText())

	params.PreviousResponseID = openai.String(resp.ID)

	params.Input = responses.ResponseNewParamsInputUnion{
		OfString: openai.String("Who created it?"),
	}

	resp, err = oaiClient.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(resp.OutputText())
}
