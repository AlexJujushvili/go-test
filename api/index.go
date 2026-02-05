package handler

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
		http.Error(w, "კონფიგურაციის შეცდომა", 500)
		return
	}

	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	jsonData := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": "მოიძიე და ქართულად შეაჯამე ბოლო 1 საათის სიახლეები ევროპიდან. გამოიყენე პუნქტები. იყავი კორექტული.",
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
		http.Error(w, "API Error", 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	// მიმდინარე დრო ქართული ფორმატით
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
			:root {
				--primary: #1a73e8;
				--bg: #f8f9fa;
				--card: #ffffff;
				--text: #202124;
				--secondary: #5f6368;
			}
			body { 
				background-color: var(--bg); 
				font-family: 'Segoe UI', Roboto, Helvetica, sans-serif; 
				margin: 0; 
				padding: 16px; 
				color: var(--text);
				display: flex;
				justify-content: center;
			}
			.container { 
				max-width: 650px; 
				width: 100%%; 
				background: var(--card); 
				padding: 24px; 
				border-radius: 20px; 
				box-shadow: 0 4px 20px rgba(0,0,0,0.08); 
				animation: fadeIn 0.6s ease-out;
			}
			@keyframes fadeIn {
				from { opacity: 0; transform: translateY(10px); }
				to { opacity: 1; transform: translateY(0); }
			}
			.header {
				display: flex;
				justify-content: space-between;
				align-items: center;
				margin-bottom: 20px;
				border-bottom: 1px solid #eee;
				padding-bottom: 15px;
			}
			h1 { 
				background: linear-gradient(45deg, #1a73e8, #8ab4f8);
				-webkit-background-clip: text;
				-webkit-text-fill-color: transparent;
				font-size: 1.6rem; 
				margin: 0;
			}
			.time-badge {
				background: #e8f0fe;
				color: var(--primary);
				padding: 4px 10px;
				border-radius: 20px;
				font-size: 0.85rem;
				font-weight: bold;
			}
			.news-content { 
				white-space: pre-wrap; 
				font-size: 1.05rem; 
				line-height: 1.8; 
				word-wrap: break-word;
			}
			/* პუნქტების (bullet points) სტილიზაცია */
			.news-content ul { padding-left: 20px; }
			.news-content li { margin-bottom: 10px; }
			
			.refresh-btn {
				display: block;
				width: 100%%;
				text-align: center;
				background: var(--primary);
				color: white;
				text-decoration: none;
				padding: 12px;
				border-radius: 12px;
				margin-top: 25px;
				font-weight: 500;
				transition: background 0.3s;
			}
			.refresh-btn:active { background: #174ea6; }

			@media (max-width: 480px) {
				.container { padding: 20px; border-radius: 16px; }
				h1 { font-size: 1.3rem; }
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>🇪🇺 ევროპა</h1>
				<span class="time-badge">🕒 %s</span>
			</div>
			<div class="news-content">`, currentTime)

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		fmt.Fprint(w, geminiResp.Candidates[0].Content.Parts[0].Text)
	} else {
		fmt.Fprint(w, "<p style='text-align:center; color:gray;'>სიახლეები ვერ მოიძებნა. სცადეთ მოგვიანებით.</p>")
	}

	fmt.Fprintf(w, `
			</div>
			<a href="/" class="refresh-btn">განახლება</a>
		</div>
	</body>
	</html>`)
}

func main() {
	http.HandleFunc("/", Handler)
	fmt.Println("სერვერი მზად არის: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
