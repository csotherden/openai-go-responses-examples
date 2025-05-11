package main

import (
	"context"
	"encoding/json"
	"github.com/invopop/jsonschema"
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
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "TIOBE-Index",
					Schema:      TIOBEIndexSchema,
					Strict:      openai.Bool(true),
					Description: openai.String("JSON Schema for the TIOBE Programming Community Index results"),
					Type:        "json_schema",
				},
			},
		},
		Model:           openai.ChatModelGPT4o,
		Temperature:     openai.Float(0.7),
		MaxOutputTokens: openai.Int(2048),
		Tools:           agentTools,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Please provide me with the top ten results from the latest TIOBE index"),
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
	{OfWebSearch: &responses.WebSearchToolParam{Type: "web_search_preview"}},
}

type TIOBEIndex struct {
	IndexVersion string                       `json:"index_version" jsonschema_description:"The index version the results are from in MMMM yyyy format"`
	Rankings     []ProgrammingLanguageRanking `json:"rankings" jsonschema_description:"The ordered ranking results"`
}

type ProgrammingLanguageRanking struct {
	Name             string  `json:"name" jsonschema_description:"Programming language name"`
	CurrentRanking   int     `json:"current_ranking" jsonschema_description:"Where the language ranks in the current index"`
	PriorYearRanking int     `json:"prior_year_ranking" jsonschema_description:"Where the language ranked in the index 12 months prior"`
	Rating           float64 `json:"rating" jsonschema_description:"The popularity share for the programming language"`
	RatingChange     float64 `json:"rating_change" jsonschema_description:"The year over year ratings change"`
}

var TIOBEIndexSchema = GenerateSchema[TIOBEIndex]()

func GenerateSchema[T any]() map[string]interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	schemaJson, err := schema.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var schemaObj map[string]interface{}
	err = json.Unmarshal(schemaJson, &schemaObj)
	if err != nil {
		panic(err)
	}

	return schemaObj
}
