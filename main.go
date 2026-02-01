package main

import (
	"log"
	handler "my-gemini-app/api" // ყურადღება: გამოიყენე შენი go.mod-ის სახელი
	"net/http"
)

func main() {
	// დაუკავშირე შენი Handler მთავარ მისამართს
	http.HandleFunc("/", handler.Handler)

	log.Println("სერვერი გაეშვა: http://localhost:8080")
	// სერვერის გაშვება
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
