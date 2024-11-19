package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"integration-service/proto/pb"
)

const ApiKey = "AIzaSyBKuxWI1SxM0MXDjFCDWbCyj662-ydmHiE"

func createResponseData(rawData map[string]interface{}) (*pb.WritingTaskAbsResponse, error) {
	responseData := &pb.WritingTaskAbsResponse{}

	if feedback, ok := rawData["feedback"].(string); ok {
		responseData.Feedback = feedback
	} else {
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
			*ptr = float32(score)
		} else {
			*ptr = 0
		}
	}

	return responseData, nil
}
func processEssay(essayText string) (*pb.WritingTaskAbsResponse, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(ApiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	configureModel(model)

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(essayText))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return parseResponse(resp)
}

func configureModel(model *genai.GenerativeModel) {
	model.SetTemperature(1)
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

func parseResponse(resp *genai.GenerateContentResponse) (*pb.WritingTaskAbsResponse, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	var result pb.WritingTaskAbsResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// Clamp the scores to the range 0 to 7
	result.CoherenceScore = clamp(result.CoherenceScore, 0, 7)
	result.GrammarScore = clamp(result.GrammarScore, 0, 7)
	result.LexicalResourceScore = clamp(result.LexicalResourceScore, 0, 7)
	result.TaskAchievementScore = clamp(result.TaskAchievementScore, 0, 7)
	result.TaskBandScore = clamp(result.TaskBandScore, 0, 7)

	return &result, nil
}
