package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"main/pkg/game"
)

type Game struct {
	Word     string
	Letters  []string
	Wrong    int
	Alphabet []string
	Corrects int
	Finished bool
	Message  string
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
			word := strings.ToLower(game.Getword("easy"))
			letters := make([]string, len(word))
			for i := range word {
				letters[i] = "_"
			}

			alphabet := []string{
				"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
				"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			}

			gameState[sessionID] = &Game{
				Word:     word,
				Letters:  letters,
				Wrong:    0,
				Alphabet: alphabet,
			}
		}

		g := gameState[sessionID]

		if r.Method == http.MethodPost {
			guess := strings.ToLower(r.FormValue("letter"))
			correct := false
			for i, ch := range g.Word {
				if string(ch) == guess {
					g.Letters[i] = guess
					correct = true
					g.Corrects++
				}
			}
			if !correct {
				g.Wrong++
			}
			for i, letter := range g.Alphabet {
				if letter == guess {
					g.Alphabet = append(g.Alphabet[:i], g.Alphabet[i+1:]...)
					break
				}
			}

			if g.Corrects == len(g.Word) {
				fmt.Printf("You win! The word was %s", g.Word)
				g.Finished = true
			}

			if g.Wrong == 8 {
				fmt.Printf("You lose! The word was %s", g.Word)
				g.Finished = true
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
