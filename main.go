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
	// ეს ლოგი დაგვეხმარება დავინახოთ, შემოვიდა თუ არა მოთხოვნა
	fmt.Printf("მოთხოვნა შემოვიდა მისამართზე: %s\n", r.URL.Path)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		http.Error(w, "Error: GEMINI_API_KEY missing", 500)
		return
	}

	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	jsonData := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{"text": "შეაჯამე ბოლო 1 საათის სიახლეები ევროპიდან ქართულად."},
				},
			},
		},
		"tools": []interface{}{
			map[string]interface{}{"google_search": map[string]interface{}{}},
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
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body { background: #f0f2f5; font-family: sans-serif; padding: 20px; display: flex; justify-content: center; }
			.card { max-width: 600px; width: 100%%; background: white; padding: 20px; border-radius: 15px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
			h1 { color: #1a73e8; font-size: 1.4rem; border-bottom: 1px solid #eee; padding-bottom: 10px; }
			.text { white-space: pre-wrap; line-height: 1.6; }
			.btn { display: block; text-align: center; background: #1a73e8; color: white; text-decoration: none; padding: 12px; border-radius: 10px; margin-top: 20px; }
		</style>
	</head>
	<body>
		<div class="card">
			<h1>🇪🇺 ევროპა (%s)</h1>
			<div class="text">`, currentTime)

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// მნიშვნელოვანია: გამოიყენე http.DefaultServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	fmt.Printf("სერვერი გაეშვა პორტზე: %s\n", port)
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("სერვერის შეცდომა: %v\n", err)
	}
}
