package stt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var providerDefaults = map[string][2]string{
	"groq":   {"https://api.groq.com/openai/v1", "whisper-large-v3-turbo"},
	"openai": {"https://api.openai.com/v1", "whisper-1"},
}

// AudioExtensions lists file extensions accepted for STT.
var AudioExtensions = map[string]bool{
	".ogg": true, ".mp3": true, ".wav": true, ".flac": true,
	".m4a": true, ".mp4": true, ".mpeg": true, ".mpga": true, ".webm": true,
}

// IsAudioFile returns true if the path has a known audio extension.
func IsAudioFile(path string) bool {
	return AudioExtensions[strings.ToLower(filepath.Ext(path))]
}

// Client calls an OpenAI-compatible /v1/audio/transcriptions endpoint.
type Client struct {
	baseURL  string
	apiKey   string
	model    string
	language string
	http     *http.Client
}

// New creates an STT client. Provider determines default baseURL and model.
func New(provider, apiKey, model, language string) *Client {
	defaults := providerDefaults[provider]
	if defaults == [2]string{} {
		defaults = providerDefaults["groq"]
	}
	baseURL := defaults[0]
	if model == "" {
		model = defaults[1]
	}
	return &Client{
		baseURL:  baseURL,
		apiKey:   apiKey,
		model:    model,
		language: language,
		http:     &http.Client{},
	}
}

type transcriptionResponse struct {
	Text string `json:"text"`
}

// Transcribe sends an audio file to the STT API and returns the transcribed text.
func (c *Client) Transcribe(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open audio: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, f); err != nil {
		return "", err
	}
	w.WriteField("model", c.model)
	w.WriteField("response_format", "json")
	if c.language != "" {
		w.WriteField("language", c.language)
	}
	w.Close()

	req, err := http.NewRequest("POST", c.baseURL+"/audio/transcriptions", &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("stt request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("stt %d: %s", resp.StatusCode, string(body))
	}

	var result transcriptionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("stt decode: %w", err)
	}
	return strings.TrimSpace(result.Text), nil
}
