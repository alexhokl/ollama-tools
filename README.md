# ollama-tools

An sample CLI application checking weather of a city. This application utilises
tools calling of Ollama.

## Example

```sh
go run main.go generate -m gpt-oss:20b -f tools.json -q "what is the weather of London?"
```

The following is the content of `tools.json`.

```json
[
  {
    "type": "function",
    "function": {
      "name": "wttr.in",
      "description": "Get the current weather for a city",
      "parameters": {
        "type": "object",
        "required": [
          "city"
        ],
        "properties": {
          "city": {
            "type": "string",
            "description": "The name of the city"
          }
        }
      }
    }
  }
]
```

## Coding

To hard-code the tools, in can be done with the following code.

```go
functionProperties := map[string]api.ToolProperty{
  "city": {
    Type:        api.PropertyType{"string"},
    Description: "The name of the city",
  },
}

tools := []api.Tool{
  {
    Type: "function",
    Function: api.ToolFunction{
      Name:        "wttr.in",
      Description: "Get the current weather for a city",
      Parameters: api.ToolFunctionParameters{
        Type:       "object",
        Required:   []string{"city"},
        Properties: functionProperties,
      },
    },
  },
}
```

