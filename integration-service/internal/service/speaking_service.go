package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"integration-service/proto/pb"
	"io"
	"os"
)

// Define the clamp function to ensure the score stays within the desired range (0 to 7)
func clamp(value float32, min, max float32) float32 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func processPartOfSpeaking(question string, message []byte) (*pb.SpeakingPartAbsResponse, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(ApiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")
	model.SetTemperature(10)
	model.SetTopK(1000)
	model.SetTopP(10)
	model.SetMaxOutputTokens(100000)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type:     genai.TypeObject,
		Required: []string{"response"},
		Properties: map[string]*genai.Schema{
			"response": {
				Type:     genai.TypeObject,
				Required: []string{"fluency_score", "grammar_score", "vocabulary_score", "coherence_score", "topic_dev_score", "relevance_score", "word_count", "transcription", "part_band_score"},
				Properties: map[string]*genai.Schema{
					"fluency_score":    {Type: genai.TypeNumber},
					"grammar_score":    {Type: genai.TypeNumber},
					"vocabulary_score": {Type: genai.TypeNumber},
					"coherence_score":  {Type: genai.TypeNumber},
					"topic_dev_score":  {Type: genai.TypeNumber},
					"relevance_score":  {Type: genai.TypeNumber},
					"word_count":       {Type: genai.TypeInteger},
					"transcription": {
						Type:     genai.TypeObject,
						Required: []string{"feedback", "transcription"},
						Properties: map[string]*genai.Schema{
							"feedback":      {Type: genai.TypeString},
							"transcription": {Type: genai.TypeString},
						},
					},
					"part_band_score": {Type: genai.TypeNumber},
				},
			},
		},
	}

	tempFile, err := os.CreateTemp("", "audio_*.mp3")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write(message); err != nil {
		return nil, fmt.Errorf("error writing to temp file: %v", err)
	}

	fileURI, err := uploadToGemini(ctx, client, tempFile.Name(), "audio/mpeg")
	if err != nil {
		return nil, err
	}

	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.FileData{URI: fileURI},
				genai.Text(fmt.Sprintf("Analyze the audio for the following question as IELTS speaking score 0 to 7, give exactly score as IELTS: %s.", question)),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("Please provide valid JSON format with all required schema fields as IELTS band scores for speaking. Rate 0 to 7, don't increase beyond this"))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	var result pb.SpeakingPartAbsResponse

	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			jsonStr := string(text)

			var parsedResult struct {
				Response struct {
					FluencyScore    float32 `json:"fluency_score"`
					GrammarScore    float32 `json:"grammar_score"`
					VocabularyScore float32 `json:"vocabulary_score"`
					CoherenceScore  float32 `json:"coherence_score"`
					TopicDevScore   float32 `json:"topic_dev_score"`
					RelevanceScore  float32 `json:"relevance_score"`
					WordCount       int32   `json:"word_count"`
					PartBandScore   float32 `json:"part_band_score"`
					Transcription   struct {
						Feedback      string `json:"feedback"`
						Transcription string `json:"transcription"`
					} `json:"transcription"`
				} `json:"response"`
			}

			if err := json.Unmarshal([]byte(jsonStr), &parsedResult); err == nil {
				// Apply clamp to ensure scores are between 0 and 7
				result = pb.SpeakingPartAbsResponse{
					FluencyScore:    clamp(parsedResult.Response.FluencyScore, 0, 7),
					GrammarScore:    clamp(parsedResult.Response.GrammarScore, 0, 7),
					VocabularyScore: clamp(parsedResult.Response.VocabularyScore, 0, 7),
					CoherenceScore:  clamp(parsedResult.Response.CoherenceScore, 0, 7),
					TopicDevScore:   clamp(parsedResult.Response.TopicDevScore, 0, 7),
					RelevanceScore:  clamp(parsedResult.Response.RelevanceScore, 0, 7),
					WordCount:       parsedResult.Response.WordCount,
					PartBandScore:   clamp(parsedResult.Response.PartBandScore, 0, 7),
					Transcription: &pb.Transcription{
						Question:      question,
						Feedback:      parsedResult.Response.Transcription.Feedback,
						Transcription: parsedResult.Response.Transcription.Transcription,
					},
				}
				break
			}
		}
	}

	return &result, nil
}

func uploadToGemini(ctx context.Context, client *genai.Client, path, mimeType string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	options := &genai.UploadFileOptions{
		DisplayName: path,
		MIMEType:    mimeType,
	}

	resp, err := client.UploadFile(ctx, "", bytes.NewReader(fileContent), options)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}

	return resp.URI, nil
}
