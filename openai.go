package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func gptv(prompt string, imageBase64 string) (string, error) {
	requestURL := "https://api.openai.com/v1/chat/completions"
	request := GPTVRequest{
		Model: "gpt-4-vision-preview",
		Messages: []ReqMessage{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "image_url",
						ImageURL: ImageURL{
							URL: fmt.Sprintf("data:image/jpeg;base64,%s", imageBase64),
						},
					},
				},
			},
		},
		MaxTokens: 1024,
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	bodyReader := bytes.NewReader(requestJson)

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))

	client := http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	response := GPTVResponse{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

// structs

type GPTVRequest struct {
	Model     string       `json:"model"`
	Messages  []ReqMessage `json:"messages"`
	MaxTokens int          `json:"max_tokens"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Content struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

type ReqMessage struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type GPTVResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type FinishDetails struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}

type ResMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message       ResMessage    `json:"message"`
	FinishDetails FinishDetails `json:"finish_details"`
	Index         int           `json:"index"`
}
