package server

import (
	"html/template"
	"net/http"

	"main/pkg/game"
)

type Game struct {
	Word    string
	Letters []string
	Wrong   int
}

var (
	gameState = make(map[string]*Game)
)

func CreateServer() *http.Server {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sessionID := "player1"

		if gameState[sessionID] == nil {
			word := game.Getword("easy")
			letters := make([]string, len(word))
			for i := range word {
				letters[i] = "_"
			}
			gameState[sessionID] = &Game{
				Word:    word,
				Letters: letters,
				Wrong:   0,
			}
		}

		g := gameState[sessionID]

		if r.Method == http.MethodPost {
			guess := r.FormValue("letter")
			correct := false
			for i, ch := range g.Word {
				if string(ch) == guess {
					g.Letters[i] = guess
					correct = true
				}
			}
			if !correct {
				g.Wrong++
			}
		}

		tmpl := template.Must(template.ParseFiles("./static/index.html"))
		tmpl.Execute(w, g)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
}
