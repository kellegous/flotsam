# Filtering MCP Tools with jq

Example code for the post [Filtering MCP Tools with jq](https://kellegous.com/j/). The program runs a small in-process MCP server (plain and jq-filtering variants) and an interactive agent that can call those tools, so you can compare behavior with and without jq-shaped tool inputs.

## Requirements

- [Go](https://go.dev/dl/) 1.26.1 or newer (see `go.mod`)

## Build

From this directory:

```sh
make
```

That produces `bin/agent`. Alternatively:

```sh
go build -o bin/agent ./cmd/agent
```

## Run

```sh
./bin/agent
```

The process listens on `--http.addr` (default `localhost:4000`) for MCP over HTTP and reads prompts from stdin. Type a question at the `>` prompt; use `/exit` or `/quit` to stop.

Useful flags:

| Flag                    | Purpose                                                                                  |
| ----------------------- | ---------------------------------------------------------------------------------------- |
| `--use-jq`              | Point the agent at the jq-aware MCP surface (`/jq`) instead of the plain one (`/plain`). |
| `--http.addr host:port` | Bind address for the MCP HTTP server.                                                    |
| `--model …`             | LLM provider and credentials (see below).                                                |
| `--log-file path`       | After exit, write a JSON log of events to `path`.                                        |

The default `--model` is tuned to the author’s environment. For your machine, pass an explicit `--model` (see next section).

## Choosing an API provider (`--model`)

`--model` uses the form `provider:option=value,option=value,…`. Supported `provider` values are `anthropic`, `openai`, and `local` (OpenAI-compatible HTTP API).

### Anthropic

`api-key` is required. `model` defaults to `claude-sonnet-4-6` if omitted.

```sh
./bin/agent --model 'anthropic:api-key=YOUR_ANTHROPIC_API_KEY'
./bin/agent --model 'anthropic:api-key=YOUR_KEY,model=claude-sonnet-4-6'
```

To avoid embedding a literal key in your shell history, you can pass through a variable you set in the same session, for example `--model "anthropic:api-key=$ANTHROPIC_API_KEY"` after `export ANTHROPIC_API_KEY=…`.

### OpenAI

`api-key` is required. `model` defaults to `gpt-4o-mini`. Optional `base-url` targets the [OpenAI API](https://platform.openai.com/) or any OpenAI-compatible endpoint (many hosted and local stacks use this shape).

```sh
./bin/agent --model 'openai:api-key=YOUR_OPENAI_API_KEY'
./bin/agent --model 'openai:api-key=YOUR_KEY,model=gpt-4o-mini'
./bin/agent --model 'openai:api-key=YOUR_KEY,model=MODEL_NAME,base-url=https://api.example.com/v1'
```

### Local / OpenAI-compatible (Ollama, LM Studio, vLLM, etc.)

`local` uses the same client as `openai` but with different defaults in code. Override `base-url`, `api-key`, and `model` for your server. For [Ollama](https://ollama.com/)’s OpenAI-compatible API, a typical setup is:

```sh
./bin/agent --model 'local:base-url=http://127.0.0.1:11434/v1,api-key=ollama,model=llama3.2'
```

Adjust `base-url` and `model` to match your stack’s documentation.

### Example: jq vs plain

```sh
./bin/agent --use-jq --model 'openai:api-key=YOUR_KEY'
```

Without `--use-jq`, the agent uses the plain MCP tools; with it, tools accept jq expressions as in the blog post.
