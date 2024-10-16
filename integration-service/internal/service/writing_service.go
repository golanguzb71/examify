package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"integration-service/proto/pb"
)

func processEssayWithRetry(essayText string) (*pb.WritingTaskAbsResponse, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		response, err := processEssay(essayText)
		if err == nil {
			return response, nil
		}

		if shouldRetry(err) {
			delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
			log.Printf("Attempt %d failed, retrying in %v: %v", attempt+1, delay, err)
			time.Sleep(delay)
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("max retries reached")
}

func shouldRetry(err error) bool {
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		return apiErr.Code == 429 || (apiErr.Code >= 500 && apiErr.Code < 600)
	}
	return false
}

func processEssay(essayText string) (*pb.WritingTaskAbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	client, err := getClientWithRetry(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	configureModel(model)

	session := model.StartChat()
	session.History = getInitialChatHistory()

	resp, err := session.SendMessage(ctx, genai.Text(essayText))
	if err != nil {
		if shouldRetry(err) {
			return nil, fmt.Errorf("rate limit or server error: %w", err)
		}
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return parseResponse(resp)
}

func configureModel(model *genai.GenerativeModel) {
	model.SetTemperature(0.7)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"feedback":               {Type: genai.TypeString},
			"coherence_score":        {Type: genai.TypeNumber},
			"grammar_score":          {Type: genai.TypeNumber},
			"lexical_resource_score": {Type: genai.TypeNumber},
			"task_achievement_score": {Type: genai.TypeNumber},
			"task_band_score":        {Type: genai.TypeNumber},
		},
	}
}

func getInitialChatHistory() []*genai.Content {
	return []*genai.Content{
		{
			Role:  "user",
			Parts: []genai.Part{genai.Text("Some people think that parents should teach their children how to be good members of society. Others, however, believe that school is the best place to learn this. Discuss both views and give your own opinion.\n\n[Sample essay content...]")},
		},
		{
			Role:  "model",
			Parts: []genai.Part{genai.Text("```json\n{\"coherence_score\": 6, \"feedback\": \"The essay has a clear structure and a well-defined thesis statement. The examples used to support the arguments are relevant and well-chosen. However, the essay could be improved by providing more specific examples and further developing the arguments. \", \"grammar_score\": 6, \"lexical_resource_score\": 6, \"task_achievement_score\": 6, \"task_band_score\": 6}\n\n```")},
		},
	}
}

func parseResponse(resp *genai.GenerateContentResponse) (*pb.WritingTaskAbsResponse, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	log.Printf("Received response: %s", text)

	jsonStr := extractJSON(string(text))
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	log.Printf("Extracted JSON: %s", jsonStr)

	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &rawData); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return createResponseData(rawData)
}

func extractJSON(text string) string {
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}") + 1
	if jsonStart < 0 || jsonEnd <= jsonStart {
		return ""
	}
	return text[jsonStart:jsonEnd]
}

func createResponseData(rawData map[string]interface{}) (*pb.WritingTaskAbsResponse, error) {
	responseData := &pb.WritingTaskAbsResponse{}

	if feedback, ok := rawData["feedback"].(string); ok {
		responseData.Feedback = feedback
	} else {
		log.Printf("Warning: 'feedback' field is missing or not a string")
		responseData.Feedback = "No feedback provided"
	}

	floatFields := map[string]*float32{
		"coherence_score":        &responseData.CoherenceScore,
		"grammar_score":          &responseData.GrammarScore,
		"lexical_resource_score": &responseData.LexicalResourceScore,
		"task_achievement_score": &responseData.TaskAchievementScore,
		"task_band_score":        &responseData.TaskBandScore,
	}

	for field, ptr := range floatFields {
		if score, ok := rawData[field].(float64); ok {
			*ptr = float32(clampScore(score))
		} else {
			log.Printf("Warning: '%s' field is missing or not a number", field)
			*ptr = 0
		}
	}

	return responseData, nil
}

func clampScore(score float64) float64 {
	if score > 7.5 {
		return 7.5
	}
	return score
}

func getClientWithRetry(ctx context.Context) (*genai.Client, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		apiKey := getNextAPIKey()

		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err == nil {
			return client, nil
		}

		if !shouldRetry(err) {
			return nil, fmt.Errorf("error creating client with API key: %w", err)
		}

		delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
		log.Printf("Attempt %d to create client failed, retrying in %v: %v", attempt+1, delay, err)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("max retries reached for creating client")
}
