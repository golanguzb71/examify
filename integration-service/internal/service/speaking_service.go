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

func processPartOfSpeaking(question string, message []byte) (*pb.SpeakingPartAbsResponse, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(ApiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")
	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
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

	// Create a temporary file to write the audio message
	tempFile, err := os.CreateTemp("", "audio_*.mp3")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write(message); err != nil {
		return nil, fmt.Errorf("error writing to temp file: %v", err)
	}

	// Upload the audio file to Gemini
	fileURI, err := uploadToGemini(ctx, client, tempFile.Name(), "audio/mpeg")
	if err != nil {
		return nil, err
	}

	// Send the message with the audio file and question
	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.FileData{URI: fileURI},
				genai.Text(fmt.Sprintf("Analyze the audio for the following question: %s.", question)),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("Please provide valid JSON format with all required schema fields."))
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
				result = pb.SpeakingPartAbsResponse{
					FluencyScore:    parsedResult.Response.FluencyScore,
					GrammarScore:    parsedResult.Response.GrammarScore,
					VocabularyScore: parsedResult.Response.VocabularyScore,
					CoherenceScore:  parsedResult.Response.CoherenceScore,
					TopicDevScore:   parsedResult.Response.TopicDevScore,
					RelevanceScore:  parsedResult.Response.RelevanceScore,
					WordCount:       parsedResult.Response.WordCount,
					PartBandScore:   parsedResult.Response.PartBandScore,
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
