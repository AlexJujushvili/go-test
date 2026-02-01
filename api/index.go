package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	prompt := "მოიძიე და ქართულად შეაჯამე ბოლო 1 საათის სიახლეები ევროპის რეგიონიდან. გამოიყენე პუნქტები."

	jsonData := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{"text": prompt},
				},
			},
		},
	}

	jsonBytes, _ := json.Marshal(jsonData)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, "API Error", 500)
		fmt.Println("სიახლეები ვერ მოიძებნა.")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<body style='background:#f0f2f5; font-family:sans-serif; padding:20px;'>")
	fmt.Fprintf(w, "<h1 style='color:#1a73e8;'>🇪🇺 ევროპის სიახლეები</h1>")

	if len(geminiResp.Candidates) > 0 {
		txt := geminiResp.Candidates[0].Content.Parts[0].Text
		fmt.Fprintf(w, "<div style='background:white; padding:20px; border-radius:10px;'>%s</div>", txt)
	} else {
		fmt.Fprintf(w, "სიახლეები ვერ მოიძებნა.")
		fmt.Println("სიახლეები ვერ მოიძებნა.")
	}
	fmt.Fprintf(w, "</body>")
}
