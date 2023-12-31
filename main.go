package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		panic("Falta la API KEY")
	}

	ctx := context.Background()
	client := openai.NewClient(apiKey)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3TextDavinci003,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You are an experienced lawyer",
			},
		},
		MaxTokens: 100,
	}

	rootCmd := &cobra.Command{
		Use: "chatgpt",
		Short: "Chat con ChatGPT en consola.",
		Run: func(cmd *cobra.Command, args []string){
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit{
				fmt.Print("Diga algo 'quit' para finalizar: ")
				if !scanner.Scan(){
					break
				}
				question := scanner.Text()
				switch question {
				case "quit":
					quit = true
				default:
					GetResponse(client, req, ctx, question)
				}

			}
		},
	}

	rootCmd.Execute()
}

func GetResponse(client *openai.Client, req openai.ChatCompletionRequest, ctx context.Context, question string, ) {
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		Content: question,
	})
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("CreateChatCompletion error: %v\n", err)
		return
	}

	fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
	req.Messages = append(req.Messages, resp.Choices[0].Message)
	fmt.Print("> ")
}