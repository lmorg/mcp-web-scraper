# mcp-web-scraper

This package uses Google Chrome's headless APIs to scrape web pages.

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

## Fallback Modes

### If Google Chrome is not installed

If you do not have Google Chrome installed, then mcp-web-scraper will fallback to use Go's HTTP user agent.

This will work in the majority of cases, however you might not get any content for sites that requires Javascript to render.

### If the webpage cannot be converted into Markdown 

By default, this tool converts pages to Markdown.

It does this to reduce the token count, which is valuable when your LLM provider charges you by tokens used. 

However if the page cannot be converted into markdown for any reason, the HTML document is returned instead. This failure mode ensures that this tool will work the vast majority of the time, regardless of what weirdness site owners return.

To reduce the token count of any HTML which is sent, this tool strips SVGs, Javascript and other HTML tags which are obviously irrelevant.