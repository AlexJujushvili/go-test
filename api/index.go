package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Postgres დრაივერი
)

// Handler არის Vercel-ის მთავარი ფუნქცია
func Handler(w http.ResponseWriter, r *http.Request) {
	// 1. ავიღოთ ბაზის მისამართი Vercel-ის Environment Variables-იდან
	// (დარწმუნდი, რომ Vercel-ის პანელში დაამატე DATABASE_URL)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		http.Error(w, "DATABASE_URL არ არის მითითებული", http.StatusInternalServerError)
		return
	}

	// 2. დაკავშირება მონაცემთა ბაზასთან
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		http.Error(w, "ბაზასთან კავშირის შეცდომა", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 3. შევამოწმოთ კავშირი (Ping)
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(w, "Vercel მუშაობს, მაგრამ Neon-თან კავშირი ვერ დამყარდა: %v", err)
		return
	}

	// 4. მარტივი როუტინგი (Routing)
	path := r.URL.Path
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch path {
	case "/api/about":
		fmt.Fprintf(w, "ეს არის /about გვერდი Go-ზე! 🚀")
	default:
		fmt.Fprintf(w, "წარმატება! სერვერიც მუშაობს და Neon-ის ბაზაც დაკავშირებულია. ✅")
	}
}
