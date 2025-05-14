package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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
		Tools:           agentTools,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("What's the current stock price for Apple?"),
		},
	}

	resp, err := oaiClient.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	params.PreviousResponseID = openai.String(resp.ID)

	params.Input = responses.ResponseNewParamsInputUnion{}

	for _, output := range resp.Output {
		if output.Type == "function_call" {
			toolCall := output.AsFunctionCall()

			result, err := processToolCall(ctx, toolCall)
			if err != nil {
				params.Input.OfInputItemList = append(params.Input.OfInputItemList, responses.ResponseInputItemParamOfFunctionCallOutput(toolCall.CallID, err.Error()))
			} else {
				params.Input.OfInputItemList = append(params.Input.OfInputItemList, responses.ResponseInputItemParamOfFunctionCallOutput(toolCall.CallID, result))
			}
		}
	}

	// No tools calls made, we already have our final response
	if len(params.Input.OfInputItemList) == 0 {
		log.Println(resp.OutputText())
		return
	}

	// Make a final call with our tools results and no tools to get the final output
	params.Tools = nil
	resp, err = oaiClient.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(resp.OutputText())
}

// getStockTool defines the OpenAI tool for getting a single Stock by ID
var getStockTool = responses.ToolUnionParam{
	OfFunction: &responses.FunctionToolParam{
		Name:        "get_stock_price",
		Description: openai.String("The get_stock_price tool retrieves the current price of a single stock by it's ticker symbol"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"symbol": map[string]string{
					"type":        "string",
					"description": "The ticker symbol of the stock to retrieve",
				},
			},
			"required": []string{"symbol"},
		},
	},
}

// GetStockPriceArgs represents the arguments for the get_stock_price function
type GetStockPriceArgs struct {
	StockSymbol string `json:"symbol"`
}

// GetStockPrice is a mockup implementation of the get_stock_price function
func GetStockPrice(ctx context.Context, args []byte) (string, error) {
	var getArgs GetStockPriceArgs
	if err := json.Unmarshal(args, &getArgs); err != nil {
		return "", fmt.Errorf("failed to parse get_stock_price arguments: %w", err)
	}

	// Validate the stock symbol
	if strings.TrimSpace(getArgs.StockSymbol) == "" {
		return "", fmt.Errorf("stock symbol is required")
	}

	// Return a static placeholder
	return "$198.53 USD", nil
}

// agentTools is the list of all tools available to the agent
var agentTools = []responses.ToolUnionParam{
	{OfWebSearch: &responses.WebSearchToolParam{Type: "web_search_preview"}},
	getStockTool,
}

// processToolCall handles a tool call from the OpenAI API
func processToolCall(ctx context.Context, toolCall responses.ResponseFunctionToolCall) (string, error) {
	switch toolCall.Name {
	case "get_stock_price":
		return GetStockPrice(ctx, []byte(toolCall.Arguments))
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Name)
	}
}
