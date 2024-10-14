package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"integration-service/proto/pb"
)

func uploadToGemini(ctx context.Context, client *genai.Client, path, mimeType string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	options := genai.UploadFileOptions{
		DisplayName: path,
		MIMEType:    mimeType,
	}
	fileData, err := client.UploadFile(ctx, "", file, &options)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	log.Printf("Uploaded file %s as: %s", fileData.DisplayName, fileData.URI)
	return fileData.URI
}

func processPartOfSpeaking(question string, message []byte) (*pb.SpeakingPartAbsResponse, error) {
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		return nil, fmt.Errorf("Environment variable GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type:     genai.TypeObject,
		Enum:     []string{},
		Required: []string{"fluency_score", "grammar_score", "vocabulary_score", "coherence_score", "topic_dev_score", "relevance_score", "word_count", "transcription", "part_band_score"},
		Properties: map[string]*genai.Schema{
			"fluency_score":    {Type: genai.TypeNumber},
			"grammar_score":    {Type: genai.TypeNumber},
			"vocabulary_score": {Type: genai.TypeNumber},
			"coherence_score":  {Type: genai.TypeNumber},
			"topic_dev_score":  {Type: genai.TypeNumber},
			"relevance_score":  {Type: genai.TypeNumber},
			"word_count":       {Type: genai.TypeNumber},
			"transcription":    {Type: genai.TypeString},
			"part_band_score":  {Type: genai.TypeNumber},
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
	fileURI := uploadToGemini(ctx, client, tempFile.Name(), "audio/mpeg")

	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.FileData{URI: fileURI},
				genai.Text(fmt.Sprintf("question : %s", question)),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("Analyze the audio and provide scores and transcription"))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}

	var result pb.SpeakingPartAbsResponse
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			jsonStr := text[8 : len(text)-4]
			err = json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
			}
			break
		}
	}

	return &result, nil
}
