package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	prompt := "მოიძიე და ქართულად შეაჯამე ბოლო 1 საათის სიახლეები ევროპიდან. გამოიყენე პუნქტები (bullet points)."

	// სინტაქსურად გამართული JSON სტრუქტურა
	jsonData := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{"text": prompt},
				},
			},
		},
		"tools": []interface{}{
			map[string]interface{}{
				"google_search": map[string]interface{}{},
			},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		http.Error(w, "JSON Marshal Error", 500)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, "Network Error", 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API Error. რაღაც შეცდომა: %s\n", string(body))
		http.Error(w, "API returned error: "+resp.Status, resp.StatusCode)
		return
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		http.Error(w, "JSON Unmarshal Error", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<body style='background:#f0f2f5; font-family:sans-serif; padding:20px;'>")
	fmt.Fprintf(w, "<div style='max-width:800px; margin:auto; background:white; padding:30px; border-radius:12px;'>")
	fmt.Fprintf(w, "<h1 style='color:#1a73e8;'>🇪🇺 ევროპის სიახლეები</h1>")

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		txt := geminiResp.Candidates[0].Content.Parts[0].Text
		fmt.Fprintf(w, "<div style='white-space: pre-wrap; font-size: 16px; line-height: 1.6;'>%s</div>", txt)
	} else {
		fmt.Fprintf(w, "<p>სიახლეები ვერ მოიძებნა.</p>")
	}
	fmt.Fprintf(w, "</div></body>")
}
