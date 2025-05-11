# OpenAI Responses API Examples in Go

This repository contains a collection of fully working Go examples demonstrating how to use the new [OpenAI Responses API](https://platform.openai.com/docs/guides/responses) with the official [openai-go SDK](https://github.com/openai/openai-go).

Each directory contains a self-contained example that maps to a specific concept covered in the blog post:

**Read the full post here:**  
[https://chris.sotherden.io/openai-responses-api-using-go](https://chris.sotherden.io/openai-responses-api-using-go)

---

## Examples Included

| Directory                       | Description                                                                            |
|--------------------------------|----------------------------------------------------------------------------------------|
| `responses`                    | Basic single-prompt usage with the Responses API                                       |
| `responses-file-search`        | Uploads files to a vector store and perform semantic search using the file_search tool |
| `responses-input-file`         | Shows how to include input files as part of the prompt context                         |
| `responses-state-management`   | Demonstrates multi-turn conversation state management using `PreviousResponseID`       |
| `responses-structured-output`  | Uses structured output with JSON Schema to control the response format                 |
| `responses-tool-calling`       | Shows how to enable and handle tool calling and web search                             |

---

## Prerequisites

- Go 1.18 or higher
- OpenAI API key (set `OPENAI_API_KEY` as an environment variable)
- An existing Vector Store for the file search example (set `VECTOR_STORE_ID` as an environment variable)
- Modules initialized via `go mod tidy` 

---

## Usage

Clone this repository:

```bash
git clone https://github.com/csotherden/openai-responses-go-examples.git
cd openai-responses-go-examples
go mod tidy
```

Run any individual example:

```bash
cd responses-structured-output
go run main.go
```

---

## Related Reading

- [OpenAI Responses API Overview](https://platform.openai.com/docs/guides/responses)
- [Official OpenAI Go SDK](https://github.com/openai/openai-go)
- [My blog post: Using the OpenAI Responses API with Go](https://chris.sotherden.io/openai-responses-api-using-go)

---

## Feedback & Contributions

Feel free to open an issue or PR if you’d like to expand or improve the examples!

---

© 2025 Chris Sotherden
