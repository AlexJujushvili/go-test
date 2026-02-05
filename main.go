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

// GeminiResponse სტრუქტურა API-დან პასუხის მისაღებად
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
	// 1. API გასაღების წაკითხვა გარემო ცვლადიდან
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: GEMINI_API_KEY is missing")
		http.Error(w, "სერვერის კონფიგურაციის შეცდომა", 500)
		return
	}

	// 2. API URL (v1beta ვერსია Gemini 2.5 Flash-ისთვის)
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	// 3. JSON მოთხოვნა Google Search-ის ჩართვით
	jsonData := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": "მოიძიე და ქართულად შეაჯამე ბოლო 1 საათის სიახლეები ევროპიდან. გამოიყენე პუნქტები (bullet points).",
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
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Google API Error: %s\n", string(body))
		http.Error(w, "Gemini API Error", resp.StatusCode)
		return
	}

	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	// მიმდინარე დრო
	currentTime := time.Now().Format("15:04")

	// 4. HTML პასუხი მობილურისთვის (ანდროიდისთვის)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="ka">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>ევროპის სიახლეები</title>
		<style>
			:root { --primary: #1a73e8; --bg: #f1f3f4; }
			body { 
				background-color: var(--bg); 
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; 
				margin: 0; padding: 10px; display: flex; justify-content: center; 
			}
			.container { 
				max-width: 600px; width: 100%%; background: white; 
				padding: 20px; border-radius: 18px; 
				box-shadow: 0 10px 25px rgba(0,0,0,0.05); box-sizing: border-box; 
			}
			.header { 
				display: flex; justify-content: space-between; align-items: center; 
				margin-bottom: 20px; border-bottom: 2px solid #f0f2f5; padding-bottom: 12px; 
			}
			h1 { 
				background: linear-gradient(45deg, #1a73e8, #4285f4); 
				-webkit-background-clip: text; -webkit-text-fill-color: transparent; 
				font-size: 1.4rem; margin: 0; 
			}
			.time { background: #e8f0fe; color: var(--primary); padding: 4px 10px; border-radius: 20px; font-size: 0.8rem; font-weight: bold; }
			.content { white-space: pre-wrap; font-size: 16px; line-height: 1.7; color: #3c4043; }
			.refresh-btn { 
				display: block; width: 100%%; text-align: center; background: var(--primary); 
				color: white; text-decoration: none; padding: 14px; border-radius: 12px; 
				margin-top: 20px; font-weight: 600; box-sizing: border-box;
			}
			@media (max-width: 480px) {
				.container { padding: 15px; }
				h1 { font-size: 1.2rem; }
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>🇪🇺 ევროპა</h1>
				<span class="time">🕒 %s</span>
			</div>
			<div class="content">`, currentTime)

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		fmt.Fprint(w, geminiResp.Candidates[0].Content.Parts[0].Text)
	} else {
		fmt.Fprint(w, "სიახლეები ვერ მოიძებნა. სცადეთ მოგვიანებით.")
	}

	fmt.Fprintf(w, `</div>
			<a href="/" class="refresh-btn">განახლება</a>
		</div>
	</body>
	</html>`)
}

func main() {
	// Render-ისთვის აუცილებელია პორტის წაკითხვა გარემო ცვლადიდან
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ლოკალური ტესტისთვის
	}

	// მთავარი როუტი
	http.HandleFunc("/", Handler)

	fmt.Println("სერვერი გაეშვა პორტზე:", port)
	// აუცილებლად ":" + port
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("სერვერის გაშვების შეცდომა: %v\n", err)
	}
}
