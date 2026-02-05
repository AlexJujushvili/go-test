package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		http.Error(w, "სერვერის შეცდომა: GEMINI_API_KEY არ არის კონფიგურირებული.", 500)
		return
	}

	// 2. API URL (Gemini 2.5 Flash - შენი სიის მიხედვით)
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	// 3. მოთხოვნის მომზადება (JSON)
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
		http.Error(w, "კავშირის შეცდომა API-სთან", 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API Error: %s\n", string(body))
		http.Error(w, "Gemini API-მ დააბრუნა შეცდომა", resp.StatusCode)
		return
	}

	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	// 4. HTML პასუხი მობილურზე მორგებული დიზაინით
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="ka">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>ევროპის სიახლეები</title>
		<style>
			body { 
				background-color: #f4f7f9; 
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; 
				margin: 0; 
				padding: 10px; 
				color: #333;
			}
			.container { 
				max-width: 600px; 
				margin: 20px auto; 
				background: white; 
				padding: 20px; 
				border-radius: 16px; 
				box-shadow: 0 10px 25px rgba(0,0,0,0.05); 
				box-sizing: border-box;
			}
			h1 { 
				color: #1a73e8; 
				font-size: 1.4rem; 
				margin-top: 0; 
				border-bottom: 2px solid #f0f2f5;
				padding-bottom: 12px;
				display: flex;
				align-items: center;
			}
			.news-content { 
				white-space: pre-wrap; 
				font-size: 16px; 
				line-height: 1.7; 
				word-wrap: break-word;
			}
			.footer {
				margin-top: 20px;
				font-size: 11px;
				color: #999;
				text-align: center;
			}
			@media (max-width: 480px) {
				body { padding: 5px; }
				.container { margin: 10px auto; border-radius: 12px; padding: 15px; }
				h1 { font-size: 1.2rem; }
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>🇪🇺 ევროპის სიახლეები</h1>
			<div class="news-content">`)

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		fmt.Fprint(w, geminiResp.Candidates[0].Content.Parts[0].Text)
	} else {
		fmt.Fprint(w, "ამ წამს სიახლეები ხელმისაწვდომი არ არის.")
	}

	fmt.Fprintf(w, `
			</div>
			<div class="footer">წყარო: Gemini 2.5 Flash • Real-time Search</div>
		</div>
	</body>
	</html>`)
}

func main() {
	http.HandleFunc("/", Handler)
	port := "8080"
	fmt.Println("სერვერი ჩაირთო პორტზე :" + port)
	http.ListenAndServe(":"+port, nil)
}
