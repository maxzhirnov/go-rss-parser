package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TranslationRequest struct {
	FolderID           string   `json:"folderId"`
	Texts              []string `json:"texts"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
}

type Translation struct {
	Text                 string `json:"text"`
	DetectedLanguageCode string `json:"detectedLanguageCode"`
}

type TranslationResponse struct {
	Translations []Translation `json:"translations"`
}

type Translator struct {
	FolderID string
	Token    string
}

func NewTranslator(folderID, token string) (*Translator, error) {
	return &Translator{
		FolderID: folderID,
		Token:    token,
	}, nil
}

func (t *Translator) Translate(text, targetLanguageCode string) (*TranslationResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("text is empty")
	}
	request := &TranslationRequest{
		FolderID:           t.FolderID,
		Texts:              []string{text},
		TargetLanguageCode: targetLanguageCode,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://translate.api.cloud.yandex.net/translate/v2/translate", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+t.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	translationResponse := &TranslationResponse{}
	err = json.Unmarshal(responseBody, translationResponse)
	if err != nil {
		return nil, err
	}

	return translationResponse, nil
}
