package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

type OpenAIAdapter struct {
	APIKey string
}

func (adapter *OpenAIAdapter) FetchCompletionStream(req openai.CompletionRequest, respCh chan<- string) error {
	defer close(respCh)

	client := openai.NewClient(adapter.APIKey)
	ctx := context.Background()

	stream, err := client.CreateCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("error creating completion stream: %v", err)
	}

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("stream error: %v", err)
		}

		respCh <- response.Choices[0].Text
	}
}
