package cmd

import (
	"fmt"
	"os/exec"

	"github.com/alexhokl/helper/cli"
	"github.com/alexhokl/helper/iohelper"
	"github.com/alexhokl/helper/jsonhelper"
	"github.com/goforj/godump"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

type generateOptions struct {
	modelName               string
	toolsDefinitionFilePath string
	question                string
	debug                   bool
}

var generateOps generateOptions

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate function/method calls by asking a question",
	RunE:  runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	flags := generateCmd.Flags()
	flags.StringVarP(&generateOps.modelName, "model", "m", "llama3.1:8b", "Model to use")
	flags.StringVarP(&generateOps.toolsDefinitionFilePath, "file", "f", "", "Path to tools definition file")
	flags.StringVarP(&generateOps.question, "question", "q", "", "Question to ask")
	flags.BoolVarP(&generateOps.debug, "debug", "d", false, "Enable debug mode")

	if err := generateCmd.MarkFlagRequired("file"); err != nil {
		cli.LogUnableToMarkFlagAsRequired("file", err)
	}
	if err := generateCmd.MarkFlagRequired("question"); err != nil {
		cli.LogUnableToMarkFlagAsRequired("question", err)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	toolDefinitionString, err := iohelper.ReadStringFromFile(generateOps.toolsDefinitionFilePath)
	if err != nil {
		return fmt.Errorf("unable to read from tools definition file %w", err)
	}

	var tools []api.Tool
	err = jsonhelper.ParseJSONString(toolDefinitionString, &tools)
	if err != nil {
		return fmt.Errorf("unable to parse tools definition file %w", err)
	}

	if generateOps.debug {
		godump.Dump(tools)
	}

	// create the API client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	messages := []api.Message{
		{
			Role:    "user",
			Content: generateOps.question,
		},
	}

	req := &api.ChatRequest{
		Model:    generateOps.modelName,
		Messages: messages,
		Tools:    tools,
	}

	var toolCalls = make([]api.ToolCall, 0)

	ctx := cmd.Context()
	respFunc := func(resp api.ChatResponse) error {
		toolCalls = append(toolCalls, resp.Message.ToolCalls...)
		return nil
	}

	if err := client.Chat(ctx, req, respFunc); err != nil {
		return err
	}

	if len(toolCalls) == 0 {
		fmt.Println("Sorry I don't understand your query on weather")
		return nil
	}

	for _, toolCall := range toolCalls {
		fmt.Printf("%s/%s\n", toolCall.Function.Name, toolCall.Function.Arguments["city"])

		byteOutput, err := exec.Command("curl", fmt.Sprintf("%s/%s", toolCall.Function.Name, toolCall.Function.Arguments["city"])).Output()
		if err != nil {
			return fmt.Errorf("unable to execute tool call %w", err)
		}
		fmt.Println(string(byteOutput))
	}

	return nil
}
