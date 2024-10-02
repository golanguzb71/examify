package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"integration-service/proto/pb"
)

func processEssay(essayText string) (*pb.WritingTaskAbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	apiKey := "AIzaSyBeQQyXZL0Duo-K36pDbTRM4EDi6thAMjo"
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
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

	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role:  "user",
			Parts: []genai.Part{genai.Text("Some people think that parents should teach their children how to be good members of society. Others, however, believe that school is the best place to learn this. Discuss both views and give your own opinion.\n\n[Sample essay content...]")},
		},
		{
			Role:  "model",
			Parts: []genai.Part{genai.Text("```json\n{\"coherence_score\": 6, \"feedback\": \"The essay has a clear structure and a well-defined thesis statement. The examples used to support the arguments are relevant and well-chosen. However, the essay could be improved by providing more specific examples and further developing the arguments. \", \"grammar_score\": 6, \"lexical_resource_score\": 6, \"task_achievement_score\": 6, \"task_band_score\": 6}\n\n```")},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text(essayText))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	log.Printf("Received response: %s", text)

	jsonStart := strings.Index(string(text), "{")
	jsonEnd := strings.LastIndex(string(text), "}") + 1
	if jsonStart < 0 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := string(text)[jsonStart:jsonEnd]
	log.Printf("Extracted JSON: %s", jsonStr)

	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &rawData); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	requiredFields := []string{"feedback", "coherence_score", "grammar_score", "lexical_resource_score", "task_achievement_score", "task_band_score"}
	for _, field := range requiredFields {
		if _, ok := rawData[field]; !ok {
			return nil, fmt.Errorf("required field '%s' is missing from the response", field)
		}
	}

	responseData := &pb.WritingTaskAbsResponse{}

	if feedback, ok := rawData["feedback"].(string); ok {
		responseData.Feedback = feedback
	} else {
		return nil, fmt.Errorf("feedback field is not a string")
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
			return nil, fmt.Errorf("%s field is not a number", field)
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
