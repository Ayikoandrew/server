package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	mathrand "math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	visitorCount = 0
	visitors     = make(map[string]bool)
	mutex        = &sync.Mutex{}
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func VisitorCounterHandler(w http.ResponseWriter, r *http.Request) {
	var count int
	sessionCookie, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		sessionID := generateSessionID()

		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)
		mutex.Lock()
		if _, exists := visitors[sessionID]; !exists {
			visitors[sessionID] = true
			visitorCount++
		}
		count = visitorCount
		mutex.Unlock()
	} else {
		sessionID := sessionCookie.Value

		mutex.Lock()
		count = visitorCount
		mutex.Unlock()
		if _, exists := visitors[sessionID]; !exists {
			mutex.Lock()
			visitors[sessionID] = true
			visitorCount++
			count = visitorCount
			mutex.Unlock()
		}
	}

	quotes := []Quote{
		{Text: "The only way to do great work is to love what you do.", Author: "Steve Jobs"},
		{Text: "Life is what happens when you're busy making other plans.", Author: "John Lennon"},
		{Text: "The best way to predict the future is to invent it.", Author: "Alan Kay"},
		{Text: "Code is like humor. When you have to explain it, it's bad.", Author: "Cory House"},
		{Text: "Simplicity is the soul of efficiency.", Author: "Austin Freeman"}}

	randomQuote := quotes[mathrand.Intn(len(quotes))]

	data := struct {
		Count int
		Quote Quote
	}{
		Count: count,
		Quote: randomQuote,
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Visitor Counter</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding-top: 100px; }
        .counter { font-size: 100px; margin: 50px; color: #333; }
        .message { font-size: 24px; color: #666; }
        .quote { margin: 50px; padding: 20px; background:rgb(104, 219, 51); border-radius: 10px; }
        .quote-text { font-style: italic; }
        .quote-author { color: #666; margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Welcome to my server! Hello there 👋</h1>
    <div class="counter">{{.Count}}</div>
    <p class="message">You are visitor number {{.Count}} to this page!</p>
    <div class="quote">
        <p class="quote-text">"{{.Quote.Text}}"</p>
        <p class="quote-author">- {{.Quote.Author}}</p>
    </div>
</body>
</html>`

	t, _ := template.New("counter").Parse(tmpl)
	t.Execute(w, data)
}
