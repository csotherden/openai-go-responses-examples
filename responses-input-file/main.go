package main

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	ctx := context.Background()

	args := os.Args
	if len(args) < 2 {
		log.Fatalln("Usage: ./responses-input-file <file name>")
	}

	filePath := args[1]
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(fileName)

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

	params := responses.ResponseNewParams{
		Model:           openai.ChatModelGPT4o,
		Temperature:     openai.Float(0.7),
		MaxOutputTokens: openai.Int(512),
	}

	params.Input = responses.ResponseNewParamsInputUnion{
		OfInputItemList: responses.ResponseInputParam{
			responses.ResponseInputItemParamOfMessage(
				responses.ResponseInputMessageContentListParam{
					responses.ResponseInputContentUnionParam{
						OfInputFile: &responses.ResponseInputFileParam{
							FileID: openai.String(storedFile.ID),
							Type:   "input_file",
						},
					},
					responses.ResponseInputContentUnionParam{
						OfInputText: &responses.ResponseInputTextParam{
							Text: "Provide a one paragraph summary of the provided document.",
							Type: "input_text",
						},
					},
				},
				"user",
			),
		},
	}

	resp, err := oaiClient.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(resp.OutputText())
}
