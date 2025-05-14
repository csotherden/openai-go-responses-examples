package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

func main() {
	ctx := context.Background()

	args := os.Args
	if len(args) < 3 {
		log.Fatalln("Usage: ./responses-file-search <file name> <user prompt>")
	}

	filePath := args[1]
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(fileName)

	userPrompt := args[2]

	if strings.ToLower(fileExt) != ".pdf" {
		log.Fatalln("Input file must be .pdf")
	}

	oaiClient := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	fileReader, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	inputFile := openai.File(fileReader, fileName, "application/pdf")

	storedFile, err := oaiClient.Files.New(ctx, openai.FileNewParams{
		File:    inputFile,
		Purpose: openai.FilePurposeUserData,
	})
	if err != nil {
		log.Fatalln(fmt.Sprintf("error uploading file to OpenAI: %s", err.Error()))
		return
	}

	_, err = oaiClient.VectorStores.Files.New(ctx, os.Getenv("VECTOR_STORE_ID"), openai.VectorStoreFileNewParams{
		FileID: storedFile.ID,
		Attributes: map[string]openai.VectorStoreFileNewParamsAttributeUnion{
			"file_type": {OfString: openai.String("application/pdf")},
		},
	})
	if err != nil {
		log.Fatalln("failed to associate file with Vector Store", err.Error())
	}

	params := responses.ResponseNewParams{
		Model:           openai.ChatModelGPT4o,
		Temperature:     openai.Float(0.7),
		MaxOutputTokens: openai.Int(512),
		Tools:           agentTools,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(userPrompt),
		},
	}

	resp, err := oaiClient.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(resp.OutputText())
}

// agentTools is the list of all tools available to the agent
var agentTools = []responses.ToolUnionParam{
	{
		OfFileSearch: &responses.FileSearchToolParam{
			VectorStoreIDs: []string{os.Getenv("VECTOR_STORE_ID")},
			MaxNumResults:  openai.Int(10),
			Filters: responses.FileSearchToolFiltersUnionParam{
				OfComparisonFilter: &shared.ComparisonFilterParam{
					Key:  "file_type",
					Type: "eq",
					Value: shared.ComparisonFilterValueUnionParam{
						OfString: openai.String("application/pdf"),
					},
				},
			},
			Type: "file_search",
		},
	},
}
