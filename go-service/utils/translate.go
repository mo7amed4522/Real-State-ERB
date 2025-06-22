package utils

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

type TranslationRequest struct {
    Text        string   `json:"text"`
    SourceLang  string   `json:"source_lang"`
    TargetLangs []string `json:"target_langs"`
}

func TranslateText(text, sourceLang string, targetLangs []string) (map[string]string, error) {
    reqBody := TranslationRequest{
        Text:        text,
        SourceLang:  sourceLang,
        TargetLangs: targetLangs,
    }
    body, _ := json.Marshal(reqBody)
    resp, err := http.Post("http://python-translate-service:8010/translate", "application/json", bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)
    var translations map[string]string
    if err := json.Unmarshal(respBody, &translations); err != nil {
        return nil, err
    }
    return translations, nil
} 