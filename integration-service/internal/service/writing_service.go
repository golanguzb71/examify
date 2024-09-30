package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"integration-service/proto/pb"
	"strings"
	"time"
)

func processEssay(essayText string) (*pb.WritingTaskAbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	apiKey := "AIzaSyBeQQyXZL0Duo-K36pDbTRM4EDi6thAMjo"

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
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

	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("Some people think that parents should teach their children how to be good members of society. Others, however, believe that school is the best place to learn this. Discuss both views and give your own opinion.\n\n[Sample essay content...]"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\"coherence_score\": 6, \"feedback\": \"The essay has a clear structure and a well-defined thesis statement. The examples used to support the arguments are relevant and well-chosen. However, the essay could be improved by providing more specific examples and further developing the arguments. \", \"grammar_score\": 6, \"lexical_resource_score\": 6, \"task_achievement_score\": 6, \"task_band_score\": 6}\n\n```"),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text(essayText))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}

	var responseData pb.WritingTaskAbsResponse
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		text := resp.Candidates[0].Content.Parts[0].(genai.Text)

		jsonStart := strings.Index(string(text), "{")
		jsonEnd := strings.LastIndex(string(text), "}") + 1
		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonStr := string(text)[jsonStart:jsonEnd]
			var rawData map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &rawData)
			if err != nil {
				return nil, fmt.Errorf("error parsing response: %v", err)
			}

			responseData.Feedback = rawData["feedback"].(string)
			responseData.CoherenceScore = float32(clampScore(rawData["coherence_score"].(float64)))
			responseData.GrammarScore = float32(clampScore(rawData["grammar_score"].(float64)))
			responseData.LexicalResourceScore = float32(clampScore(rawData["lexical_resource_score"].(float64)))
			responseData.TaskAchievementScore = float32(clampScore(rawData["task_achievement_score"].(float64)))
			responseData.TaskBandScore = float32(clampScore(rawData["task_band_score"].(float64)))
		} else {
			return nil, fmt.Errorf("no valid JSON found in response")
		}
	} else {
		return nil, fmt.Errorf("no content in response")
	}

	return &responseData, nil
}

func clampScore(score float64) float64 {
	if score < 1 {
		return 1
	}
	if score > 7.5 {
		return 7.5
	}
	return score
}
