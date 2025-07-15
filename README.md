# mcp-web-scraper

This package uses Google Chrome's headless APIs to scrape web pages for AI/LLM agents.

Because it uses Chrome as its default user agent, any sites that require Javascript (for example, single page applications) should also be parsable with this tool.

It supports being called either from Go (go lang) via [LangChainGo](https://github.com/tmc/langchaingo/tree/main), or as an MCP server.

## MCP Server

First compile the code using `go`:

```sh
go build .
```

### Claude Desktop

```json
{
  "mcpServers": {
    "mcp-web-scraper": {
      "command": "/path/to/mcp-web-scraper",
      "args": []
    }
  }
}
```

### Visual Studio Code

```json
{
  "mcp": {
    "servers": {
      "mcp-web-scraper": {
        "command": "/path/to/mcp-web-scraper",
        "args": []
      }
    }
  }
}
```

## LangChainGo tool

Integration into langchain is easy:

```go
import 	"github.com/lmorg/mcp-web-scraper/langchain"

func example() {
    scraper := langchain.NewScraper()
}
```

Please consult the [langchaingo docs](https://tmc.github.io/langchaingo/docs/) for how to use tools with their libraries.

## Fallback Modes

### If Google Chrome is not installed

If you do not have Google Chrome installed, then mcp-web-scraper will fallback to use Go's HTTP user agent.

This will work in the majority of cases, however you might not get any content for sites that requires Javascript to render.
