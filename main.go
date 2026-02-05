package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
	if apiKey == "" {
		http.Error(w, "Config Error: API Key missing", 500)
		return
	}

	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	jsonData := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": "მოიძიე და ქართულად შეაჯამე ბოლო 1 საათის სიახლეები ევროპიდან. გამოიყენე პუნქტები.",
					},
				},
			},
		},
		"tools": []interface{}{
			map[string]interface{}{
				"google_search": map[string]interface{}{},
			},
		},
	}

	jsonBytes, _ := json.Marshal(jsonData)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, "Network Error", 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	currentTime := time.Now().Format("15:04")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="ka">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>ევროპის სიახლეები</title>
		<style>
			:root { --primary: #1a73e8; --bg: #f8f9fa; --card: #ffffff; }
			body { background-color: var(--bg); font-family: system-ui, -apple-system, sans-serif; margin: 0; padding: 15px; display: flex; justify-content: center; }
			.container { max-width: 600px; width: 100%%; background: var(--card); padding: 25px; border-radius: 20px; box-shadow: 0 10px 30px rgba(0,0,0,0.05); }
			.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; border-bottom: 2px solid #f0f2f5; padding-bottom: 15px; }
			h1 { background: linear-gradient(45deg, #1a73e8, #4285f4); -webkit-background-clip: text; -webkit-text-fill-color: transparent; font-size: 1.5rem; margin: 0; }
			.time-badge { background: #e8f0fe; color: var(--primary); padding: 5px 12px; border-radius: 50px; font-size: 0.8rem; font-weight: bold; }
			.content { white-space: pre-wrap; font-size: 1.05rem; line-height: 1.8; color: #3c4043; }
			.btn { display: block; width: 100%%; text-align: center; background: var(--primary); color: white; text-decoration: none; padding: 14px; border-radius: 12px; margin-top: 25px; font-weight: 600; box-sizing: border-box; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>🇪🇺 ევროპა</h1>
				<span class="time-badge">🕒 %s</span>
			</div>
			<div class="content">`, currentTime)

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		fmt.Fprint(w, geminiResp.Candidates[0].Content.Parts[0].Text)
	} else {
		fmt.Fprint(w, "სიახლეები ვერ მოიძებნა.")
	}

	fmt.Fprintf(w, `</div>
			<a href="/" class="btn">განახლება</a>
		</div>
	</body>
	</html>`)
}

func main() {
	// Render-ს სჭირდება პორტის წაკითხვა გარემო ცვლადიდან
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", Handler)
	fmt.Printf("სერვერი გაეშვა პორტზე %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
