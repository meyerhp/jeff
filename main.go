package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	API_KEY = ""
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatCompletion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

func main() {
	loop()
}

func loop() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Q: ")
		scanner.Scan()
		userInput := scanner.Text()
		request(userInput)
	}
}

func request(prompt string) {
	apiEndpoint := "https://api.openai.com/v1/chat/completions"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", API_KEY),
	}

	chatRequest := ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a terminal assistant called Jeff."},
			{Role: "user", Content: "What is the command to list the files in the current dir?"},
			{Role: "system", Content: "ls"},
			{Role: "user", Content: prompt},
		},
	}

	p, err := json.Marshal(chatRequest)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(p))
	if err != nil {
		return
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var chatCompletion ChatCompletion
	err = json.Unmarshal(responseBody, &chatCompletion)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	fmt.Println("\n==============================================================")
	fmt.Println(chatCompletion.Choices[0].Message.Content)
	fmt.Println("==============================================================\n")
}
