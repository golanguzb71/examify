package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"integration-service/proto/pb"
)

func processSpeakingWithRetry(question string, audioData []byte) (*pb.SpeakingPartAbsResponse, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		response, err := processSpeaking(question, audioData)
		if err == nil {
			return response, nil
		}

		if shouldRetryS(err) {
			delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
			log.Printf("Attempt %d failed, retrying in %v: %v", attempt+1, delay, err)
			time.Sleep(delay)
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("max retries reached")
}

func shouldRetryS(err error) bool {
	if apiErr, ok := err.(*googleapi.Error); ok {
		return apiErr.Code == 429 || (apiErr.Code >= 500 && apiErr.Code < 600)
	}
	return false
}

func processSpeaking(question string, audioData []byte) (*pb.SpeakingPartAbsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	apiKey := "AIzaSyCDa-dcBGtOVdh4ClJuJg8jK4pvTP03T-E"
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")
	configureModelS(model)

	fileURI, err := uploadAudioToGemini(ctx, client, audioData)
	if err != nil {
		return nil, fmt.Errorf("error uploading audio: %w", err)
	}

	session := model.StartChat()
	session.History = getInitialChatHistoryS(question, fileURI)

	resp, err := session.SendMessage(ctx, genai.Text("Please provide valid JSON format with all required schema fields."))
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok && apiErr.Code == 429 {
			return nil, fmt.Errorf("rate limit exceeded: %w", err)
		}
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return parseResponseS(resp)
}

func configureModelS(model *genai.GenerativeModel) {
	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"fluency_score":    {Type: genai.TypeNumber},
			"grammar_score":    {Type: genai.TypeNumber},
			"vocabulary_score": {Type: genai.TypeNumber},
			"coherence_score":  {Type: genai.TypeNumber},
			"topic_dev_score":  {Type: genai.TypeNumber},
			"relevance_score":  {Type: genai.TypeNumber},
			"word_count":       {Type: genai.TypeInteger},
			"part_band_score":  {Type: genai.TypeNumber},
			"transcription": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"feedback":      {Type: genai.TypeString},
					"transcription": {Type: genai.TypeString},
				},
			},
		},
	}
}

func getInitialChatHistoryS(question string, fileURI string) []*genai.Content {
	return []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.FileData{URI: fileURI},
				genai.Text(fmt.Sprintf("Analyze the audio for the following question: %s. Please provide valid JSON format with all required schema fields. Provide a detailed analysis, including fluency_score, grammar_score, vocabulary_score, coherence_score, topic_dev_score, relevance_score, word_count, transcription, part_band_score, in valid JSON format. All scores should be between 1 and 8 and you should give feedback for giving answer more better. Do not leave any field null.", question)),
			},
		},
	}
}

func uploadAudioToGemini(ctx context.Context, client *genai.Client, audioData []byte) (string, error) {
	tempFile, err := os.CreateTemp("", "audio_*.mp3")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write(audioData); err != nil {
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}

	options := genai.UploadFileOptions{
		DisplayName: tempFile.Name(),
		MIMEType:    "audio/mpeg",
	}
	fileData, err := client.UploadFile(ctx, "", tempFile, &options)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %w", err)
	}

	log.Printf("Uploaded file %s as: %s", fileData.DisplayName, fileData.URI)
	return fileData.URI, nil
}

func parseResponseS(resp *genai.GenerateContentResponse) (*pb.SpeakingPartAbsResponse, error) {
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

	return createResponseDataS(rawData)
}

func extractJSONS(text string) string {
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}") + 1
	if jsonStart < 0 || jsonEnd <= jsonStart {
		return ""
	}
	return text[jsonStart:jsonEnd]
}

func createResponseDataS(rawData map[string]interface{}) (*pb.SpeakingPartAbsResponse, error) {
	responseData := &pb.SpeakingPartAbsResponse{}

	floatFields := map[string]*float32{
		"fluency_score":    &responseData.FluencyScore,
		"grammar_score":    &responseData.GrammarScore,
		"vocabulary_score": &responseData.VocabularyScore,
		"coherence_score":  &responseData.CoherenceScore,
		"topic_dev_score":  &responseData.TopicDevScore,
		"relevance_score":  &responseData.RelevanceScore,
		"part_band_score":  &responseData.PartBandScore,
	}

	for field, ptr := range floatFields {
		if score, ok := rawData[field].(float64); ok {
			*ptr = float32(clampScoreS(score))
		} else {
			log.Printf("Warning: '%s' field is missing or not a number", field)
			*ptr = 0
		}
	}

	if wordCount, ok := rawData["word_count"].(float64); ok {
		responseData.WordCount = int32(wordCount)
	} else {
		log.Printf("Warning: 'word_count' field is missing or not a number")
		responseData.WordCount = 0
	}

	if transcription, ok := rawData["transcription"].(map[string]interface{}); ok {
		responseData.Transcription = &pb.Transcription{}
		if feedback, ok := transcription["feedback"].(string); ok {
			responseData.Transcription.Feedback = feedback
		} else {
			log.Printf("Warning: 'feedback' field is missing or not a string")
			responseData.Transcription.Feedback = "No feedback provided"
		}
		if transcriptionText, ok := transcription["transcription"].(string); ok {
			responseData.Transcription.Transcription = transcriptionText
		} else {
			log.Printf("Warning: 'transcription' field is missing or not a string")
			responseData.Transcription.Transcription = "No transcription provided"
		}
	} else {
		log.Printf("Warning: 'transcription' field is missing or not an object")
		responseData.Transcription = &pb.Transcription{
			Feedback:      "No feedback provided",
			Transcription: "No transcription provided",
		}
	}

	return responseData, nil
}

func clampScoreS(score float64) float64 {
	if score < 1 {
		return 1
	}
	if score > 8 {
		return 8
	}
	return score
}
