package handlers

import (
	"html/template"
	"math/rand"
	"net/http"
	"sync"
)

var (
	visitorCount = 0
	mutex        = &sync.Mutex{}
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func VisitorCounterHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	visitorCount++
	count := visitorCount
	mutex.Unlock()

	quotes := []Quote{
		{Text: "The only way to do great work is to love what you do.", Author: "Steve Jobs"},
		{Text: "Life is what happens when you're busy making other plans.", Author: "John Lennon"},
		{Text: "The future belongs to those who believe in the beauty of their dreams.", Author: "Eleanor Roosevelt"},
		{Text: "The best way to predict the future is to invent it.", Author: "Alan Kay"},
		{Text: "Code is like humor. When you have to explain it, it's bad.", Author: "Cory House"},
		{Text: "Simplicity is the soul of efficiency.", Author: "Austin Freeman"},
	}

	randomQuote := quotes[rand.Intn(len(quotes))]

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
        .quote { margin: 50px; padding: 20px; background: #f5f5f5; border-radius: 10px; }
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
